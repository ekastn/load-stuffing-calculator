import 'package:json_annotation/json_annotation.dart';

part 'dashboard_dto.g.dart';

@JsonSerializable()
class DashboardStatsDto {
  final AdminStatsDto? admin;
  final PlannerStatsDto? planner;
  final OperatorStatsDto? operator;

  DashboardStatsDto({
    this.admin,
    this.planner,
    this.operator,
  });

  factory DashboardStatsDto.fromJson(Map<String, dynamic> json) =>
      _$DashboardStatsDtoFromJson(json);
}

@JsonSerializable()
class AdminStatsDto {
  @JsonKey(name: 'total_users')
  final int totalUsers;
  @JsonKey(name: 'active_shipments')
  final int activeShipments;
  @JsonKey(name: 'container_types')
  final int containerTypes;
  @JsonKey(name: 'success_rate')
  final double successRate;

  AdminStatsDto({
    required this.totalUsers,
    required this.activeShipments,
    required this.containerTypes,
    required this.successRate,
  });

  factory AdminStatsDto.fromJson(Map<String, dynamic> json) =>
      _$AdminStatsDtoFromJson(json);
}

@JsonSerializable()
class PlannerStatsDto {
  @JsonKey(name: 'pending_plans')
  final int pendingPlans;
  @JsonKey(name: 'completed_today')
  final int completedToday;
  @JsonKey(name: 'avg_utilization')
  final double avgUtilization;
  @JsonKey(name: 'items_processed')
  final int itemsProcessed;

  PlannerStatsDto({
    required this.pendingPlans,
    required this.completedToday,
    required this.avgUtilization,
    required this.itemsProcessed,
  });

  factory PlannerStatsDto.fromJson(Map<String, dynamic> json) =>
      _$PlannerStatsDtoFromJson(json);
}

@JsonSerializable()
class OperatorStatsDto {
  @JsonKey(name: 'active_loads')
  final int activeLoads;
  @JsonKey(name: 'completed')
  final int completed;
  @JsonKey(name: 'failed_validations')
  final int failedValidations;
  @JsonKey(name: 'avg_time_per_load')
  final String avgTimePerLoad;

  OperatorStatsDto({
    required this.activeLoads,
    required this.completed,
    required this.failedValidations,
    required this.avgTimePerLoad,
  });

  factory OperatorStatsDto.fromJson(Map<String, dynamic> json) =>
      _$OperatorStatsDtoFromJson(json);
}
