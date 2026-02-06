import 'package:flutter/material.dart';
import '../config/theme.dart';

class UiHelpers {
  static Color getStatusColor(String status) {
    switch (status.toUpperCase()) {
      case 'COMPLETED':
        return AppColors.success.withOpacity(0.2);
      case 'DRAFT':
        return AppColors.textTertiary.withOpacity(0.2);
      case 'IN_PROGRESS':
        return AppColors.info.withOpacity(0.2);
      case 'FAILED':
        return AppColors.error.withOpacity(0.2);
      default:
        return AppColors.surface;
    }
  }
}
