import 'package:flutter/material.dart';
import 'package:provider/provider.dart';
import 'package:go_router/go_router.dart';
import '../../components/buttons/quick_action_button.dart';
import '../../components/cards/stat_card.dart';
import '../../dtos/dashboard_dto.dart';
import '../../providers/auth_provider.dart';
import '../../providers/dashboard_provider.dart';
import '../../utils/ui_helpers.dart';

class DashboardPage extends StatefulWidget {
  const DashboardPage({super.key});

  @override
  State<DashboardPage> createState() => _DashboardPageState();
}

class _DashboardPageState extends State<DashboardPage> {
  @override
  void initState() {
    super.initState();
    // Fetch data on fresh load
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
                      style: const TextStyle(color: Colors.red)),
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
              padding: const EdgeInsets.all(16.0),
              children: [
                if (provider.stats != null) ...[
                  _buildStatsSection(context, provider.stats!),
                  const SizedBox(height: 24),
                  _buildQuickActions(context),
                  const SizedBox(height: 24),
                ],
                Text(
                  'Recent Plans',
                  style: Theme.of(context).textTheme.titleLarge,
                ),
                const SizedBox(height: 8),
                if (provider.recentPlans.isEmpty)
                  const Padding(
                    padding: EdgeInsets.all(16.0),
                    child: Center(child: Text('No plans found.')),
                  )
                else
                  ...provider.recentPlans.map((plan) => Card(
                        child: ListTile(
                          title: Text(plan.title),
                          subtitle: Text(plan.planCode),
                          trailing: Chip(
                            label: Text(
                              plan.status,
                              style: const TextStyle(fontSize: 10),
                            ),
                            backgroundColor: UiHelpers.getStatusColor(plan.status),
                          ),
                          onTap: () {
                            // TODO: Navigate to details
                            // context.push('/plans/${plan.planId}');
                          },
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
       // Show Admin-derived stats if available
       if (stats.admin != null) {
          cards.add(StatCard(
            title: 'Active Shipments',
            value: stats.admin!.activeShipments.toString(),
            icon: Icons.local_shipping,
            color: Colors.blue,
          ));
       }
       
       // Show Planner-derived stats if available
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
            color: Colors.orange,
          ));
       }

       // Show Admin-derived container stats if available
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
        color: Colors.orange,
      ));
      cards.add(StatCard(
        title: 'Completed Today',
        value: stats.planner!.completedToday.toString(),
        icon: Icons.check_circle,
        color: Colors.green,
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
        color: Colors.orange,
      ));
    } else if (userRole == 'operator' && stats.operator != null) {
       cards.add(StatCard(
        title: 'Active Loads',
        value: stats.operator!.activeLoads.toString(),
        icon: Icons.forklift,
        color: Colors.amber.shade900,
      ));
       cards.add(StatCard(
        title: 'Completed',
        value: stats.operator!.completed.toString(),
        icon: Icons.task_alt,
        color: Colors.teal,
      ));
    }
    

    // Fallback if cards is empty
    if (cards.isEmpty) {
       // If absolutely no stats were matched above, try to show at least something
       if (stats.admin != null) {
          // This duplicates logic above but is a safety net if roles didn't match
          cards.add(StatCard(title: 'Active Shipments', value: stats.admin!.activeShipments.toString(), icon: Icons.local_shipping, color: Colors.blue));
       }
    }

    if (cards.isEmpty) return const SizedBox.shrink();

    return GridView.count(
      shrinkWrap: true,
      physics: const NeverScrollableScrollPhysics(),
      crossAxisCount: 2,
      crossAxisSpacing: 16,
      mainAxisSpacing: 16,
      childAspectRatio: 1.5,
      children: cards,
    );
  }

  Widget _buildQuickActions(BuildContext context) {
      return Card(
        child: Padding(
          padding: const EdgeInsets.all(16.0),
          child: Column(
            crossAxisAlignment: CrossAxisAlignment.start,
            children: [
              Text('Quick Actions', style: Theme.of(context).textTheme.titleMedium),
              const SizedBox(height: 16),
              Row(
                mainAxisAlignment: MainAxisAlignment.spaceAround,
                children: [
                   QuickActionButton(label: 'Products', icon: Icons.inventory, onTap: () {
                    context.push('/products');
                  }),
                  QuickActionButton(label: 'Containers', icon: Icons.view_in_ar, onTap: () {
                    context.push('/containers');
                  }),
                   QuickActionButton(label: 'New Plan', icon: Icons.add_circle_outline, onTap: () {
                    // Navigate to Create Plan
                    // context.push('/plans/new');
                  }),
                ],
              ),
            ],
          ),
        ),
      );
  }


}
