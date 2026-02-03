// GENERATED CODE - DO NOT MODIFY BY HAND

part of 'product_dto.dart';

// **************************************************************************
// JsonSerializableGenerator
// **************************************************************************

ProductResponseDto _$ProductResponseDtoFromJson(Map<String, dynamic> json) =>
    ProductResponseDto(
      id: json['id'] as String,
      name: json['name'] as String,
      lengthMm: (json['length_mm'] as num).toDouble(),
      widthMm: (json['width_mm'] as num).toDouble(),
      heightMm: (json['height_mm'] as num).toDouble(),
      weightKg: (json['weight_kg'] as num).toDouble(),
      colorHex: json['color_hex'] as String?,
    );

Map<String, dynamic> _$ProductResponseDtoToJson(ProductResponseDto instance) =>
    <String, dynamic>{
      'id': instance.id,
      'name': instance.name,
      'length_mm': instance.lengthMm,
      'width_mm': instance.widthMm,
      'height_mm': instance.heightMm,
      'weight_kg': instance.weightKg,
      'color_hex': instance.colorHex,
    };

CreateProductRequestDto _$CreateProductRequestDtoFromJson(
  Map<String, dynamic> json,
) => CreateProductRequestDto(
  name: json['name'] as String,
  lengthMm: (json['length_mm'] as num).toDouble(),
  widthMm: (json['width_mm'] as num).toDouble(),
  heightMm: (json['height_mm'] as num).toDouble(),
  weightKg: (json['weight_kg'] as num).toDouble(),
  colorHex: json['color_hex'] as String?,
);

Map<String, dynamic> _$CreateProductRequestDtoToJson(
  CreateProductRequestDto instance,
) => <String, dynamic>{
  'name': instance.name,
  'length_mm': instance.lengthMm,
  'width_mm': instance.widthMm,
  'height_mm': instance.heightMm,
  'weight_kg': instance.weightKg,
  'color_hex': instance.colorHex,
};

UpdateProductRequestDto _$UpdateProductRequestDtoFromJson(
  Map<String, dynamic> json,
) => UpdateProductRequestDto(
  name: json['name'] as String,
  lengthMm: (json['length_mm'] as num).toDouble(),
  widthMm: (json['width_mm'] as num).toDouble(),
  heightMm: (json['height_mm'] as num).toDouble(),
  weightKg: (json['weight_kg'] as num).toDouble(),
  colorHex: json['color_hex'] as String?,
);

Map<String, dynamic> _$UpdateProductRequestDtoToJson(
  UpdateProductRequestDto instance,
) => <String, dynamic>{
  'name': instance.name,
  'length_mm': instance.lengthMm,
  'width_mm': instance.widthMm,
  'height_mm': instance.heightMm,
  'weight_kg': instance.weightKg,
  'color_hex': instance.colorHex,
};
