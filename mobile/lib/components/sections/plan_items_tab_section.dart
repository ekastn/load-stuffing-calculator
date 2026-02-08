import 'package:flutter/material.dart';
import '../../models/plan_detail_model.dart';
import '../widgets/empty_state.dart';
import '../cards/plan_item_card.dart';

class PlanItemsTabSection extends StatelessWidget {
  final PlanDetailModel plan;

  const PlanItemsTabSection({super.key, required this.plan});

  @override
  Widget build(BuildContext context) {
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
        return PlanItemCard(item: item);
      },
    );
  }
}
