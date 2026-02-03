import 'package:flutter/material.dart';
import 'package:go_router/go_router.dart';

class MainShellPage extends StatelessWidget {
  final Widget child;

  const MainShellPage({super.key, required this.child});

  @override
  Widget build(BuildContext context) {
    return Scaffold(
      body: child,
      bottomNavigationBar: NavigationBar(
        selectedIndex: _calculateSelectedIndex(context),
        onDestinationSelected: (index) => _onItemTapped(index, context),
        destinations: const [
          NavigationDestination(
            icon: Icon(Icons.dashboard_outlined),
            selectedIcon: Icon(Icons.dashboard),
            label: 'Dashboard',
          ),
          NavigationDestination(
            icon: Icon(Icons.list_alt_outlined),
            selectedIcon: Icon(Icons.list_alt),
            label: 'Plans',
          ),
          NavigationDestination(
            icon: Icon(Icons.inventory_outlined),
            selectedIcon: Icon(Icons.inventory),
            label: 'Products',
          ),
          NavigationDestination(
            icon: Icon(Icons.view_in_ar_outlined),
            selectedIcon: Icon(Icons.view_in_ar),
            label: 'Containers',
          ),
          NavigationDestination(
            icon: Icon(Icons.person_outline),
            selectedIcon: Icon(Icons.person),
            label: 'Profile',
          ),
        ],
      ),
    );
  }

  static int _calculateSelectedIndex(BuildContext context) {
    final String location = GoRouterState.of(context).uri.path;
    if (location.startsWith('/plans')) {
      return 1;
    }
    if (location.startsWith('/products')) {
      return 2;
    }
    if (location.startsWith('/containers')) {
      return 3;
    }
    if (location.startsWith('/profile')) {
      return 4;
    }
    return 0; // Dashboard
  }

  void _onItemTapped(int index, BuildContext context) {
    switch (index) {
      case 0:
        context.go('/');
        break;
      case 1:
        context.go('/plans');
        break;
      case 2:
        context.go('/products');
        break;
      case 3:
        context.go('/containers');
        break;
      case 4:
        context.go('/profile');
        break;
    }
  }
}
