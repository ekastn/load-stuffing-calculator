import 'package:flutter/material.dart';
import '../../models/plan_detail_model.dart';
import '../cards/app_card.dart';
import '../cards/info_row.dart';

class ContainerInfoSection extends StatelessWidget {
  final ContainerInfo container;

  const ContainerInfoSection({super.key, required this.container});

  @override
  Widget build(BuildContext context) {
    return Column(
      crossAxisAlignment: CrossAxisAlignment.start,
      children: [
        Text(
          'Container Information',
          style: Theme.of(
            context,
          ).textTheme.titleLarge?.copyWith(fontWeight: FontWeight.bold),
        ),
        const SizedBox(height: 16),
        AppCard(
          child: Column(
            crossAxisAlignment: CrossAxisAlignment.start,
            children: [
              if (container.name != null) ...[
                InfoRow(label: 'Name', value: container.name!),
                const SizedBox(height: 8),
              ],
              InfoRow(
                label: 'Dimensions',
                value:
                    '${container.lengthMm.toInt()} × ${container.widthMm.toInt()} × ${container.heightMm.toInt()} mm',
              ),
              const SizedBox(height: 8),
              InfoRow(
                label: 'Volume',
                value: '${container.volumeM3.toStringAsFixed(2)} m³',
              ),
              const SizedBox(height: 8),
              InfoRow(
                label: 'Max Weight',
                value: '${container.maxWeightKg.toStringAsFixed(1)} kg',
              ),
            ],
          ),
        ),
      ],
    );
  }
}
