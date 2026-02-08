import 'package:flutter/material.dart';
import '../../components/inputs/app_text_field.dart';
import '../../components/cards/app_card.dart';

class PlanBasicInfoSection extends StatelessWidget {
  final TextEditingController titleController;
  final TextEditingController notesController;

  const PlanBasicInfoSection({
    super.key,
    required this.titleController,
    required this.notesController,
  });

  @override
  Widget build(BuildContext context) {
    return AppCard(
      padding: const EdgeInsets.all(20.0),
      child: Column(
        crossAxisAlignment: CrossAxisAlignment.start,
        children: [
          Text(
            '1. Basic Information',
            style: Theme.of(
              context,
            ).textTheme.titleMedium?.copyWith(fontWeight: FontWeight.bold),
          ),
          const SizedBox(height: 20),
          AppTextField(
            controller: titleController,
            label: 'Plan Title',
            hint: 'e.g., Export to Japan - December 2025',
            required: true,
          ),
          const SizedBox(height: 16),
          AppTextField(
            controller: notesController,
            label: 'Notes',
            hint: 'Optional notes',
            maxLines: 3,
          ),
        ],
      ),
    );
  }
}
