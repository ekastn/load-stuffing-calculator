import '../dtos/product_dto.dart';
import '../models/product_model.dart';

class ProductMapper {
  static ProductModel toModel(ProductResponseDto dto) {
    return ProductModel(
      id: dto.id,
      name: dto.name,
      lengthMm: dto.lengthMm,
      widthMm: dto.widthMm,
      heightMm: dto.heightMm,
      weightKg: dto.weightKg,
      colorHex: dto.colorHex,
    );
  }

  static CreateProductRequestDto toCreateDto(ProductModel model) {
    return CreateProductRequestDto(
      name: model.name,
      lengthMm: model.lengthMm,
      widthMm: model.widthMm,
      heightMm: model.heightMm,
      weightKg: model.weightKg,
      colorHex: model.colorHex,
    );
  }

  static UpdateProductRequestDto toUpdateDto(ProductModel model) {
    return UpdateProductRequestDto(
      name: model.name,
      lengthMm: model.lengthMm,
      widthMm: model.widthMm,
      heightMm: model.heightMm,
      weightKg: model.weightKg,
      colorHex: model.colorHex,
    );
  }
}
