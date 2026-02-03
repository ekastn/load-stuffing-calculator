import 'package:flutter/foundation.dart';
import '../models/plan_detail_model.dart';
import '../services/plan_service.dart';

class PlanDetailProvider extends ChangeNotifier {
  final PlanService _planService;

  PlanDetailProvider(this._planService);

  PlanDetailModel? _plan;
  PlanDetailModel? get plan => _plan;

  bool _isLoading = false;
  bool get isLoading => _isLoading;

  bool _isCalculating = false;
  bool get isCalculating => _isCalculating;

  String? _error;
  String? get error => _error;

  Future<void> fetchPlanDetail(String planId) async {
    _isLoading = true;
    _error = null;
    notifyListeners();

    try {
      _plan = await _planService.getPlanDetail(planId);
      _error = null;
    } catch (e) {
      _error = e.toString();
      _plan = null;
    } finally {
      _isLoading = false;
      notifyListeners();
    }
  }

  Future<void> recalculate(
    String planId, {
    String strategy = 'bestfitdecreasing',
    String? goal,
    bool gravity = true,
  }) async {
    if (_isCalculating) return; // Prevent duplicate calls

    _isCalculating = true;
    _error = null;
    notifyListeners();

    try {
      _plan = await _planService.recalculate(
        planId,
        strategy: strategy,
        goal: goal,
        gravity: gravity,
      );
      _error = null;
    } catch (e) {
      _error = e.toString();
    } finally {
      _isCalculating = false;
      notifyListeners();
    }
  }

  void clearError() {
    _error = null;
    notifyListeners();
  }

  void reset() {
    _plan = null;
    _isLoading = false;
    _isCalculating = false;
    _error = null;
    notifyListeners();
  }
}
