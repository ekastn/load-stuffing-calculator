import 'package:flutter/material.dart';
import '../../components/cards/stat_card.dart';
import '../../dtos/dashboard_dto.dart';
import '../../config/theme.dart';

class DashboardStatsSection extends StatelessWidget {
  final DashboardStatsDto stats;
  final String userRole;

  const DashboardStatsSection({
    super.key,
    required this.stats,
    required this.userRole,
  });

  @override
  Widget build(BuildContext context) {
    final cards = _buildStatCards();

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

  List<Widget> _buildStatCards() {
    final role = userRole.toLowerCase();
    final List<Widget> cards = [];

    // Personal/Admin/Owner Logic (matches Web Client)
    if (role == 'personal' || role == 'admin' || role == 'owner') {
      if (stats.admin != null) {
        cards.add(
          StatCard(
            title: 'Active Shipments',
            value: stats.admin!.activeShipments.toString(),
            icon: Icons.local_shipping,
            color: AppColors.info,
          ),
        );
      }

      if (stats.planner != null) {
        cards.add(
          StatCard(
            title: 'Avg Utilization',
            value: '${stats.planner!.avgUtilization.toStringAsFixed(1)}%',
            icon: Icons.pie_chart,
            color: Colors.purple,
          ),
        );
        cards.add(
          StatCard(
            title: 'Items Processed',
            value: stats.planner!.itemsProcessed.toString(),
            icon: Icons.inventory_2,
            color: AppColors.secondary,
          ),
        );
      }

      if (stats.admin != null) {
        cards.add(
          StatCard(
            title: 'Container Types',
            value: stats.admin!.containerTypes.toString(),
            icon: Icons.view_in_ar,
            color: Colors.teal,
          ),
        );
      }
    } else if (role == 'planner' && stats.planner != null) {
      cards.add(
        StatCard(
          title: 'Pending Plans',
          value: stats.planner!.pendingPlans.toString(),
          icon: Icons.pending_actions,
          color: AppColors.warning,
        ),
      );
      cards.add(
        StatCard(
          title: 'Completed Today',
          value: stats.planner!.completedToday.toString(),
          icon: Icons.check_circle,
          color: AppColors.success,
        ),
      );
      cards.add(
        StatCard(
          title: 'Avg Utilization',
          value: '${stats.planner!.avgUtilization.toStringAsFixed(1)}%',
          icon: Icons.pie_chart,
          color: Colors.purple,
        ),
      );
      cards.add(
        StatCard(
          title: 'Items Processed',
          value: stats.planner!.itemsProcessed.toString(),
          icon: Icons.inventory_2,
          color: AppColors.secondary,
        ),
      );
    } else if (role == 'operator' && stats.operator != null) {
      cards.add(
        StatCard(
          title: 'Active Loads',
          value: stats.operator!.activeLoads.toString(),
          icon: Icons.local_shipping,
          color: AppColors.warning,
        ),
      );
      cards.add(
        StatCard(
          title: 'Completed',
          value: stats.operator!.completed.toString(),
          icon: Icons.task_alt,
          color: Colors.teal,
        ),
      );
    }

    // Fallback if no role-specific stats
    if (cards.isEmpty && stats.admin != null) {
      cards.add(
        StatCard(
          title: 'Active Shipments',
          value: stats.admin!.activeShipments.toString(),
          icon: Icons.local_shipping,
          color: AppColors.info,
        ),
      );
    }

    return cards;
  }
}
