import 'package:flutter/material.dart';
import 'package:provider/provider.dart';
import 'package:go_router/go_router.dart';
import '../../providers/product_provider.dart';
import '../../models/product_model.dart';
import '../../dtos/product_dto.dart';
import '../../components/sections/resource_form_section.dart';
import '../../components/forms/product_form_content.dart';

class ProductFormPage extends StatefulWidget {
  final String? productId;

  const ProductFormPage({super.key, this.productId});

  @override
  State<ProductFormPage> createState() => _ProductFormPageState();
}

class _ProductFormPageState extends State<ProductFormPage> {
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
    return ResourceFormSection<ProductModel>(
      title: widget.productId == null ? 'New Product' : 'Edit Product',
      isLoading: _isLoading,
      isEditMode: widget.productId != null,
      submitLabel: widget.productId == null
          ? 'Create Product'
          : 'Update Product',
      formContent: (context, formKey) => ProductFormContent(
        nameController: _nameController,
        lengthController: _lengthController,
        widthController: _widthController,
        heightController: _heightController,
        weightController: _weightController,
        colorController: _colorController,
        onColorChanged: (color) {
          setState(() {
            _colorController.text = color;
          });
        },
        formKey: formKey,
      ),
      onSubmit: (context) => _submit(),
    );
  }

  Future<void> _submit() async {
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
        await provider.createProduct(
          CreateProductRequestDto(
            name: name,
            lengthMm: length,
            widthMm: width,
            heightMm: height,
            weightKg: weight,
            colorHex: color,
          ),
        );
      } else {
        await provider.updateProduct(
          widget.productId!,
          UpdateProductRequestDto(
            name: name,
            lengthMm: length,
            widthMm: width,
            heightMm: height,
            weightKg: weight,
            colorHex: color,
          ),
        );
      }

      if (mounted) {
        context.pop();
        ScaffoldMessenger.of(context).showSnackBar(
          SnackBar(
            content: Text(
              widget.productId == null
                  ? 'Product created successfully'
                  : 'Product updated successfully',
            ),
          ),
        );
      }
    } catch (e) {
      if (mounted) {
        ScaffoldMessenger.of(
          context,
        ).showSnackBar(SnackBar(content: Text('Error: $e')));
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
