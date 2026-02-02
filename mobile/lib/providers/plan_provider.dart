import 'package:flutter/material.dart';
import '../models/plan_model.dart';
import '../services/plan_service.dart';
import '../dtos/plan_dto.dart';

class PlanProvider with ChangeNotifier {
  final PlanService _service;

  PlanProvider(this._service);

  List<PlanModel> _plans = [];
  bool _isLoading = false;
  String? _error;

  List<PlanModel> get plans => _plans;
  bool get isLoading => _isLoading;
  String? get error => _error;

  Future<void> fetchPlans() async {
    _isLoading = true;
    _error = null;
    notifyListeners();

    try {
      _plans = await _service.getPlans(page: 1, limit: 100);
      _error = null;
    } catch (e) {
      _error = e.toString();
      _plans = [];
    } finally {
      _isLoading = false;
      notifyListeners();
    }
  }

  Future<void> createPlan(CreatePlanRequestDto request) async {
    _isLoading = true;
    notifyListeners();

    try {
      final newPlan = await _service.createPlan(request);
      _plans.insert(0, newPlan); // Add to beginning
      _error = null;
    } catch (e) {
      _error = e.toString();
      rethrow;
    } finally {
      _isLoading = false;
      notifyListeners();
    }
  }

  Future<void> deletePlan(String planId) async {
    try {
      await _service.deletePlan(planId);
      _plans.removeWhere((p) => p.id == planId);
      _error = null;
      notifyListeners();
    } catch (e) {
      _error = e.toString();
      notifyListeners();
      rethrow;
    }
  }
}
