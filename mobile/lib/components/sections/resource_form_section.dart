import 'package:flutter/material.dart';
import 'package:go_router/go_router.dart';
import '../widgets/loading_state.dart';
import '../buttons/app_button.dart';

class ResourceFormSection<T> extends StatefulWidget {
  final String title;
  final bool isLoading;
  final bool isEditMode;
  final Widget Function(BuildContext, GlobalKey<FormState>) formContent;
  final Future<void> Function(BuildContext) onSubmit;
  final String? submitLabel;
  final IconData? submitIcon;

  const ResourceFormSection({
    super.key,
    required this.title,
    required this.isLoading,
    required this.isEditMode,
    required this.formContent,
    required this.onSubmit,
    this.submitLabel,
    this.submitIcon,
  });

  @override
  State<ResourceFormSection<T>> createState() => _ResourceFormSectionState<T>();
}

class _ResourceFormSectionState<T> extends State<ResourceFormSection<T>> {
  final _formKey = GlobalKey<FormState>();
  bool _isSubmitting = false;

  Future<void> _handleSubmit() async {
    if (!_formKey.currentState!.validate()) {
      return;
    }

    setState(() => _isSubmitting = true);
    try {
      await widget.onSubmit(context);
    } finally {
      if (mounted) {
        setState(() => _isSubmitting = false);
      }
    }
  }

  @override
  Widget build(BuildContext context) {
    return Scaffold(
      appBar: AppBar(
        title: Text(widget.title),
        leading: IconButton(
          icon: const Icon(Icons.arrow_back),
          onPressed: () => context.pop(),
        ),
      ),
      body: widget.isLoading
          ? const LoadingState()
          : Form(
              key: _formKey,
              child: SingleChildScrollView(
                padding: const EdgeInsets.all(20.0),
                child: Column(
                  children: [
                    widget.formContent(context, _formKey),
                    const SizedBox(height: 32),
                    AppButton(
                      onPressed: _isSubmitting ? null : _handleSubmit,
                      label:
                          widget.submitLabel ??
                          (widget.isEditMode ? 'Update' : 'Create'),
                      icon: widget.submitIcon,
                      isLoading: _isSubmitting,
                      isFullWidth: true,
                    ),
                    const SizedBox(height: 20),
                  ],
                ),
              ),
            ),
    );
  }
}
