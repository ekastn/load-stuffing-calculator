import 'package:flutter/material.dart';
import '../../config/theme.dart';
import 'app_card.dart';

class ResourceListItem extends StatelessWidget {
  final Widget? leading;
  final String? code;
  final String title;
  final Widget? subtitle;
  final Widget? trailing; // Badge or Color Indicator
  final Widget? content; // Middle section (Progress bars or stats)
  final List<Widget>? actions; // Bottom row buttons
  final VoidCallback? onTap;

  const ResourceListItem({
    super.key,
    this.leading, // Optional now, as designs might not have a big left icon
    this.code,
    required this.title,
    this.subtitle,
    this.trailing,
    this.content,
    this.actions,
    this.onTap,
  });

  @override
  Widget build(BuildContext context) {
    return Padding(
      padding: const EdgeInsets.symmetric(horizontal: 16, vertical: 8), // Increased vertical spacing
      child: AppCard(
        padding: const EdgeInsets.all(20), // Increased padding
        onTap: onTap,
        child: Column(
          crossAxisAlignment: CrossAxisAlignment.start,
          children: [
            // Header Row: Code (optional)
            if (code != null) ...[
              Text(
                code!,
                style: const TextStyle(
                  fontSize: 12,
                  fontWeight: FontWeight.w500,
                  color: AppColors.textSecondary,
                  letterSpacing: 0.5,
                ),
              ),
              const SizedBox(height: 4),
            ],

            // Main Header: Title + Trailing
            Row(
              crossAxisAlignment: CrossAxisAlignment.start,
              children: [
                Expanded(
                  child: Column(
                    crossAxisAlignment: CrossAxisAlignment.start,
                    children: [
                      Text(
                        title,
                        style: const TextStyle(
                          fontSize: 18,
                          fontWeight: FontWeight.bold,
                          color: AppColors.textPrimary,
                        ),
                      ),
                      if (subtitle != null) ...[
                        const SizedBox(height: 4),
                        DefaultTextStyle(
                          style: const TextStyle(
                            fontSize: 14,
                            color: AppColors.textSecondary,
                          ),
                          child: subtitle!,
                        ),
                      ],
                    ],
                  ),
                ),
                if (trailing != null) ...[
                  const SizedBox(width: 16),
                  trailing!,
                ],
              ],
            ),

            // Content Section
            if (content != null) ...[
              const SizedBox(height: 16),
              content!,
            ],

            // Actions Section
            if (actions != null && actions!.isNotEmpty) ...[
              const SizedBox(height: 20), // Spacing before actions
              Row(
                mainAxisAlignment: MainAxisAlignment.end, // Align to right like cards usually do, or spread?
                // Designs show:
                // Plan: View (Left), Calculate (Left-Center), Delete (Right)
                // Product: Edit (Left), Delete (Right)
                // Let's use MainAxisAlignment.spaceBetween if 2 items, or custom. 
                // Actually easier to just let the parent decide, but standardizing is better.
                // Let's assume the passed 'actions' List are the buttons themselves.
                children: _buildActions(actions!),
              ),
            ],
          ],
        ),
      ),
    );
  }

  // Helper to space out actions. 
  // If 2 actions, typically Edit/Delete -> Space Between?
  // If 3 actions -> ?
  // Simple: MainAxisAlignment.spaceBetween if count == 2 (common for Edit/Delete)
  // Else: Use spacing.
  // Actually, looking at the images:
  // Image 1: View (Left) ... Calculate ... Delete (Right) -> Hard to tell exact alignment without more context.
  // Image 2: Edit (Left) ... Delete (Right)
  // Image 3: Edit (Left) ... Delete (Right)
  // So SpaceBetween seems appropriate for the "Edit/Delete" pair.
  List<Widget> _buildActions(List<Widget> widgets) {
    if (widgets.isEmpty) return [];
    
    // If only 1 action, align end? or start (as per Image 1 View)?
    // Let's just return them. The Row above handles alignment.
    // Wait, the Row above has MainAxis.end. 
    // Image 2 has Edit on LEFT, Delete on RIGHT.
    // So the Row should probably be MainAxisAlignment.spaceBetween.
    if (widgets.length > 1) {
       // For > 1 items, let's try to space them between if it's 2 items.
       // If 3 items, maybe space evenly?
       // Let's check Image 1 again. "View" is far left. "Calculate" is next to it. "Delete" is far right?
       // Or "View" "Calculate" are left-aligned group, Delete is right-aligned.
       
       // To support flexible layouts, generic is better. 
       // But simplified for now:
       return [
         Expanded(
           child: Row(
             children: [
               widgets.first,
               if (widgets.length > 2) ...[
                 const SizedBox(width: 24),
                 widgets[1],
               ]
             ],
           ),
         ),
         widgets.last, // Delete usually last
       ];
    }
    
    return widgets;
  }
}
