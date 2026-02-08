import 'package:flutter/material.dart';
import 'app_text_field.dart';

class NumberField extends StatelessWidget {
  final TextEditingController controller;
  final String label;
  final bool required;
  final String? unit;
  final double? min;
  final double? max;
  final String hintText;
  final double borderRadius;
  final String? Function(String?)? validator;

  const NumberField({
    super.key,
    required this.controller,
    required this.label,
    this.required = false,
    this.unit,
    this.min,
    this.max,
    this.hintText = '0',
    this.borderRadius = 8.0,
    this.validator,
  });

  @override
  Widget build(BuildContext context) {
    return AppTextField(
      controller: controller,
      label: label,
      hint: hintText,
      required: required,
      keyboardType: const TextInputType.numberWithOptions(decimal: true),
      suffix: unit != null ? Text(unit!) : null,
      validator: validator ?? _defaultValidator,
    );
  }

  String? _defaultValidator(String? value) {
    if (value == null || value.isEmpty) {
      return required ? 'Required' : null;
    }

    final number = double.tryParse(value);
    if (number == null) {
      return 'Invalid number';
    }

    if (min != null && number < min!) {
      return 'Must be ≥ $min';
    }

    if (max != null && number > max!) {
      return 'Must be ≤ $max';
    }

    return null;
  }
}
