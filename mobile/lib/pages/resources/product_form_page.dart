import 'package:flutter/material.dart';
import 'package:provider/provider.dart';
import 'package:go_router/go_router.dart';
import '../../providers/product_provider.dart';
import '../../models/product_model.dart';
import '../../dtos/product_dto.dart';
import '../../components/inputs/app_text_field.dart';
import '../../components/inputs/number_field.dart';
import '../../components/widgets/loading_state.dart';
import '../../components/cards/app_card.dart';
import '../../components/buttons/app_button.dart';
import '../../config/theme.dart';

class ProductFormPage extends StatefulWidget {
  final String? productId;

  const ProductFormPage({super.key, this.productId});

  @override
  State<ProductFormPage> createState() => _ProductFormPageState();
}

class _ProductFormPageState extends State<ProductFormPage> {
  final _formKey = GlobalKey<FormState>();
  
  // Controllers
  final _nameController = TextEditingController();
  final _lengthController = TextEditingController();
  final _widthController = TextEditingController();
  final _heightController = TextEditingController();
  final _weightController = TextEditingController();
  final _colorController = TextEditingController(text: '#3498db');

  bool _isLoading = false;

  @override
  void initState() {
    super.initState();
    if (widget.productId != null) {
      _loadProductData();
    }
  }

  Future<void> _loadProductData() async {
    setState(() => _isLoading = true);
    try {
      final provider = context.read<ProductProvider>();
      var product = provider.products.cast<ProductModel?>().firstWhere(
            (p) => p?.id == widget.productId,
            orElse: () => null,
          );
      
      if (product != null) {
        _nameController.text = product.name;
        _lengthController.text = product.lengthMm.toString();
        _widthController.text = product.widthMm.toString();
        _heightController.text = product.heightMm.toString();
        _weightController.text = product.weightKg.toString();
        _colorController.text = product.colorHex ?? '#3498db';
      }
    } finally {
      setState(() => _isLoading = false);
    }
  }

  @override
  Widget build(BuildContext context) {
    return Scaffold(
      appBar: AppBar(
        title: Text(widget.productId == null ? 'New Product' : 'Edit Product'),
      ),
      body: _isLoading 
          ? const LoadingState()
          : SingleChildScrollView(
              padding: const EdgeInsets.all(20.0),
              child: Form(
                key: _formKey,
                child: Column(
                  children: [
                    AppCard(
                        padding: const EdgeInsets.all(20),
                        child: Column(
                            crossAxisAlignment: CrossAxisAlignment.start,
                            children: [
                              Text(
                                'Product Details',
                                style: Theme.of(context).textTheme.titleMedium?.copyWith(
                                    fontWeight: FontWeight.bold
                                ),
                              ),
                              const SizedBox(height: 20),
                              AppTextField(
                                controller: _nameController,
                                label: 'Product Name',
                                required: true,
                              ),
                              const SizedBox(height: 16),
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
                                      controller: _weightController,
                                      label: 'Weight',
                                      unit: 'kg',
                                      required: true,
                                      min: 0,
                                    ),
                                  ),
                                ],
                              ),
                              const SizedBox(height: 16),
                              _buildColorPicker(),
                            ]
                        )
                    ),
                    const SizedBox(height: 32),
                    AppButton(
                        onPressed: _submit,
                        label: widget.productId == null ? 'Create Product' : 'Update Product',
                        isFullWidth: true,
                        isLoading: _isLoading,
                    ),
                  ],
                ),
              ),
            ),
    );
  }

  Future<void> _submit() async {
    if (!_formKey.currentState!.validate()) return;

    setState(() => _isLoading = true);
    try {
      final provider = context.read<ProductProvider>();
      
      final name = _nameController.text;
      final length = double.parse(_lengthController.text);
      final width = double.parse(_widthController.text);
      final height = double.parse(_heightController.text);
      final weight = double.parse(_weightController.text);
      final color = _colorController.text;

      if (widget.productId == null) {
        await provider.createProduct(CreateProductRequestDto(
          name: name,
          lengthMm: length,
          widthMm: width,
          heightMm: height,
          weightKg: weight,
          colorHex: color,
        ));
      } else {
        await provider.updateProduct(widget.productId!, UpdateProductRequestDto(
          name: name,
          lengthMm: length,
          widthMm: width,
          heightMm: height,
          weightKg: weight,
          colorHex: color,
        ));
      }

      if (mounted) {
        context.pop();
        ScaffoldMessenger.of(context).showSnackBar(
          SnackBar(
            content: Text(widget.productId == null 
                ? 'Product created successfully' 
                : 'Product updated successfully'),
          ),
        );
      }
    } catch (e) {
      if (mounted) {
        ScaffoldMessenger.of(context).showSnackBar(
          SnackBar(content: Text('Error: $e')),
        );
      }
    } finally {
      if (mounted) setState(() => _isLoading = false);
    }
  }

  Widget _buildColorPicker() {
    final List<String> palette = [
      '#3498db', // Blue
      '#e74c3c', // Red
      '#2ecc71', // Green
      '#f1c40f', // Yellow
      '#9b59b6', // Purple
      '#e67e22', // Orange
      '#1abc9c', // Teal
      '#34495e', // Navy
      '#7f8c8d', // Grey
    ];

    return Column(
      crossAxisAlignment: CrossAxisAlignment.start,
      children: [
         const Text(
          'Color',
          style: TextStyle(fontSize: 14, fontWeight: FontWeight.w500, color: AppColors.textPrimary),
        ),
        const SizedBox(height: 12),
        Wrap(
          spacing: 12,
          runSpacing: 12,
          children: palette.map((colorHex) {
            final isSelected = _colorController.text.toLowerCase() == colorHex.toLowerCase();
            return GestureDetector(
              onTap: () {
                setState(() {
                  _colorController.text = colorHex;
                });
              },
              child: Container(
                width: 40,
                height: 40,
                decoration: BoxDecoration(
                  color: Color(int.parse(colorHex.replaceFirst('#', '0xFF'))),
                  shape: BoxShape.circle,
                  border: isSelected 
                      ? Border.all(color: AppColors.primary, width: 3)
                      : Border.all(color: Colors.transparent),
                  boxShadow: [
                    BoxShadow(
                      color: Colors.black.withValues(alpha: 0.1),
                      blurRadius: 4,
                      offset: const Offset(0, 2),
                    ),
                  ],
                ),
                child: isSelected 
                    ? const Icon(Icons.check, color: Colors.white, size: 24)
                    : null,
              ),
            );
          }).toList(),
        ),
        const SizedBox(height: 16),
        AppTextField(
          controller: _colorController,
          label: 'Custom Hex Color',
          hint: '#3498db',
          validator: (v) => v == null || !v.startsWith('#') || v.length != 7
              ? 'Invalid Hex (e.g. #3498db)' 
              : null,
          onChanged: (val) => setState(() {}),
        ),
      ],
    );
  }

  @override
  void dispose() {
    _nameController.dispose();
    _lengthController.dispose();
    _widthController.dispose();
    _heightController.dispose();
    _weightController.dispose();
    _colorController.dispose();
    super.dispose();
  }
}
