import 'package:flutter/material.dart';
import 'package:provider/provider.dart';
import 'package:go_router/go_router.dart';
import '../../components/buttons/quick_action_button.dart';
import '../../components/cards/stat_card.dart';
import '../../components/cards/app_card.dart';
import '../../components/widgets/status_badge.dart';
import '../../dtos/dashboard_dto.dart';
import '../../providers/auth_provider.dart';
import '../../providers/dashboard_provider.dart';
import '../../config/theme.dart';

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
      appBar: AppBar(
        title: const Text('Dashboard'),
      ),
      body: Consumer<DashboardProvider>(
        builder: (context, provider, child) {
          if (provider.isLoading && provider.stats == null) {
            return const Center(child: CircularProgressIndicator());
          }

          if (provider.error != null && provider.stats == null) {
            return Center(
              child: Column(
                mainAxisAlignment: MainAxisAlignment.center,
                children: [
                  Text('Error: ${provider.error}',
                      style: const TextStyle(color: AppColors.error)),
                  const SizedBox(height: 16),
                  ElevatedButton(
                    onPressed: () => provider.loadData(),
                    child: const Text('Retry'),
                  ),
                ],
              ),
            );
          }

          return RefreshIndicator(
            onRefresh: provider.loadData,
            child: ListView(
              padding: const EdgeInsets.all(20.0),
              children: [
                if (provider.stats != null) ...[
                  _buildStatsSection(context, provider.stats!),
                  const SizedBox(height: 32),
                  _buildQuickActions(context),
                  const SizedBox(height: 32),
                ],
                Text(
                  'Recent Plans',
                  style: Theme.of(context).textTheme.titleLarge?.copyWith(
                    fontWeight: FontWeight.bold,
                  ),
                ),
                const SizedBox(height: 16),
                if (provider.recentPlans.isEmpty)
                  const Padding(
                    padding: EdgeInsets.all(32.0),
                    child: Center(
                      child: Text(
                        'No plans found',
                        style: TextStyle(color: AppColors.textSecondary),
                      ),
                    ),
                  )
                else
                  ...provider.recentPlans.map((plan) => Padding(
                        padding: const EdgeInsets.only(bottom: 12),
                        child: AppCard(
                          onTap: () {
                             context.push('/plans/${plan.id}');
                          },
                          child: Row(
                            children: [
                              Container(
                                padding: const EdgeInsets.all(12),
                                decoration: BoxDecoration(
                                  color: AppColors.primary.withValues(alpha: 0.1),
                                  borderRadius: BorderRadius.circular(12),
                                ),
                                child: const Icon(Icons.inventory_2_outlined, color: AppColors.primary),
                              ),
                              const SizedBox(width: 16),
                              Expanded(
                                child: Column(
                                  crossAxisAlignment: CrossAxisAlignment.start,
                                  children: [
                                    Text(
                                      plan.title.isEmpty ? plan.code : plan.title,
                                      style: const TextStyle(fontWeight: FontWeight.w600, fontSize: 16),
                                    ),
                                    const SizedBox(height: 4),
                                    Text(
                                      plan.code,
                                      style: const TextStyle(
                                        fontFamily: 'monospace',
                                        fontSize: 12,
                                        color: AppColors.textSecondary,
                                      ),
                                    ),
                                  ],
                                ),
                              ),
                              StatusBadge.fromStatus(plan.status),
                            ],
                          ),
                        ),
                      )),
              ],
            ),
          );
        },
      ),
    );
  }

  Widget _buildStatsSection(BuildContext context, DashboardStatsDto stats) {
    // Get user role
    final userRole = context.read<AuthProvider>().user?.role.toLowerCase() ?? '';
    final List<Widget> cards = [];

    // Personal User Logic (matches Web Client)
    if (userRole == 'personal' || userRole == 'admin' || userRole == 'owner') {
       if (stats.admin != null) {
          cards.add(StatCard(
            title: 'Active Shipments',
            value: stats.admin!.activeShipments.toString(),
            icon: Icons.local_shipping,
            color: AppColors.info,
          ));
       }
       
       if (stats.planner != null) {
          cards.add(StatCard(
            title: 'Avg Utilization',
            value: '${stats.planner!.avgUtilization.toStringAsFixed(1)}%',
            icon: Icons.pie_chart,
            color: Colors.purple,
          ));
          cards.add(StatCard(
            title: 'Items Processed',
            value: stats.planner!.itemsProcessed.toString(),
            icon: Icons.inventory_2,
            color: AppColors.secondary,
          ));
       }

       if (stats.admin != null) {
          cards.add(StatCard(
            title: 'Container Types',
            value: stats.admin!.containerTypes.toString(),
            icon: Icons.view_in_ar,
            color: Colors.teal,
          ));
       }
    } else if (userRole == 'planner' && stats.planner != null) {
       cards.add(StatCard(
        title: 'Pending Plans',
        value: stats.planner!.pendingPlans.toString(),
        icon: Icons.pending_actions,
        color: AppColors.warning,
      ));
      cards.add(StatCard(
        title: 'Completed Today',
        value: stats.planner!.completedToday.toString(),
        icon: Icons.check_circle,
        color: AppColors.success,
      ));
      cards.add(StatCard(
        title: 'Avg Utilization',
        value: '${stats.planner!.avgUtilization.toStringAsFixed(1)}%',
        icon: Icons.pie_chart,
        color: Colors.purple,
      ));
       cards.add(StatCard(
        title: 'Items Processed',
        value: stats.planner!.itemsProcessed.toString(),
        icon: Icons.inventory_2,
        color: AppColors.secondary,
      ));
    } else if (userRole == 'operator' && stats.operator != null) {
       cards.add(StatCard(
        title: 'Active Loads',
        value: stats.operator!.activeLoads.toString(),
        icon: Icons.local_shipping,
        color: AppColors.warning,
      ));
       cards.add(StatCard(
        title: 'Completed',
        value: stats.operator!.completed.toString(),
        icon: Icons.task_alt,
        color: Colors.teal,
      ));
    }
    
    if (cards.isEmpty) {
       if (stats.admin != null) {
          cards.add(StatCard(title: 'Active Shipments', value: stats.admin!.activeShipments.toString(), icon: Icons.local_shipping, color: AppColors.info));
       }
    }

    if (cards.isEmpty) return const SizedBox.shrink();

    return GridView.count(
      shrinkWrap: true,
      physics: const NeverScrollableScrollPhysics(),
      crossAxisCount: 2,
      crossAxisSpacing: 16,
      mainAxisSpacing: 16,
      childAspectRatio: 1.1,
      children: cards,
    );
  }

  Widget _buildQuickActions(BuildContext context) {
      return Column(
        crossAxisAlignment: CrossAxisAlignment.start,
        children: [
          Text(
            'Quick Actions', 
            style: Theme.of(context).textTheme.titleLarge?.copyWith(fontWeight: FontWeight.bold)
          ),
          const SizedBox(height: 16),
          Row(
            mainAxisAlignment: MainAxisAlignment.spaceBetween,
            children: [
               QuickActionButton(label: 'Products', icon: Icons.inventory, onTap: () {
                context.go('/products');
              }),
              QuickActionButton(label: 'Containers', icon: Icons.view_in_ar, onTap: () {
                context.go('/containers');
              }),
               QuickActionButton(label: 'New Plan', icon: Icons.add_circle_outline, onTap: () {
                context.push('/plans/new');
              }),
            ],
          ),
        ],
      );
  }
}
