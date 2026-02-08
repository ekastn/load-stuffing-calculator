import 'package:flutter/material.dart';
import 'package:provider/provider.dart';
import 'package:go_router/go_router.dart';
import '../../providers/container_provider.dart';
import '../../models/container_model.dart';
import '../../components/sections/resource_list_section.dart';
import '../../config/theme.dart';

class ContainerListPage extends StatefulWidget {
  final bool isEmbedded;
  const ContainerListPage({super.key, this.isEmbedded = false});

  @override
  State<ContainerListPage> createState() => _ContainerListPageState();
}

class _ContainerListPageState extends State<ContainerListPage> {
  @override
  void initState() {
    super.initState();
    WidgetsBinding.instance.addPostFrameCallback((_) {
      context.read<ContainerProvider>().fetchContainers();
    });
  }

  @override
  Widget build(BuildContext context) {
    final content = Consumer<ContainerProvider>(
      builder: (context, provider, child) {
        return ResourceListSection<ContainerModel>(
          resources: provider.containers,
          isLoading: provider.isLoading,
          error: provider.error,
          emptyMessage: 'No containers found.\nTap + to add a new container.',
          emptyIcon: Icons.view_in_ar_outlined,
          title: (container) => container.name,
          id: (container) => container.id,
          trailing: (context, container) => Container(
            padding: const EdgeInsets.symmetric(horizontal: 8, vertical: 4),
            decoration: BoxDecoration(
              color: AppColors.textSecondary.withValues(alpha: 0.1),
              borderRadius: BorderRadius.circular(12),
            ),
            child: const Text(
              'Container',
              style: TextStyle(fontSize: 12, color: AppColors.textSecondary),
            ),
          ),
          details: (context, container) => Column(
            crossAxisAlignment: CrossAxisAlignment.start,
            children: [
              Text(
                'Dimensions: ${container.innerLengthMm.toInt()} × ${container.innerWidthMm.toInt()} × ${container.innerHeightMm.toInt()} mm',
                style: const TextStyle(
                  color: AppColors.textSecondary,
                  fontSize: 14,
                ),
              ),
              const SizedBox(height: 4),
              Text(
                'Max Weight: ${container.maxWeightKg.toStringAsFixed(1)} kg',
                style: const TextStyle(
                  color: AppColors.textSecondary,
                  fontSize: 14,
                ),
              ),
            ],
          ),
          onEdit: (context, container) =>
              context.push('/containers/${container.id}'),
          onDelete: (context, container) =>
              context.read<ContainerProvider>().deleteContainer(container.id),
          onRetry: () => context.read<ContainerProvider>().fetchContainers(),
        );
      },
    );

    if (widget.isEmbedded) {
      return Scaffold(
        body: content,
        floatingActionButton: FloatingActionButton(
          onPressed: () => context.push('/containers/new'),
          child: const Icon(Icons.add),
        ),
      );
    }

    return Scaffold(
      appBar: AppBar(
        title: const Text('Containers'),
        actions: [
          IconButton(
            icon: const Icon(Icons.refresh),
            onPressed: () =>
                context.read<ContainerProvider>().fetchContainers(),
          ),
        ],
      ),
      floatingActionButton: FloatingActionButton(
        onPressed: () => context.push('/containers/new'),
        child: const Icon(Icons.add),
      ),
      body: content,
    );
  }
}
