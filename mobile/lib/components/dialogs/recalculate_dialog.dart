import 'package:flutter/material.dart';

/// A dialog for recalculating a load plan with strategy selection.
///
/// Displays strategy options and gravity settling toggle, then calls the
/// provided callback with the selected parameters.
///
/// Example:
/// ```dart
/// showDialog(
///   context: context,
///   builder: (context) => RecalculateDialog(
///     onRecalculate: (strategy, goal, gravity) {
///       provider.recalculate(planId, strategy: strategy, gravity: gravity);
///     },
///   ),
/// )
/// ```
class RecalculateDialog extends StatefulWidget {
  final Function(String strategy, String? goal, bool gravity) onRecalculate;
  final String initialStrategy;
  final bool initialGravity;

  const RecalculateDialog({
    required this.onRecalculate,
    this.initialStrategy = 'bestfitdecreasing',
    this.initialGravity = true,
    super.key,
  });

  @override
  State<RecalculateDialog> createState() => _RecalculateDialogState();
}

class _RecalculateDialogState extends State<RecalculateDialog> {
  late String _strategy;
  late bool _gravity;
  String? _goal;

  @override
  void initState() {
    super.initState();
    _strategy = widget.initialStrategy;
    _gravity = widget.initialGravity;
  }

  @override
  Widget build(BuildContext context) {
    return AlertDialog(
      title: const Text('Recalculate Plan'),
      content: Column(
        mainAxisSize: MainAxisSize.min,
        crossAxisAlignment: CrossAxisAlignment.start,
        children: [
          const Text(
            'Choose packing strategy:',
            style: TextStyle(fontSize: 14, fontWeight: FontWeight.w500),
          ),
          const SizedBox(height: 12),
          DropdownButtonFormField<String>(
            value: _strategy,
            decoration: const InputDecoration(
              labelText: 'Strategy',
              border: OutlineInputBorder(),
              isDense: true,
            ),
            items: const [
              DropdownMenuItem(
                value: 'bestfitdecreasing',
                child: Text('Best Fit Decreasing'),
              ),
              DropdownMenuItem(
                value: 'minimizeboxes',
                child: Text('Minimize Boxes'),
              ),
              DropdownMenuItem(value: 'greedy', child: Text('Greedy')),
              DropdownMenuItem(
                value: 'parallel',
                child: Text('Parallel (Auto)'),
              ),
            ],
            onChanged: (value) {
              if (value != null) {
                setState(() {
                  _strategy = value;
                });
              }
            },
          ),
          const SizedBox(height: 12),
          CheckboxListTile(
            title: const Text('Gravity settling'),
            subtitle: const Text('Drop items to reduce floating'),
            value: _gravity,
            onChanged: (value) {
              setState(() {
                _gravity = value ?? true;
              });
            },
            contentPadding: EdgeInsets.zero,
          ),
        ],
      ),
      actions: [
        TextButton(
          onPressed: () => Navigator.pop(context),
          child: const Text('Cancel'),
        ),
        FilledButton(
          onPressed: () {
            Navigator.pop(context);
            widget.onRecalculate(_strategy, _goal, _gravity);
          },
          child: const Text('Recalculate'),
        ),
      ],
    );
  }
}
