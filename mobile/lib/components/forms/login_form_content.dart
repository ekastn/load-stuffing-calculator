import 'package:flutter/material.dart';
import '../inputs/app_text_field.dart';
import '../buttons/app_button.dart';
import '../cards/app_card.dart';

class LoginFormContent extends StatelessWidget {
  final TextEditingController usernameController;
  final TextEditingController passwordController;
  final VoidCallback onSubmit;
  final bool isLoading;

  const LoginFormContent({
    super.key,
    required this.usernameController,
    required this.passwordController,
    required this.onSubmit,
    this.isLoading = false,
  });

  @override
  Widget build(BuildContext context) {
    return AppCard(
      padding: const EdgeInsets.all(24),
      child: Column(
        children: [
          AppTextField(
            controller: usernameController,
            label: 'Username',
            prefixIcon: Icons.person_outline,
            validator: (value) {
              if (value?.isEmpty ?? true) return 'Username is required';
              return null;
            },
          ),
          const SizedBox(height: 16),
          AppTextField(
            controller: passwordController,
            label: 'Password',
            obscureText: true,
            prefixIcon: Icons.lock_outline,
            validator: (value) {
              if (value?.isEmpty ?? true) return 'Password is required';
              return null;
            },
          ),
          const SizedBox(height: 32),
          AppButton(
            onPressed: onSubmit,
            label: 'Sign In',
            isLoading: isLoading,
            isFullWidth: true,
          ),
        ],
      ),
    );
  }
}
