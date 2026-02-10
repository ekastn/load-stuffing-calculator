// GENERATED CODE - DO NOT MODIFY BY HAND
// coverage:ignore-file
// ignore_for_file: type=lint
// ignore_for_file: unused_element, deprecated_member_use, deprecated_member_use_from_same_package, use_function_type_syntax_for_parameters, unnecessary_const, avoid_init_to_null, invalid_override_different_default_values_named, prefer_expression_function_bodies, annotate_overrides, invalid_annotation_target, unnecessary_question_mark

part of 'loading_session.dart';

// **************************************************************************
// FreezedGenerator
// **************************************************************************

// dart format off
T _$identity<T>(T value) => value;

/// @nodoc
mixin _$LoadingSession {

 String get sessionId; String get planId; DateTime get startedAt; DateTime? get completedAt; int get totalItems; int get validatedCount; int get skippedCount; int get currentStepIndex; List<ItemValidation> get validations; String get status;
/// Create a copy of LoadingSession
/// with the given fields replaced by the non-null parameter values.
@JsonKey(includeFromJson: false, includeToJson: false)
@pragma('vm:prefer-inline')
$LoadingSessionCopyWith<LoadingSession> get copyWith => _$LoadingSessionCopyWithImpl<LoadingSession>(this as LoadingSession, _$identity);

  /// Serializes this LoadingSession to a JSON map.
  Map<String, dynamic> toJson();


@override
bool operator ==(Object other) {
  return identical(this, other) || (other.runtimeType == runtimeType&&other is LoadingSession&&(identical(other.sessionId, sessionId) || other.sessionId == sessionId)&&(identical(other.planId, planId) || other.planId == planId)&&(identical(other.startedAt, startedAt) || other.startedAt == startedAt)&&(identical(other.completedAt, completedAt) || other.completedAt == completedAt)&&(identical(other.totalItems, totalItems) || other.totalItems == totalItems)&&(identical(other.validatedCount, validatedCount) || other.validatedCount == validatedCount)&&(identical(other.skippedCount, skippedCount) || other.skippedCount == skippedCount)&&(identical(other.currentStepIndex, currentStepIndex) || other.currentStepIndex == currentStepIndex)&&const DeepCollectionEquality().equals(other.validations, validations)&&(identical(other.status, status) || other.status == status));
}

@JsonKey(includeFromJson: false, includeToJson: false)
@override
int get hashCode => Object.hash(runtimeType,sessionId,planId,startedAt,completedAt,totalItems,validatedCount,skippedCount,currentStepIndex,const DeepCollectionEquality().hash(validations),status);

@override
String toString() {
  return 'LoadingSession(sessionId: $sessionId, planId: $planId, startedAt: $startedAt, completedAt: $completedAt, totalItems: $totalItems, validatedCount: $validatedCount, skippedCount: $skippedCount, currentStepIndex: $currentStepIndex, validations: $validations, status: $status)';
}


}

/// @nodoc
abstract mixin class $LoadingSessionCopyWith<$Res>  {
  factory $LoadingSessionCopyWith(LoadingSession value, $Res Function(LoadingSession) _then) = _$LoadingSessionCopyWithImpl;
@useResult
$Res call({
 String sessionId, String planId, DateTime startedAt, DateTime? completedAt, int totalItems, int validatedCount, int skippedCount, int currentStepIndex, List<ItemValidation> validations, String status
});




}
/// @nodoc
class _$LoadingSessionCopyWithImpl<$Res>
    implements $LoadingSessionCopyWith<$Res> {
  _$LoadingSessionCopyWithImpl(this._self, this._then);

  final LoadingSession _self;
  final $Res Function(LoadingSession) _then;

/// Create a copy of LoadingSession
/// with the given fields replaced by the non-null parameter values.
@pragma('vm:prefer-inline') @override $Res call({Object? sessionId = null,Object? planId = null,Object? startedAt = null,Object? completedAt = freezed,Object? totalItems = null,Object? validatedCount = null,Object? skippedCount = null,Object? currentStepIndex = null,Object? validations = null,Object? status = null,}) {
  return _then(_self.copyWith(
sessionId: null == sessionId ? _self.sessionId : sessionId // ignore: cast_nullable_to_non_nullable
as String,planId: null == planId ? _self.planId : planId // ignore: cast_nullable_to_non_nullable
as String,startedAt: null == startedAt ? _self.startedAt : startedAt // ignore: cast_nullable_to_non_nullable
as DateTime,completedAt: freezed == completedAt ? _self.completedAt : completedAt // ignore: cast_nullable_to_non_nullable
as DateTime?,totalItems: null == totalItems ? _self.totalItems : totalItems // ignore: cast_nullable_to_non_nullable
as int,validatedCount: null == validatedCount ? _self.validatedCount : validatedCount // ignore: cast_nullable_to_non_nullable
as int,skippedCount: null == skippedCount ? _self.skippedCount : skippedCount // ignore: cast_nullable_to_non_nullable
as int,currentStepIndex: null == currentStepIndex ? _self.currentStepIndex : currentStepIndex // ignore: cast_nullable_to_non_nullable
as int,validations: null == validations ? _self.validations : validations // ignore: cast_nullable_to_non_nullable
as List<ItemValidation>,status: null == status ? _self.status : status // ignore: cast_nullable_to_non_nullable
as String,
  ));
}

}


/// Adds pattern-matching-related methods to [LoadingSession].
extension LoadingSessionPatterns on LoadingSession {
/// A variant of `map` that fallback to returning `orElse`.
///
/// It is equivalent to doing:
/// ```dart
/// switch (sealedClass) {
///   case final Subclass value:
///     return ...;
///   case _:
///     return orElse();
/// }
/// ```

@optionalTypeArgs TResult maybeMap<TResult extends Object?>(TResult Function( _LoadingSession value)?  $default,{required TResult orElse(),}){
final _that = this;
switch (_that) {
case _LoadingSession() when $default != null:
return $default(_that);case _:
  return orElse();

}
}
/// A `switch`-like method, using callbacks.
///
/// Callbacks receives the raw object, upcasted.
/// It is equivalent to doing:
/// ```dart
/// switch (sealedClass) {
///   case final Subclass value:
///     return ...;
///   case final Subclass2 value:
///     return ...;
/// }
/// ```

@optionalTypeArgs TResult map<TResult extends Object?>(TResult Function( _LoadingSession value)  $default,){
final _that = this;
switch (_that) {
case _LoadingSession():
return $default(_that);case _:
  throw StateError('Unexpected subclass');

}
}
/// A variant of `map` that fallback to returning `null`.
///
/// It is equivalent to doing:
/// ```dart
/// switch (sealedClass) {
///   case final Subclass value:
///     return ...;
///   case _:
///     return null;
/// }
/// ```

@optionalTypeArgs TResult? mapOrNull<TResult extends Object?>(TResult? Function( _LoadingSession value)?  $default,){
final _that = this;
switch (_that) {
case _LoadingSession() when $default != null:
return $default(_that);case _:
  return null;

}
}
/// A variant of `when` that fallback to an `orElse` callback.
///
/// It is equivalent to doing:
/// ```dart
/// switch (sealedClass) {
///   case Subclass(:final field):
///     return ...;
///   case _:
///     return orElse();
/// }
/// ```

@optionalTypeArgs TResult maybeWhen<TResult extends Object?>(TResult Function( String sessionId,  String planId,  DateTime startedAt,  DateTime? completedAt,  int totalItems,  int validatedCount,  int skippedCount,  int currentStepIndex,  List<ItemValidation> validations,  String status)?  $default,{required TResult orElse(),}) {final _that = this;
switch (_that) {
case _LoadingSession() when $default != null:
return $default(_that.sessionId,_that.planId,_that.startedAt,_that.completedAt,_that.totalItems,_that.validatedCount,_that.skippedCount,_that.currentStepIndex,_that.validations,_that.status);case _:
  return orElse();

}
}
/// A `switch`-like method, using callbacks.
///
/// As opposed to `map`, this offers destructuring.
/// It is equivalent to doing:
/// ```dart
/// switch (sealedClass) {
///   case Subclass(:final field):
///     return ...;
///   case Subclass2(:final field2):
///     return ...;
/// }
/// ```

@optionalTypeArgs TResult when<TResult extends Object?>(TResult Function( String sessionId,  String planId,  DateTime startedAt,  DateTime? completedAt,  int totalItems,  int validatedCount,  int skippedCount,  int currentStepIndex,  List<ItemValidation> validations,  String status)  $default,) {final _that = this;
switch (_that) {
case _LoadingSession():
return $default(_that.sessionId,_that.planId,_that.startedAt,_that.completedAt,_that.totalItems,_that.validatedCount,_that.skippedCount,_that.currentStepIndex,_that.validations,_that.status);case _:
  throw StateError('Unexpected subclass');

}
}
/// A variant of `when` that fallback to returning `null`
///
/// It is equivalent to doing:
/// ```dart
/// switch (sealedClass) {
///   case Subclass(:final field):
///     return ...;
///   case _:
///     return null;
/// }
/// ```

@optionalTypeArgs TResult? whenOrNull<TResult extends Object?>(TResult? Function( String sessionId,  String planId,  DateTime startedAt,  DateTime? completedAt,  int totalItems,  int validatedCount,  int skippedCount,  int currentStepIndex,  List<ItemValidation> validations,  String status)?  $default,) {final _that = this;
switch (_that) {
case _LoadingSession() when $default != null:
return $default(_that.sessionId,_that.planId,_that.startedAt,_that.completedAt,_that.totalItems,_that.validatedCount,_that.skippedCount,_that.currentStepIndex,_that.validations,_that.status);case _:
  return null;

}
}

}

/// @nodoc
@JsonSerializable()

class _LoadingSession implements LoadingSession {
  const _LoadingSession({required this.sessionId, required this.planId, required this.startedAt, this.completedAt, required this.totalItems, this.validatedCount = 0, this.skippedCount = 0, this.currentStepIndex = 0, final  List<ItemValidation> validations = const [], this.status = 'IN_PROGRESS'}): _validations = validations;
  factory _LoadingSession.fromJson(Map<String, dynamic> json) => _$LoadingSessionFromJson(json);

@override final  String sessionId;
@override final  String planId;
@override final  DateTime startedAt;
@override final  DateTime? completedAt;
@override final  int totalItems;
@override@JsonKey() final  int validatedCount;
@override@JsonKey() final  int skippedCount;
@override@JsonKey() final  int currentStepIndex;
 final  List<ItemValidation> _validations;
@override@JsonKey() List<ItemValidation> get validations {
  if (_validations is EqualUnmodifiableListView) return _validations;
  // ignore: implicit_dynamic_type
  return EqualUnmodifiableListView(_validations);
}

@override@JsonKey() final  String status;

/// Create a copy of LoadingSession
/// with the given fields replaced by the non-null parameter values.
@override @JsonKey(includeFromJson: false, includeToJson: false)
@pragma('vm:prefer-inline')
_$LoadingSessionCopyWith<_LoadingSession> get copyWith => __$LoadingSessionCopyWithImpl<_LoadingSession>(this, _$identity);

@override
Map<String, dynamic> toJson() {
  return _$LoadingSessionToJson(this, );
}

@override
bool operator ==(Object other) {
  return identical(this, other) || (other.runtimeType == runtimeType&&other is _LoadingSession&&(identical(other.sessionId, sessionId) || other.sessionId == sessionId)&&(identical(other.planId, planId) || other.planId == planId)&&(identical(other.startedAt, startedAt) || other.startedAt == startedAt)&&(identical(other.completedAt, completedAt) || other.completedAt == completedAt)&&(identical(other.totalItems, totalItems) || other.totalItems == totalItems)&&(identical(other.validatedCount, validatedCount) || other.validatedCount == validatedCount)&&(identical(other.skippedCount, skippedCount) || other.skippedCount == skippedCount)&&(identical(other.currentStepIndex, currentStepIndex) || other.currentStepIndex == currentStepIndex)&&const DeepCollectionEquality().equals(other._validations, _validations)&&(identical(other.status, status) || other.status == status));
}

@JsonKey(includeFromJson: false, includeToJson: false)
@override
int get hashCode => Object.hash(runtimeType,sessionId,planId,startedAt,completedAt,totalItems,validatedCount,skippedCount,currentStepIndex,const DeepCollectionEquality().hash(_validations),status);

@override
String toString() {
  return 'LoadingSession(sessionId: $sessionId, planId: $planId, startedAt: $startedAt, completedAt: $completedAt, totalItems: $totalItems, validatedCount: $validatedCount, skippedCount: $skippedCount, currentStepIndex: $currentStepIndex, validations: $validations, status: $status)';
}


}

/// @nodoc
abstract mixin class _$LoadingSessionCopyWith<$Res> implements $LoadingSessionCopyWith<$Res> {
  factory _$LoadingSessionCopyWith(_LoadingSession value, $Res Function(_LoadingSession) _then) = __$LoadingSessionCopyWithImpl;
@override @useResult
$Res call({
 String sessionId, String planId, DateTime startedAt, DateTime? completedAt, int totalItems, int validatedCount, int skippedCount, int currentStepIndex, List<ItemValidation> validations, String status
});




}
/// @nodoc
class __$LoadingSessionCopyWithImpl<$Res>
    implements _$LoadingSessionCopyWith<$Res> {
  __$LoadingSessionCopyWithImpl(this._self, this._then);

  final _LoadingSession _self;
  final $Res Function(_LoadingSession) _then;

/// Create a copy of LoadingSession
/// with the given fields replaced by the non-null parameter values.
@override @pragma('vm:prefer-inline') $Res call({Object? sessionId = null,Object? planId = null,Object? startedAt = null,Object? completedAt = freezed,Object? totalItems = null,Object? validatedCount = null,Object? skippedCount = null,Object? currentStepIndex = null,Object? validations = null,Object? status = null,}) {
  return _then(_LoadingSession(
sessionId: null == sessionId ? _self.sessionId : sessionId // ignore: cast_nullable_to_non_nullable
as String,planId: null == planId ? _self.planId : planId // ignore: cast_nullable_to_non_nullable
as String,startedAt: null == startedAt ? _self.startedAt : startedAt // ignore: cast_nullable_to_non_nullable
as DateTime,completedAt: freezed == completedAt ? _self.completedAt : completedAt // ignore: cast_nullable_to_non_nullable
as DateTime?,totalItems: null == totalItems ? _self.totalItems : totalItems // ignore: cast_nullable_to_non_nullable
as int,validatedCount: null == validatedCount ? _self.validatedCount : validatedCount // ignore: cast_nullable_to_non_nullable
as int,skippedCount: null == skippedCount ? _self.skippedCount : skippedCount // ignore: cast_nullable_to_non_nullable
as int,currentStepIndex: null == currentStepIndex ? _self.currentStepIndex : currentStepIndex // ignore: cast_nullable_to_non_nullable
as int,validations: null == validations ? _self._validations : validations // ignore: cast_nullable_to_non_nullable
as List<ItemValidation>,status: null == status ? _self.status : status // ignore: cast_nullable_to_non_nullable
as String,
  ));
}


}


/// @nodoc
mixin _$ItemValidation {

 String get itemId; String get itemLabel; int get stepNumber; String? get scannedBarcode; String? get expectedBarcode; String get status; DateTime get validatedAt; String? get notes;
/// Create a copy of ItemValidation
/// with the given fields replaced by the non-null parameter values.
@JsonKey(includeFromJson: false, includeToJson: false)
@pragma('vm:prefer-inline')
$ItemValidationCopyWith<ItemValidation> get copyWith => _$ItemValidationCopyWithImpl<ItemValidation>(this as ItemValidation, _$identity);

  /// Serializes this ItemValidation to a JSON map.
  Map<String, dynamic> toJson();


@override
bool operator ==(Object other) {
  return identical(this, other) || (other.runtimeType == runtimeType&&other is ItemValidation&&(identical(other.itemId, itemId) || other.itemId == itemId)&&(identical(other.itemLabel, itemLabel) || other.itemLabel == itemLabel)&&(identical(other.stepNumber, stepNumber) || other.stepNumber == stepNumber)&&(identical(other.scannedBarcode, scannedBarcode) || other.scannedBarcode == scannedBarcode)&&(identical(other.expectedBarcode, expectedBarcode) || other.expectedBarcode == expectedBarcode)&&(identical(other.status, status) || other.status == status)&&(identical(other.validatedAt, validatedAt) || other.validatedAt == validatedAt)&&(identical(other.notes, notes) || other.notes == notes));
}

@JsonKey(includeFromJson: false, includeToJson: false)
@override
int get hashCode => Object.hash(runtimeType,itemId,itemLabel,stepNumber,scannedBarcode,expectedBarcode,status,validatedAt,notes);

@override
String toString() {
  return 'ItemValidation(itemId: $itemId, itemLabel: $itemLabel, stepNumber: $stepNumber, scannedBarcode: $scannedBarcode, expectedBarcode: $expectedBarcode, status: $status, validatedAt: $validatedAt, notes: $notes)';
}


}

/// @nodoc
abstract mixin class $ItemValidationCopyWith<$Res>  {
  factory $ItemValidationCopyWith(ItemValidation value, $Res Function(ItemValidation) _then) = _$ItemValidationCopyWithImpl;
@useResult
$Res call({
 String itemId, String itemLabel, int stepNumber, String? scannedBarcode, String? expectedBarcode, String status, DateTime validatedAt, String? notes
});




}
/// @nodoc
class _$ItemValidationCopyWithImpl<$Res>
    implements $ItemValidationCopyWith<$Res> {
  _$ItemValidationCopyWithImpl(this._self, this._then);

  final ItemValidation _self;
  final $Res Function(ItemValidation) _then;

/// Create a copy of ItemValidation
/// with the given fields replaced by the non-null parameter values.
@pragma('vm:prefer-inline') @override $Res call({Object? itemId = null,Object? itemLabel = null,Object? stepNumber = null,Object? scannedBarcode = freezed,Object? expectedBarcode = freezed,Object? status = null,Object? validatedAt = null,Object? notes = freezed,}) {
  return _then(_self.copyWith(
itemId: null == itemId ? _self.itemId : itemId // ignore: cast_nullable_to_non_nullable
as String,itemLabel: null == itemLabel ? _self.itemLabel : itemLabel // ignore: cast_nullable_to_non_nullable
as String,stepNumber: null == stepNumber ? _self.stepNumber : stepNumber // ignore: cast_nullable_to_non_nullable
as int,scannedBarcode: freezed == scannedBarcode ? _self.scannedBarcode : scannedBarcode // ignore: cast_nullable_to_non_nullable
as String?,expectedBarcode: freezed == expectedBarcode ? _self.expectedBarcode : expectedBarcode // ignore: cast_nullable_to_non_nullable
as String?,status: null == status ? _self.status : status // ignore: cast_nullable_to_non_nullable
as String,validatedAt: null == validatedAt ? _self.validatedAt : validatedAt // ignore: cast_nullable_to_non_nullable
as DateTime,notes: freezed == notes ? _self.notes : notes // ignore: cast_nullable_to_non_nullable
as String?,
  ));
}

}


/// Adds pattern-matching-related methods to [ItemValidation].
extension ItemValidationPatterns on ItemValidation {
/// A variant of `map` that fallback to returning `orElse`.
///
/// It is equivalent to doing:
/// ```dart
/// switch (sealedClass) {
///   case final Subclass value:
///     return ...;
///   case _:
///     return orElse();
/// }
/// ```

@optionalTypeArgs TResult maybeMap<TResult extends Object?>(TResult Function( _ItemValidation value)?  $default,{required TResult orElse(),}){
final _that = this;
switch (_that) {
case _ItemValidation() when $default != null:
return $default(_that);case _:
  return orElse();

}
}
/// A `switch`-like method, using callbacks.
///
/// Callbacks receives the raw object, upcasted.
/// It is equivalent to doing:
/// ```dart
/// switch (sealedClass) {
///   case final Subclass value:
///     return ...;
///   case final Subclass2 value:
///     return ...;
/// }
/// ```

@optionalTypeArgs TResult map<TResult extends Object?>(TResult Function( _ItemValidation value)  $default,){
final _that = this;
switch (_that) {
case _ItemValidation():
return $default(_that);case _:
  throw StateError('Unexpected subclass');

}
}
/// A variant of `map` that fallback to returning `null`.
///
/// It is equivalent to doing:
/// ```dart
/// switch (sealedClass) {
///   case final Subclass value:
///     return ...;
///   case _:
///     return null;
/// }
/// ```

@optionalTypeArgs TResult? mapOrNull<TResult extends Object?>(TResult? Function( _ItemValidation value)?  $default,){
final _that = this;
switch (_that) {
case _ItemValidation() when $default != null:
return $default(_that);case _:
  return null;

}
}
/// A variant of `when` that fallback to an `orElse` callback.
///
/// It is equivalent to doing:
/// ```dart
/// switch (sealedClass) {
///   case Subclass(:final field):
///     return ...;
///   case _:
///     return orElse();
/// }
/// ```

@optionalTypeArgs TResult maybeWhen<TResult extends Object?>(TResult Function( String itemId,  String itemLabel,  int stepNumber,  String? scannedBarcode,  String? expectedBarcode,  String status,  DateTime validatedAt,  String? notes)?  $default,{required TResult orElse(),}) {final _that = this;
switch (_that) {
case _ItemValidation() when $default != null:
return $default(_that.itemId,_that.itemLabel,_that.stepNumber,_that.scannedBarcode,_that.expectedBarcode,_that.status,_that.validatedAt,_that.notes);case _:
  return orElse();

}
}
/// A `switch`-like method, using callbacks.
///
/// As opposed to `map`, this offers destructuring.
/// It is equivalent to doing:
/// ```dart
/// switch (sealedClass) {
///   case Subclass(:final field):
///     return ...;
///   case Subclass2(:final field2):
///     return ...;
/// }
/// ```

@optionalTypeArgs TResult when<TResult extends Object?>(TResult Function( String itemId,  String itemLabel,  int stepNumber,  String? scannedBarcode,  String? expectedBarcode,  String status,  DateTime validatedAt,  String? notes)  $default,) {final _that = this;
switch (_that) {
case _ItemValidation():
return $default(_that.itemId,_that.itemLabel,_that.stepNumber,_that.scannedBarcode,_that.expectedBarcode,_that.status,_that.validatedAt,_that.notes);case _:
  throw StateError('Unexpected subclass');

}
}
/// A variant of `when` that fallback to returning `null`
///
/// It is equivalent to doing:
/// ```dart
/// switch (sealedClass) {
///   case Subclass(:final field):
///     return ...;
///   case _:
///     return null;
/// }
/// ```

@optionalTypeArgs TResult? whenOrNull<TResult extends Object?>(TResult? Function( String itemId,  String itemLabel,  int stepNumber,  String? scannedBarcode,  String? expectedBarcode,  String status,  DateTime validatedAt,  String? notes)?  $default,) {final _that = this;
switch (_that) {
case _ItemValidation() when $default != null:
return $default(_that.itemId,_that.itemLabel,_that.stepNumber,_that.scannedBarcode,_that.expectedBarcode,_that.status,_that.validatedAt,_that.notes);case _:
  return null;

}
}

}

/// @nodoc
@JsonSerializable()

class _ItemValidation implements ItemValidation {
  const _ItemValidation({required this.itemId, required this.itemLabel, required this.stepNumber, this.scannedBarcode, this.expectedBarcode, required this.status, required this.validatedAt, this.notes});
  factory _ItemValidation.fromJson(Map<String, dynamic> json) => _$ItemValidationFromJson(json);

@override final  String itemId;
@override final  String itemLabel;
@override final  int stepNumber;
@override final  String? scannedBarcode;
@override final  String? expectedBarcode;
@override final  String status;
@override final  DateTime validatedAt;
@override final  String? notes;

/// Create a copy of ItemValidation
/// with the given fields replaced by the non-null parameter values.
@override @JsonKey(includeFromJson: false, includeToJson: false)
@pragma('vm:prefer-inline')
_$ItemValidationCopyWith<_ItemValidation> get copyWith => __$ItemValidationCopyWithImpl<_ItemValidation>(this, _$identity);

@override
Map<String, dynamic> toJson() {
  return _$ItemValidationToJson(this, );
}

@override
bool operator ==(Object other) {
  return identical(this, other) || (other.runtimeType == runtimeType&&other is _ItemValidation&&(identical(other.itemId, itemId) || other.itemId == itemId)&&(identical(other.itemLabel, itemLabel) || other.itemLabel == itemLabel)&&(identical(other.stepNumber, stepNumber) || other.stepNumber == stepNumber)&&(identical(other.scannedBarcode, scannedBarcode) || other.scannedBarcode == scannedBarcode)&&(identical(other.expectedBarcode, expectedBarcode) || other.expectedBarcode == expectedBarcode)&&(identical(other.status, status) || other.status == status)&&(identical(other.validatedAt, validatedAt) || other.validatedAt == validatedAt)&&(identical(other.notes, notes) || other.notes == notes));
}

@JsonKey(includeFromJson: false, includeToJson: false)
@override
int get hashCode => Object.hash(runtimeType,itemId,itemLabel,stepNumber,scannedBarcode,expectedBarcode,status,validatedAt,notes);

@override
String toString() {
  return 'ItemValidation(itemId: $itemId, itemLabel: $itemLabel, stepNumber: $stepNumber, scannedBarcode: $scannedBarcode, expectedBarcode: $expectedBarcode, status: $status, validatedAt: $validatedAt, notes: $notes)';
}


}

/// @nodoc
abstract mixin class _$ItemValidationCopyWith<$Res> implements $ItemValidationCopyWith<$Res> {
  factory _$ItemValidationCopyWith(_ItemValidation value, $Res Function(_ItemValidation) _then) = __$ItemValidationCopyWithImpl;
@override @useResult
$Res call({
 String itemId, String itemLabel, int stepNumber, String? scannedBarcode, String? expectedBarcode, String status, DateTime validatedAt, String? notes
});




}
/// @nodoc
class __$ItemValidationCopyWithImpl<$Res>
    implements _$ItemValidationCopyWith<$Res> {
  __$ItemValidationCopyWithImpl(this._self, this._then);

  final _ItemValidation _self;
  final $Res Function(_ItemValidation) _then;

/// Create a copy of ItemValidation
/// with the given fields replaced by the non-null parameter values.
@override @pragma('vm:prefer-inline') $Res call({Object? itemId = null,Object? itemLabel = null,Object? stepNumber = null,Object? scannedBarcode = freezed,Object? expectedBarcode = freezed,Object? status = null,Object? validatedAt = null,Object? notes = freezed,}) {
  return _then(_ItemValidation(
itemId: null == itemId ? _self.itemId : itemId // ignore: cast_nullable_to_non_nullable
as String,itemLabel: null == itemLabel ? _self.itemLabel : itemLabel // ignore: cast_nullable_to_non_nullable
as String,stepNumber: null == stepNumber ? _self.stepNumber : stepNumber // ignore: cast_nullable_to_non_nullable
as int,scannedBarcode: freezed == scannedBarcode ? _self.scannedBarcode : scannedBarcode // ignore: cast_nullable_to_non_nullable
as String?,expectedBarcode: freezed == expectedBarcode ? _self.expectedBarcode : expectedBarcode // ignore: cast_nullable_to_non_nullable
as String?,status: null == status ? _self.status : status // ignore: cast_nullable_to_non_nullable
as String,validatedAt: null == validatedAt ? _self.validatedAt : validatedAt // ignore: cast_nullable_to_non_nullable
as DateTime,notes: freezed == notes ? _self.notes : notes // ignore: cast_nullable_to_non_nullable
as String?,
  ));
}


}

// dart format on
