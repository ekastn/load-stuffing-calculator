class Constants {
  // Use 10.0.2.2 for Android Emulator to access localhost
  static const String apiBaseUrl = String.fromEnvironment(
    'API_URL',
    defaultValue: 'http://10.0.2.2:8080/api/v1',
  );
  
  // Web client base URL (for WebView embed routes)
  static const String webBaseUrl = String.fromEnvironment(
    'WEB_URL',
    defaultValue: 'http://10.0.2.2:3000',
  );
  
  static const String accessTokenKey = 'access_token';
  static const String refreshTokenKey = 'refresh_token';
  static const String activeWorkspaceIdKey = 'active_workspace_id';
}

// Backwards compatibility alias
@Deprecated('Use Constants instead')
typedef AppConstants = Constants;
