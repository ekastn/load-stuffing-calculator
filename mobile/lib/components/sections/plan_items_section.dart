import 'package:flutter/material.dart';
import 'package:provider/provider.dart';
import '../../providers/product_provider.dart';
import '../../components/cards/app_card.dart';
import '../../components/inputs/app_text_field.dart';
import '../../components/inputs/number_field.dart';
import '../../components/inputs/dimension_inputs.dart';
import '../../components/buttons/app_button.dart';
import '../../config/theme.dart';

class PlanItemsSection extends StatefulWidget {
  final String itemMode;
  final String? selectedProductId;
  final TextEditingController quantityController;
  final TextEditingController itemLabelController;
  final TextEditingController itemLengthController;
  final TextEditingController itemWidthController;
  final TextEditingController itemHeightController;
  final TextEditingController itemWeightController;
  final TextEditingController itemQuantityController;
  final ValueChanged<String> onItemModeChanged;
  final ValueChanged<String?> onProductSelected;
  final VoidCallback onAddCatalogItem;
  final VoidCallback onAddManualItem;

  const PlanItemsSection({
    super.key,
    required this.itemMode,
    required this.selectedProductId,
    required this.quantityController,
    required this.itemLabelController,
    required this.itemLengthController,
    required this.itemWidthController,
    required this.itemHeightController,
    required this.itemWeightController,
    required this.itemQuantityController,
    required this.onItemModeChanged,
    required this.onProductSelected,
    required this.onAddCatalogItem,
    required this.onAddManualItem,
  });

  @override
  State<PlanItemsSection> createState() => _PlanItemsSectionState();
}

class _PlanItemsSectionState extends State<PlanItemsSection> {
  @override
  Widget build(BuildContext context) {
    return AppCard(
      padding: const EdgeInsets.all(20.0),
      child: Column(
        crossAxisAlignment: CrossAxisAlignment.start,
        children: [
          Text(
            '3. Add Items',
            style: Theme.of(
              context,
            ).textTheme.titleMedium?.copyWith(fontWeight: FontWeight.bold),
          ),
          const SizedBox(height: 20),
          SizedBox(
            width: double.infinity,
            child: SegmentedButton<String>(
              segments: const [
                ButtonSegment(
                  value: 'catalog',
                  label: Text('Catalog'),
                  icon: Icon(Icons.inventory_outlined),
                ),
                ButtonSegment(
                  value: 'manual',
                  label: Text('Manual'),
                  icon: Icon(Icons.edit),
                ),
              ],
              selected: {widget.itemMode},
              onSelectionChanged: (Set<String> newSelection) {
                widget.onItemModeChanged(newSelection.first);
              },
            ),
          ),
          const SizedBox(height: 20),
          if (widget.itemMode == 'catalog') ...[
            _buildCatalogItemForm(),
          ] else ...[
            _buildManualItemForm(),
          ],
        ],
      ),
    );
  }

  Widget _buildCatalogItemForm() {
    return Consumer<ProductProvider>(
      builder: (context, provider, child) {
        if (provider.products.isEmpty && !provider.isLoading) {
          return const Padding(
            padding: EdgeInsets.symmetric(vertical: 8.0),
            child: Text(
              'No products available. Use manual mode.',
              style: TextStyle(color: AppColors.textSecondary),
            ),
          );
        }

        return Column(
          children: [
            DropdownButtonFormField<String>(
              isExpanded: true,
              initialValue: widget.selectedProductId,
              decoration: const InputDecoration(
                labelText: 'Select Product',
                border: OutlineInputBorder(),
              ),
              items: provider.products.map((product) {
                return DropdownMenuItem(
                  value: product.id,
                  child: Text(product.name),
                );
              }).toList(),
              onChanged: widget.onProductSelected,
            ),
            const SizedBox(height: 16),
            NumberField(
              controller: widget.quantityController,
              label: 'Quantity',
              required: true,
              min: 1,
            ),
            const SizedBox(height: 16),
            AppButton(
              onPressed: widget.selectedProductId == null
                  ? null
                  : widget.onAddCatalogItem,
              icon: Icons.add,
              label: 'Add to List',
              variant: AppButtonVariant.secondary,
            ),
          ],
        );
      },
    );
  }

  Widget _buildManualItemForm() {
    return Column(
      children: [
        AppTextField(
          controller: widget.itemLabelController,
          label: 'Item Label',
          hint: 'e.g., Special Cargo Box',
          required: true,
        ),
        const SizedBox(height: 16),
        DimensionInputs(
          lengthController: widget.itemLengthController,
          widthController: widget.itemWidthController,
          heightController: widget.itemHeightController,
          weightController: widget.itemWeightController,
          lengthLabel: 'Length',
          widthLabel: 'Width',
          heightLabel: 'Height',
          weightLabel: 'Weight',
        ),
        const SizedBox(height: 16),
        NumberField(
          controller: widget.itemQuantityController,
          label: 'Quantity',
          required: true,
          min: 1,
        ),
        const SizedBox(height: 16),
        AppButton(
          onPressed: widget.onAddManualItem,
          icon: Icons.add,
          label: 'Add to List',
          variant: AppButtonVariant.secondary,
        ),
      ],
    );
  }
}
