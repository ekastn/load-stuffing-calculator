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
