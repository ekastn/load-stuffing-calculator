import 'package:flutter_secure_storage/flutter_secure_storage.dart';
import '../config/constants.dart';

class StorageService {
  final FlutterSecureStorage _storage;

  StorageService({FlutterSecureStorage? storage})
      : _storage = storage ?? const FlutterSecureStorage();

  Future<void> setAccessToken(String token) async {
    await _storage.write(key: AppConstants.accessTokenKey, value: token);
  }

  Future<String?> getAccessToken() async {
    return await _storage.read(key: AppConstants.accessTokenKey);
  }

  Future<void> setRefreshToken(String token) async {
    await _storage.write(key: AppConstants.refreshTokenKey, value: token);
  }

  Future<String?> getRefreshToken() async {
    return await _storage.read(key: AppConstants.refreshTokenKey);
  }

  Future<void> setActiveWorkspaceId(String id) async {
    await _storage.write(key: AppConstants.activeWorkspaceIdKey, value: id);
  }

  Future<String?> getActiveWorkspaceId() async {
    return await _storage.read(key: AppConstants.activeWorkspaceIdKey);
  }

  Future<void> clearAll() async {
    await _storage.deleteAll();
  }

  Future<void> saveLoadingSession(String payload) async {
    await _storage.write(key: AppConstants.loadingSessionKey, value: payload);
  }

  Future<String?> getLoadingSession() async {
    return await _storage.read(key: AppConstants.loadingSessionKey);
  }

  Future<void> deleteLoadingSession() async {
    await _storage.delete(key: AppConstants.loadingSessionKey);
  }

  Future<void> saveLoadingExpectedItems(String planId, String payload) async {
    await _storage.write(
      key: '${AppConstants.loadingSessionKey}_items_$planId',
      value: payload,
    );
  }

  Future<String?> getLoadingExpectedItems(String planId) async {
    return await _storage.read(key: '${AppConstants.loadingSessionKey}_items_$planId');
  }

  Future<void> deleteLoadingExpectedItems(String planId) async {
    await _storage.delete(key: '${AppConstants.loadingSessionKey}_items_$planId');
  }
}
