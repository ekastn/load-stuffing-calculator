import 'package:flutter/material.dart';
import 'number_field.dart';

/// A widget that displays a 2x2 grid of dimension input fields.
///
/// Provides consistent layout for entering Length, Width, Height, and Weight.
/// Used in plan form for both container and item dimension inputs.
///
/// Example:
/// ```dart
/// DimensionInputs(
///   lengthController: _lengthController,
///   widthController: _widthController,
///   heightController: _heightController,
///   weightController: _weightController,
/// )
/// ```
class DimensionInputs extends StatelessWidget {
  final TextEditingController lengthController;
  final TextEditingController widthController;
  final TextEditingController heightController;
  final TextEditingController weightController;
  final String lengthLabel;
  final String widthLabel;
  final String heightLabel;
  final String weightLabel;
  final bool required;

  const DimensionInputs({
    required this.lengthController,
    required this.widthController,
    required this.heightController,
    required this.weightController,
    this.lengthLabel = 'Length',
    this.widthLabel = 'Width',
    this.heightLabel = 'Height',
    this.weightLabel = 'Max Weight',
    this.required = true,
    super.key,
  });

  @override
  Widget build(BuildContext context) {
    return Column(
      children: [
        // First row: Length and Width
        Row(
          children: [
            Expanded(
              child: NumberField(
                controller: lengthController,
                label: lengthLabel,
                unit: 'mm',
                required: required,
                min: 0,
              ),
            ),
            const SizedBox(width: 12),
            Expanded(
              child: NumberField(
                controller: widthController,
                label: widthLabel,
                unit: 'mm',
                required: required,
                min: 0,
              ),
            ),
          ],
        ),
        const SizedBox(height: 16),
        // Second row: Height and Weight
        Row(
          children: [
            Expanded(
              child: NumberField(
                controller: heightController,
                label: heightLabel,
                unit: 'mm',
                required: required,
                min: 0,
              ),
            ),
            const SizedBox(width: 12),
            Expanded(
              child: NumberField(
                controller: weightController,
                label: weightLabel,
                unit: 'kg',
                required: required,
                min: 0,
              ),
            ),
          ],
        ),
      ],
    );
  }
}
