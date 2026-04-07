// GENERATED CODE - DO NOT MODIFY BY HAND

part of 'container_dto.dart';

// **************************************************************************
// JsonSerializableGenerator
// **************************************************************************

ContainerResponseDto _$ContainerResponseDtoFromJson(
  Map<String, dynamic> json,
) => ContainerResponseDto(
  id: json['id'] as String,
  name: json['name'] as String,
  innerLengthMm: (json['inner_length_mm'] as num).toDouble(),
  innerWidthMm: (json['inner_width_mm'] as num).toDouble(),
  innerHeightMm: (json['inner_height_mm'] as num).toDouble(),
  maxWeightKg: (json['max_weight_kg'] as num).toDouble(),
  description: json['description'] as String?,
);

Map<String, dynamic> _$ContainerResponseDtoToJson(
  ContainerResponseDto instance,
) => <String, dynamic>{
  'id': instance.id,
  'name': instance.name,
  'inner_length_mm': instance.innerLengthMm,
  'inner_width_mm': instance.innerWidthMm,
  'inner_height_mm': instance.innerHeightMm,
  'max_weight_kg': instance.maxWeightKg,
  'description': instance.description,
};

CreateContainerRequestDto _$CreateContainerRequestDtoFromJson(
  Map<String, dynamic> json,
) => CreateContainerRequestDto(
  name: json['name'] as String,
  innerLengthMm: (json['inner_length_mm'] as num).toDouble(),
  innerWidthMm: (json['inner_width_mm'] as num).toDouble(),
  innerHeightMm: (json['inner_height_mm'] as num).toDouble(),
  maxWeightKg: (json['max_weight_kg'] as num).toDouble(),
  description: json['description'] as String?,
);

Map<String, dynamic> _$CreateContainerRequestDtoToJson(
  CreateContainerRequestDto instance,
) => <String, dynamic>{
  'name': instance.name,
  'inner_length_mm': instance.innerLengthMm,
  'inner_width_mm': instance.innerWidthMm,
  'inner_height_mm': instance.innerHeightMm,
  'max_weight_kg': instance.maxWeightKg,
  'description': instance.description,
};

UpdateContainerRequestDto _$UpdateContainerRequestDtoFromJson(
  Map<String, dynamic> json,
) => UpdateContainerRequestDto(
  name: json['name'] as String,
  innerLengthMm: (json['inner_length_mm'] as num).toDouble(),
  innerWidthMm: (json['inner_width_mm'] as num).toDouble(),
  innerHeightMm: (json['inner_height_mm'] as num).toDouble(),
  maxWeightKg: (json['max_weight_kg'] as num).toDouble(),
  description: json['description'] as String?,
);

Map<String, dynamic> _$UpdateContainerRequestDtoToJson(
  UpdateContainerRequestDto instance,
) => <String, dynamic>{
  'name': instance.name,
  'inner_length_mm': instance.innerLengthMm,
  'inner_width_mm': instance.innerWidthMm,
  'inner_height_mm': instance.innerHeightMm,
  'max_weight_kg': instance.maxWeightKg,
  'description': instance.description,
};
