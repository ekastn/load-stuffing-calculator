import 'package:flutter/material.dart';
import 'package:provider/provider.dart';
import 'package:flutter_secure_storage/flutter_secure_storage.dart';
import 'config/routes.dart';
import 'config/theme.dart';
import 'providers/auth_provider.dart';
import 'providers/dashboard_provider.dart';
import 'services/api_service.dart';
import 'services/auth_service.dart';
import 'services/dashboard_service.dart';
import 'services/plan_service.dart';
import 'services/storage_service.dart';

void main() {
  WidgetsFlutterBinding.ensureInitialized();

  // 1. Initialize Infrastucture (Services)
  final storageService = StorageService(storage: const FlutterSecureStorage());
  final apiService = ApiService(storageService);
  final authService = AuthService(apiService, storageService);
  final dashboardService = DashboardService(apiService);
  final planService = PlanService(apiService);

  runApp(MyApp(
    authService: authService,
    storageService: storageService,
    dashboardService: dashboardService,
    planService: planService,
  ));
}

class MyApp extends StatelessWidget {
  final AuthService authService;
  final StorageService storageService;
  final DashboardService dashboardService;
  final PlanService planService;

  const MyApp({
    super.key,
    required this.authService,
    required this.storageService,
    required this.dashboardService,
    required this.planService,
  });

  @override
  Widget build(BuildContext context) {
    // 2. Setup Provider Graph with Dependency Injection
    return MultiProvider(
      providers: [
        ChangeNotifierProvider(
          create: (_) => AuthProvider(authService),
        ),
        ChangeNotifierProvider(
          create: (_) => DashboardProvider(dashboardService, planService),
        ),
      ],
      child: const AppRoot(),
    );
  }
}

class AppRoot extends StatelessWidget {
  const AppRoot({super.key});

  @override
  Widget build(BuildContext context) {
    // 3. Consume Providers for High-Level Logic (e.g. Auth Status for Router)
    // We create the router here, listening to AuthProvider changes.
    final authProvider = context.read<AuthProvider>();
    final router = createRouter(authProvider);

    return MaterialApp.router(
      title: 'Load Stuffing Calculator',
      theme: AppTheme.lightTheme,
      routerConfig: router,
      debugShowCheckedModeBanner: false,
    );
  }
}
