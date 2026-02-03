import 'package:flutter/material.dart';

class ResourceListItem extends StatelessWidget {
  final Widget leading;
  final String title;
  final String subtitle;
  final VoidCallback? onTap;
  final VoidCallback? onEdit;
  final VoidCallback? onDelete;
  final List<Widget>? additionalActions;

  const ResourceListItem({
    super.key,
    required this.leading,
    required this.title,
    required this.subtitle,
    this.onTap,
    this.onEdit,
    this.onDelete,
    this.additionalActions,
  });

  @override
  Widget build(BuildContext context) {
    final actions = <Widget>[];
    
    if (onEdit != null) {
      actions.add(
        IconButton(
          icon: const Icon(Icons.edit, color: Colors.blue),
          onPressed: onEdit,
          tooltip: 'Edit',
        ),
      );
    }
    
    if (additionalActions != null) {
      actions.addAll(additionalActions!);
    }
    
    if (onDelete != null) {
      actions.add(
        IconButton(
          icon: const Icon(Icons.delete, color: Colors.red),
          onPressed: onDelete,
          tooltip: 'Delete',
        ),
      );
    }

    return Card(
      margin: const EdgeInsets.symmetric(horizontal: 16, vertical: 8),
      child: ListTile(
        leading: leading,
        title: Text(title),
        subtitle: Text(subtitle),
        onTap: onTap,
        trailing: actions.isNotEmpty
            ? Row(
                mainAxisSize: MainAxisSize.min,
                children: actions,
              )
            : null,
      ),
    );
  }
}
