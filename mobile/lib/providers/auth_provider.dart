import 'package:flutter/material.dart';
import '../models/user_model.dart';
import '../services/auth_service.dart';

class AuthProvider extends ChangeNotifier {
  final AuthService _authService;

  AuthProvider(this._authService);

  UserModel? _user;
  UserModel? get user => _user;

  String? _accessToken;
  String? get accessToken => _accessToken;

  bool get isAuthenticated => _user != null;

  bool _isLoading = true;
  bool get isLoading => _isLoading;

  String? _error;
  String? get error => _error;

  /// Initialize auth state by checking for existing session
  Future<void> initialize() async {
    _isLoading = true;
    notifyListeners();

    try {
      _user = await _authService.getCurrentUser();
      if (_user != null) {
        _accessToken = await _authService.getAccessToken();
      }
    } catch (e) {
      _user = null;
      // Silent failure - user will just see login screen
    } finally {
      _isLoading = false;
      notifyListeners();
    }
  }

  Future<void> login(String username, String password) async {
    _isLoading = true;
    _error = null;
    notifyListeners();

    try {
      _user = await _authService.login(username, password);
      _accessToken = await _authService.getAccessToken();
    } catch (e) {
      _error = e.toString();
      _user = null;
      rethrow;
    } finally {
      _isLoading = false;
      notifyListeners();
    }
  }

  Future<void> logout() async {
    await _authService.logout();
    _user = null;
    _accessToken = null;
    notifyListeners();
  }
}
