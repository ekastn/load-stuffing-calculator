import 'package:flutter/material.dart';
import '../widgets/empty_state.dart';
import '../widgets/loading_state.dart';
import '../widgets/error_state.dart';
import '../cards/resource_list_item.dart';
import '../dialogs/confirm_dialog.dart';
import '../../config/theme.dart';

class ResourceListSection<T> extends StatelessWidget {
  final List<T> resources;
  final bool isLoading;
  final String? error;
  final String emptyMessage;
  final IconData emptyIcon;

  final String Function(T) title;
  final String Function(T) id;
  final Widget Function(BuildContext, T) trailing;
  final Widget Function(BuildContext, T) details;

  final Function(BuildContext, T) onEdit;
  final Future<void> Function(BuildContext, T) onDelete;
  final VoidCallback onRetry;
  final VoidCallback? onTap;

  const ResourceListSection({
    super.key,
    required this.resources,
    required this.isLoading,
    required this.error,
    required this.emptyMessage,
    required this.emptyIcon,
    required this.title,
    required this.id,
    required this.trailing,
    required this.details,
    required this.onEdit,
    required this.onDelete,
    required this.onRetry,
    this.onTap,
  });

  @override
  Widget build(BuildContext context) {
    if (isLoading && resources.isEmpty) {
      return const LoadingState();
    }

    if (error != null && resources.isEmpty) {
      return ErrorState(message: error!, onRetry: onRetry);
    }

    if (resources.isEmpty) {
      return EmptyState(message: emptyMessage, icon: emptyIcon);
    }

    return ListView.builder(
      itemCount: resources.length,
      itemBuilder: (context, index) {
        final resource = resources[index];
        return ResourceListItem(
          title: title(resource),
          trailing: trailing(context, resource),
          content: details(context, resource),
          actions: [
            _buildActionButton(
              context,
              label: 'Edit',
              icon: Icons.edit_outlined,
              color: AppColors.primary,
              onTap: () => onEdit(context, resource),
            ),
            _buildActionButton(
              context,
              label: 'Delete',
              icon: Icons.delete_outline,
              color: AppColors.error,
              onTap: () => _confirmDelete(context, resource),
            ),
          ],
          onTap: onTap,
        );
      },
    );
  }

  void _confirmDelete(BuildContext context, T resource) {
    ConfirmDialog.show(
      context: context,
      title: 'Delete',
      content: 'Are you sure you want to delete ${title(resource)}?',
      confirmText: 'Delete',
      isDangerous: true,
      onConfirm: () async {
        try {
          await onDelete(context, resource);
          if (context.mounted) {
            ScaffoldMessenger.of(context).showSnackBar(
              const SnackBar(content: Text('Deleted successfully')),
            );
          }
        } catch (e) {
          if (context.mounted) {
            ScaffoldMessenger.of(
              context,
            ).showSnackBar(SnackBar(content: Text('Failed to delete: $e')));
          }
        }
      },
    );
  }

  Widget _buildActionButton(
    BuildContext context, {
    required String label,
    required IconData icon,
    required Color color,
    required VoidCallback onTap,
  }) {
    return TextButton.icon(
      onPressed: onTap,
      icon: Icon(icon, size: 18, color: color),
      label: Text(
        label,
        style: TextStyle(color: color, fontWeight: FontWeight.w600),
      ),
      style: TextButton.styleFrom(
        padding: const EdgeInsets.symmetric(horizontal: 12, vertical: 8),
        minimumSize: Size.zero,
        tapTargetSize: MaterialTapTargetSize.shrinkWrap,
      ),
    );
  }
}
