import 'package:flutter/material.dart';
import 'package:provider/provider.dart';
import 'package:go_router/go_router.dart';
import '../../providers/container_provider.dart';
import '../../models/container_model.dart';
// import '../../dtos/container_dto.dart';

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
        onPressed: () {
            // Navigate to create form
            context.push('/containers/new');
        },
        child: const Icon(Icons.add),
      ),
      body: Consumer<ContainerProvider>(
        builder: (context, provider, child) {
          if (provider.isLoading && provider.containers.isEmpty) {
            return const Center(child: CircularProgressIndicator());
          }

          if (provider.error != null && provider.containers.isEmpty) {
            return Center(child: Text('Error: ${provider.error}'));
          }

          if (provider.containers.isEmpty) {
            return const Center(child: Text('No containers found.'));
          }

          return ListView.builder(
            itemCount: provider.containers.length,
            itemBuilder: (context, index) {
              final container = provider.containers[index];
              return Card(
                margin: const EdgeInsets.symmetric(horizontal: 16, vertical: 8),
                child: ListTile(
                  leading: const Icon(Icons.view_in_ar, size: 40),
                  title: Text(container.name),
                  subtitle: Text('${container.innerLengthMm}x${container.innerWidthMm}x${container.innerHeightMm}mm\nMax: ${container.maxWeightKg}kg'),
                  isThreeLine: true,
                  trailing: Row(
                    mainAxisSize: MainAxisSize.min,
                    children: [
                       IconButton(
                        icon: const Icon(Icons.edit, color: Colors.blue),
                        onPressed: () {
                            // Navigation to edit
                            context.push('/containers/${container.id}');
                        },
                      ),
                      IconButton(
                        icon: const Icon(Icons.delete, color: Colors.red),
                        onPressed: () => _confirmDelete(context, container),
                      ),
                    ],
                  ),
                ),
              );
            },
          );
        },
      ),
    );
  }

  void _confirmDelete(BuildContext context, ContainerModel container) {
    showDialog(
      context: context,
      builder: (context) => AlertDialog(
        title: const Text('Delete Container'),
        content: Text('Are you sure you want to delete ${container.name}?'),
        actions: [
          TextButton(
            onPressed: () => Navigator.pop(context),
            child: const Text('Cancel'),
          ),
          TextButton(
            onPressed: () async {
              Navigator.pop(context);
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
            child: const Text('Delete', style: TextStyle(color: Colors.red)),
          ),
        ],
      ),
    );
  }
}
