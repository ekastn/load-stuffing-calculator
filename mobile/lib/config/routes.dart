import 'package:flutter/material.dart';
import 'package:go_router/go_router.dart';
import 'package:provider/provider.dart';
import '../pages/auth/login_page.dart';
import '../providers/auth_provider.dart';

// Simple placeholder page until Home is implemented
class PlaceholderHomePage extends StatelessWidget {
  const PlaceholderHomePage({super.key});
  @override
  Widget build(BuildContext context) {
    return const Scaffold(body: Center(child: Text('Home Page Placeholder')));
  }
}

GoRouter createRouter(AuthProvider authProvider) {
  return GoRouter(
    initialLocation: '/',
    refreshListenable: authProvider,
    redirect: (context, state) {
      final isLoggingIn = state.uri.toString() == '/login';

      if (!authProvider.isAuthenticated && !isLoggingIn) {
        return '/login';
      }
      if (authProvider.isAuthenticated && isLoggingIn) {
        return '/';
      }
      return null;
    },
    routes: [
      GoRoute(
        path: '/',
        builder: (context, state) => const PlaceholderHomePage(),
      ),
      GoRoute(
        path: '/login',
        builder: (context, state) => const LoginPage(),
      ),
    ],
  );
}
