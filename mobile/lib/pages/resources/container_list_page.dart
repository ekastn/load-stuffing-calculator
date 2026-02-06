import 'package:flutter/material.dart';
import 'package:provider/provider.dart';
import 'package:go_router/go_router.dart';
import '../../providers/container_provider.dart';
import '../../models/container_model.dart';
import '../../components/widgets/empty_state.dart';
import '../../components/widgets/loading_state.dart';
import '../../components/widgets/error_state.dart';
import '../../components/cards/resource_list_item.dart';
import '../../components/dialogs/confirm_dialog.dart';
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
          if (provider.isLoading && provider.containers.isEmpty) {
            return const LoadingState();
          }

          if (provider.error != null && provider.containers.isEmpty) {
            return ErrorState(
              message: provider.error!,
              onRetry: () => provider.fetchContainers(),
            );
          }

          if (provider.containers.isEmpty) {
            return const EmptyState(
              message: 'No containers found.\nTap + to add a new container.',
              icon: Icons.view_in_ar_outlined,
            );
          }

          return ListView.builder(
            itemCount: provider.containers.length,
            itemBuilder: (context, index) {
              final container = provider.containers[index];
              return ResourceListItem(
                title: container.name,
                trailing: Container(
                  padding: const EdgeInsets.symmetric(horizontal: 8, vertical: 4),
                  decoration: BoxDecoration(
                    color: AppColors.textSecondary.withValues(alpha: 0.1),
                    borderRadius: BorderRadius.circular(12),
                  ),
                  child: const Text('Container', style: TextStyle(fontSize: 12, color: AppColors.textSecondary)),
                ),
                content: Column(
                  crossAxisAlignment: CrossAxisAlignment.start,
                  children: [
                    Text(
                      'Dimensions: ${container.innerLengthMm.toInt()} × ${container.innerWidthMm.toInt()} × ${container.innerHeightMm.toInt()} mm',
                      style: const TextStyle(color: AppColors.textSecondary, fontSize: 14),
                    ),
                    const SizedBox(height: 4),
                    Text(
                      'Max Weight: ${container.maxWeightKg.toStringAsFixed(1)} kg',
                      style: const TextStyle(color: AppColors.textSecondary, fontSize: 14),
                    ),
                  ],
                ),
                actions: [
                  _buildActionButton(
                    context,
                    label: 'Edit',
                    icon: Icons.edit_outlined,
                    color: AppColors.primary,
                    onTap: () => context.push('/containers/${container.id}'),
                  ),
                  _buildActionButton(
                    context,
                    label: 'Delete',
                    icon: Icons.delete_outline,
                    color: AppColors.error,
                    onTap: () => _confirmDelete(context, container),
                  ),
                ],
                onTap: () => context.push('/containers/${container.id}'),
              );
            },
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
            onPressed: () => context.read<ContainerProvider>().fetchContainers(),
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

  void _confirmDelete(BuildContext context, ContainerModel container) {
    ConfirmDialog.show(
      context: context,
      title: 'Delete Container',
      content: 'Are you sure you want to delete ${container.name}?',
      confirmText: 'Delete',
      isDangerous: true,
      onConfirm: () async {
        try {
          await context.read<ContainerProvider>().deleteContainer(container.id);
          if (context.mounted) {
            ScaffoldMessenger.of(context).showSnackBar(
              const SnackBar(content: Text('Container deleted')),
            );
          }
        } catch (e) {
          if (context.mounted) {
            ScaffoldMessenger.of(context).showSnackBar(
              SnackBar(content: Text('Failed to delete: $e')),
            );
          }
        }
      },
    );
  }

  Widget _buildActionButton(BuildContext context, {
    required String label,
    required IconData icon,
    required Color color,
    required VoidCallback onTap,
  }) {
    return TextButton.icon(
      onPressed: onTap,
      icon: Icon(icon, size: 18, color: color),
      label: Text(
        label,
        style: TextStyle(
          color: color,
          fontWeight: FontWeight.w600,
        ),
      ),
      style: TextButton.styleFrom(
        padding: const EdgeInsets.symmetric(horizontal: 12, vertical: 8),
        minimumSize: Size.zero,
        tapTargetSize: MaterialTapTargetSize.shrinkWrap,
      ),
    );
  }
}
