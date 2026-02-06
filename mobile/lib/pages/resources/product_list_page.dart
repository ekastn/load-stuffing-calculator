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
import '../../config/theme.dart';

class ProductListPage extends StatefulWidget {
  final bool isEmbedded;
  const ProductListPage({super.key, this.isEmbedded = false});

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
    final content = Consumer<ProductProvider>(
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
                code: product.id.length > 8 ? product.id.substring(0, 8).toUpperCase() : product.id.toUpperCase(),
                title: product.name,
                trailing: Container(
                  width: 24,
                  height: 24,
                  decoration: BoxDecoration(
                    color: product.colorHex != null
                        ? Color(int.parse(product.colorHex!.replaceFirst('#', '0xFF')))
                        : Colors.grey,
                    borderRadius: BorderRadius.circular(6),
                  ),
                ),
                content: Column(
                  crossAxisAlignment: CrossAxisAlignment.start,
                  children: [
                    Text(
                      'Dimensions: ${product.lengthMm.toInt()} × ${product.widthMm.toInt()} × ${product.heightMm.toInt()} mm',
                      style: const TextStyle(color: AppColors.textSecondary, fontSize: 14),
                    ),
                    const SizedBox(height: 4),
                    Text(
                      'Weight: ${product.weightKg.toStringAsFixed(1)} kg',
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
                    onTap: () => context.push('/products/${product.id}'),
                  ),
                  _buildActionButton(
                    context,
                    label: 'Delete',
                    icon: Icons.delete_outline,
                    color: AppColors.error,
                    onTap: () => _confirmDelete(context, product),
                  ),
                ],
                onTap: () => context.push('/products/${product.id}'),
              );
            },
          );
        },
      );

    if (widget.isEmbedded) {
      return Scaffold(
        body: content,
        floatingActionButton: FloatingActionButton(
          onPressed: () => context.push('/products/new'),
          child: const Icon(Icons.add),
        ),
      );
    }

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
      body: content,
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
