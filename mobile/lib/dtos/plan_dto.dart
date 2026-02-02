import 'package:json_annotation/json_annotation.dart';

part 'plan_dto.g.dart';

@JsonSerializable()
class PlanListItemDto {
  @JsonKey(name: 'plan_id')
  final String planId;
  @JsonKey(name: 'plan_code')
  final String planCode;
  final String title;
  final String status;
  @JsonKey(name: 'total_items')
  final int totalItems;
  @JsonKey(name: 'total_weight_kg')
  final double totalWeightKg;
  @JsonKey(name: 'volume_utilization_pct')
  final double? volumeUtilizationPct;
  @JsonKey(name: 'created_by')
  final String createdBy;
  @JsonKey(name: 'created_at')
  final String createdAt;

  PlanListItemDto({
    required this.planId,
    required this.planCode,
    required this.title,
    required this.status,
    required this.totalItems,
    required this.totalWeightKg,
    this.volumeUtilizationPct,
    required this.createdBy,
    required this.createdAt,
  });

  factory PlanListItemDto.fromJson(Map<String, dynamic> json) =>
      _$PlanListItemDtoFromJson(json);
}

@JsonSerializable()
class CreatePlanRequestDto {
  final String title;
  final String? notes;
  @JsonKey(name: 'auto_calculate')
  final bool? autoCalculate;
  final CreatePlanContainerDto container;
  final List<CreatePlanItemDto> items;

  CreatePlanRequestDto({
    required this.title,
    this.notes,
    this.autoCalculate = true,
    required this.container,
    required this.items,
  });

  Map<String, dynamic> toJson() => _$CreatePlanRequestDtoToJson(this);
}

@JsonSerializable()
class CreatePlanContainerDto {
  @JsonKey(name: 'container_id')
  final String? containerId;
  @JsonKey(name: 'length_mm')
  final double? lengthMm;
  @JsonKey(name: 'width_mm')
  final double? widthMm;
  @JsonKey(name: 'height_mm')
  final double? heightMm;
  @JsonKey(name: 'max_weight_kg')
  final double? maxWeightKg;

  CreatePlanContainerDto({
    this.containerId,
    this.lengthMm,
    this.widthMm,
    this.heightMm,
    this.maxWeightKg,
  });

  factory CreatePlanContainerDto.fromJson(Map<String, dynamic> json) =>
      _$CreatePlanContainerDtoFromJson(json);
  Map<String, dynamic> toJson() => _$CreatePlanContainerDtoToJson(this);
}

@JsonSerializable()
class CreatePlanItemDto {
  @JsonKey(name: 'product_sku')
  final String? productSku;
  final String? label;
  @JsonKey(name: 'length_mm')
  final double lengthMm;
  @JsonKey(name: 'width_mm')
  final double widthMm;
  @JsonKey(name: 'height_mm')
  final double heightMm;
  @JsonKey(name: 'weight_kg')
  final double weightKg;
  final int quantity;
  @JsonKey(name: 'allow_rotation')
  final bool? allowRotation;
  @JsonKey(name: 'color_hex')
  final String? colorHex;

  CreatePlanItemDto({
    this.productSku,
    this.label,
    required this.lengthMm,
    required this.widthMm,
    required this.heightMm,
    required this.weightKg,
    required this.quantity,
    this.allowRotation = true,
    this.colorHex,
  });

  factory CreatePlanItemDto.fromJson(Map<String, dynamic> json) =>
      _$CreatePlanItemDtoFromJson(json);
  Map<String, dynamic> toJson() => _$CreatePlanItemDtoToJson(this);
}

@JsonSerializable()
class CreatePlanResponseDto {
  @JsonKey(name: 'plan_id')
  final String planId;
  @JsonKey(name: 'plan_code')
  final String planCode;
  final String title;
  final String status;
  @JsonKey(name: 'total_items')
  final int totalItems;
  @JsonKey(name: 'total_weight_kg')
  final double totalWeightKg;
  @JsonKey(name: 'total_volume_m3')
  final double totalVolumeM3;
  @JsonKey(name: 'calculation_job_id')
  final String? calculationJobId;
  @JsonKey(name: 'created_at')
  final String createdAt;

  CreatePlanResponseDto({
    required this.planId,
    required this.planCode,
    required this.title,
    required this.status,
    required this.totalItems,
    required this.totalWeightKg,
    required this.totalVolumeM3,
    this.calculationJobId,
    required this.createdAt,
  });

  factory CreatePlanResponseDto.fromJson(Map<String, dynamic> json) =>
      _$CreatePlanResponseDtoFromJson(json);
}
