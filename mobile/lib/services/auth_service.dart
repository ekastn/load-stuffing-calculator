import 'package:dio/dio.dart';

import '../dtos/api_response_dto.dart';
import '../dtos/auth_dto.dart';
import '../mappers/auth_mapper.dart';
import '../models/user_model.dart';
import '../exceptions/login_exception.dart';
import 'api_service.dart';
import 'storage_service.dart';

class AuthService {
  final ApiService _api;
  final StorageService _storage;

  AuthService(this._api, this._storage);

  Future<UserModel> login(String username, String password) async {
    try {
      final response = await _api.post(
        '/auth/login',
        data: {'username': username, 'password': password},
      );

      final apiResponse = ApiResponseDto<LoginResponseDto>.fromJson(
        response.data,
        (json) => LoginResponseDto.fromJson(json as Map<String, dynamic>),
      );

      if (!apiResponse.success || apiResponse.data == null) {
        throw LoginException(
          apiResponse.errors?.firstOrNull?.message ??
              'Invalid username or password',
        );
      }

      final dto = apiResponse.data!;

      await _storage.setAccessToken(dto.accessToken);
      await _storage.setRefreshToken(dto.refreshToken);
      if (dto.activeWorkspaceId != null) {
        await _storage.setActiveWorkspaceId(dto.activeWorkspaceId!);
      }

      return AuthMapper.toUserModel(dto.user);
    } catch (e) {
      if (e is LoginException) rethrow;

      // Handle DioException (API errors)
      if (e is DioException) {
        // Try to extract error message from response
        if (e.response != null) {
          try {
            final apiResponse = ApiResponseDto<dynamic>.fromJson(
              e.response!.data,
              (json) => json,
            );
            final errorMessage =
                apiResponse.errors?.firstOrNull?.message ??
                'Invalid username or password';
            throw LoginException(errorMessage);
          } catch (_) {
            // If parsing fails, use status-based message
            if (e.response!.statusCode == 401) {
              throw LoginException('Invalid username or password');
            }
          }
        }
        // Network/connection error
        throw LoginException(
          'Unable to connect. Please check your internet connection.',
        );
      }

      // Unknown error
      throw LoginException('Login failed: ${e.toString()}');
    }
  }

  Future<void> logout() async {
    await _storage.clearAll();
  }

  Future<String?> getAccessToken() async {
    return await _storage.getAccessToken();
  }

  /// Attempts to restore user session from stored token
  Future<UserModel?> getCurrentUser() async {
    try {
      final token = await _storage.getAccessToken();

      if (token == null || token.isEmpty) {
        return null;
      }

      // Verify token is still valid by fetching current user info
      final response = await _api.get('/auth/me');

      final apiResponse = ApiResponseDto<AuthMeResponseDto>.fromJson(
        response.data,
        (json) => AuthMeResponseDto.fromJson(json as Map<String, dynamic>),
      );

      if (!apiResponse.success || apiResponse.data == null) {
        // Token is invalid, clear storage
        await _storage.clearAll();
        return null;
      }

      return AuthMapper.toUserModel(apiResponse.data!.user);
    } catch (e) {
      // If token validation fails, clear storage
      await _storage.clearAll();
      return null;
    }
  }
}
