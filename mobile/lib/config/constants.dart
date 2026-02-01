class AppConstants {
  // Use 10.0.2.2 for Android Emulator to access localhost
  static const String apiBaseUrl = String.fromEnvironment(
    'API_URL',
    defaultValue: 'http://10.0.2.2:8080/api/v1',
  );
  
  static const String accessTokenKey = 'access_token';
  static const String refreshTokenKey = 'refresh_token';
  static const String activeWorkspaceIdKey = 'active_workspace_id';
}
