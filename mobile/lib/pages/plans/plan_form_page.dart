import 'package:flutter/material.dart';
import 'package:provider/provider.dart';
import 'package:go_router/go_router.dart';
import '../../providers/plan_provider.dart';
import '../../providers/product_provider.dart';
import '../../providers/container_provider.dart';
import '../../dtos/plan_dto.dart';
import '../../components/inputs/app_text_field.dart';
import '../../components/inputs/number_field.dart';
import '../../components/cards/app_card.dart';
import '../../components/buttons/app_button.dart';
import '../../config/theme.dart';

class PlanFormPage extends StatefulWidget {
  const PlanFormPage({super.key});

  @override
  State<PlanFormPage> createState() => _PlanFormPageState();
}

class _PlanFormPageState extends State<PlanFormPage> {
  final _formKey = GlobalKey<FormState>();
  
  // Basic Info
  final _titleController = TextEditingController();
  final _notesController = TextEditingController();
  
  // Container Selection
  String _containerMode = 'preset'; // 'preset' or 'custom'
  String? _selectedContainerId;
  final _lengthController = TextEditingController();
  final _widthController = TextEditingController();
  final _heightController = TextEditingController();
  final _maxWeightController = TextEditingController();
  
  // Items
  String _itemMode = 'catalog'; // 'catalog' or 'manual'
  String? _selectedProductId;
  final _quantityController = TextEditingController(text: '1');
  
  // Manual item
  final _itemLabelController = TextEditingController();
  final _itemLengthController = TextEditingController();
  final _itemWidthController = TextEditingController();
  final _itemHeightController = TextEditingController();
  final _itemWeightController = TextEditingController();
  final _itemQuantityController = TextEditingController(text: '1');
  
  final List<CreatePlanItemDto> _items = [];
  bool _isSubmitting = false;

  @override
  void initState() {
    super.initState();
    WidgetsBinding.instance.addPostFrameCallback((_) {
      context.read<ProductProvider>().fetchProducts();
      context.read<ContainerProvider>().fetchContainers();
    });
  }

  @override
  Widget build(BuildContext context) {
    return Scaffold(
      appBar: AppBar(
        title: const Text('New Plan'),
      ),
      body: _isSubmitting
          ? const Center(child: CircularProgressIndicator())
          : SingleChildScrollView(
              padding: const EdgeInsets.all(20.0),
              child: Form(
                key: _formKey,
                child: Column(
                  crossAxisAlignment: CrossAxisAlignment.start,
                  children: [
                    _buildBasicInfoSection(),
                    const SizedBox(height: 24),
                    _buildContainerSection(),
                    const SizedBox(height: 24),
                    _buildItemsSection(),
                    const SizedBox(height: 24),
                    _buildItemsList(),
                    const SizedBox(height: 32),
                    _buildSubmitButton(),
                  ],
                ),
              ),
            ),
    );
  }

  Widget _buildBasicInfoSection() {
    return AppCard(
      padding: const EdgeInsets.all(20.0),
      child: Column(
        crossAxisAlignment: CrossAxisAlignment.start,
        children: [
          Text(
            '1. Basic Information',
            style: Theme.of(context).textTheme.titleMedium?.copyWith(
              fontWeight: FontWeight.bold
            ),
          ),
          const SizedBox(height: 20),
          AppTextField(
            controller: _titleController,
            label: 'Plan Title',
            hint: 'e.g., Export to Japan - December 2025',
            required: true,
          ),
          const SizedBox(height: 16),
          AppTextField(
            controller: _notesController,
            label: 'Notes',
            hint: 'Optional notes',
            maxLines: 3,
          ),
        ],
      ),
    );
  }

  Widget _buildContainerSection() {
    return AppCard(
      padding: const EdgeInsets.all(20.0),
      child: Column(
        crossAxisAlignment: CrossAxisAlignment.start,
        children: [
          Text(
            '2. Container Selection',
            style: Theme.of(context).textTheme.titleMedium?.copyWith(
              fontWeight: FontWeight.bold
            ),
          ),
          const SizedBox(height: 20),
          SizedBox(
            width: double.infinity,
            child: SegmentedButton<String>(
              segments: const [
                ButtonSegment(
                  value: 'preset',
                  label: Text('Preset'),
                  icon: Icon(Icons.inventory),
                ),
                ButtonSegment(
                  value: 'custom',
                  label: Text('Custom'),
                  icon: Icon(Icons.tune),
                ),
              ],
              selected: {_containerMode},
              onSelectionChanged: (Set<String> newSelection) {
                setState(() {
                  _containerMode = newSelection.first;
                });
              },
            ),
          ),
          const SizedBox(height: 20),
          if (_containerMode == 'preset') ...[
            _buildPresetContainerSelector(),
          ] else ...[
            _buildCustomContainerInputs(),
          ],
        ],
      ),
    );
  }

  Widget _buildPresetContainerSelector() {
    return Consumer<ContainerProvider>(
      builder: (context, provider, child) {
        if (provider.containers.isEmpty && !provider.isLoading) {
          return const Padding(
            padding: EdgeInsets.symmetric(vertical: 8.0),
            child: Text('No containers available. Use custom mode.', style: TextStyle(color: AppColors.textSecondary)),
          );
        }

        return DropdownButtonFormField<String>(
          isExpanded: true,
          initialValue: _selectedContainerId,
          decoration: const InputDecoration(
            labelText: 'Select Container',
            border: OutlineInputBorder(),
          ),
          items: provider.containers.map((container) {
            return DropdownMenuItem(
              value: container.id,
              child: Text(
                '${container.name} (${container.innerLengthMm}×${container.innerWidthMm}×${container.innerHeightMm}mm)',
                overflow: TextOverflow.ellipsis,
              ),
            );
          }).toList(),
          onChanged: (value) {
            setState(() {
              _selectedContainerId = value;
            });
          },
          validator: (value) =>
              value == null ? 'Please select a container' : null,
        );
      },
    );
  }

  Widget _buildCustomContainerInputs() {
    return Column(
      children: [
        Row(
          children: [
            Expanded(
              child: NumberField(
                controller: _lengthController,
                label: 'Length',
                unit: 'mm',
                required: true,
                min: 0,
              ),
            ),
            const SizedBox(width: 12),
            Expanded(
              child: NumberField(
                controller: _widthController,
                label: 'Width',
                unit: 'mm',
                required: true,
                min: 0,
              ),
            ),
          ],
        ),
        const SizedBox(height: 16),
        Row(
          children: [
            Expanded(
              child: NumberField(
                controller: _heightController,
                label: 'Height',
                unit: 'mm',
                required: true,
                min: 0,
              ),
            ),
            const SizedBox(width: 12),
            Expanded(
              child: NumberField(
                controller: _maxWeightController,
                label: 'Max Weight',
                unit: 'kg',
                required: true,
                min: 0,
              ),
            ),
          ],
        ),
      ],
    );
  }

  Widget _buildItemsSection() {
    return AppCard(
      padding: const EdgeInsets.all(20.0),
      child: Column(
        crossAxisAlignment: CrossAxisAlignment.start,
        children: [
          Text(
            '3. Add Items',
            style: Theme.of(context).textTheme.titleMedium?.copyWith(
              fontWeight: FontWeight.bold
            ),
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
              selected: {_itemMode},
              onSelectionChanged: (Set<String> newSelection) {
                setState(() {
                  _itemMode = newSelection.first;
                });
              },
            ),
          ),
          const SizedBox(height: 20),
          if (_itemMode == 'catalog') ...[
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
            child: Text('No products available. Use manual mode.', style: TextStyle(color: AppColors.textSecondary)),
          );
        }

        return Column(
          children: [
            DropdownButtonFormField<String>(
              isExpanded: true,
              initialValue: _selectedProductId,
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
              onChanged: (value) {
                setState(() {
                  _selectedProductId = value;
                });
              },
            ),
            const SizedBox(height: 16),
            NumberField(
              controller: _quantityController,
              label: 'Quantity',
              required: true,
              min: 1,
            ),
            const SizedBox(height: 16),
            AppButton(
              onPressed: _selectedProductId == null ? null : _addCatalogItem,
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
          controller: _itemLabelController,
          label: 'Item Label',
          hint: 'e.g., Special Cargo Box',
          required: true,
        ),
        const SizedBox(height: 16),
        Row(
          children: [
            Expanded(
              child: NumberField(
                controller: _itemLengthController,
                label: 'Length',
                unit: 'mm',
                required: true,
                min: 0,
              ),
            ),
            const SizedBox(width: 12),
            Expanded(
              child: NumberField(
                controller: _itemWidthController,
                label: 'Width',
                unit: 'mm',
                required: true,
                min: 0,
              ),
            ),
          ],
        ),
        const SizedBox(height: 16),
        Row(
          children: [
            Expanded(
              child: NumberField(
                controller: _itemHeightController,
                label: 'Height',
                unit: 'mm',
                required: true,
                min: 0,
              ),
            ),
            const SizedBox(width: 12),
            Expanded(
              child: NumberField(
                controller: _itemWeightController,
                label: 'Weight',
                unit: 'kg',
                required: true,
                min: 0,
              ),
            ),
          ],
        ),
        const SizedBox(height: 16),
        NumberField(
          controller: _itemQuantityController,
          label: 'Quantity',
          required: true,
          min: 1,
        ),
        const SizedBox(height: 16),
        AppButton(
          onPressed: _addManualItem,
          icon: Icons.add,
          label: 'Add to List',
          variant: AppButtonVariant.secondary,
        ),
      ],
    );
  }

  Widget _buildItemsList() {
    if (_items.isEmpty) {
      return AppCard(
        padding: const EdgeInsets.all(32.0),
        child: Center(
          child: Column(
            children: [
               const Icon(Icons.playlist_add, size: 48, color: AppColors.textTertiary),
               const SizedBox(height: 8),
               Text(
                'No items added yet',
                style: TextStyle(color: AppColors.textSecondary),
              ),
            ],
          ),
        ),
      );
    }

    return AppCard(
      padding: const EdgeInsets.all(16.0),
      child: Column(
        crossAxisAlignment: CrossAxisAlignment.start,
        children: [
          Row(
            mainAxisAlignment: MainAxisAlignment.spaceBetween,
            children: [
              Text(
                'Items',
                style: Theme.of(context).textTheme.titleMedium?.copyWith(fontWeight: FontWeight.bold),
              ),
              Container(
                padding: const EdgeInsets.symmetric(horizontal: 8, vertical: 2),
                decoration: BoxDecoration(
                  color: AppColors.primary,
                  borderRadius: BorderRadius.circular(12),
                ),
                child: Text(
                  '${_items.length}',
                  style: const TextStyle(color: Colors.white, fontSize: 12, fontWeight: FontWeight.bold),
                ),
              ),
            ],
          ),
          const SizedBox(height: 12),
          ListView.separated(
            shrinkWrap: true,
            physics: const NeverScrollableScrollPhysics(),
            itemCount: _items.length,
            separatorBuilder: (ctx, i) => const Divider(height: 1),
            itemBuilder: (ctx, index) {
               final item = _items[index];
               return ListTile(
                  contentPadding: EdgeInsets.zero,
                  title: Text(item.label ?? 'Item ${index + 1}', style: const TextStyle(fontWeight: FontWeight.w500)),
                  subtitle: Text(
                    '${item.quantity}x • ${item.lengthMm.toInt()}×${item.widthMm.toInt()}×${item.heightMm.toInt()}mm',
                    style: const TextStyle(fontSize: 12),
                  ),
                  trailing: IconButton(
                    icon: const Icon(Icons.delete_outline, color: AppColors.error),
                    onPressed: () {
                      setState(() {
                         _items.removeAt(index);
                      });
                    },
                  ),
               );
            },
          ),
        ],
      ),
    );
  }

  Widget _buildSubmitButton() {
    return AppButton(
      label: 'Create Plan',
      onPressed: _isSubmitting ? null : _submit,
      isLoading: _isSubmitting,
      isFullWidth: true,
    );
  }

  void _addCatalogItem() {
    if (_selectedProductId == null) return;

    final productProvider = context.read<ProductProvider>();
    final product = productProvider.products
        .firstWhere((p) => p.id == _selectedProductId);

    final quantity = int.tryParse(_quantityController.text) ?? 1;

    setState(() {
      _items.add(CreatePlanItemDto(
        label: product.name,
        productSku: product.id, // Using ID as SKU for reference
        lengthMm: product.lengthMm,
        widthMm: product.widthMm,
        heightMm: product.heightMm,
        weightKg: product.weightKg,
        quantity: quantity,
        colorHex: product.colorHex ?? '#3498db',
        allowRotation: true,
      ));
      _selectedProductId = null;
      _quantityController.text = '1';
    });

    ScaffoldMessenger.of(context).showSnackBar(
      SnackBar(content: Text('Added ${quantity}x ${product.name}')),
    );
  }

  void _addManualItem() {
    if (_itemLabelController.text.isEmpty ||
        _itemLengthController.text.isEmpty ||
        _itemWidthController.text.isEmpty ||
        _itemHeightController.text.isEmpty ||
        _itemWeightController.text.isEmpty) {
      ScaffoldMessenger.of(context).showSnackBar(
        const SnackBar(content: Text('Please fill in all fields')),
      );
      return;
    }

    setState(() {
      _items.add(CreatePlanItemDto(
        label: _itemLabelController.text,
        lengthMm: double.parse(_itemLengthController.text),
        widthMm: double.parse(_itemWidthController.text),
        heightMm: double.parse(_itemHeightController.text),
        weightKg: double.parse(_itemWeightController.text),
        quantity: int.parse(_itemQuantityController.text),
        colorHex: '#3498db',
        allowRotation: true,
      ));

      // Clear form
      _itemLabelController.clear();
      _itemLengthController.clear();
      _itemWidthController.clear();
      _itemHeightController.clear();
      _itemWeightController.clear();
      _itemQuantityController.text = '1';
    });

    ScaffoldMessenger.of(context).showSnackBar(
      const SnackBar(content: Text('Item added')),
    );
  }

  Future<void> _submit() async {
    if (!_formKey.currentState!.validate()) return;

    if (_items.isEmpty) {
      ScaffoldMessenger.of(context).showSnackBar(
        const SnackBar(content: Text('Please add at least one item')),
      );
      return;
    }

    // Validate container
    CreatePlanContainerDto container;
    if (_containerMode == 'preset') {
      if (_selectedContainerId == null) {
        ScaffoldMessenger.of(context).showSnackBar(
          const SnackBar(content: Text('Please select a container')),
        );
        return;
      }
      container = CreatePlanContainerDto(containerId: _selectedContainerId);
    } else {
      if (_lengthController.text.isEmpty ||
          _widthController.text.isEmpty ||
          _heightController.text.isEmpty ||
          _maxWeightController.text.isEmpty) {
        ScaffoldMessenger.of(context).showSnackBar(
          const SnackBar(content: Text('Please fill in all container fields')),
        );
        return;
      }
      container = CreatePlanContainerDto(
        lengthMm: double.parse(_lengthController.text),
        widthMm: double.parse(_widthController.text),
        heightMm: double.parse(_heightController.text),
        maxWeightKg: double.parse(_maxWeightController.text),
      );
    }

    setState(() => _isSubmitting = true);

    try {
      final request = CreatePlanRequestDto(
        title: _titleController.text,
        notes: _notesController.text.isEmpty ? null : _notesController.text,
        autoCalculate: true,
        container: container,
        items: _items,
      );

      await context.read<PlanProvider>().createPlan(request);

      if (mounted) {
        context.pop();
        ScaffoldMessenger.of(context).showSnackBar(
          const SnackBar(content: Text('Plan created successfully')),
        );
      }
    } catch (e) {
      if (mounted) {
        ScaffoldMessenger.of(context).showSnackBar(
          SnackBar(content: Text('Error: $e')),
        );
      }
    } finally {
      if (mounted) setState(() => _isSubmitting = false);
    }
  }

  @override
  void dispose() {
    _titleController.dispose();
    _notesController.dispose();
    _lengthController.dispose();
    _widthController.dispose();
    _heightController.dispose();
    _maxWeightController.dispose();
    _quantityController.dispose();
    _itemLabelController.dispose();
    _itemLengthController.dispose();
    _itemWidthController.dispose();
    _itemHeightController.dispose();
    _itemWeightController.dispose();
    _itemQuantityController.dispose();
    super.dispose();
  }
}
