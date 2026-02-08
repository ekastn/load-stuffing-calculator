import 'package:flutter/material.dart';
import '../../components/cards/app_card.dart';
import '../../components/widgets/status_badge.dart';
import '../../models/plan_model.dart';
import '../../config/theme.dart';

class RecentPlansList extends StatelessWidget {
  final List<PlanModel> plans;
  final Function(String planId) onPlanTap;

  const RecentPlansList({
    super.key,
    required this.plans,
    required this.onPlanTap,
  });

  @override
  Widget build(BuildContext context) {
    return Column(
      crossAxisAlignment: CrossAxisAlignment.start,
      children: [
        Text(
          'Recent Plans',
          style: Theme.of(
            context,
          ).textTheme.titleLarge?.copyWith(fontWeight: FontWeight.bold),
        ),
        const SizedBox(height: 16),
        if (plans.isEmpty)
          const Padding(
            padding: EdgeInsets.all(32.0),
            child: Center(
              child: Text(
                'No plans found',
                style: TextStyle(color: AppColors.textSecondary),
              ),
            ),
          )
        else
          ...plans.map(
            (plan) => Padding(
              padding: const EdgeInsets.only(bottom: 12),
              child: AppCard(
                onTap: () => onPlanTap(plan.id),
                child: Row(
                  children: [
                    Container(
                      padding: const EdgeInsets.all(12),
                      decoration: BoxDecoration(
                        color: AppColors.primary.withValues(alpha: 0.1),
                        borderRadius: BorderRadius.circular(12),
                      ),
                      child: const Icon(
                        Icons.inventory_2_outlined,
                        color: AppColors.primary,
                      ),
                    ),
                    const SizedBox(width: 16),
                    Expanded(
                      child: Column(
                        crossAxisAlignment: CrossAxisAlignment.start,
                        children: [
                          Text(
                            plan.title.isEmpty ? plan.code : plan.title,
                            style: const TextStyle(
                              fontWeight: FontWeight.w600,
                              fontSize: 16,
                            ),
                          ),
                          const SizedBox(height: 4),
                          Text(
                            plan.code,
                            style: const TextStyle(
                              fontFamily: 'monospace',
                              fontSize: 12,
                              color: AppColors.textSecondary,
                            ),
                          ),
                        ],
                      ),
                    ),
                    StatusBadge.fromStatus(plan.status),
                  ],
                ),
              ),
            ),
          ),
      ],
    );
  }
}
