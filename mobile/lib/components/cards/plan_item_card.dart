import 'package:flutter/material.dart';
import '../../config/theme.dart';
import '../../models/plan_detail_model.dart';
import 'app_card.dart';

/// A card widget displaying a single plan item with its specifications.
///
/// Shows item color dot, name, quantity, dimensions, weight, and total volume.
/// Used in plan detail and form pages to display items in a consistent format.
///
/// Example:
/// ```dart
/// PlanItemCard(
///   item: myLoadItem,
///   onDelete: () => removeItem(myLoadItem.id),
/// )
/// ```
class PlanItemCard extends StatelessWidget {
  final PlanItem item;
  final VoidCallback? onDelete;
  final bool editable;

  const PlanItemCard({
    required this.item,
    this.onDelete,
    this.editable = false,
    super.key,
  });

  Color _parseColor(String? colorHex) {
    if (colorHex == null) return Colors.grey;
    try {
      return Color(int.parse(colorHex.substring(1), radix: 16) + 0xFF000000);
    } catch (_) {
      return Colors.grey;
    }
  }

  @override
  Widget build(BuildContext context) {
    return AppCard(
      padding: const EdgeInsets.symmetric(horizontal: 16, vertical: 12),
      child: Row(
        children: [
          // Color indicator
          Container(
            width: 16,
            height: 16,
            decoration: BoxDecoration(
              color: _parseColor(item.colorHex),
              shape: BoxShape.circle,
            ),
          ),
          const SizedBox(width: 16),
          // Item details
          Expanded(
            child: Column(
              crossAxisAlignment: CrossAxisAlignment.start,
              children: [
                Text(
                  item.label ?? item.itemId,
                  style: const TextStyle(
                    fontWeight: FontWeight.w600,
                    fontSize: 16,
                  ),
                ),
                const SizedBox(height: 4),
                Text(
                  '${item.quantity}× • ${item.lengthMm.toInt()}×${item.widthMm.toInt()}×${item.heightMm.toInt()} mm • ${item.weightKg} kg',
                  style: TextStyle(
                    fontSize: 13,
                    color: AppColors.textSecondary,
                  ),
                ),
              ],
            ),
          ),
          // Weight and volume
          Column(
            mainAxisAlignment: MainAxisAlignment.center,
            crossAxisAlignment: CrossAxisAlignment.end,
            children: [
              Text(
                '${item.totalWeightKg.toStringAsFixed(1)} kg',
                style: const TextStyle(
                  fontWeight: FontWeight.bold,
                  fontSize: 14,
                ),
              ),
              Text(
                '${item.totalVolumeM3.toStringAsFixed(3)} m³',
                style: TextStyle(fontSize: 12, color: AppColors.textSecondary),
              ),
            ],
          ),
          // Delete button (if editable)
          if (editable && onDelete != null) ...[
            const SizedBox(width: 12),
            InkWell(
              onTap: onDelete,
              borderRadius: BorderRadius.circular(8),
              child: const Padding(
                padding: EdgeInsets.all(8.0),
                child: Icon(Icons.close, size: 20, color: AppColors.error),
              ),
            ),
          ],
        ],
      ),
    );
  }
}
