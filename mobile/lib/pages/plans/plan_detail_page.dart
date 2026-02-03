import 'package:flutter/material.dart';
import 'package:provider/provider.dart';
import '../../models/plan_detail_model.dart';
import '../../providers/plan_detail_provider.dart';
import '../../components/viewers/plan_visualizer_view.dart';
import '../../components/widgets/loading_state.dart';
import '../../components/widgets/error_state.dart';
import '../../components/widgets/empty_state.dart';

class PlanDetailPage extends StatefulWidget {
  final String planId;

  const PlanDetailPage({super.key, required this.planId});

  @override
  State<PlanDetailPage> createState() => _PlanDetailPageState();
}

class _PlanDetailPageState extends State<PlanDetailPage>
    with SingleTickerProviderStateMixin {
  late TabController _tabController;

  @override
  void initState() {
    super.initState();
    _tabController = TabController(length: 3, vsync: this);

    WidgetsBinding.instance.addPostFrameCallback((_) {
      context.read<PlanDetailProvider>().fetchPlanDetail(widget.planId);
    });
  }

  @override
  void dispose() {
    _tabController.dispose();
    super.dispose();
  }

  @override
  Widget build(BuildContext context) {
    return Scaffold(
      appBar: AppBar(
        title: const Text('Plan Details'),
        bottom: TabBar(
          controller: _tabController,
          tabs: const [
            Tab(icon: Icon(Icons.analytics), text: 'Overview'),
            Tab(icon: Icon(Icons.view_in_ar), text: '3D View'),
            Tab(icon: Icon(Icons.list), text: 'Items'),
          ],
        ),
      ),
      body: Consumer<PlanDetailProvider>(
        builder: (context, provider, _) {
          if (provider.isLoading) {
            return const LoadingState(message: 'Loading plan details...');
          }

          if (provider.error != null) {
            return ErrorState(
              message: provider.error!,
              onRetry: () =>
                  provider.fetchPlanDetail(widget.planId),
            );
          }

          if (provider.plan == null) {
            return const EmptyState(
              icon: Icons.inbox_outlined,
              message: 'Plan not found',
            );
          }

          return TabBarView(
            controller: _tabController,
            children: [
              _buildOverviewTab(provider.plan!),
              _buildVisualizationTab(provider.plan!),
              _buildItemsTab(provider.plan!),
            ],
          );
        },
      ),
      floatingActionButton: Consumer<PlanDetailProvider>(
        builder: (context, provider, _) {
          if (provider.plan == null) return const SizedBox.shrink();

          return FloatingActionButton.extended(
            onPressed: provider.isCalculating
                ? null
                : () => _showRecalculateDialog(context),
            icon: provider.isCalculating
                ? const SizedBox(
                    width: 20,
                    height: 20,
                    child: CircularProgressIndicator(
                      strokeWidth: 2,
                      color: Colors.white,
                    ),
                  )
                : const Icon(Icons.calculate),
            label: Text(provider.isCalculating ? 'Calculating...' : 'Recalculate'),
          );
        },
      ),
    );
  }

  Widget _buildOverviewTab(PlanDetailModel plan) {
    return SingleChildScrollView(
      padding: const EdgeInsets.all(16),
      child: Column(
        crossAxisAlignment: CrossAxisAlignment.start,
        children: [
          // Plan Header
          _buildPlanHeader(plan),
          const SizedBox(height: 24),

          // Stats Cards
          _buildStatsCards(plan.stats),
          const SizedBox(height: 24),

          // Container Info
          _buildContainerInfo(plan.container),
        ],
      ),
    );
  }

  Widget _buildPlanHeader(PlanDetailModel plan) {
    return Card(
      child: Padding(
        padding: const EdgeInsets.all(16),
        child: Column(
          crossAxisAlignment: CrossAxisAlignment.start,
          children: [
            Row(
              children: [
                Expanded(
                  child: Column(
                    crossAxisAlignment: CrossAxisAlignment.start,
                    children: [
                      Text(
                        plan.title,
                        style: const TextStyle(
                          fontSize: 20,
                          fontWeight: FontWeight.bold,
                        ),
                      ),
                      const SizedBox(height: 4),
                      Text(
                        plan.code,
                        style: TextStyle(
                          fontSize: 14,
                          color: Colors.grey[600],
                          fontFamily: 'monospace',
                        ),
                      ),
                    ],
                  ),
                ),
                _buildStatusBadge(plan.status),
              ],
            ),
            if (plan.notes != null) ...[
              const SizedBox(height: 12),
              const Divider(),
              const SizedBox(height: 8),
              Text(
                plan.notes!,
                style: TextStyle(
                  fontSize: 14,
                  color: Colors.grey[700],
                ),
              ),
            ],
          ],
        ),
      ),
    );
  }

  Widget _buildStatusBadge(String status) {
    Color color;
    switch (status.toUpperCase()) {
      case 'COMPLETED':
        color = Colors.green;
        break;
      case 'IN_PROGRESS':
        color = Colors.blue;
        break;
      case 'FAILED':
        color = Colors.red;
        break;
      case 'DRAFT':
        color = Colors.orange;
        break;
      default:
        color = Colors.grey;
    }

    return Container(
      padding: const EdgeInsets.symmetric(horizontal: 12, vertical: 6),
      decoration: BoxDecoration(
        color: color.withOpacity(0.1),
        borderRadius: BorderRadius.circular(12),
        border: Border.all(color: color.withOpacity(0.3)),
      ),
      child: Text(
        status.replaceAll('_', ' ').toUpperCase(),
        style: TextStyle(
          color: color,
          fontSize: 12,
          fontWeight: FontWeight.bold,
        ),
      ),
    );
  }

  Widget _buildStatsCards(PlanStats stats) {
    return Column(
      crossAxisAlignment: CrossAxisAlignment.start,
      children: [
        const Text(
          'Statistics',
          style: TextStyle(
            fontSize: 18,
            fontWeight: FontWeight.bold,
          ),
        ),
        const SizedBox(height: 12),
        Row(
          children: [
            Expanded(
              child: _buildStatCard(
                'Total Items',
                stats.totalItems.toString(),
                Icons.inventory_2_outlined,
                Colors.blue,
              ),
            ),
            const SizedBox(width: 12),
            Expanded(
              child: _buildStatCard(
                'Total Weight',
                '${stats.totalWeightKg.toStringAsFixed(1)} kg',
                Icons.scale_outlined,
                Colors.orange,
              ),
            ),
          ],
        ),
        const SizedBox(height: 12),
        _buildUtilizationCard(
          'Volume Utilization',
          stats.volumeUtilizationPct,
          '${stats.totalVolumeM3.toStringAsFixed(2)} m³',
          Colors.green,
        ),
        const SizedBox(height: 12),
        _buildUtilizationCard(
          'Weight Utilization',
          stats.weightUtilizationPct,
          '${stats.totalWeightKg.toStringAsFixed(1)} kg',
          Colors.purple,
        ),
      ],
    );
  }

  Widget _buildStatCard(String label, String value, IconData icon, Color color) {
    return Card(
      child: Padding(
        padding: const EdgeInsets.all(16),
        child: Column(
          crossAxisAlignment: CrossAxisAlignment.start,
          children: [
            Row(
              children: [
                Icon(icon, size: 20, color: color),
                const SizedBox(width: 8),
                Text(
                  label,
                  style: TextStyle(
                    fontSize: 12,
                    color: Colors.grey[600],
                  ),
                ),
              ],
            ),
            const SizedBox(height: 8),
            Text(
              value,
              style: const TextStyle(
                fontSize: 24,
                fontWeight: FontWeight.bold,
              ),
            ),
          ],
        ),
      ),
    );
  }

  Widget _buildUtilizationCard(
      String label, double percentage, String subtitle, Color color) {
    return Card(
      child: Padding(
        padding: const EdgeInsets.all(16),
        child: Column(
          crossAxisAlignment: CrossAxisAlignment.start,
          children: [
            Row(
              mainAxisAlignment: MainAxisAlignment.spaceBetween,
              children: [
                Text(
                  label,
                  style: const TextStyle(
                    fontSize: 14,
                    fontWeight: FontWeight.w500,
                  ),
                ),
                Text(
                  '${percentage.toStringAsFixed(1)}%',
                  style: TextStyle(
                    fontSize: 16,
                    fontWeight: FontWeight.bold,
                    color: color,
                  ),
                ),
              ],
            ),
            const SizedBox(height: 8),
            ClipRRect(
              borderRadius: BorderRadius.circular(4),
              child: LinearProgressIndicator(
                value: percentage / 100,
                minHeight: 8,
                backgroundColor: Colors.grey[200],
                valueColor: AlwaysStoppedAnimation<Color>(color),
              ),
            ),
            const SizedBox(height: 4),
            Text(
              subtitle,
              style: TextStyle(
                fontSize: 12,
                color: Colors.grey[600],
              ),
            ),
          ],
        ),
      ),
    );
  }

  Widget _buildContainerInfo(ContainerInfo container) {
    return Column(
      crossAxisAlignment: CrossAxisAlignment.start,
      children: [
        const Text(
          'Container Information',
          style: TextStyle(
            fontSize: 18,
            fontWeight: FontWeight.bold,
          ),
        ),
        const SizedBox(height: 12),
        Card(
          child: Padding(
            padding: const EdgeInsets.all(16),
            child: Column(
              crossAxisAlignment: CrossAxisAlignment.start,
              children: [
                if (container.name != null) ...[
                  _buildInfoRow('Name', container.name!),
                  const SizedBox(height: 8),
                ],
                _buildInfoRow(
                  'Dimensions',
                  '${container.lengthMm.toInt()} × ${container.widthMm.toInt()} × ${container.heightMm.toInt()} mm',
                ),
                const SizedBox(height: 8),
                _buildInfoRow(
                  'Volume',
                  '${container.volumeM3.toStringAsFixed(2)} m³',
                ),
                const SizedBox(height: 8),
                _buildInfoRow(
                  'Max Weight',
                  '${container.maxWeightKg.toStringAsFixed(1)} kg',
                ),
              ],
            ),
          ),
        ),
      ],
    );
  }

  Widget _buildInfoRow(String label, String value) {
    return Row(
      mainAxisAlignment: MainAxisAlignment.spaceBetween,
      children: [
        Text(
          label,
          style: TextStyle(
            fontSize: 14,
            color: Colors.grey[600],
          ),
        ),
        Text(
          value,
          style: const TextStyle(
            fontSize: 14,
            fontWeight: FontWeight.w500,
          ),
        ),
      ],
    );
  }

  Widget _buildVisualizationTab(PlanDetailModel plan) {
    return PlanVisualizerView(planId: plan.id);
  }

  Widget _buildItemsTab(PlanDetailModel plan) {
    if (plan.items.isEmpty) {
      return const EmptyState(
        icon: Icons.inbox_outlined,
        message: 'No items in this plan',
      );
    }

    return ListView.separated(
      padding: const EdgeInsets.all(16),
      itemCount: plan.items.length,
      separatorBuilder: (context, index) => const SizedBox(height: 8),
      itemBuilder: (context, index) {
        final item = plan.items[index];
        return _buildItemCard(item);
      },
    );
  }

  Widget _buildItemCard(PlanItem item) {
    return Card(
      child: ListTile(
        leading: Container(
          width: 16,
          height: 16,
          decoration: BoxDecoration(
            color: item.colorHex != null
                ? Color(int.parse(item.colorHex!.substring(1), radix: 16) +
                    0xFF000000)
                : Colors.grey,
            shape: BoxShape.circle,
          ),
        ),
        title: Text(
          item.label ?? item.itemId,
          style: const TextStyle(fontWeight: FontWeight.w500),
        ),
        subtitle: Text(
          '${item.quantity}× • ${item.lengthMm.toInt()}×${item.widthMm.toInt()}×${item.heightMm.toInt()} mm • ${item.weightKg} kg',
          style: TextStyle(fontSize: 12, color: Colors.grey[600]),
        ),
        trailing: Column(
          mainAxisAlignment: MainAxisAlignment.center,
          crossAxisAlignment: CrossAxisAlignment.end,
          children: [
            Text(
              '${item.totalWeightKg.toStringAsFixed(1)} kg',
              style: const TextStyle(
                fontWeight: FontWeight.bold,
                fontSize: 14,
              ),
            ),
            Text(
              '${item.totalVolumeM3.toStringAsFixed(3)} m³',
              style: TextStyle(
                fontSize: 11,
                color: Colors.grey[600],
              ),
            ),
          ],
        ),
      ),
    );
  }

  void _showRecalculateDialog(BuildContext context) {
    String strategy = 'bestfitdecreasing';
    String? goal;
    bool gravity = true;

    showDialog(
      context: context,
      builder: (context) => StatefulBuilder(
        builder: (context, setState) => AlertDialog(
          title: const Text('Recalculate Plan'),
          content: Column(
            mainAxisSize: MainAxisSize.min,
            crossAxisAlignment: CrossAxisAlignment.start,
            children: [
              const Text(
                'Choose packing strategy:',
                style: TextStyle(fontSize: 14, fontWeight: FontWeight.w500),
              ),
              const SizedBox(height: 12),
              DropdownButtonFormField<String>(
                value: strategy,
                decoration: const InputDecoration(
                  labelText: 'Strategy',
                  border: OutlineInputBorder(),
                  isDense: true,
                ),
                items: const [
                  DropdownMenuItem(
                      value: 'bestfitdecreasing', child: Text('Best Fit Decreasing')),
                  DropdownMenuItem(
                      value: 'minimizeboxes', child: Text('Minimize Boxes')),
                  DropdownMenuItem(value: 'greedy', child: Text('Greedy')),
                  DropdownMenuItem(value: 'parallel', child: Text('Parallel (Auto)')),
                ],
                onChanged: (value) {
                  setState(() {
                    strategy = value!;
                  });
                },
              ),
              const SizedBox(height: 12),
              CheckboxListTile(
                title: const Text('Gravity settling'),
                subtitle: const Text('Drop items to reduce floating'),
                value: gravity,
                onChanged: (value) {
                  setState(() {
                    gravity = value ?? true;
                  });
                },
                contentPadding: EdgeInsets.zero,
              ),
            ],
          ),
          actions: [
            TextButton(
              onPressed: () => Navigator.pop(context),
              child: const Text('Cancel'),
            ),
            FilledButton(
              onPressed: () {
                Navigator.pop(context);
                context.read<PlanDetailProvider>().recalculate(
                      widget.planId,
                      strategy: strategy,
                      goal: goal,
                      gravity: gravity,
                    );
              },
              child: const Text('Recalculate'),
            ),
          ],
        ),
      ),
    );
  }
}
