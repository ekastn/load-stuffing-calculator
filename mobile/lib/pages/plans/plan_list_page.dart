import 'package:flutter/material.dart';
import 'package:provider/provider.dart';
import 'package:go_router/go_router.dart';
import '../../providers/plan_provider.dart';
import '../../models/plan_model.dart';
import '../../components/widgets/empty_state.dart';
import '../../components/widgets/loading_state.dart';
import '../../components/widgets/error_state.dart';
import '../../components/cards/resource_list_item.dart';
import '../../components/dialogs/confirm_dialog.dart';
import '../../config/theme.dart';
import '../../components/widgets/status_badge.dart';

class PlanListPage extends StatefulWidget {
  final bool isEmbedded;
  const PlanListPage({super.key, this.isEmbedded = false});

  @override
  State<PlanListPage> createState() => _PlanListPageState();
}

class _PlanListPageState extends State<PlanListPage> {
  @override
  void initState() {
    super.initState();
    WidgetsBinding.instance.addPostFrameCallback((_) {
      context.read<PlanProvider>().fetchPlans();
    });
  }

  @override
  Widget build(BuildContext context) {
    final content = Consumer<PlanProvider>(
        builder: (context, provider, child) {
          if (provider.isLoading && provider.plans.isEmpty) {
            return const LoadingState();
          }

          if (provider.error != null && provider.plans.isEmpty) {
            return ErrorState(
              message: provider.error!,
              onRetry: () => provider.fetchPlans(),
            );
          }

          if (provider.plans.isEmpty) {
            return const EmptyState(
              message: 'No plans found.\nTap + to create a new plan.',
              icon: Icons.inventory_2_outlined,
            );
          }

          return ListView.builder(
            itemCount: provider.plans.length,
            itemBuilder: (context, index) {
              final plan = provider.plans[index];
              return ResourceListItem(
                title: plan.title.isEmpty ? plan.code : plan.title,
                subtitle: Text(
                  '${plan.title.isEmpty ? '' : '${plan.code} • '}${plan.totalItems} items',
                  style: const TextStyle(color: AppColors.textSecondary),
                ),
                trailing: StatusBadge.fromStatus(plan.status),
                content: Column(
                  children: [
                    _buildProgressBar(
                      context,
                      label: 'Volume',
                      percentage: plan.volumeUtilizationPct ?? 0,
                      color: AppColors.info,
                    ),
                  ],
                ),
                actions: [
                  _buildActionButton(
                    context,
                    label: 'View',
                    icon: Icons.visibility_outlined,
                    color: AppColors.primary,
                    onTap: () => context.push('/plans/${plan.id}'),
                  ),
                  _buildActionButton(
                    context,
                    label: 'Calculate',
                    icon: Icons.bolt_outlined,
                    color: AppColors.primary,
                    onTap: () {
                      // TODO: Implement calculation trigger
                    },
                  ),
                  _buildActionButton(
                    context,
                    icon: Icons.delete_outline,
                    color: AppColors.error,
                    onTap: () => _confirmDelete(context, plan),
                    isIconOnly: true,
                  ),
                ],
                onTap: () => context.push('/plans/${plan.id}'),
              );
            },
          );
        },
      );

    if (widget.isEmbedded) {
      return Scaffold(
        body: content,
        floatingActionButton: FloatingActionButton(
          onPressed: () => context.push('/plans/new'),
          child: const Icon(Icons.add),
        ),
      );
    }

    return Scaffold(
      appBar: AppBar(
        title: const Text('Plans'),
        actions: [
          IconButton(
            icon: const Icon(Icons.refresh),
            onPressed: () => context.read<PlanProvider>().fetchPlans(),
          ),
        ],
      ),
      floatingActionButton: FloatingActionButton(
        onPressed: () => context.push('/plans/new'),
        child: const Icon(Icons.add),
      ),
      body: content,
    );
  }



  void _confirmDelete(BuildContext context, PlanModel plan) {
    ConfirmDialog.show(
      context: context,
      title: 'Delete Plan',
      content: 'Are you sure you want to delete "${plan.title}"?',
      confirmText: 'Delete',
      isDangerous: true,
      onConfirm: () async {
        try {
          await context.read<PlanProvider>().deletePlan(plan.id);
          if (context.mounted) {
            ScaffoldMessenger.of(context).showSnackBar(
              const SnackBar(content: Text('Plan deleted')),
            );
          }
        } catch (e) {
          if (context.mounted) {
            ScaffoldMessenger.of(context).showSnackBar(
              SnackBar(content: Text('Failed to delete: $e')),
            );
          }
        }
      },
    );
  }

  Widget _buildProgressBar(BuildContext context, {
    required String label,
    required double percentage,
    required Color color,
  }) {
    return Column(
      children: [
        Row(
          mainAxisAlignment: MainAxisAlignment.spaceBetween,
          children: [
            Text(
              label,
              style: const TextStyle(
                fontSize: 14,
                color: AppColors.textSecondary,
              ),
            ),
            Text(
              '${percentage.toStringAsFixed(0)}%',
              style: const TextStyle(
                fontSize: 14,
                fontWeight: FontWeight.bold,
                color: AppColors.textPrimary,
              ),
            ),
          ],
        ),
        const SizedBox(height: 8),
        ClipRRect(
          borderRadius: BorderRadius.circular(4),
          child: LinearProgressIndicator(
            value: percentage / 100,
            backgroundColor: AppColors.border,
            valueColor: AlwaysStoppedAnimation<Color>(color),
            minHeight: 8,
          ),
        ),
      ],
    );
  }

  Widget _buildActionButton(BuildContext context, {
    String? label,
    required IconData icon,
    required Color color,
    required VoidCallback onTap,
    bool isIconOnly = false,
  }) {
    if (isIconOnly) {
      return InkWell(
        onTap: onTap,
        borderRadius: BorderRadius.circular(8),
        child: Padding(
          padding: const EdgeInsets.all(8.0),
          child: Icon(icon, color: color, size: 20),
        ),
      );
    }

    return TextButton.icon(
      onPressed: onTap,
      icon: Icon(icon, size: 18, color: color),
      label: Text(
        label!,
        style: TextStyle(
          color: color,
          fontWeight: FontWeight.w600,
        ),
      ),
      style: TextButton.styleFrom(
        padding: const EdgeInsets.symmetric(horizontal: 12, vertical: 8),
        minimumSize: Size.zero,
        tapTargetSize: MaterialTapTargetSize.shrinkWrap,
      ),
    );
  }
}
