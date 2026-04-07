import 'package:flutter/material.dart';
import '../../config/theme.dart';

class AppTextField extends StatelessWidget {
  final TextEditingController? controller;
  final String label;
  final String? hint;
  final bool required;
  final TextInputType? keyboardType;
  final String? Function(String?)? validator;
  final int? maxLines;
  final Widget? suffix;
  final bool obscureText;
  final IconData? prefixIcon;
  final ValueChanged<String>? onChanged;
  final String? initialValue;

  const AppTextField({
    super.key,
    this.controller,
    required this.label,
    this.hint,
    this.required = false,
    this.keyboardType,
    this.validator,
    this.maxLines = 1,
    this.suffix,
    this.obscureText = false,
    this.prefixIcon,
    this.onChanged,
    this.initialValue,
  });

  @override
  Widget build(BuildContext context) {
    return Column(
      crossAxisAlignment: CrossAxisAlignment.start,
      children: [
        Text(
          required ? '$label *' : label,
          style: Theme.of(context).textTheme.bodyMedium?.copyWith(
                fontWeight: FontWeight.w500,
                color: AppColors.textPrimary,
              ),
        ),
        const SizedBox(height: 8),
        TextFormField(
          controller: controller,
          initialValue: initialValue,
          onChanged: onChanged,
          decoration: InputDecoration(
            hintText: hint,
            suffixIcon: suffix,
            prefixIcon: prefixIcon != null ? Icon(prefixIcon, size: 20) : null,
          ),
          keyboardType: keyboardType,
          maxLines: maxLines,
          obscureText: obscureText,
          validator: validator ??
              (required
                  ? (value) => value == null || value.isEmpty ? '$label is required' : null
                  : null),
        ),
      ],
    );
  }
}
