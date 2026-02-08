import 'package:flutter/material.dart';
import '../cards/app_card.dart';
import '../fields/color_picker_field.dart';
import '../inputs/dimension_inputs.dart';
import '../inputs/app_text_field.dart';

class ProductFormContent extends StatelessWidget {
  final TextEditingController nameController;
  final TextEditingController lengthController;
  final TextEditingController widthController;
  final TextEditingController heightController;
  final TextEditingController weightController;
  final TextEditingController colorController;
  final Function(String) onColorChanged;
  final GlobalKey<FormState>? formKey;

  const ProductFormContent({
    super.key,
    required this.nameController,
    required this.lengthController,
    required this.widthController,
    required this.heightController,
    required this.weightController,
    required this.colorController,
    required this.onColorChanged,
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
            'Product Details',
            style: Theme.of(
              context,
            ).textTheme.titleMedium?.copyWith(fontWeight: FontWeight.bold),
          ),
          const SizedBox(height: 20),
          AppTextField(
            controller: nameController,
            label: 'Product Name',
            hint: 'Enter product name',
            required: true,
            validator: (value) {
              if (value == null || value.trim().isEmpty) {
                return 'Product name is required';
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
            lengthLabel: 'Length',
            widthLabel: 'Width',
            heightLabel: 'Height',
            weightLabel: 'Weight',
            dimensionSuffix: 'mm',
            weightSuffix: 'kg',
            hintText: '0',
            borderRadius: 8,
            spacing: 12,
            required: true,
          ),
          const SizedBox(height: 16),
          ColorPickerField(
            initialColor: colorController.text,
            onColorChanged: onColorChanged,
            label: 'Color',
          ),
        ],
      ),
    );
  }
}
