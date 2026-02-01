import 'package:freezed_annotation/freezed_annotation.dart';

part 'product_model.freezed.dart';

@freezed
abstract class ProductModel with _$ProductModel {
  const factory ProductModel({
    required String id,
    required String name,
    required double lengthMm,
    required double widthMm,
    required double heightMm,
    required double weightKg,
    String? colorHex,
  }) = _ProductModel;
}
