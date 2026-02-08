import 'package:flutter/material.dart';
import '../cards/app_card.dart';
import '../inputs/dimension_inputs.dart';
import '../inputs/app_text_field.dart';

class ContainerFormContent extends StatelessWidget {
  final TextEditingController nameController;
  final TextEditingController lengthController;
  final TextEditingController widthController;
  final TextEditingController heightController;
  final TextEditingController weightController;
  final TextEditingController descController;
  final GlobalKey<FormState>? formKey;

  const ContainerFormContent({
    super.key,
    required this.nameController,
    required this.lengthController,
    required this.widthController,
    required this.heightController,
    required this.weightController,
    required this.descController,
    this.formKey,
  });

  @override
  Widget build(BuildContext context) {
    return AppCard(
      padding: const EdgeInsets.all(20),
      child: Column(
        crossAxisAlignment: CrossAxisAlignment.start,
        children: [
          Text(
            'Container Properties',
            style: Theme.of(
              context,
            ).textTheme.titleMedium?.copyWith(fontWeight: FontWeight.bold),
          ),
          const SizedBox(height: 20),
          AppTextField(
            controller: nameController,
            label: 'Container Name',
            hint: 'Enter container name',
            required: true,
            validator: (value) {
              if (value == null || value.trim().isEmpty) {
                return 'Container name is required';
              }
              return null;
            },
          ),
          const SizedBox(height: 16),
          DimensionInputs(
            lengthController: lengthController,
            widthController: widthController,
            heightController: heightController,
            weightController: weightController,
            lengthLabel: 'Inner Length',
            widthLabel: 'Inner Width',
            heightLabel: 'Inner Height',
            weightLabel: 'Max Weight',
            dimensionSuffix: 'mm',
            weightSuffix: 'kg',
            hintText: '0',
            borderRadius: 8,
            spacing: 12,
            required: true,
          ),
          const SizedBox(height: 16),
          AppTextField(
            controller: descController,
            label: 'Description',
            hint: 'Optional description',
            maxLines: 3,
          ),
        ],
      ),
    );
  }
}
