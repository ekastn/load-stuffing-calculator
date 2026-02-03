import '../dtos/api_response_dto.dart';
import '../dtos/auth_dto.dart';
import '../mappers/auth_mapper.dart';
import '../models/user_model.dart';
import 'api_service.dart';
import 'storage_service.dart';

class AuthService {
  final ApiService _api;
  final StorageService _storage;

  AuthService(this._api, this._storage);

  Future<UserModel> login(String username, String password) async {
    try {
      final response = await _api.post('/auth/login', data: {
        'username': username,
        'password': password,
      });

      final apiResponse = ApiResponseDto<LoginResponseDto>.fromJson(
        response.data,
        (json) => LoginResponseDto.fromJson(json as Map<String, dynamic>),
      );

      if (!apiResponse.success || apiResponse.data == null) {
        throw Exception(apiResponse.errors?.firstOrNull?.message ?? 'Login failed');
      }

      final dto = apiResponse.data!;

      await _storage.setAccessToken(dto.accessToken);
      await _storage.setRefreshToken(dto.refreshToken);
      if (dto.activeWorkspaceId != null) {
        await _storage.setActiveWorkspaceId(dto.activeWorkspaceId!);
      }

      return AuthMapper.toUserModel(dto.user);
    } catch (e) {
      // In a real app we would map DioException to a domain exception here
      rethrow; 
    }
  }

  Future<void> logout() async {
    await _storage.clearAll();
  }
}
