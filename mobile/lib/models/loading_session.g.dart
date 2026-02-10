// GENERATED CODE - DO NOT MODIFY BY HAND

part of 'loading_session.dart';

// **************************************************************************
// JsonSerializableGenerator
// **************************************************************************

_LoadingSession _$LoadingSessionFromJson(Map<String, dynamic> json) =>
    _LoadingSession(
      sessionId: json['sessionId'] as String,
      planId: json['planId'] as String,
      startedAt: DateTime.parse(json['startedAt'] as String),
      completedAt: json['completedAt'] == null
          ? null
          : DateTime.parse(json['completedAt'] as String),
      totalItems: (json['totalItems'] as num).toInt(),
      validatedCount: (json['validatedCount'] as num?)?.toInt() ?? 0,
      skippedCount: (json['skippedCount'] as num?)?.toInt() ?? 0,
      currentStepIndex: (json['currentStepIndex'] as num?)?.toInt() ?? 0,
      validations:
          (json['validations'] as List<dynamic>?)
              ?.map((e) => ItemValidation.fromJson(e as Map<String, dynamic>))
              .toList() ??
          const [],
      status: json['status'] as String? ?? 'IN_PROGRESS',
    );

Map<String, dynamic> _$LoadingSessionToJson(_LoadingSession instance) =>
    <String, dynamic>{
      'sessionId': instance.sessionId,
      'planId': instance.planId,
      'startedAt': instance.startedAt.toIso8601String(),
      'completedAt': instance.completedAt?.toIso8601String(),
      'totalItems': instance.totalItems,
      'validatedCount': instance.validatedCount,
      'skippedCount': instance.skippedCount,
      'currentStepIndex': instance.currentStepIndex,
      'validations': instance.validations,
      'status': instance.status,
    };

_ItemValidation _$ItemValidationFromJson(Map<String, dynamic> json) =>
    _ItemValidation(
      itemId: json['itemId'] as String,
      itemLabel: json['itemLabel'] as String,
      stepNumber: (json['stepNumber'] as num).toInt(),
      scannedBarcode: json['scannedBarcode'] as String?,
      expectedBarcode: json['expectedBarcode'] as String?,
      status: json['status'] as String,
      validatedAt: DateTime.parse(json['validatedAt'] as String),
      notes: json['notes'] as String?,
    );

Map<String, dynamic> _$ItemValidationToJson(_ItemValidation instance) =>
    <String, dynamic>{
      'itemId': instance.itemId,
      'itemLabel': instance.itemLabel,
      'stepNumber': instance.stepNumber,
      'scannedBarcode': instance.scannedBarcode,
      'expectedBarcode': instance.expectedBarcode,
      'status': instance.status,
      'validatedAt': instance.validatedAt.toIso8601String(),
      'notes': instance.notes,
    };
