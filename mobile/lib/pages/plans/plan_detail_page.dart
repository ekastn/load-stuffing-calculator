import 'package:flutter/material.dart';
import 'package:provider/provider.dart';
import '../../providers/plan_detail_provider.dart';
import '../../components/viewers/plan_visualizer_view.dart';
import '../../components/widgets/loading_state.dart';
import '../../components/widgets/error_state.dart';
import '../../components/widgets/empty_state.dart';
import '../../components/dialogs/recalculate_dialog.dart';
import '../../components/sections/plan_overview_tab_section.dart';
import '../../components/sections/plan_items_tab_section.dart';

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
              onRetry: () => provider.fetchPlanDetail(widget.planId),
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
              PlanOverviewTabSection(plan: provider.plan!),
              PlanVisualizerView(planId: provider.plan!.id),
              PlanItemsTabSection(plan: provider.plan!),
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
            label: Text(
              provider.isCalculating ? 'Calculating...' : 'Recalculate',
            ),
          );
        },
      ),
    );
  }

  void _showRecalculateDialog(BuildContext context) {
    showDialog(
      context: context,
      builder: (context) => RecalculateDialog(
        onRecalculate: (strategy, goal, gravity) {
          context.read<PlanDetailProvider>().recalculate(
            widget.planId,
            strategy: strategy,
            goal: goal,
            gravity: gravity,
          );
        },
      ),
    );
  }
}
