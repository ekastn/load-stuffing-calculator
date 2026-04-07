import 'package:flutter/material.dart';
import 'package:provider/provider.dart';
import 'package:go_router/go_router.dart';
import '../../providers/plan_provider.dart';
import '../../providers/product_provider.dart';
import '../../providers/container_provider.dart';
import '../../dtos/plan_dto.dart';
import '../../components/buttons/app_button.dart';
import '../../components/sections/plan_basic_info_section.dart';
import '../../components/sections/plan_container_section.dart';
import '../../components/sections/plan_items_section.dart';
import '../../components/sections/plan_items_list_section.dart';

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
      appBar: AppBar(title: const Text('New Plan')),
      body: _isSubmitting
          ? const Center(child: CircularProgressIndicator())
          : SingleChildScrollView(
              padding: const EdgeInsets.all(20.0),
              child: Form(
                key: _formKey,
                child: Column(
                  crossAxisAlignment: CrossAxisAlignment.start,
                  children: [
                    PlanBasicInfoSection(
                      titleController: _titleController,
                      notesController: _notesController,
                    ),
                    const SizedBox(height: 24),
                    PlanContainerSection(
                      containerMode: _containerMode,
                      selectedContainerId: _selectedContainerId,
                      lengthController: _lengthController,
                      widthController: _widthController,
                      heightController: _heightController,
                      maxWeightController: _maxWeightController,
                      onContainerModeChanged: (mode) {
                        setState(() => _containerMode = mode);
                      },
                      onContainerSelected: (id) {
                        setState(() => _selectedContainerId = id);
                      },
                    ),
                    const SizedBox(height: 24),
                    PlanItemsSection(
                      itemMode: _itemMode,
                      selectedProductId: _selectedProductId,
                      quantityController: _quantityController,
                      itemLabelController: _itemLabelController,
                      itemLengthController: _itemLengthController,
                      itemWidthController: _itemWidthController,
                      itemHeightController: _itemHeightController,
                      itemWeightController: _itemWeightController,
                      itemQuantityController: _itemQuantityController,
                      onItemModeChanged: (mode) {
                        setState(() => _itemMode = mode);
                      },
                      onProductSelected: (id) {
                        setState(() => _selectedProductId = id);
                      },
                      onAddCatalogItem: _addCatalogItem,
                      onAddManualItem: _addManualItem,
                    ),
                    const SizedBox(height: 24),
                    PlanItemsListSection(
                      items: _items,
                      onRemoveItem: (index) {
                        setState(() => _items.removeAt(index));
                      },
                    ),
                    const SizedBox(height: 32),
                    AppButton(
                      label: 'Create Plan',
                      onPressed: _isSubmitting ? null : _submit,
                      isLoading: _isSubmitting,
                      isFullWidth: true,
                    ),
                  ],
                ),
              ),
            ),
    );
  }

  void _addCatalogItem() {
    if (_selectedProductId == null) return;

    final productProvider = context.read<ProductProvider>();
    final product = productProvider.products.firstWhere(
      (p) => p.id == _selectedProductId,
    );

    final quantity = int.tryParse(_quantityController.text) ?? 1;

    setState(() {
      _items.add(
        CreatePlanItemDto(
          label: product.name,
          productSku: product.id, // Using ID as SKU for reference
          lengthMm: product.lengthMm,
          widthMm: product.widthMm,
          heightMm: product.heightMm,
          weightKg: product.weightKg,
          quantity: quantity,
          colorHex: product.colorHex ?? '#3498db',
          allowRotation: true,
        ),
      );
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
      _items.add(
        CreatePlanItemDto(
          label: _itemLabelController.text,
          lengthMm: double.parse(_itemLengthController.text),
          widthMm: double.parse(_itemWidthController.text),
          heightMm: double.parse(_itemHeightController.text),
          weightKg: double.parse(_itemWeightController.text),
          quantity: int.parse(_itemQuantityController.text),
          colorHex: '#3498db',
          allowRotation: true,
        ),
      );

      // Clear form
      _itemLabelController.clear();
      _itemLengthController.clear();
      _itemWidthController.clear();
      _itemHeightController.clear();
      _itemWeightController.clear();
      _itemQuantityController.text = '1';
    });

    ScaffoldMessenger.of(
      context,
    ).showSnackBar(const SnackBar(content: Text('Item added')));
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
        ScaffoldMessenger.of(
          context,
        ).showSnackBar(SnackBar(content: Text('Error: $e')));
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
