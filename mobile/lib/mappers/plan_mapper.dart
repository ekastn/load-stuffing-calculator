import '../dtos/plan_dto.dart';
import '../models/plan_model.dart';

class PlanMapper {
  static PlanModel fromDto(PlanListItemDto dto) {
    return PlanModel(
      id: dto.planId,
      code: dto.planCode,
      title: dto.title,
      status: dto.status,
      totalItems: dto.totalItems,
      totalWeightKg: dto.totalWeightKg,
      volumeUtilizationPct: dto.volumeUtilizationPct,
      createdBy: dto.createdBy,
      createdAt: DateTime.parse(dto.createdAt),
    );
  }
}
