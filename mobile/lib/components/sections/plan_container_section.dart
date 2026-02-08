import 'package:flutter/material.dart';
import 'package:provider/provider.dart';
import '../../providers/container_provider.dart';
import '../../components/cards/app_card.dart';
import '../../components/inputs/dimension_inputs.dart';
import '../../config/theme.dart';

class PlanContainerSection extends StatefulWidget {
  final String containerMode;
  final String? selectedContainerId;
  final TextEditingController lengthController;
  final TextEditingController widthController;
  final TextEditingController heightController;
  final TextEditingController maxWeightController;
  final ValueChanged<String> onContainerModeChanged;
  final ValueChanged<String?> onContainerSelected;

  const PlanContainerSection({
    super.key,
    required this.containerMode,
    required this.selectedContainerId,
    required this.lengthController,
    required this.widthController,
    required this.heightController,
    required this.maxWeightController,
    required this.onContainerModeChanged,
    required this.onContainerSelected,
  });

  @override
  State<PlanContainerSection> createState() => _PlanContainerSectionState();
}

class _PlanContainerSectionState extends State<PlanContainerSection> {
  @override
  Widget build(BuildContext context) {
    return AppCard(
      padding: const EdgeInsets.all(20.0),
      child: Column(
        crossAxisAlignment: CrossAxisAlignment.start,
        children: [
          Text(
            '2. Container Selection',
            style: Theme.of(
              context,
            ).textTheme.titleMedium?.copyWith(fontWeight: FontWeight.bold),
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
              selected: {widget.containerMode},
              onSelectionChanged: (Set<String> newSelection) {
                widget.onContainerModeChanged(newSelection.first);
              },
            ),
          ),
          const SizedBox(height: 20),
          if (widget.containerMode == 'preset') ...[
            _buildPresetContainerSelector(),
          ] else ...[
            DimensionInputs(
              lengthController: widget.lengthController,
              widthController: widget.widthController,
              heightController: widget.heightController,
              weightController: widget.maxWeightController,
              lengthLabel: 'Length',
              widthLabel: 'Width',
              heightLabel: 'Height',
              weightLabel: 'Max Weight',
            ),
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
            child: Text(
              'No containers available. Use custom mode.',
              style: TextStyle(color: AppColors.textSecondary),
            ),
          );
        }

        return DropdownButtonFormField<String>(
          isExpanded: true,
          initialValue: widget.selectedContainerId,
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
          onChanged: widget.onContainerSelected,
          validator: (value) =>
              value == null ? 'Please select a container' : null,
        );
      },
    );
  }
}
