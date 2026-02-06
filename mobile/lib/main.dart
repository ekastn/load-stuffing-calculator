import 'package:flutter/material.dart';
import 'package:provider/provider.dart';
import 'package:flutter_secure_storage/flutter_secure_storage.dart';
import 'package:go_router/go_router.dart';
import 'config/routes.dart';
import 'config/theme.dart';
import 'providers/auth_provider.dart';
import 'providers/container_provider.dart';
import 'providers/dashboard_provider.dart';
import 'providers/plan_provider.dart';
import 'providers/plan_detail_provider.dart';
import 'providers/product_provider.dart';
import 'services/api_service.dart';
import 'services/auth_service.dart';
import 'services/container_service.dart';
import 'services/dashboard_service.dart';
import 'services/plan_service.dart';
import 'services/product_service.dart';
import 'services/storage_service.dart';

void main() {
  WidgetsFlutterBinding.ensureInitialized();

  // 1. Initialize Infrastucture (Services)
  final storageService = StorageService(storage: const FlutterSecureStorage());
  final apiService = ApiService(storageService);
  final authService = AuthService(apiService, storageService);
  final dashboardService = DashboardService(apiService);
  final planService = PlanService(apiService);
  final productService = ProductService(apiService);
  final containerService = ContainerService(apiService);

  runApp(MyApp(
    authService: authService,
    storageService: storageService,
    dashboardService: dashboardService,
    planService: planService,
    productService: productService,
    containerService: containerService,
  ));
}

class MyApp extends StatelessWidget {
  final AuthService authService;
  final StorageService storageService;
  final DashboardService dashboardService;
  final PlanService planService;
  final ProductService productService;
  final ContainerService containerService;

  const MyApp({
    super.key,
    required this.authService,
    required this.storageService,
    required this.dashboardService,
    required this.planService,
    required this.productService,
    required this.containerService,
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
        ChangeNotifierProvider(
          create: (_) => PlanProvider(planService),
        ),
        ChangeNotifierProvider(
          create: (_) => ProductProvider(productService),
        ),
        ChangeNotifierProvider(
          create: (_) => ContainerProvider(containerService),
        ),
        ChangeNotifierProvider(
          create: (_) => PlanDetailProvider(planService),
        ),
      ],
      child: const AppRoot(),
    );
  }
}

class AppRoot extends StatefulWidget {
  const AppRoot({super.key});

  @override
  State<AppRoot> createState() => _AppRootState();
}

class _AppRootState extends State<AppRoot> {
  bool _initialized = false;
  late final GoRouter _router;
  
  @override
  void didChangeDependencies() {
    super.didChangeDependencies();
    // Initialize auth state only once
    if (!_initialized) {
      _initialized = true;
      final authProvider = context.read<AuthProvider>();
      _router = createRouter(authProvider);
      authProvider.initialize();
    }
  }

  @override
  Widget build(BuildContext context) {
    return MaterialApp.router(
      title: 'Load Stuffing Calculator',
      theme: AppTheme.lightTheme,
      routerConfig: _router,
      debugShowCheckedModeBanner: false,
    );
  }
}
