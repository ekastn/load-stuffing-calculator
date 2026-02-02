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

class ContainerListPage extends StatefulWidget {
  const ContainerListPage({super.key});

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
      body: Consumer<ContainerProvider>(
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
                leading: Icon(
                  Icons.view_in_ar,
                  size: 40,
                  color: Theme.of(context).colorScheme.primary,
                ),
                title: container.name,
                subtitle:
                    '${container.innerLengthMm}×${container.innerWidthMm}×${container.innerHeightMm}mm\nMax: ${container.maxWeightKg}kg',
                onEdit: () => context.push('/containers/${container.id}'),
                onDelete: () => _confirmDelete(context, container),
              );
            },
          );
        },
      ),
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
}
