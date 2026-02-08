import 'package:flutter/material.dart';
import 'package:provider/provider.dart';
import 'package:go_router/go_router.dart';
import '../../providers/container_provider.dart';
import '../../models/container_model.dart';
import '../../dtos/container_dto.dart';
import '../../components/sections/resource_form_section.dart';
import '../../components/forms/container_form_content.dart';

class ContainerFormPage extends StatefulWidget {
  final String? containerId;

  const ContainerFormPage({super.key, this.containerId});

  @override
  State<ContainerFormPage> createState() => _ContainerFormPageState();
}

class _ContainerFormPageState extends State<ContainerFormPage> {
  // Controllers
  final _nameController = TextEditingController();
  final _lengthController = TextEditingController();
  final _widthController = TextEditingController();
  final _heightController = TextEditingController();
  final _weightController = TextEditingController();
  final _descController = TextEditingController();

  bool _isLoading = false;

  @override
  void initState() {
    super.initState();
    if (widget.containerId != null) {
      _loadContainerData();
    }
  }

  Future<void> _loadContainerData() async {
    setState(() => _isLoading = true);
    try {
      final provider = context.read<ContainerProvider>();
      var container = provider.containers.cast<ContainerModel?>().firstWhere(
        (c) => c?.id == widget.containerId,
        orElse: () => null,
      );

      if (container != null) {
        _nameController.text = container.name;
        _lengthController.text = container.innerLengthMm.toString();
        _widthController.text = container.innerWidthMm.toString();
        _heightController.text = container.innerHeightMm.toString();
        _weightController.text = container.maxWeightKg.toString();
        _descController.text = container.description ?? '';
      }
    } finally {
      setState(() => _isLoading = false);
    }
  }

  @override
  Widget build(BuildContext context) {
    return ResourceFormSection<ContainerModel>(
      title: widget.containerId == null ? 'New Container' : 'Edit Container',
      isLoading: _isLoading,
      isEditMode: widget.containerId != null,
      submitLabel: widget.containerId == null
          ? 'Create Container'
          : 'Update Container',
      formContent: (context, formKey) => ContainerFormContent(
        nameController: _nameController,
        lengthController: _lengthController,
        widthController: _widthController,
        heightController: _heightController,
        weightController: _weightController,
        descController: _descController,
        formKey: formKey,
      ),
      onSubmit: (context) => _submit(),
    );
  }

  Future<void> _submit() async {
    setState(() => _isLoading = true);
    try {
      final provider = context.read<ContainerProvider>();

      final name = _nameController.text;
      final length = double.parse(_lengthController.text);
      final width = double.parse(_widthController.text);
      final height = double.parse(_heightController.text);
      final weight = double.parse(_weightController.text);
      final desc = _descController.text.isEmpty ? null : _descController.text;

      if (widget.containerId == null) {
        await provider.createContainer(
          CreateContainerRequestDto(
            name: name,
            innerLengthMm: length,
            innerWidthMm: width,
            innerHeightMm: height,
            maxWeightKg: weight,
            description: desc,
          ),
        );
      } else {
        await provider.updateContainer(
          widget.containerId!,
          UpdateContainerRequestDto(
            name: name,
            innerLengthMm: length,
            innerWidthMm: width,
            innerHeightMm: height,
            maxWeightKg: weight,
            description: desc,
          ),
        );
      }

      if (mounted) {
        context.pop();
        ScaffoldMessenger.of(context).showSnackBar(
          SnackBar(
            content: Text(
              widget.containerId == null
                  ? 'Container created successfully'
                  : 'Container updated successfully',
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
    _descController.dispose();
    super.dispose();
  }
}
