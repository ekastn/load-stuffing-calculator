import 'package:flutter/material.dart';
import 'package:provider/provider.dart';
import 'package:go_router/go_router.dart';
import '../../providers/auth_provider.dart';
import '../../providers/dashboard_provider.dart';
import '../../components/sections/dashboard_stats_section.dart';
import '../../components/sections/recent_plans_list.dart';
import '../../components/sections/quick_actions_section.dart';
import '../../components/widgets/loading_state.dart';
import '../../components/widgets/error_state.dart';

class DashboardPage extends StatefulWidget {
  const DashboardPage({super.key});

  @override
  State<DashboardPage> createState() => _DashboardPageState();
}

class _DashboardPageState extends State<DashboardPage> {
  @override
  void initState() {
    super.initState();
    WidgetsBinding.instance.addPostFrameCallback((_) {
      context.read<DashboardProvider>().loadData();
    });
  }

  @override
  Widget build(BuildContext context) {
    return Scaffold(
      appBar: AppBar(title: const Text('Dashboard')),
      body: Consumer<DashboardProvider>(
        builder: (context, provider, child) {
          if (provider.isLoading && provider.stats == null) {
            return const LoadingState();
          }

          if (provider.error != null && provider.stats == null) {
            return ErrorState(
              message: provider.error!,
              onRetry: () => provider.loadData(),
            );
          }

          return RefreshIndicator(
            onRefresh: provider.loadData,
            child: ListView(
              padding: const EdgeInsets.all(20.0),
              children: [
                if (provider.stats != null) ...[
                  DashboardStatsSection(
                    stats: provider.stats!,
                    userRole: context.watch<AuthProvider>().user?.role ?? '',
                  ),
                  const SizedBox(height: 32),
                  QuickActionsSection(
                    onProductsTap: () => context.go('/products'),
                    onContainersTap: () => context.go('/containers'),
                    onNewPlanTap: () => context.push('/plans/new'),
                  ),
                  const SizedBox(height: 32),
                ],
                RecentPlansList(
                  plans: provider.recentPlans,
                  onPlanTap: (planId) => context.push('/plans/$planId'),
                ),
              ],
            ),
          );
        },
      ),
    );
  }
}
