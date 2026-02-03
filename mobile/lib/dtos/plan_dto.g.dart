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

CreatePlanRequestDto _$CreatePlanRequestDtoFromJson(
  Map<String, dynamic> json,
) => CreatePlanRequestDto(
  title: json['title'] as String,
  notes: json['notes'] as String?,
  autoCalculate: json['auto_calculate'] as bool? ?? true,
  container: CreatePlanContainerDto.fromJson(
    json['container'] as Map<String, dynamic>,
  ),
  items: (json['items'] as List<dynamic>)
      .map((e) => CreatePlanItemDto.fromJson(e as Map<String, dynamic>))
      .toList(),
);

Map<String, dynamic> _$CreatePlanRequestDtoToJson(
  CreatePlanRequestDto instance,
) => <String, dynamic>{
  'title': instance.title,
  'notes': instance.notes,
  'auto_calculate': instance.autoCalculate,
  'container': instance.container,
  'items': instance.items,
};

CreatePlanContainerDto _$CreatePlanContainerDtoFromJson(
  Map<String, dynamic> json,
) => CreatePlanContainerDto(
  containerId: json['container_id'] as String?,
  lengthMm: (json['length_mm'] as num?)?.toDouble(),
  widthMm: (json['width_mm'] as num?)?.toDouble(),
  heightMm: (json['height_mm'] as num?)?.toDouble(),
  maxWeightKg: (json['max_weight_kg'] as num?)?.toDouble(),
);

Map<String, dynamic> _$CreatePlanContainerDtoToJson(
  CreatePlanContainerDto instance,
) => <String, dynamic>{
  'container_id': instance.containerId,
  'length_mm': instance.lengthMm,
  'width_mm': instance.widthMm,
  'height_mm': instance.heightMm,
  'max_weight_kg': instance.maxWeightKg,
};

CreatePlanItemDto _$CreatePlanItemDtoFromJson(Map<String, dynamic> json) =>
    CreatePlanItemDto(
      productSku: json['product_sku'] as String?,
      label: json['label'] as String?,
      lengthMm: (json['length_mm'] as num).toDouble(),
      widthMm: (json['width_mm'] as num).toDouble(),
      heightMm: (json['height_mm'] as num).toDouble(),
      weightKg: (json['weight_kg'] as num).toDouble(),
      quantity: (json['quantity'] as num).toInt(),
      allowRotation: json['allow_rotation'] as bool? ?? true,
      colorHex: json['color_hex'] as String?,
    );

Map<String, dynamic> _$CreatePlanItemDtoToJson(CreatePlanItemDto instance) =>
    <String, dynamic>{
      'product_sku': instance.productSku,
      'label': instance.label,
      'length_mm': instance.lengthMm,
      'width_mm': instance.widthMm,
      'height_mm': instance.heightMm,
      'weight_kg': instance.weightKg,
      'quantity': instance.quantity,
      'allow_rotation': instance.allowRotation,
      'color_hex': instance.colorHex,
    };

CreatePlanResponseDto _$CreatePlanResponseDtoFromJson(
  Map<String, dynamic> json,
) => CreatePlanResponseDto(
  planId: json['plan_id'] as String,
  planCode: json['plan_code'] as String,
  title: json['title'] as String,
  status: json['status'] as String,
  totalItems: (json['total_items'] as num).toInt(),
  totalWeightKg: (json['total_weight_kg'] as num).toDouble(),
  totalVolumeM3: (json['total_volume_m3'] as num).toDouble(),
  calculationJobId: json['calculation_job_id'] as String?,
  createdAt: json['created_at'] as String,
);

Map<String, dynamic> _$CreatePlanResponseDtoToJson(
  CreatePlanResponseDto instance,
) => <String, dynamic>{
  'plan_id': instance.planId,
  'plan_code': instance.planCode,
  'title': instance.title,
  'status': instance.status,
  'total_items': instance.totalItems,
  'total_weight_kg': instance.totalWeightKg,
  'total_volume_m3': instance.totalVolumeM3,
  'calculation_job_id': instance.calculationJobId,
  'created_at': instance.createdAt,
};
