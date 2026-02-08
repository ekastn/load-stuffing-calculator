import 'package:flutter/material.dart';
import 'package:provider/provider.dart';
import '../../providers/auth_provider.dart';
import '../../components/buttons/app_button.dart';
import '../../components/sections/settings_section.dart';
import '../../components/widgets/settings_list_item.dart';
import '../../components/widgets/profile_header.dart';
import '../../config/theme.dart';

class ProfilePage extends StatelessWidget {
  const ProfilePage({super.key});

  @override
  Widget build(BuildContext context) {
    final authProvider = context.watch<AuthProvider>();
    final user = authProvider.user;

    if (user == null) {
      return const Scaffold(body: Center(child: Text('Not logged in')));
    }

    return Scaffold(
      appBar: AppBar(
        title: const Text('Profile'),
        centerTitle: true,
        elevation: 0,
        backgroundColor: Colors.transparent,
      ),
      extendBodyBehindAppBar: true,
      body: SingleChildScrollView(
        padding: EdgeInsets.only(
          top: MediaQuery.of(context).padding.top + kToolbarHeight + 16,
          left: 20,
          right: 20,
          bottom: 24,
        ),
        child: Column(
          children: [
            // Header
            ProfileHeader(username: user.username, role: user.role),
            const SizedBox(height: 32),

            // Account Settings
            SettingsSection(
              title: 'Account Settings',
              children: [
                SettingsListItem(
                  icon: Icons.person_outline,
                  title: 'Edit Profile',
                  onTap: () {},
                ),
                SettingsListItem(
                  icon: Icons.lock_outline,
                  title: 'Change Password',
                  onTap: () {},
                ),
                SettingsListItem(
                  icon: Icons.notifications_outlined,
                  title: 'Notifications',
                  onTap: () {},
                ),
              ],
            ),
            const SizedBox(height: 32),

            // General Settings
            SettingsSection(
              title: 'General',
              children: [
                SettingsListItem(
                  icon: Icons.description_outlined,
                  title: 'Terms of Service',
                  onTap: () {},
                ),
                SettingsListItem(
                  icon: Icons.privacy_tip_outlined,
                  title: 'Privacy Policy',
                  onTap: () {},
                ),
                SettingsListItem(
                  icon: Icons.info_outline,
                  title: 'About App',
                  trailing: const Text(
                    'v1.0.0',
                    style: TextStyle(
                      color: AppColors.textSecondary,
                      fontSize: 13,
                    ),
                  ),
                  onTap: () {},
                ),
              ],
            ),

            const SizedBox(height: 48),

            // Logout
            AppButton(
              label: 'Logout',
              variant: AppButtonVariant.destructive,
              icon: Icons.logout,
              isFullWidth: true,
              onPressed: () => _showLogoutConfirmation(context, authProvider),
            ),
          ],
        ),
      ),
    );
  }

  void _showLogoutConfirmation(
    BuildContext context,
    AuthProvider authProvider,
  ) {
    showDialog(
      context: context,
      builder: (context) => AlertDialog(
        title: const Text('Logout'),
        content: const Text('Are you sure you want to logout?'),
        icon: const Icon(Icons.logout, color: AppColors.error, size: 32),
        actions: [
          TextButton(
            onPressed: () => Navigator.of(context).pop(),
            child: const Text('Cancel'),
          ),
          FilledButton(
            onPressed: () async {
              Navigator.of(context).pop();
              await authProvider.logout();
            },
            style: FilledButton.styleFrom(
              backgroundColor: AppColors.error,
              foregroundColor: Colors.white,
            ),
            child: const Text('Logout'),
          ),
        ],
      ),
    );
  }
}
