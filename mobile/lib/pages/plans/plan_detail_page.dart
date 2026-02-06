import 'package:flutter/material.dart';
import 'package:provider/provider.dart';
import '../../models/plan_detail_model.dart';
import '../../providers/plan_detail_provider.dart';
import '../../components/viewers/plan_visualizer_view.dart';
import '../../components/widgets/loading_state.dart';
import '../../components/widgets/error_state.dart';
import '../../components/widgets/empty_state.dart';
import '../../components/cards/app_card.dart';
import '../../components/cards/stat_card.dart';
import '../../components/widgets/status_badge.dart';
import '../../config/theme.dart';

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
    return AppCard(
      child: Column(
        crossAxisAlignment: CrossAxisAlignment.start,
        children: [
          Row(
            crossAxisAlignment: CrossAxisAlignment.start,
            children: [
              Expanded(
                child: Column(
                  crossAxisAlignment: CrossAxisAlignment.start,
                  children: [
                    Text(
                      plan.title.isEmpty ? plan.code : plan.title,
                      style: Theme.of(context).textTheme.headlineSmall?.copyWith(
                        fontWeight: FontWeight.bold,
                      ),
                    ),
                    if (plan.title.isNotEmpty) ...[
                      const SizedBox(height: 4),
                      Text(
                        plan.code,
                        style: TextStyle(
                          fontSize: 14,
                          color: AppColors.textSecondary,
                          fontFamily: 'monospace',
                        ),
                      ),
                    ],
                  ],
                ),
              ),
              const SizedBox(width: 8),
              StatusBadge.fromStatus(plan.status),
            ],
          ),
          if (plan.notes != null && plan.notes!.isNotEmpty) ...[
            const SizedBox(height: 16),
            const Divider(),
            const SizedBox(height: 12),
            Text(
              'Notes:',
              style: Theme.of(context).textTheme.labelMedium?.copyWith(
                fontWeight: FontWeight.bold,
                color: AppColors.textPrimary,
              ),
            ),
             const SizedBox(height: 4),
            Text(
              plan.notes!,
              style: TextStyle(
                fontSize: 14,
                color: AppColors.textSecondary,
                height: 1.5,
              ),
            ),
          ],
        ],
      ),
    );
  }

  Widget _buildStatsCards(PlanStats stats) {
    return Column(
      crossAxisAlignment: CrossAxisAlignment.start,
      children: [
        Text(
          'Statistics',
          style: Theme.of(context).textTheme.titleLarge?.copyWith(
            fontWeight: FontWeight.bold,
          ),
        ),
        const SizedBox(height: 16),
        Row(
          children: [
            Expanded(
              child: StatCard(
                title: 'Total Items',
                value: stats.totalItems.toString(),
                icon: Icons.inventory_2_outlined,
                color: AppColors.info,
              ),
            ),
            const SizedBox(width: 16),
            Expanded(
              child: StatCard(
                title: 'Total Weight',
                value: '${stats.totalWeightKg.toStringAsFixed(1)} kg',
                icon: Icons.scale_outlined,
                color: AppColors.secondary,
              ),
            ),
          ],
        ),
        const SizedBox(height: 16),
        _buildUtilizationCard(
          'Volume Utilization',
          stats.volumeUtilizationPct,
          '${stats.totalVolumeM3.toStringAsFixed(2)} m³',
          AppColors.success,
        ),
         const SizedBox(height: 16),
        _buildUtilizationCard(
          'Weight Utilization',
          stats.weightUtilizationPct,
          '${stats.totalWeightKg.toStringAsFixed(1)} kg',
          Colors.purple,
        ),
      ],
    );
  }

  Widget _buildUtilizationCard(
      String label, double percentage, String subtitle, Color color) {
    return AppCard(
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
              backgroundColor: AppColors.background,
              valueColor: AlwaysStoppedAnimation<Color>(color),
            ),
          ),
          const SizedBox(height: 4),
          Text(
            subtitle,
            style: TextStyle(
              fontSize: 12,
              color: AppColors.textSecondary,
            ),
          ),
        ],
      ),
    );
  }

  Widget _buildContainerInfo(ContainerInfo container) {
    return Column(
      crossAxisAlignment: CrossAxisAlignment.start,
      children: [
        Text(
          'Container Information',
          style: Theme.of(context).textTheme.titleLarge?.copyWith(
            fontWeight: FontWeight.bold,
          ),
        ),
        const SizedBox(height: 16),
        AppCard(
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
    return AppCard(
      padding: const EdgeInsets.symmetric(horizontal: 16, vertical: 12),
      child: Row(
        children: [
            Container(
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
            const SizedBox(width: 16),
            Expanded(
              child: Column(
                crossAxisAlignment: CrossAxisAlignment.start,
                children: [
                   Text(
                    item.label ?? item.itemId,
                    style: const TextStyle(fontWeight: FontWeight.w600, fontSize: 16),
                  ),
                   const SizedBox(height: 4),
                  Text(
                    '${item.quantity}× • ${item.lengthMm.toInt()}×${item.widthMm.toInt()}×${item.heightMm.toInt()} mm • ${item.weightKg} kg',
                    style: TextStyle(fontSize: 13, color: AppColors.textSecondary),
                  ),
                ],
              ),
            ),
            Column(
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
                    fontSize: 12,
                    color: AppColors.textSecondary,
                  ),
                ),
              ],
            ),
        ],
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
