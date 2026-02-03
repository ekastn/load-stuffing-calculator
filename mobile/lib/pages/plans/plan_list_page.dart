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
import '../../utils/ui_helpers.dart';

class PlanListPage extends StatefulWidget {
  const PlanListPage({super.key});

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
      body: Consumer<PlanProvider>(
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
                leading: CircleAvatar(
                  backgroundColor: UiHelpers.getStatusColor(plan.status),
                  child: const Icon(Icons.inventory_2, color: Colors.black87),
                ),
                title: plan.title,
                subtitle:
                    '${plan.code} • ${plan.totalItems} items • ${plan.totalWeightKg.toStringAsFixed(1)}kg\n${plan.status.replaceAll('_', ' ').toUpperCase()}${plan.volumeUtilizationPct != null ? " • ${plan.volumeUtilizationPct!.toStringAsFixed(1)}% utilized" : ""}',
                onTap: () => context.push('/plans/${plan.id}'),
                onDelete: () => _confirmDelete(context, plan),
              );
            },
          );
        },
      ),
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
}
