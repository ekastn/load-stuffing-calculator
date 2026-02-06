import 'package:flutter/material.dart';
import 'package:go_router/go_router.dart';
import '../plans/plan_list_page.dart';
import '../resources/product_list_page.dart';
import '../resources/container_list_page.dart';

class ResourcesPage extends StatefulWidget {
  final int initialIndex;

  const ResourcesPage({super.key, this.initialIndex = 0});

  @override
  State<ResourcesPage> createState() => _ResourcesPageState();
}

class _ResourcesPageState extends State<ResourcesPage> with SingleTickerProviderStateMixin {
  late TabController _tabController;

  @override
  void initState() {
    super.initState();
    _tabController = TabController(length: 3, vsync: this, initialIndex: widget.initialIndex);
  }

  @override
  void didUpdateWidget(ResourcesPage oldWidget) {
    super.didUpdateWidget(oldWidget);
    if (oldWidget.initialIndex != widget.initialIndex) {
      _tabController.animateTo(widget.initialIndex);
    }
  }

  @override
  void dispose() {
    _tabController.dispose();
    super.dispose();
  }

  void _handleTabSelection(int index) {
    switch (index) {
      case 0:
        context.go('/plans');
        break;
      case 1:
        context.go('/products');
        break;
      case 2:
        context.go('/containers');
        break;
    }
  }

  @override
  Widget build(BuildContext context) {
    return Scaffold(
      appBar: AppBar(
        title: const Text('Resources'),
        bottom: TabBar(
          controller: _tabController,
          onTap: _handleTabSelection,
          tabs: const [
            Tab(text: 'Plans'),
            Tab(text: 'Products'),
            Tab(text: 'Containers'),
          ],
        ),
      ),
      body: TabBarView(
        controller: _tabController,
        children: const [
          PlanListPage(isEmbedded: true),
          ProductListPage(isEmbedded: true),
          ContainerListPage(isEmbedded: true),
        ],
      ),
    );
  }
}
