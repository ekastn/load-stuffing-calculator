import 'package:flutter/material.dart';
import '../../config/theme.dart';
import 'app_card.dart';

/// A card widget displaying progress utilization with percentage and subtitle.
///
/// Combines a label, progress bar, percentage display, and optional subtitle
/// in a consistent card layout. Used for showing volume and weight utilization
/// metrics in plan detail views.
///
/// Example:
/// ```dart
/// UtilizationProgressCard(
///   label: 'Volume Utilization',
///   percentage: 75.5,
///   subtitle: '2.45 m³',
///   color: Colors.green,
/// )
/// ```
class UtilizationProgressCard extends StatelessWidget {
  final String label;
  final double percentage;
  final String subtitle;
  final Color color;

  const UtilizationProgressCard({
    required this.label,
    required this.percentage,
    required this.subtitle,
    required this.color,
    super.key,
  });

  @override
  Widget build(BuildContext context) {
    return AppCard(
      child: Column(
        crossAxisAlignment: CrossAxisAlignment.start,
        children: [
          Row(
            mainAxisAlignment: MainAxisAlignment.spaceBetween,
            children: [
              Text(
                label,
                style: const TextStyle(
                  fontSize: 14,
                  fontWeight: FontWeight.w500,
                ),
              ),
              Text(
                '${percentage.toStringAsFixed(1)}%',
                style: TextStyle(
                  fontSize: 16,
                  fontWeight: FontWeight.bold,
                  color: color,
                ),
              ),
            ],
          ),
          const SizedBox(height: 8),
          ClipRRect(
            borderRadius: BorderRadius.circular(4),
            child: LinearProgressIndicator(
              value: (percentage / 100).clamp(0.0, 1.0),
              minHeight: 8,
              backgroundColor: AppColors.background,
              valueColor: AlwaysStoppedAnimation<Color>(color),
            ),
          ),
          const SizedBox(height: 4),
          Text(
            subtitle,
            style: TextStyle(fontSize: 12, color: AppColors.textSecondary),
          ),
        ],
      ),
    );
  }
}
