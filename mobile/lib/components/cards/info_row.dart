import 'package:flutter/material.dart';
import '../../config/theme.dart';

/// A reusable row widget that displays a label and value pair.
///
/// Used throughout the app to show information in a consistent format,
/// typically within cards or detail pages.
///
/// Example:
/// ```dart
/// InfoRow(label: 'Name', value: 'My Container')
/// InfoRow(label: 'Dimensions', value: '1000 × 800 × 600 mm')
/// ```
class InfoRow extends StatelessWidget {
  final String label;
  final String value;
  final IconData? icon;
  final TextStyle? valueStyle;

  const InfoRow({
    required this.label,
    required this.value,
    this.icon,
    this.valueStyle,
    super.key,
  });

  @override
  Widget build(BuildContext context) {
    return Row(
      mainAxisAlignment: MainAxisAlignment.spaceBetween,
      children: [
        Expanded(
          child: Row(
            children: [
              if (icon != null) ...[
                Icon(icon, size: 16, color: AppColors.textSecondary),
                const SizedBox(width: 8),
              ],
              Text(
                label,
                style: TextStyle(fontSize: 14, color: AppColors.textSecondary),
              ),
            ],
          ),
        ),
        Text(
          value,
          style:
              valueStyle ??
              const TextStyle(fontSize: 14, fontWeight: FontWeight.w500),
          textAlign: TextAlign.end,
        ),
      ],
    );
  }
}
