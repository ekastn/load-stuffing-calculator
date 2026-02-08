import 'package:flutter/material.dart';

/// A button widget for plan actions with optional label and icon.
///
/// Can display as icon-only or with label, providing consistent styling
/// for action buttons throughout the app (e.g., delete, recalculate, view).
///
/// Example:
/// ```dart
/// PlanActionButton(
///   icon: Icons.delete,
///   color: Colors.red,
///   onTap: () => deletePlan(),
///   label: 'Delete',
/// )
/// ```
class PlanActionButton extends StatelessWidget {
  final IconData icon;
  final Color color;
  final VoidCallback onTap;
  final String? label;
  final bool isIconOnly;

  const PlanActionButton({
    required this.icon,
    required this.color,
    required this.onTap,
    this.label,
    this.isIconOnly = false,
    super.key,
  });

  @override
  Widget build(BuildContext context) {
    if (isIconOnly) {
      return InkWell(
        onTap: onTap,
        borderRadius: BorderRadius.circular(8),
        child: Padding(
          padding: const EdgeInsets.all(8.0),
          child: Icon(icon, color: color, size: 20),
        ),
      );
    }

    return TextButton.icon(
      onPressed: onTap,
      icon: Icon(icon, size: 18, color: color),
      label: Text(
        label ?? '',
        style: TextStyle(color: color, fontWeight: FontWeight.w600),
      ),
      style: TextButton.styleFrom(
        padding: const EdgeInsets.symmetric(horizontal: 12, vertical: 8),
        minimumSize: Size.zero,
        tapTargetSize: MaterialTapTargetSize.shrinkWrap,
      ),
    );
  }
}
