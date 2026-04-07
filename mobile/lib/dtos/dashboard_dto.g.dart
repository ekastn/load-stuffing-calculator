// GENERATED CODE - DO NOT MODIFY BY HAND

part of 'dashboard_dto.dart';

// **************************************************************************
// JsonSerializableGenerator
// **************************************************************************

DashboardStatsDto _$DashboardStatsDtoFromJson(Map<String, dynamic> json) =>
    DashboardStatsDto(
      admin: json['admin'] == null
          ? null
          : AdminStatsDto.fromJson(json['admin'] as Map<String, dynamic>),
      planner: json['planner'] == null
          ? null
          : PlannerStatsDto.fromJson(json['planner'] as Map<String, dynamic>),
      operator: json['operator'] == null
          ? null
          : OperatorStatsDto.fromJson(json['operator'] as Map<String, dynamic>),
    );

Map<String, dynamic> _$DashboardStatsDtoToJson(DashboardStatsDto instance) =>
    <String, dynamic>{
      'admin': instance.admin,
      'planner': instance.planner,
      'operator': instance.operator,
    };

AdminStatsDto _$AdminStatsDtoFromJson(Map<String, dynamic> json) =>
    AdminStatsDto(
      totalUsers: (json['total_users'] as num).toInt(),
      activeShipments: (json['active_shipments'] as num).toInt(),
      containerTypes: (json['container_types'] as num).toInt(),
      successRate: (json['success_rate'] as num).toDouble(),
    );

Map<String, dynamic> _$AdminStatsDtoToJson(AdminStatsDto instance) =>
    <String, dynamic>{
      'total_users': instance.totalUsers,
      'active_shipments': instance.activeShipments,
      'container_types': instance.containerTypes,
      'success_rate': instance.successRate,
    };

PlannerStatsDto _$PlannerStatsDtoFromJson(Map<String, dynamic> json) =>
    PlannerStatsDto(
      pendingPlans: (json['pending_plans'] as num).toInt(),
      completedToday: (json['completed_today'] as num).toInt(),
      avgUtilization: (json['avg_utilization'] as num).toDouble(),
      itemsProcessed: (json['items_processed'] as num).toInt(),
    );

Map<String, dynamic> _$PlannerStatsDtoToJson(PlannerStatsDto instance) =>
    <String, dynamic>{
      'pending_plans': instance.pendingPlans,
      'completed_today': instance.completedToday,
      'avg_utilization': instance.avgUtilization,
      'items_processed': instance.itemsProcessed,
    };

OperatorStatsDto _$OperatorStatsDtoFromJson(Map<String, dynamic> json) =>
    OperatorStatsDto(
      activeLoads: (json['active_loads'] as num).toInt(),
      completed: (json['completed'] as num).toInt(),
      failedValidations: (json['failed_validations'] as num).toInt(),
      avgTimePerLoad: json['avg_time_per_load'] as String,
    );

Map<String, dynamic> _$OperatorStatsDtoToJson(OperatorStatsDto instance) =>
    <String, dynamic>{
      'active_loads': instance.activeLoads,
      'completed': instance.completed,
      'failed_validations': instance.failedValidations,
      'avg_time_per_load': instance.avgTimePerLoad,
    };
