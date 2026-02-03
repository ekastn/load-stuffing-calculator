import 'package:json_annotation/json_annotation.dart';

part 'container_dto.g.dart';

@JsonSerializable()
class ContainerResponseDto {
  final String id;
  final String name;
  @JsonKey(name: 'inner_length_mm')
  final double innerLengthMm;
  @JsonKey(name: 'inner_width_mm')
  final double innerWidthMm;
  @JsonKey(name: 'inner_height_mm')
  final double innerHeightMm;
  @JsonKey(name: 'max_weight_kg')
  final double maxWeightKg;
  final String? description;

  ContainerResponseDto({
    required this.id,
    required this.name,
    required this.innerLengthMm,
    required this.innerWidthMm,
    required this.innerHeightMm,
    required this.maxWeightKg,
    this.description,
  });

  factory ContainerResponseDto.fromJson(Map<String, dynamic> json) => _$ContainerResponseDtoFromJson(json);
}

@JsonSerializable()
class CreateContainerRequestDto {
  final String name;
  @JsonKey(name: 'inner_length_mm')
  final double innerLengthMm;
  @JsonKey(name: 'inner_width_mm')
  final double innerWidthMm;
  @JsonKey(name: 'inner_height_mm')
  final double innerHeightMm;
  @JsonKey(name: 'max_weight_kg')
  final double maxWeightKg;
  final String? description;

  CreateContainerRequestDto({
    required this.name,
    required this.innerLengthMm,
    required this.innerWidthMm,
    required this.innerHeightMm,
    required this.maxWeightKg,
    this.description,
  });

  Map<String, dynamic> toJson() => _$CreateContainerRequestDtoToJson(this);
}

@JsonSerializable()
class UpdateContainerRequestDto {
  final String name;
  @JsonKey(name: 'inner_length_mm')
  final double innerLengthMm;
  @JsonKey(name: 'inner_width_mm')
  final double innerWidthMm;
  @JsonKey(name: 'inner_height_mm')
  final double innerHeightMm;
  @JsonKey(name: 'max_weight_kg')
  final double maxWeightKg;
  final String? description;

  UpdateContainerRequestDto({
    required this.name,
    required this.innerLengthMm,
    required this.innerWidthMm,
    required this.innerHeightMm,
    required this.maxWeightKg,
    this.description,
  });

  Map<String, dynamic> toJson() => _$UpdateContainerRequestDtoToJson(this);
}
