import 'package:json_annotation/json_annotation.dart';

part 'product_dto.g.dart';

@JsonSerializable()
class ProductResponseDto {
  final String id;
  final String name;
  @JsonKey(name: 'length_mm')
  final double lengthMm;
  @JsonKey(name: 'width_mm')
  final double widthMm;
  @JsonKey(name: 'height_mm')
  final double heightMm;
  @JsonKey(name: 'weight_kg')
  final double weightKg;
  @JsonKey(name: 'color_hex')
  final String? colorHex;

  ProductResponseDto({
    required this.id,
    required this.name,
    required this.lengthMm,
    required this.widthMm,
    required this.heightMm,
    required this.weightKg,
    this.colorHex,
  });

  factory ProductResponseDto.fromJson(Map<String, dynamic> json) => _$ProductResponseDtoFromJson(json);
}

@JsonSerializable()
class CreateProductRequestDto {
  final String name;
  @JsonKey(name: 'length_mm')
  final double lengthMm;
  @JsonKey(name: 'width_mm')
  final double widthMm;
  @JsonKey(name: 'height_mm')
  final double heightMm;
  @JsonKey(name: 'weight_kg')
  final double weightKg;
  @JsonKey(name: 'color_hex')
  final String? colorHex;

  CreateProductRequestDto({
    required this.name,
    required this.lengthMm,
    required this.widthMm,
    required this.heightMm,
    required this.weightKg,
    this.colorHex,
  });

  Map<String, dynamic> toJson() => _$CreateProductRequestDtoToJson(this);
}

@JsonSerializable()
class UpdateProductRequestDto {
  final String name;
  @JsonKey(name: 'length_mm')
  final double lengthMm;
  @JsonKey(name: 'width_mm')
  final double widthMm;
  @JsonKey(name: 'height_mm')
  final double heightMm;
  @JsonKey(name: 'weight_kg')
  final double weightKg;
  @JsonKey(name: 'color_hex')
  final String? colorHex;

  UpdateProductRequestDto({
    required this.name,
    required this.lengthMm,
    required this.widthMm,
    required this.heightMm,
    required this.weightKg,
    this.colorHex,
  });

  Map<String, dynamic> toJson() => _$UpdateProductRequestDtoToJson(this);
}
