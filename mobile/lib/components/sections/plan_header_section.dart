import 'package:flutter/material.dart';
import '../../models/plan_detail_model.dart';
import '../cards/app_card.dart';
import '../cards/stat_card.dart';
import '../cards/utilization_progress_card.dart';
import '../widgets/status_badge.dart';
import '../../config/theme.dart';

class PlanHeaderSection extends StatelessWidget {
  final PlanDetailModel plan;

  const PlanHeaderSection({super.key, required this.plan});

  @override
  Widget build(BuildContext context) {
    return Column(
      crossAxisAlignment: CrossAxisAlignment.start,
      children: [
        // Plan Header
        _buildPlanHeader(context),
        const SizedBox(height: 24),

        // Stats Cards
        _buildStatsCards(context),
      ],
    );
  }

  Widget _buildPlanHeader(BuildContext context) {
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
                      style: Theme.of(context).textTheme.headlineSmall
                          ?.copyWith(fontWeight: FontWeight.bold),
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

  Widget _buildStatsCards(BuildContext context) {
    final stats = plan.stats;
    return Column(
      crossAxisAlignment: CrossAxisAlignment.start,
      children: [
        Text(
          'Statistics',
          style: Theme.of(
            context,
          ).textTheme.titleLarge?.copyWith(fontWeight: FontWeight.bold),
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
        UtilizationProgressCard(
          label: 'Volume Utilization',
          percentage: stats.volumeUtilizationPct,
          subtitle: '${stats.totalVolumeM3.toStringAsFixed(2)} m³',
          color: AppColors.success,
        ),
        const SizedBox(height: 16),
        UtilizationProgressCard(
          label: 'Weight Utilization',
          percentage: stats.weightUtilizationPct,
          subtitle: '${stats.totalWeightKg.toStringAsFixed(1)} kg',
          color: Colors.purple,
        ),
      ],
    );
  }
}
