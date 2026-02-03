import '../dtos/container_dto.dart';
import '../dtos/api_response_dto.dart';
import '../models/container_model.dart';
import '../mappers/container_mapper.dart';
import 'api_service.dart';

class ContainerService {
  final ApiService _apiService;

  ContainerService(this._apiService);

  Future<List<ContainerModel>> getContainers() async {
    final response = await _apiService.get('/containers');
    final apiResponse = ApiResponseDto<List<dynamic>>.fromJson(
      response.data,
      (json) => json as List<dynamic>,
    );
    
    if (!apiResponse.success) {
        throw Exception(apiResponse.errors?.firstOrNull?.message ?? 'Failed to fetch containers');
    }

    if (apiResponse.data == null) return [];

    return apiResponse.data!
        .map((e) => ContainerResponseDto.fromJson(e as Map<String, dynamic>))
        .map((dto) => ContainerMapper.toModel(dto))
        .toList();
  }

  Future<ContainerModel> getContainer(String id) async {
      final response = await _apiService.get('/containers/$id');
      final apiResponse = ApiResponseDto<ContainerResponseDto>.fromJson(
        response.data,
        (json) => ContainerResponseDto.fromJson(json as Map<String, dynamic>),
      );

      if (!apiResponse.success) {
          throw Exception(apiResponse.errors?.firstOrNull?.message ?? 'Failed to fetch container');
      }
      return ContainerMapper.toModel(apiResponse.data!);
  }

  Future<ContainerModel> createContainer(CreateContainerRequestDto request) async {
    final response = await _apiService.post('/containers', data: request.toJson());
    final apiResponse = ApiResponseDto<ContainerResponseDto>.fromJson(
      response.data,
      (json) => ContainerResponseDto.fromJson(json as Map<String, dynamic>),
    );
    
    if (!apiResponse.success) {
      throw Exception(apiResponse.errors?.firstOrNull?.message ?? 'Failed to create container');
    }

    return ContainerMapper.toModel(apiResponse.data!);
  }

  Future<ContainerModel> updateContainer(String id, UpdateContainerRequestDto request) async {
    final response = await _apiService.put('/containers/$id', data: request.toJson());
    final apiResponse = ApiResponseDto<ContainerResponseDto>.fromJson(
      response.data,
      (json) => ContainerResponseDto.fromJson(json as Map<String, dynamic>),
    );
    
    if (!apiResponse.success) {
      throw Exception(apiResponse.errors?.firstOrNull?.message ?? 'Failed to update container');
    }

    return ContainerMapper.toModel(apiResponse.data!);
  }

  Future<void> deleteContainer(String id) async {
    final response = await _apiService.delete('/containers/$id');
    final apiResponse = ApiResponseDto<void>.fromJson(
      response.data,
      (json) {}, // Void callback
    );

    if (!apiResponse.success) {
      throw Exception(apiResponse.errors?.firstOrNull?.message ?? 'Failed to delete container');
    }
  }
}
