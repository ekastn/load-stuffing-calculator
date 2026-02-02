import 'package:flutter/material.dart';
import 'package:provider/provider.dart';
import 'package:go_router/go_router.dart';
import '../../providers/container_provider.dart';
import '../../models/container_model.dart';
import '../../dtos/container_dto.dart';
import '../../components/inputs/app_text_field.dart';
import '../../components/inputs/number_field.dart';
import '../../components/widgets/loading_state.dart';

class ContainerFormPage extends StatefulWidget {
  final String? containerId;

  const ContainerFormPage({super.key, this.containerId});

  @override
  State<ContainerFormPage> createState() => _ContainerFormPageState();
}

class _ContainerFormPageState extends State<ContainerFormPage> {
  final _formKey = GlobalKey<FormState>();
  
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
    return Scaffold(
      appBar: AppBar(
        title: Text(widget.containerId == null ? 'New Container' : 'Edit Container'),
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
                      label: 'Container Name',
                      required: true,
                    ),
                    const SizedBox(height: 16),
                    Row(
                      children: [
                        Expanded(
                          child: NumberField(
                            controller: _lengthController,
                            label: 'Inner Length',
                            unit: 'mm',
                            required: true,
                            min: 0,
                          ),
                        ),
                        const SizedBox(width: 12),
                        Expanded(
                          child: NumberField(
                            controller: _widthController,
                            label: 'Inner Width',
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
                            label: 'Inner Height',
                            unit: 'mm',
                            required: true,
                            min: 0,
                          ),
                        ),
                        const SizedBox(width: 12),
                        Expanded(
                          child: NumberField(
                            controller: _weightController,
                            label: 'Max Weight',
                            unit: 'kg',
                            required: true,
                            min: 0,
                          ),
                        ),
                      ],
                    ),
                    const SizedBox(height: 16),
                    AppTextField(
                      controller: _descController,
                      label: 'Description',
                      hint: 'Optional description',
                      maxLines: 3,
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
                          widget.containerId == null ? 'Create Container' : 'Update Container',
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
      final provider = context.read<ContainerProvider>();
      
      final name = _nameController.text;
      final length = double.parse(_lengthController.text);
      final width = double.parse(_widthController.text);
      final height = double.parse(_heightController.text);
      final weight = double.parse(_weightController.text);
      final desc = _descController.text.isEmpty ? null : _descController.text;

      if (widget.containerId == null) {
        await provider.createContainer(CreateContainerRequestDto(
          name: name,
          innerLengthMm: length,
          innerWidthMm: width,
          innerHeightMm: height,
          maxWeightKg: weight,
          description: desc,
        ));
      } else {
        await provider.updateContainer(widget.containerId!, UpdateContainerRequestDto(
          name: name,
          innerLengthMm: length,
          innerWidthMm: width,
          innerHeightMm: height,
          maxWeightKg: weight,
          description: desc,
        ));
      }

      if (mounted) {
        context.pop();
        ScaffoldMessenger.of(context).showSnackBar(
          SnackBar(
            content: Text(widget.containerId == null 
                ? 'Container created successfully' 
                : 'Container updated successfully'),
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
    _descController.dispose();
    super.dispose();
  }
}
