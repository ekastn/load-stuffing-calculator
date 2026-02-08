import 'package:flutter/material.dart';
import 'package:provider/provider.dart';
import 'package:go_router/go_router.dart';
import '../../config/theme.dart';
import '../../providers/auth_provider.dart';
import '../../exceptions/login_exception.dart';
import '../../components/sections/auth_form_layout.dart';
import '../../components/forms/login_form_content.dart';

class LoginPage extends StatefulWidget {
  const LoginPage({super.key});

  @override
  State<LoginPage> createState() => _LoginPageState();
}

class _LoginPageState extends State<LoginPage> {
  final _usernameController = TextEditingController();
  final _passwordController = TextEditingController();

  @override
  void dispose() {
    _usernameController.dispose();
    _passwordController.dispose();
    super.dispose();
  }

  Future<void> _handleLogin() async {
    final username = _usernameController.text.trim();
    final password = _passwordController.text.trim();

    if (username.isEmpty || password.isEmpty) {
      ScaffoldMessenger.of(context).showSnackBar(
        const SnackBar(content: Text('Please fill in all fields')),
      );
      return;
    }

    try {
      await context.read<AuthProvider>().login(username, password);
      if (mounted) {
        context.go('/');
      }
    } catch (e) {
      if (mounted) {
        final message = e is LoginException
            ? e.toString()
            : 'Login failed. Please try again.';
        ScaffoldMessenger.of(context).showSnackBar(
          SnackBar(content: Text(message), backgroundColor: Colors.red),
        );
      }
    }
  }

  @override
  Widget build(BuildContext context) {
    return Scaffold(
      backgroundColor: AppColors.background,
      body: AuthFormLayout(
        title: 'Welcome Back',
        subtitle: 'Sign in to your account',
        formContent: LoginFormContent(
          usernameController: _usernameController,
          passwordController: _passwordController,
          onSubmit: _handleLogin,
          isLoading: context.watch<AuthProvider>().isLoading,
        ),
        footerActions: [
          TextButton(
            onPressed: () {},
            child: const Text('Don\'t have an account? Sign up'),
          ),
        ],
      ),
    );
  }
}
