import 'package:flutter/material.dart';

class AppColors {
  // Approximate conversions from Oklch in globals.css
  // Primary: oklch(0.3 0.15 260) -> Deep Blue/Purple
  static const Color primary = Color(0xFF1E3A8A); 
  static const Color primaryForeground = Color(0xFFFAFAFA);

  // Secondary: oklch(0.65 0.2 30) -> Orange/Red-ish
  static const Color secondary = Color(0xFFEA580C);
  static const Color secondaryForeground = Color(0xFFFAFAFA);

  // Backgrounds
  static const Color background = Color(0xFFFAFAFA); // oklch(0.98 0 0)
  static const Color surface = Color(0xFFFFFFFF);    // oklch(1 0 0) - card/popover

  // Text
  static const Color foreground = Color(0xFF171717); // oklch(0.16 0 0)
  static const Color muted = Color(0xFFF5F5F5);      // oklch(0.94 0 0)
  static const Color mutedForeground = Color(0xFF737373); // oklch(0.5 0 0)

  // Border & Input
  static const Color border = Color(0xFFE5E5E5);    // oklch(0.92 0 0)
  static const Color input = Color(0xFFF5F5F5);     // oklch(0.96 0 0)

  // Status
  static const Color destructive = Color(0xFFDC2626); // oklch(0.577 0.245 27.325)
}

class AppTheme {
  static final ThemeData lightTheme = ThemeData(
    useMaterial3: true,
    colorScheme: ColorScheme.light(
      primary: AppColors.primary,
      onPrimary: AppColors.primaryForeground,
      secondary: AppColors.secondary,
      onSecondary: AppColors.secondaryForeground,
      surface: AppColors.surface,
      onSurface: AppColors.foreground,
      error: AppColors.destructive,
      onError: Colors.white,
      background: AppColors.background,
      onBackground: AppColors.foreground,
      outline: AppColors.border,
    ),
    scaffoldBackgroundColor: AppColors.background,
    cardTheme: const CardTheme(
      color: AppColors.surface,
      elevation: 0,
      margin: EdgeInsets.all(0),
      shape: RoundedRectangleBorder(
        side: BorderSide(color: AppColors.border),
        borderRadius: BorderRadius.all(Radius.circular(8)),
      ),
    ),
    inputDecorationTheme: InputDecorationTheme(
      filled: true,
      fillColor: AppColors.surface,
      border: OutlineInputBorder(
        borderRadius: BorderRadius.circular(8),
        borderSide: const BorderSide(color: AppColors.input),
      ),
      enabledBorder: OutlineInputBorder(
        borderRadius: BorderRadius.circular(8),
        borderSide: const BorderSide(color: AppColors.input),
      ),
      focusedBorder: OutlineInputBorder(
        borderRadius: BorderRadius.circular(8),
        borderSide: const BorderSide(color: AppColors.primary, width: 2),
      ),
      contentPadding: const EdgeInsets.symmetric(horizontal: 16, vertical: 12),
    ),
    elevatedButtonTheme: ElevatedButtonThemeData(
      style: ElevatedButton.styleFrom(
        backgroundColor: AppColors.primary,
        foregroundColor: AppColors.primaryForeground,
        shape: RoundedRectangleBorder(borderRadius: BorderRadius.circular(8)),
        padding: const EdgeInsets.symmetric(horizontal: 24, vertical: 12),
        elevation: 0,
      ),
    ),
  );
}
