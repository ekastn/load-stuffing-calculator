class Constants {
  // Use 10.0.2.2 for Android Emulator to access localhost
  static const String apiBaseUrl = String.fromEnvironment(
    'API_URL',
    defaultValue: 'https://stuffing-api.irc-enter.tech/api/v1',
  );
  
  // Web client base URL (for WebView embed routes)
  static const String webBaseUrl = String.fromEnvironment(
    'WEB_URL',
    defaultValue: 'https://stuffing.irc-enter.tech',
  );
  
  static const String accessTokenKey = 'access_token';
  static const String refreshTokenKey = 'refresh_token';
  static const String activeWorkspaceIdKey = 'active_workspace_id';
  static const String loadingSessionKey = 'loading_session';
}

// Backwards compatibility alias
@Deprecated('Use Constants instead')
typedef AppConstants = Constants;
