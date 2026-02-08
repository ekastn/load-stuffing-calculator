import 'package:flutter/material.dart';
import '../inputs/app_text_field.dart';
import '../../config/theme.dart';

class ColorPickerField extends StatefulWidget {
  final String initialColor;
  final ValueChanged<String> onColorChanged;
  final String? label;
  final bool readOnly;

  const ColorPickerField({
    super.key,
    required this.initialColor,
    required this.onColorChanged,
    this.label = 'Color',
    this.readOnly = false,
  });

  @override
  State<ColorPickerField> createState() => _ColorPickerFieldState();
}

class _ColorPickerFieldState extends State<ColorPickerField> {
  late TextEditingController _colorController;

  final List<String> _palette = [
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

  @override
  void initState() {
    super.initState();
    _colorController = TextEditingController(text: widget.initialColor);
  }

  @override
  Widget build(BuildContext context) {
    return Column(
      crossAxisAlignment: CrossAxisAlignment.start,
      children: [
        if (widget.label != null)
          Text(
            widget.label!,
            style: const TextStyle(
              fontSize: 14,
              fontWeight: FontWeight.w500,
              color: AppColors.textPrimary,
            ),
          ),
        if (widget.label != null) const SizedBox(height: 12),
        Wrap(
          spacing: 12,
          runSpacing: 12,
          children: _palette.map((colorHex) {
            final isSelected =
                _colorController.text.toLowerCase() == colorHex.toLowerCase();
            return GestureDetector(
              onTap: widget.readOnly
                  ? null
                  : () {
                      setState(() {
                        _colorController.text = colorHex;
                      });
                      widget.onColorChanged(colorHex);
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
          onChanged: (val) {
            setState(() {});
            if (val.startsWith('#') && val.length == 7) {
              widget.onColorChanged(val);
            }
          },
        ),
      ],
    );
  }

  @override
  void dispose() {
    _colorController.dispose();
    super.dispose();
  }
}
