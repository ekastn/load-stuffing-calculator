import 'package:flutter/material.dart';
import 'package:provider/provider.dart';
import 'package:go_router/go_router.dart';
import '../../providers/product_provider.dart';
import '../../models/product_model.dart';
import '../../components/sections/resource_list_section.dart';
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
        return ResourceListSection<ProductModel>(
          resources: provider.products,
          isLoading: provider.isLoading,
          error: provider.error,
          emptyMessage: 'No products found.\nTap + to add a new product.',
          emptyIcon: Icons.inventory_outlined,
          title: (product) => product.name,
          id: (product) => product.id,
          trailing: (context, product) => Container(
            width: 24,
            height: 24,
            decoration: BoxDecoration(
              color: product.colorHex != null
                  ? Color(
                      int.parse(product.colorHex!.replaceFirst('#', '0xFF')),
                    )
                  : Colors.grey,
              borderRadius: BorderRadius.circular(6),
            ),
          ),
          details: (context, product) => Column(
            crossAxisAlignment: CrossAxisAlignment.start,
            children: [
              Text(
                'Dimensions: ${product.lengthMm.toInt()} × ${product.widthMm.toInt()} × ${product.heightMm.toInt()} mm',
                style: const TextStyle(
                  color: AppColors.textSecondary,
                  fontSize: 14,
                ),
              ),
              const SizedBox(height: 4),
              Text(
                'Weight: ${product.weightKg.toStringAsFixed(1)} kg',
                style: const TextStyle(
                  color: AppColors.textSecondary,
                  fontSize: 14,
                ),
              ),
            ],
          ),
          onEdit: (context, product) => context.push('/products/${product.id}'),
          onDelete: (context, product) =>
              context.read<ProductProvider>().deleteProduct(product.id),
          onRetry: () => context.read<ProductProvider>().fetchProducts(),
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
}
