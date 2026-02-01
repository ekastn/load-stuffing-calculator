import 'package:flutter/material.dart';
import 'package:provider/provider.dart';
import '../../dtos/dashboard_dto.dart';
import '../../providers/auth_provider.dart';
import '../../providers/dashboard_provider.dart';

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
                            backgroundColor: _getStatusColor(plan.status),
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
          cards.add(_buildStatCard(
            context,
            'Active Shipments',
            stats.admin!.activeShipments.toString(),
            Icons.local_shipping,
            Colors.blue,
          ));
       }
       
       // Show Planner-derived stats if available
       if (stats.planner != null) {
          cards.add(_buildStatCard(
            context,
            'Avg Utilization',
            '${stats.planner!.avgUtilization.toStringAsFixed(1)}%',
            Icons.pie_chart,
            Colors.purple,
          ));
          cards.add(_buildStatCard(
            context,
            'Items Processed',
            stats.planner!.itemsProcessed.toString(),
            Icons.inventory_2,
            Colors.orange,
          ));
       }

       // Show Admin-derived container stats if available
       if (stats.admin != null) {
          cards.add(_buildStatCard(
            context,
            'Container Types',
            stats.admin!.containerTypes.toString(),
            Icons.view_in_ar,
            Colors.teal,
          ));
       }
    } else if (userRole == 'planner' && stats.planner != null) {
       cards.add(_buildStatCard(
        context,
        'Pending Plans',
        stats.planner!.pendingPlans.toString(),
        Icons.pending_actions,
        Colors.orange,
      ));
      cards.add(_buildStatCard(
        context,
        'Completed Today',
        stats.planner!.completedToday.toString(),
        Icons.check_circle,
        Colors.green,
      ));
      cards.add(_buildStatCard(
        context,
        'Avg Utilization',
        '${stats.planner!.avgUtilization.toStringAsFixed(1)}%',
        Icons.pie_chart,
        Colors.purple,
      ));
       cards.add(_buildStatCard(
        context,
        'Items Processed',
        stats.planner!.itemsProcessed.toString(),
        Icons.inventory_2,
        Colors.orange,
      ));
    } else if (userRole == 'operator' && stats.operator != null) {
       cards.add(_buildStatCard(
        context,
        'Active Loads',
        stats.operator!.activeLoads.toString(),
        Icons.forklift,
        Colors.amber.shade900,
      ));
       cards.add(_buildStatCard(
        context,
        'Completed',
        stats.operator!.completed.toString(),
        Icons.task_alt,
        Colors.teal,
      ));
    }
    

    // Fallback if cards is empty
    if (cards.isEmpty) {
       // If absolutely no stats were matched above, try to show at least something
       if (stats.admin != null) {
          // This duplicates logic above but is a safety net if roles didn't match
          cards.add(_buildStatCard(context, 'Active Shipments', stats.admin!.activeShipments.toString(), Icons.local_shipping, Colors.blue));
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
                  _buildActionButton(context, 'New Plan', Icons.add_circle_outline, () {
                    // Navigate to Create Plan
                  }),
                  _buildActionButton(context, 'New Container', Icons.add_box_outlined, () {
                    // Navigate to Create Container
                  }),
                  _buildActionButton(context, 'View Plans', Icons.list_alt, () {
                    // Navigate to Plans
                     // Default Tab Controller or GoRouter?
                     // context.go('/plans'); 
                     // But we are in Dashboard tab, switching tab is tricky without context.go
                  }),
                ],
              ),
            ],
          ),
        ),
      );
  }

  Widget _buildActionButton(BuildContext context, String label, IconData icon, VoidCallback onTap) {
    return InkWell(
      onTap: onTap,
      child: Column(
        children: [
          Container(
            padding: const EdgeInsets.all(12),
            decoration: BoxDecoration(
              color: Theme.of(context).colorScheme.primaryContainer,
              shape: BoxShape.circle,
            ),
            child: Icon(icon, color: Theme.of(context).colorScheme.primary),
          ),
          const SizedBox(height: 8),
          Text(label, style: const TextStyle(fontSize: 12)),
        ],
      ),
    );
  }

  Widget _buildStatCard(BuildContext context, String title, String value,
      IconData icon, Color color) {
    return Card(
      elevation: 2,
      child: Padding(
        padding: const EdgeInsets.all(12.0),
        child: Column(
          mainAxisAlignment: MainAxisAlignment.center,
          children: [
            Icon(icon, size: 32, color: color),
            const SizedBox(height: 8),
            Text(
              value,
              style: Theme.of(context).textTheme.headlineSmall?.copyWith(
                    fontWeight: FontWeight.bold,
                  ),
            ),
            Text(
              title,
              style: Theme.of(context).textTheme.bodySmall,
              textAlign: TextAlign.center,
            ),
          ],
        ),
      ),
    );
  }

  Color _getStatusColor(String status) {
    switch (status.toUpperCase()) {
      case 'COMPLETED':
        return Colors.green.shade100;
      case 'DRAFT':
        return Colors.grey.shade200;
      case 'IN_PROGRESS':
        return Colors.blue.shade100;
      case 'FAILED':
        return Colors.red.shade100;
      default:
        return Colors.grey.shade100;
    }
  }
}
