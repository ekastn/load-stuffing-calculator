import 'package:flutter/material.dart';
import '../../models/plan_detail_model.dart';
import 'plan_header_section.dart';
import 'container_info_section.dart';

class PlanOverviewTabSection extends StatelessWidget {
  final PlanDetailModel plan;

  const PlanOverviewTabSection({super.key, required this.plan});

  @override
  Widget build(BuildContext context) {
    return SingleChildScrollView(
      padding: const EdgeInsets.all(16),
      child: Column(
        crossAxisAlignment: CrossAxisAlignment.start,
        children: [
          PlanHeaderSection(plan: plan),
          const SizedBox(height: 24),
          ContainerInfoSection(container: plan.container),
        ],
      ),
    );
  }
}
