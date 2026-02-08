import 'package:flutter/material.dart';
import '../../config/theme.dart';

class AuthFormLayout extends StatelessWidget {
  final String title;
  final String subtitle;
  final Widget formContent;
  final List<Widget>? footerActions;

  const AuthFormLayout({
    super.key,
    required this.title,
    required this.subtitle,
    required this.formContent,
    this.footerActions,
  });

  @override
  Widget build(BuildContext context) {
    return SafeArea(
      child: SingleChildScrollView(
        padding: const EdgeInsets.all(20.0),
        child: Column(
          mainAxisAlignment: MainAxisAlignment.center,
          children: [
            const SizedBox(height: 60),
            // Logo
            Container(
              width: 80,
              height: 80,
              decoration: BoxDecoration(
                color: AppColors.primary.withValues(alpha: 0.1),
                borderRadius: BorderRadius.circular(20),
              ),
              child: const Icon(
                Icons.view_in_ar,
                size: 40,
                color: AppColors.primary,
              ),
            ),
            const SizedBox(height: 32),
            // Title
            Text(
              title,
              style: Theme.of(
                context,
              ).textTheme.headlineMedium?.copyWith(fontWeight: FontWeight.bold),
              textAlign: TextAlign.center,
            ),
            const SizedBox(height: 8),
            // Subtitle
            Text(
              subtitle,
              style: Theme.of(
                context,
              ).textTheme.bodyLarge?.copyWith(color: AppColors.textSecondary),
              textAlign: TextAlign.center,
            ),
            const SizedBox(height: 48),
            // Form Content (passed as widget)
            formContent,
            const SizedBox(height: 24),
            // Footer Actions (optional)
            if (footerActions != null) ...footerActions!,
          ],
        ),
      ),
    );
  }
}
