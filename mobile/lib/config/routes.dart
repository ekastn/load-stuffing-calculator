import 'package:go_router/go_router.dart';
import '../pages/auth/login_page.dart';
import '../pages/splash_screen.dart';
import '../pages/main_shell_page.dart';
import '../pages/plans/plan_form_page.dart';
import '../pages/plans/plan_detail_page.dart';
import '../pages/loading/loading_session_page.dart';
import '../pages/profile/profile_page.dart';
import '../pages/resources/container_form_page.dart';
import '../pages/resources/product_form_page.dart';
import '../pages/resources/resources_page.dart';
import '../providers/auth_provider.dart';

import '../pages/dashboard/dashboard_page.dart';

GoRouter createRouter(AuthProvider authProvider) {
  return GoRouter(
    initialLocation: '/splash',
    refreshListenable: authProvider,
    redirect: (context, state) {
      final isLoggingIn = state.uri.toString() == '/login';
      final isSplash = state.uri.toString() == '/splash';

      // Still loading auth state
      if (authProvider.isLoading) {
        return '/splash';
      }

      // If we are on splash and loading is done, decide where to go
      if (isSplash && !authProvider.isLoading) {
        return authProvider.isAuthenticated ? '/' : '/login';
      }

      if (!authProvider.isAuthenticated && !isLoggingIn && !isSplash) {
        return '/login';
      }

      if (authProvider.isAuthenticated && (isLoggingIn || isSplash)) {
        return '/';
      }

      return null;
    },
    routes: [
      GoRoute(
        path: '/splash',
        builder: (context, state) => const SplashScreen(),
      ),
      GoRoute(path: '/login', builder: (context, state) => const LoginPage()),
      ShellRoute(
        builder: (context, state, child) => MainShellPage(child: child),
        routes: [
          GoRoute(
            path: '/',
            builder: (context, state) => const DashboardPage(),
          ),
          GoRoute(
            path: '/plans',
            builder: (context, state) => const ResourcesPage(initialIndex: 0),
            routes: [
              GoRoute(
                path: 'new',
                builder: (context, state) => const PlanFormPage(),
              ),
              GoRoute(
                path: ':id',
                builder: (context, state) {
                  final id = state.pathParameters['id']!;
                  return PlanDetailPage(planId: id);
                },
              ),
              GoRoute(
                path: ':id/loading',
                builder: (context, state) {
                  final id = state.pathParameters['id']!;
                  return LoadingSessionPage(planId: id);
                },
              ),
            ],
          ),
          GoRoute(
            path: '/profile',
            builder: (context, state) => const ProfilePage(),
          ),
          GoRoute(
            path: '/products',
            builder: (context, state) => const ResourcesPage(initialIndex: 1),
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
            builder: (context, state) => const ResourcesPage(initialIndex: 2),
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
