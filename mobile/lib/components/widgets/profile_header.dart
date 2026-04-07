import 'package:flutter/material.dart';
import '../../config/theme.dart';

class ProfileHeader extends StatelessWidget {
  final String username;
  final String role;

  const ProfileHeader({super.key, required this.username, required this.role});

  @override
  Widget build(BuildContext context) {
    return Column(
      children: [
        Container(
          width: 100,
          height: 100,
          decoration: BoxDecoration(
            color: AppColors.primary,
            shape: BoxShape.circle,
            border: Border.all(color: Colors.white, width: 4),
            boxShadow: [
              BoxShadow(
                color: Colors.black.withValues(alpha: 0.1),
                blurRadius: 10,
                offset: const Offset(0, 5),
              ),
            ],
          ),
          child: Center(
            child: Text(
              username.isNotEmpty
                  ? username.substring(0, 1).toUpperCase()
                  : '?',
              style: const TextStyle(
                fontSize: 48,
                fontWeight: FontWeight.bold,
                color: Colors.white,
              ),
            ),
          ),
        ),
        const SizedBox(height: 16),
        Text(
          username,
          style: Theme.of(context).textTheme.headlineSmall?.copyWith(
            fontWeight: FontWeight.bold,
            color: AppColors.textPrimary,
          ),
        ),
        const SizedBox(height: 8),
        Container(
          padding: const EdgeInsets.symmetric(horizontal: 16, vertical: 6),
          decoration: BoxDecoration(
            color: _getRoleColor(role).withValues(alpha: 0.1),
            borderRadius: BorderRadius.circular(20),
            border: Border.all(
              color: _getRoleColor(role).withValues(alpha: 0.2),
            ),
          ),
          child: Text(
            role.toUpperCase(),
            style: TextStyle(
              color: _getRoleColor(role),
              fontWeight: FontWeight.bold,
              fontSize: 12,
              letterSpacing: 0.5,
            ),
          ),
        ),
      ],
    );
  }

  Color _getRoleColor(String role) {
    switch (role.toLowerCase()) {
      case 'admin':
      case 'founder':
        return Colors.purple;
      case 'operator':
        return AppColors.info;
      case 'planner':
        return AppColors.success;
      default:
        return AppColors.textSecondary;
    }
  }
}
