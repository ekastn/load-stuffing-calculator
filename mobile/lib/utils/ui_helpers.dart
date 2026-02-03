import 'package:flutter/material.dart';

class UiHelpers {
  static Color getStatusColor(String status) {
    switch (status.toUpperCase()) {
      case 'COMPLETED':
        return Colors.green.shade100;
      case 'DRAFT':
        return Colors.grey.shade200;
      case 'IN_PROGRESS':
        return Colors.blue.shade100;
      case 'FAILED':
        return Colors.red.shade100;
      default:
        return Colors.grey.shade100;
    }
  }
}
