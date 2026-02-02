import '../dtos/product_dto.dart';
import '../dtos/api_response_dto.dart';
import '../models/product_model.dart';
import '../mappers/product_mapper.dart';
import 'api_service.dart';

class ProductService {
  final ApiService _apiService;

  ProductService(this._apiService);

  Future<List<ProductModel>> getProducts() async {
    final response = await _apiService.get('/products');
    final apiResponse = ApiResponseDto<List<dynamic>>.fromJson(
      response.data,
      (json) => json as List<dynamic>,
    );
    
    if (!apiResponse.success) {
        throw Exception(apiResponse.errors?.firstOrNull?.message ?? 'Failed to fetch products');
    }

    if (apiResponse.data == null) return [];

    return apiResponse.data!
        .map((e) => ProductResponseDto.fromJson(e as Map<String, dynamic>))
        .map((dto) => ProductMapper.toModel(dto))
        .toList();
  }

  Future<ProductModel> getProduct(String id) async {
      final response = await _apiService.get('/products/$id');
      final apiResponse = ApiResponseDto<ProductResponseDto>.fromJson(
        response.data,
        (json) => ProductResponseDto.fromJson(json as Map<String, dynamic>),
      );

      if (!apiResponse.success) {
          throw Exception(apiResponse.errors?.firstOrNull?.message ?? 'Failed to fetch product');
      }
      return ProductMapper.toModel(apiResponse.data!);
  }

  Future<ProductModel> createProduct(CreateProductRequestDto request) async {
    final response = await _apiService.post('/products', data: request.toJson());
    final apiResponse = ApiResponseDto<ProductResponseDto>.fromJson(
      response.data,
      (json) => ProductResponseDto.fromJson(json as Map<String, dynamic>),
    );
    
    if (!apiResponse.success) {
      throw Exception(apiResponse.errors?.firstOrNull?.message ?? 'Failed to create product');
    }

    return ProductMapper.toModel(apiResponse.data!);
  }

  Future<ProductModel> updateProduct(String id, UpdateProductRequestDto request) async {
    final response = await _apiService.put('/products/$id', data: request.toJson());
    final apiResponse = ApiResponseDto<ProductResponseDto>.fromJson(
      response.data,
      (json) => ProductResponseDto.fromJson(json as Map<String, dynamic>),
    );
    
    if (!apiResponse.success) {
      throw Exception(apiResponse.errors?.firstOrNull?.message ?? 'Failed to update product');
    }

    return ProductMapper.toModel(apiResponse.data!);
  }

  Future<void> deleteProduct(String id) async {
    final response = await _apiService.delete('/products/$id');
    final apiResponse = ApiResponseDto<void>.fromJson(
      response.data,
      (json) {}, // Void callback
    );

    if (!apiResponse.success) {
      throw Exception(apiResponse.errors?.firstOrNull?.message ?? 'Failed to delete product');
    }
  }
}
