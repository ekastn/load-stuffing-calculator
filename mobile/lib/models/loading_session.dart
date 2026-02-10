import 'package:freezed_annotation/freezed_annotation.dart';

part 'loading_session.freezed.dart';
part 'loading_session.g.dart';

@freezed
abstract class LoadingSession with _$LoadingSession {
  const factory LoadingSession({
    required String sessionId,
    required String planId,
    required DateTime startedAt,
    DateTime? completedAt,
    required int totalItems,
    @Default(0) int validatedCount,
    @Default(0) int skippedCount,
    @Default(0) int currentStepIndex,
    @Default([]) List<ItemValidation> validations,
    @Default('IN_PROGRESS') String status,
  }) = _LoadingSession;

  factory LoadingSession.fromJson(Map<String, dynamic> json) =>
      _$LoadingSessionFromJson(json);
}

@freezed
abstract class ItemValidation with _$ItemValidation {
  const factory ItemValidation({
    required String itemId,
    required String itemLabel,
    required int stepNumber,
    String? scannedBarcode,
    String? expectedBarcode,
    required String status,
    required DateTime validatedAt,
    String? notes,
  }) = _ItemValidation;

  factory ItemValidation.fromJson(Map<String, dynamic> json) =>
      _$ItemValidationFromJson(json);
}
