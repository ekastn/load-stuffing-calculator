import 'package:flutter/material.dart';
import '../dtos/dashboard_dto.dart';
import '../models/plan_model.dart';
import '../services/dashboard_service.dart';
import '../services/plan_service.dart';

class DashboardProvider extends ChangeNotifier {
  final DashboardService _dashboardService;
  final PlanService _planService;

  DashboardProvider(this._dashboardService, this._planService);

  DashboardStatsDto? _stats;
  List<PlanModel> _recentPlans = [];
  bool _isLoading = false;
  String? _error;

  DashboardStatsDto? get stats => _stats;
  List<PlanModel> get recentPlans => _recentPlans;
  bool get isLoading => _isLoading;
  String? get error => _error;

  Future<void> loadData() async {
    _isLoading = true;
    _error = null;
    notifyListeners();

    try {
      final results = await Future.wait([
        _dashboardService.getStats(),
        _planService.getPlans(limit: 5),
      ]);

      _stats = results[0] as DashboardStatsDto;
      _recentPlans = results[1] as List<PlanModel>;
    } catch (e) {
      _error = e.toString();
      // Keep previous data if refresh fails? Or clear? 
      // For now we keep it visible or let UI handle error overlay.
    } finally {
      _isLoading = false;
      notifyListeners();
    }
  }
}
