import 'package:flutter/material.dart';
import '../../dtos/plan_dto.dart';
import '../../components/cards/app_card.dart';
import '../../config/theme.dart';

class PlanItemsListSection extends StatelessWidget {
  final List<CreatePlanItemDto> items;
  final ValueChanged<int> onRemoveItem;

  const PlanItemsListSection({
    super.key,
    required this.items,
    required this.onRemoveItem,
  });

  @override
  Widget build(BuildContext context) {
    if (items.isEmpty) {
      return AppCard(
        padding: const EdgeInsets.all(32.0),
        child: Center(
          child: Column(
            children: [
              const Icon(
                Icons.playlist_add,
                size: 48,
                color: AppColors.textTertiary,
              ),
              const SizedBox(height: 8),
              Text(
                'No items added yet',
                style: TextStyle(color: AppColors.textSecondary),
              ),
            ],
          ),
        ),
      );
    }

    return AppCard(
      padding: const EdgeInsets.all(16.0),
      child: Column(
        crossAxisAlignment: CrossAxisAlignment.start,
        children: [
          Row(
            mainAxisAlignment: MainAxisAlignment.spaceBetween,
            children: [
              Text(
                'Items',
                style: Theme.of(
                  context,
                ).textTheme.titleMedium?.copyWith(fontWeight: FontWeight.bold),
              ),
              Container(
                padding: const EdgeInsets.symmetric(horizontal: 8, vertical: 2),
                decoration: BoxDecoration(
                  color: AppColors.primary,
                  borderRadius: BorderRadius.circular(12),
                ),
                child: Text(
                  '${items.length}',
                  style: const TextStyle(
                    color: Colors.white,
                    fontSize: 12,
                    fontWeight: FontWeight.bold,
                  ),
                ),
              ),
            ],
          ),
          const SizedBox(height: 12),
          ListView.separated(
            shrinkWrap: true,
            physics: const NeverScrollableScrollPhysics(),
            itemCount: items.length,
            separatorBuilder: (ctx, i) => const Divider(height: 1),
            itemBuilder: (ctx, index) {
              final item = items[index];
              return ListTile(
                contentPadding: EdgeInsets.zero,
                title: Text(
                  item.label ?? 'Item ${index + 1}',
                  style: const TextStyle(fontWeight: FontWeight.w500),
                ),
                subtitle: Text(
                  '${item.quantity}x • ${item.lengthMm.toInt()}×${item.widthMm.toInt()}×${item.heightMm.toInt()}mm',
                  style: const TextStyle(fontSize: 12),
                ),
                trailing: IconButton(
                  icon: const Icon(
                    Icons.delete_outline,
                    color: AppColors.error,
                  ),
                  onPressed: () => onRemoveItem(index),
                ),
              );
            },
          ),
        ],
      ),
    );
  }
}
