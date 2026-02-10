import 'dart:convert';
import 'package:flutter/foundation.dart';
import '../models/loading_session.dart';
import '../models/plan_detail_model.dart';
import '../services/plan_service.dart';
import '../services/storage_service.dart';

class LoadingProvider extends ChangeNotifier {
  final PlanService _planService;
  final StorageService _storage;

  LoadingProvider(this._planService, this._storage);

  LoadingSession? _currentSession;
  LoadingSession? get currentSession => _currentSession;

  PlanDetailModel? _currentPlan;
  List<PlacementDetail>? _placements; // Sorted by step_number

  bool _isLoading = false;
  bool get isLoading => _isLoading;

  String? _error;
  String? get error => _error;

  int get progressPercentage =>
      _currentSession == null || _currentSession!.totalItems == 0
      ? 0
      : (_currentSession!.validatedCount / _currentSession!.totalItems * 100)
            .round();

  // Start new loading session (100% local initially)
  Future<void> startSession(String planId) async {
    _isLoading = true;
    _error = null;
    notifyListeners();

    try {
      // 1. Fetch plan data from backend
      _currentPlan = await _planService.getPlanDetail(planId);

      if (_currentPlan?.calculation?.placements == null ||
          _currentPlan!.calculation!.placements!.isEmpty) {
        throw Exception("Plan has no placements to load");
      }

      _placements = List.from(_currentPlan!.calculation!.placements!)
        ..sort((a, b) => a.stepNumber.compareTo(b.stepNumber));

      // 2. Create local session
      _currentSession = LoadingSession(
        sessionId: DateTime.now().millisecondsSinceEpoch
            .toString(), // Simplified ID
        planId: planId,
        startedAt: DateTime.now(),
        totalItems: _placements!.length,
        currentStepIndex: 0,
        status: 'IN_PROGRESS',
      );

      // 3. Save to local storage
      await _saveSessionToStorage();

      _error = null;
    } catch (e) {
      _error = e.toString();
      _currentSession = null;
      _currentPlan = null;
      _placements = null;
    } finally {
      _isLoading = false;
      notifyListeners();
    }
  }

  // Resume existing session from storage
  Future<void> resumeSession() async {
    _isLoading = true;
    notifyListeners();

    try {
      final sessionJson = await _storage.getLoadingSession();
      if (sessionJson != null) {
        _currentSession = LoadingSession.fromJson(jsonDecode(sessionJson));

        // Reload plan data
        _currentPlan = await _planService.getPlanDetail(
          _currentSession!.planId,
        );
        _placements = List.from(_currentPlan!.calculation!.placements!)
          ..sort((a, b) => a.stepNumber.compareTo(b.stepNumber));
      }
    } catch (e) {
      _error = "Failed to resume session: $e";
    } finally {
      _isLoading = false;
      notifyListeners();
    }
  }

  // Get current expected item
  ExpectedItem? getCurrentExpectedItem() {
    if (_currentSession == null ||
        _placements == null ||
        _currentPlan == null) {
      return null;
    }
    if (_currentSession!.currentStepIndex >= _placements!.length) return null;

    final placement = _placements![_currentSession!.currentStepIndex];
    final item = _currentPlan!.items.firstWhere(
      (i) => i.itemId == placement.itemId,
      orElse: () => throw Exception("Item not found: ${placement.itemId}"),
    );

    final generatedBarcode = generateBarcode(
      _currentSession!.planId,
      placement.stepNumber,
      item.itemId,
    );

    return ExpectedItem(
      itemId: item.itemId,
      itemLabel: item.label ?? item.productSku ?? 'Item ${item.itemId}',
      generatedBarcode: generatedBarcode,
      stepNumber: placement.stepNumber,
      position:
          '(${placement.posX.round()}, ${placement.posY.round()}, ${placement.posZ.round()})',
      dimensions:
          '${item.lengthMm.round()} × ${item.widthMm.round()} × ${item.heightMm.round()} mm',
    );
  }

  // Generate barcode for item (matches backend logic)
  String generateBarcode(String planId, int stepNumber, String itemId) {
    final planShort = planId.length > 8 ? planId.substring(0, 8) : planId;
    final itemShort = itemId.length > 8 ? itemId.substring(0, 8) : itemId;
    final stepPadded = stepNumber.toString().padLeft(3, '0');
    return 'PLAN-$planShort-STEP-$stepPadded-$itemShort';
  }

  // Parse scanned barcode
  ParsedBarcode? parseBarcode(String barcode) {
    final parts = barcode.split('-');
    if (parts.length != 5 || parts[0] != 'PLAN' || parts[2] != 'STEP') {
      return null;
    }

    final stepNumber = int.tryParse(parts[3]);
    if (stepNumber == null) return null;

    return ParsedBarcode(
      planIdShort: parts[1],
      stepNumber: stepNumber,
      itemIdShort: parts[4],
    );
  }

  // Validate scanned barcode (client-side)
  ValidationResult validateBarcode(String scannedBarcode) {
    final expectedItem = getCurrentExpectedItem();
    if (expectedItem == null) {
      return ValidationResult(matched: false, error: 'No more items to load');
    }

    // Parse scanned barcode
    final parsed = parseBarcode(scannedBarcode);
    if (parsed == null) {
      return ValidationResult(
        matched: false,
        error: 'Invalid QR code format',
        status: 'INVALID_FORMAT',
      );
    }

    // Validate
    final planShort = _currentSession!.planId.length > 8
        ? _currentSession!.planId.substring(0, 8)
        : _currentSession!.planId;
    final itemShort = expectedItem.itemId.length > 8
        ? expectedItem.itemId.substring(0, 8)
        : expectedItem.itemId;

    bool matched = false;
    String status = 'MISMATCHED';

    if (parsed.planIdShort != planShort) {
      status = 'WRONG_PLAN';
    } else if (parsed.stepNumber == expectedItem.stepNumber &&
        parsed.itemIdShort == itemShort) {
      matched = true;
      status = 'MATCHED';
    } else if (parsed.stepNumber != expectedItem.stepNumber) {
      status = 'OUT_OF_SEQUENCE';
    }

    final validation = ItemValidation(
      itemId: expectedItem.itemId,
      itemLabel: expectedItem.itemLabel,
      stepNumber: expectedItem.stepNumber,
      scannedBarcode: scannedBarcode,
      expectedBarcode: expectedItem.generatedBarcode,
      status: status,
      validatedAt: DateTime.now(),
    );

    if (matched) {
      _moveToNextItem(validation);
    }

    return ValidationResult(
      matched: matched,
      status: status,
      expectedItem: expectedItem,
      validation: validation,
    );
  }

  // Manual confirmation
  void manualConfirm({String? notes}) {
    final expectedItem = getCurrentExpectedItem();
    if (expectedItem == null) return;

    final validation = ItemValidation(
      itemId: expectedItem.itemId,
      itemLabel: expectedItem.itemLabel,
      stepNumber: expectedItem.stepNumber,
      expectedBarcode: expectedItem.generatedBarcode,
      status: 'MANUAL_CONFIRMED',
      validatedAt: DateTime.now(),
      notes: notes,
    );

    _moveToNextItem(validation);
  }

  // Skip item
  void skipItem({String? reason}) {
    final expectedItem = getCurrentExpectedItem();
    if (expectedItem == null) return;

    final validation = ItemValidation(
      itemId: expectedItem.itemId,
      itemLabel: expectedItem.itemLabel,
      stepNumber: expectedItem.stepNumber,
      expectedBarcode: expectedItem.generatedBarcode,
      status: 'SKIPPED',
      validatedAt: DateTime.now(),
      notes: reason,
    );

    _currentSession = _currentSession!.copyWith(
      skippedCount: _currentSession!.skippedCount + 1,
      currentStepIndex: _currentSession!.currentStepIndex + 1,
      validations: [..._currentSession!.validations, validation],
    );

    _saveSessionToStorage();
    notifyListeners();
  }

  void _moveToNextItem(ItemValidation validation) {
    _currentSession = _currentSession!.copyWith(
      validatedCount: _currentSession!.validatedCount + 1,
      currentStepIndex: _currentSession!.currentStepIndex + 1,
      validations: [..._currentSession!.validations, validation],
    );

    _saveSessionToStorage();
    notifyListeners();
  }

  Future<void> completeSession({bool syncToBackend = true}) async {
    if (_currentSession == null) return;

    _currentSession = _currentSession!.copyWith(
      status: 'COMPLETED',
      completedAt: DateTime.now(),
    );

    if (syncToBackend) {
      await _syncToBackend();
    }

    await _storage.deleteLoadingSession();

    _currentSession = null;
    _currentPlan = null;
    _placements = null;

    notifyListeners();
  }

  Future<void> _syncToBackend() async {
    if (_currentSession == null || _currentPlan == null) return;

    final loadingNotes = {
      'session_id': _currentSession!.sessionId,
      'started_at': _currentSession!.startedAt.toIso8601String(),
      'completed_at': _currentSession!.completedAt?.toIso8601String(),
      'total_items': _currentSession!.totalItems,
      'validated_count': _currentSession!.validatedCount,
      'skipped_count': _currentSession!.skippedCount,
      'validations': _currentSession!.validations
          .map((v) => v.toJson())
          .toList(),
    };

    try {
      await _planService.updatePlan(
        _currentPlan!.id,
        status: 'COMPLETED',
        loadingNotes: loadingNotes,
      );
    } catch (e) {
      debugPrint('Failed to sync loading notes: $e');
    }
  }

  Future<void> _saveSessionToStorage() async {
    if (_currentSession != null) {
      await _storage.saveLoadingSession(jsonEncode(_currentSession!.toJson()));
    }
  }

  void clearError() {
    _error = null;
    notifyListeners();
  }
}

class ExpectedItem {
  final String itemId;
  final String itemLabel;
  final String generatedBarcode;
  final int stepNumber;
  final String position;
  final String dimensions;

  ExpectedItem({
    required this.itemId,
    required this.itemLabel,
    required this.generatedBarcode,
    required this.stepNumber,
    required this.position,
    required this.dimensions,
  });
}

class ValidationResult {
  final bool matched;
  final String status;
  final ExpectedItem? expectedItem;
  final ItemValidation? validation;
  final String? error;

  ValidationResult({
    required this.matched,
    this.status = '',
    this.expectedItem,
    this.validation,
    this.error,
  });
}

class ParsedBarcode {
  final String planIdShort;
  final int stepNumber;
  final String itemIdShort;

  ParsedBarcode({
    required this.planIdShort,
    required this.stepNumber,
    required this.itemIdShort,
  });
}
