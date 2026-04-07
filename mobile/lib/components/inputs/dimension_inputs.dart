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
  final String dimensionSuffix;
  final String weightSuffix;
  final String hintText;
  final double borderRadius;
  final double spacing;
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
    this.dimensionSuffix = 'mm',
    this.weightSuffix = 'kg',
    this.hintText = '0',
    this.borderRadius = 8.0,
    this.spacing = 12.0,
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
                unit: dimensionSuffix,
                required: required,
                min: 0,
                hintText: hintText,
                borderRadius: borderRadius,
              ),
            ),
            SizedBox(width: spacing),
            Expanded(
              child: NumberField(
                controller: widthController,
                label: widthLabel,
                unit: dimensionSuffix,
                required: required,
                min: 0,
                hintText: hintText,
                borderRadius: borderRadius,
              ),
            ),
          ],
        ),
        SizedBox(height: spacing + 4),
        // Second row: Height and Weight
        Row(
          children: [
            Expanded(
              child: NumberField(
                controller: heightController,
                label: heightLabel,
                unit: dimensionSuffix,
                required: required,
                min: 0,
                hintText: hintText,
                borderRadius: borderRadius,
              ),
            ),
            SizedBox(width: spacing),
            Expanded(
              child: NumberField(
                controller: weightController,
                label: weightLabel,
                unit: weightSuffix,
                required: required,
                min: 0,
                hintText: hintText,
                borderRadius: borderRadius,
              ),
            ),
          ],
        ),
      ],
    );
  }
}
