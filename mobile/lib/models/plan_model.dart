import 'package:freezed_annotation/freezed_annotation.dart';

part 'plan_model.freezed.dart';

@freezed
abstract class PlanModel with _$PlanModel {
  const factory PlanModel({
    required String id,
    required String code,
    required String title,
    required String status,
    required int totalItems,
    required double totalWeightKg,
    double? volumeUtilizationPct,
    required String createdBy,
    required DateTime createdAt,
  }) = _PlanModel;
}
