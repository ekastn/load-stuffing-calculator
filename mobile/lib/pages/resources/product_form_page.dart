import 'package:flutter/material.dart';
import 'package:provider/provider.dart';
import 'package:go_router/go_router.dart';
import '../../providers/product_provider.dart';
import '../../models/product_model.dart';
import '../../dtos/product_dto.dart';
import '../../components/inputs/app_text_field.dart';
import '../../components/inputs/number_field.dart';
import '../../components/widgets/loading_state.dart';

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
              padding: const EdgeInsets.all(16.0),
              child: Form(
                key: _formKey,
                child: Column(
                  children: [
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
                    AppTextField(
                      controller: _colorController,
                      label: 'Color (Hex)',
                      hint: '#3498db',
                      validator: (v) => v == null || !v.startsWith('#') 
                          ? 'Must start with #' 
                          : null,
                    ),
                    const SizedBox(height: 32),
                    SizedBox(
                      width: double.infinity,
                      child: ElevatedButton(
                        onPressed: _submit,
                        style: ElevatedButton.styleFrom(
                          padding: const EdgeInsets.symmetric(vertical: 16),
                        ),
                        child: Text(
                          widget.productId == null ? 'Create Product' : 'Update Product',
                          style: const TextStyle(fontSize: 16),
                        ),
                      ),
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
