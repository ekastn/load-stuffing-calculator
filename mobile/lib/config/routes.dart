import 'package:flutter/material.dart';
import 'package:go_router/go_router.dart';
import 'package:provider/provider.dart';
import '../pages/auth/login_page.dart';
import '../pages/main_shell_page.dart';
import '../pages/plans/plan_list_page.dart';
import '../pages/plans/plan_form_page.dart';
import '../pages/resources/container_form_page.dart';
import '../pages/resources/container_list_page.dart';
import '../pages/resources/product_form_page.dart';
import '../pages/resources/product_list_page.dart';
import '../providers/auth_provider.dart';

import '../pages/dashboard/dashboard_page.dart';

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
            builder: (context, state) => const PlanListPage(),
            routes: [
              GoRoute(
                path: 'new',
                builder: (context, state) => const PlanFormPage(),
              ),
            ],
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
          GoRoute(
            path: '/products',
            builder: (context, state) => const ProductListPage(),
            routes: [
              GoRoute(
                path: 'new',
                builder: (context, state) => const ProductFormPage(),
              ),
              GoRoute(
                path: ':id',
                builder: (context, state) {
                  final id = state.pathParameters['id']!;
                  return ProductFormPage(productId: id);
                },
              ),
            ],
          ),
          GoRoute(
            path: '/containers',
            builder: (context, state) => const ContainerListPage(),
            routes: [
              GoRoute(
                path: 'new',
                builder: (context, state) => const ContainerFormPage(),
              ),
              GoRoute(
                path: ':id',
                builder: (context, state) {
                  final id = state.pathParameters['id']!;
                  return ContainerFormPage(containerId: id);
                },
              ),
            ],
          ),
        ],
      ),
    ],
  );
}
