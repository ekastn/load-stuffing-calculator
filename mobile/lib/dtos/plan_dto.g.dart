// GENERATED CODE - DO NOT MODIFY BY HAND

part of 'plan_dto.dart';

// **************************************************************************
// JsonSerializableGenerator
// **************************************************************************

PlanListItemDto _$PlanListItemDtoFromJson(Map<String, dynamic> json) =>
    PlanListItemDto(
      planId: json['plan_id'] as String,
      planCode: json['plan_code'] as String,
      title: json['title'] as String,
      status: json['status'] as String,
      totalItems: (json['total_items'] as num).toInt(),
      totalWeightKg: (json['total_weight_kg'] as num).toDouble(),
      volumeUtilizationPct: (json['volume_utilization_pct'] as num?)
          ?.toDouble(),
      createdBy: json['created_by'] as String,
      createdAt: json['created_at'] as String,
    );

Map<String, dynamic> _$PlanListItemDtoToJson(PlanListItemDto instance) =>
    <String, dynamic>{
      'plan_id': instance.planId,
      'plan_code': instance.planCode,
      'title': instance.title,
      'status': instance.status,
      'total_items': instance.totalItems,
      'total_weight_kg': instance.totalWeightKg,
      'volume_utilization_pct': instance.volumeUtilizationPct,
      'created_by': instance.createdBy,
      'created_at': instance.createdAt,
    };
