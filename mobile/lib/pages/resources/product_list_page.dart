import 'package:flutter/material.dart';
import 'package:provider/provider.dart';
import 'package:go_router/go_router.dart';
import '../../providers/product_provider.dart';
import '../../models/product_model.dart';
import '../../components/widgets/empty_state.dart';
import '../../components/widgets/loading_state.dart';
import '../../components/widgets/error_state.dart';
import '../../components/cards/resource_list_item.dart';
import '../../components/dialogs/confirm_dialog.dart';

class ProductListPage extends StatefulWidget {
  const ProductListPage({super.key});

  @override
  State<ProductListPage> createState() => _ProductListPageState();
}

class _ProductListPageState extends State<ProductListPage> {
  @override
  void initState() {
    super.initState();
    WidgetsBinding.instance.addPostFrameCallback((_) {
      context.read<ProductProvider>().fetchProducts();
    });
  }

  @override
  Widget build(BuildContext context) {
    return Scaffold(
      appBar: AppBar(
        title: const Text('Products'),
        actions: [
          IconButton(
            icon: const Icon(Icons.refresh),
            onPressed: () => context.read<ProductProvider>().fetchProducts(),
          ),
        ],
      ),
      floatingActionButton: FloatingActionButton(
        onPressed: () => context.push('/products/new'),
        child: const Icon(Icons.add),
      ),
      body: Consumer<ProductProvider>(
        builder: (context, provider, child) {
          if (provider.isLoading && provider.products.isEmpty) {
            return const LoadingState();
          }

          if (provider.error != null && provider.products.isEmpty) {
            return ErrorState(
              message: provider.error!,
              onRetry: () => provider.fetchProducts(),
            );
          }

          if (provider.products.isEmpty) {
            return const EmptyState(
              message: 'No products found.\nTap + to add a new product.',
              icon: Icons.inventory_outlined,
            );
          }

          return ListView.builder(
            itemCount: provider.products.length,
            itemBuilder: (context, index) {
              final product = provider.products[index];
              return ResourceListItem(
                leading: CircleAvatar(
                  backgroundColor: product.colorHex != null
                      ? Color(int.parse(
                          product.colorHex!.replaceFirst('#', '0xFF')))
                      : Colors.grey,
                  child: Text(
                    product.name.substring(0, 1).toUpperCase(),
                    style: const TextStyle(color: Colors.white),
                  ),
                ),
                title: product.name,
                subtitle:
                    '${product.lengthMm}×${product.widthMm}×${product.heightMm}mm • ${product.weightKg}kg',
                onEdit: () => context.push('/products/${product.id}'),
                onDelete: () => _confirmDelete(context, product),
              );
            },
          );
        },
      ),
    );
  }

  void _confirmDelete(BuildContext context, ProductModel product) {
    ConfirmDialog.show(
      context: context,
      title: 'Delete Product',
      content: 'Are you sure you want to delete ${product.name}?',
      confirmText: 'Delete',
      isDangerous: true,
      onConfirm: () async {
        try {
          await context.read<ProductProvider>().deleteProduct(product.id);
          if (context.mounted) {
            ScaffoldMessenger.of(context).showSnackBar(
              const SnackBar(content: Text('Product deleted')),
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
