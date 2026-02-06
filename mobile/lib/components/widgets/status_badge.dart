import 'package:flutter/material.dart';
import '../../config/theme.dart';

enum StatusType { success, warning, error, info, neutral }

class StatusBadge extends StatelessWidget {
  final String label;
  final StatusType type;

  const StatusBadge({
    super.key,
    required this.label,
    required this.type,
  });

  factory StatusBadge.fromStatus(String status) {
    StatusType type;
    switch (status.toUpperCase()) {
      case 'COMPLETED':
        type = StatusType.success;
        break;
      case 'IN_PROGRESS':
        type = StatusType.info;
        break;
      case 'FAILED':
        type = StatusType.error;
        break;
      case 'DRAFT':
        type = StatusType.warning;
        break;
      default:
        type = StatusType.neutral;
    }
    
    return StatusBadge(
      label: status.replaceAll('_', ' ').toUpperCase(),
      type: type,
    );
  }

  @override
  Widget build(BuildContext context) {
    Color backgroundColor;
    Color textColor;

    switch (type) {
      case StatusType.success:
        backgroundColor = AppColors.success.withOpacity(0.1);
        textColor = AppColors.success;
        break;
      case StatusType.warning:
        backgroundColor = AppColors.warning.withOpacity(0.1);
        textColor = AppColors.warning;
        break;
      case StatusType.error:
        backgroundColor = AppColors.error.withOpacity(0.1);
        textColor = AppColors.error;
        break;
      case StatusType.info:
        backgroundColor = AppColors.info.withOpacity(0.1);
        textColor = AppColors.info;
        break;
      case StatusType.neutral:
        backgroundColor = AppColors.textTertiary.withOpacity(0.1);
        textColor = AppColors.textTertiary;
        break;
    }

    return Container(
      padding: const EdgeInsets.symmetric(horizontal: 10, vertical: 4),
      decoration: BoxDecoration(
        color: backgroundColor,
        borderRadius: BorderRadius.circular(16),
      ),
      child: Text(
        label,
        style: TextStyle(
          color: textColor,
          fontWeight: FontWeight.bold,
          fontSize: 11,
        ),
      ),
    );
  }
}
