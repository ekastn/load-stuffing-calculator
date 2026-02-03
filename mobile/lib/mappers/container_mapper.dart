import '../dtos/container_dto.dart';
import '../models/container_model.dart';

class ContainerMapper {
  static ContainerModel toModel(ContainerResponseDto dto) {
    return ContainerModel(
      id: dto.id,
      name: dto.name,
      innerLengthMm: dto.innerLengthMm,
      innerWidthMm: dto.innerWidthMm,
      innerHeightMm: dto.innerHeightMm,
      maxWeightKg: dto.maxWeightKg,
      description: dto.description,
    );
  }

  static CreateContainerRequestDto toCreateDto(ContainerModel model) {
    return CreateContainerRequestDto(
      name: model.name,
      innerLengthMm: model.innerLengthMm,
      innerWidthMm: model.innerWidthMm,
      innerHeightMm: model.innerHeightMm,
      maxWeightKg: model.maxWeightKg,
      description: model.description,
    );
  }

  static UpdateContainerRequestDto toUpdateDto(ContainerModel model) {
    return UpdateContainerRequestDto(
      name: model.name,
      innerLengthMm: model.innerLengthMm,
      innerWidthMm: model.innerWidthMm,
      innerHeightMm: model.innerHeightMm,
      maxWeightKg: model.maxWeightKg,
      description: model.description,
    );
  }
}
