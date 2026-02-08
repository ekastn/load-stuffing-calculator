import 'package:flutter/material.dart';
import '../../components/buttons/quick_action_button.dart';

class QuickActionsSection extends StatelessWidget {
  final VoidCallback onProductsTap;
  final VoidCallback onContainersTap;
  final VoidCallback onNewPlanTap;

  const QuickActionsSection({
    super.key,
    required this.onProductsTap,
    required this.onContainersTap,
    required this.onNewPlanTap,
  });

  @override
  Widget build(BuildContext context) {
    return Column(
      crossAxisAlignment: CrossAxisAlignment.start,
      children: [
        Text(
          'Quick Actions',
          style: Theme.of(
            context,
          ).textTheme.titleLarge?.copyWith(fontWeight: FontWeight.bold),
        ),
        const SizedBox(height: 16),
        Row(
          mainAxisAlignment: MainAxisAlignment.spaceBetween,
          children: [
            QuickActionButton(
              label: 'Products',
              icon: Icons.inventory,
              onTap: onProductsTap,
            ),
            QuickActionButton(
              label: 'Containers',
              icon: Icons.view_in_ar,
              onTap: onContainersTap,
            ),
            QuickActionButton(
              label: 'New Plan',
              icon: Icons.add_circle_outline,
              onTap: onNewPlanTap,
            ),
          ],
        ),
      ],
    );
  }
}
