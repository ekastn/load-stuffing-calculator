import 'package:flutter/material.dart';
import 'package:go_router/go_router.dart';
import 'package:provider/provider.dart';
import '../pages/auth/login_page.dart';
import '../providers/auth_provider.dart';

import '../pages/dashboard/dashboard_page.dart';
import '../pages/main_shell_page.dart';

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
        path: '/login',
        builder: (context, state) => const LoginPage(),
      ),
      ShellRoute(
        builder: (context, state, child) => MainShellPage(child: child),
        routes: [
          GoRoute(
            path: '/',
            builder: (context, state) => const DashboardPage(),
          ),
          GoRoute(
            path: '/plans',
            builder: (context, state) => const Scaffold(
              body: Center(child: Text('Plans List')),
            ),
          ),
          GoRoute(
            path: '/profile',
            builder: (context, state) {
              return Scaffold(
                body: Center(
                  child: ElevatedButton(
                    onPressed: () => context.read<AuthProvider>().logout(),
                    child: const Text('Logout'),
                  ),
                ),
              );
            },
          ),
        ],
      ),
    ],
  );
}
