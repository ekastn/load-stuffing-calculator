import 'package:flutter/material.dart';
import '../config/theme.dart';

class UiHelpers {
  static Color getStatusColor(String status) {
    switch (status.toUpperCase()) {
      case 'COMPLETED':
        return AppColors.success.withValues(alpha: 0.2);
      case 'DRAFT':
        return AppColors.textTertiary.withValues(alpha: 0.2);
      case 'IN_PROGRESS':
        return AppColors.info.withValues(alpha: 0.2);
      case 'FAILED':
        return AppColors.error.withValues(alpha: 0.2);
      default:
        return AppColors.surface;
    }
  }
}
