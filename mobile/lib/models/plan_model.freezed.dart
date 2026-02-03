// GENERATED CODE - DO NOT MODIFY BY HAND
// coverage:ignore-file
// ignore_for_file: type=lint
// ignore_for_file: unused_element, deprecated_member_use, deprecated_member_use_from_same_package, use_function_type_syntax_for_parameters, unnecessary_const, avoid_init_to_null, invalid_override_different_default_values_named, prefer_expression_function_bodies, annotate_overrides, invalid_annotation_target, unnecessary_question_mark

part of 'plan_model.dart';

// **************************************************************************
// FreezedGenerator
// **************************************************************************

// dart format off
T _$identity<T>(T value) => value;
/// @nodoc
mixin _$PlanModel {

 String get id; String get code; String get title; String get status; int get totalItems; double get totalWeightKg; double? get volumeUtilizationPct; String get createdBy; DateTime get createdAt;
/// Create a copy of PlanModel
/// with the given fields replaced by the non-null parameter values.
@JsonKey(includeFromJson: false, includeToJson: false)
@pragma('vm:prefer-inline')
$PlanModelCopyWith<PlanModel> get copyWith => _$PlanModelCopyWithImpl<PlanModel>(this as PlanModel, _$identity);



@override
bool operator ==(Object other) {
  return identical(this, other) || (other.runtimeType == runtimeType&&other is PlanModel&&(identical(other.id, id) || other.id == id)&&(identical(other.code, code) || other.code == code)&&(identical(other.title, title) || other.title == title)&&(identical(other.status, status) || other.status == status)&&(identical(other.totalItems, totalItems) || other.totalItems == totalItems)&&(identical(other.totalWeightKg, totalWeightKg) || other.totalWeightKg == totalWeightKg)&&(identical(other.volumeUtilizationPct, volumeUtilizationPct) || other.volumeUtilizationPct == volumeUtilizationPct)&&(identical(other.createdBy, createdBy) || other.createdBy == createdBy)&&(identical(other.createdAt, createdAt) || other.createdAt == createdAt));
}


@override
int get hashCode => Object.hash(runtimeType,id,code,title,status,totalItems,totalWeightKg,volumeUtilizationPct,createdBy,createdAt);

@override
String toString() {
  return 'PlanModel(id: $id, code: $code, title: $title, status: $status, totalItems: $totalItems, totalWeightKg: $totalWeightKg, volumeUtilizationPct: $volumeUtilizationPct, createdBy: $createdBy, createdAt: $createdAt)';
}


}

/// @nodoc
abstract mixin class $PlanModelCopyWith<$Res>  {
  factory $PlanModelCopyWith(PlanModel value, $Res Function(PlanModel) _then) = _$PlanModelCopyWithImpl;
@useResult
$Res call({
 String id, String code, String title, String status, int totalItems, double totalWeightKg, double? volumeUtilizationPct, String createdBy, DateTime createdAt
});




}
/// @nodoc
class _$PlanModelCopyWithImpl<$Res>
    implements $PlanModelCopyWith<$Res> {
  _$PlanModelCopyWithImpl(this._self, this._then);

  final PlanModel _self;
  final $Res Function(PlanModel) _then;

/// Create a copy of PlanModel
/// with the given fields replaced by the non-null parameter values.
@pragma('vm:prefer-inline') @override $Res call({Object? id = null,Object? code = null,Object? title = null,Object? status = null,Object? totalItems = null,Object? totalWeightKg = null,Object? volumeUtilizationPct = freezed,Object? createdBy = null,Object? createdAt = null,}) {
  return _then(_self.copyWith(
id: null == id ? _self.id : id // ignore: cast_nullable_to_non_nullable
as String,code: null == code ? _self.code : code // ignore: cast_nullable_to_non_nullable
as String,title: null == title ? _self.title : title // ignore: cast_nullable_to_non_nullable
as String,status: null == status ? _self.status : status // ignore: cast_nullable_to_non_nullable
as String,totalItems: null == totalItems ? _self.totalItems : totalItems // ignore: cast_nullable_to_non_nullable
as int,totalWeightKg: null == totalWeightKg ? _self.totalWeightKg : totalWeightKg // ignore: cast_nullable_to_non_nullable
as double,volumeUtilizationPct: freezed == volumeUtilizationPct ? _self.volumeUtilizationPct : volumeUtilizationPct // ignore: cast_nullable_to_non_nullable
as double?,createdBy: null == createdBy ? _self.createdBy : createdBy // ignore: cast_nullable_to_non_nullable
as String,createdAt: null == createdAt ? _self.createdAt : createdAt // ignore: cast_nullable_to_non_nullable
as DateTime,
  ));
}

}


/// Adds pattern-matching-related methods to [PlanModel].
extension PlanModelPatterns on PlanModel {
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

@optionalTypeArgs TResult maybeMap<TResult extends Object?>(TResult Function( _PlanModel value)?  $default,{required TResult orElse(),}){
final _that = this;
switch (_that) {
case _PlanModel() when $default != null:
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

@optionalTypeArgs TResult map<TResult extends Object?>(TResult Function( _PlanModel value)  $default,){
final _that = this;
switch (_that) {
case _PlanModel():
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

@optionalTypeArgs TResult? mapOrNull<TResult extends Object?>(TResult? Function( _PlanModel value)?  $default,){
final _that = this;
switch (_that) {
case _PlanModel() when $default != null:
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

@optionalTypeArgs TResult maybeWhen<TResult extends Object?>(TResult Function( String id,  String code,  String title,  String status,  int totalItems,  double totalWeightKg,  double? volumeUtilizationPct,  String createdBy,  DateTime createdAt)?  $default,{required TResult orElse(),}) {final _that = this;
switch (_that) {
case _PlanModel() when $default != null:
return $default(_that.id,_that.code,_that.title,_that.status,_that.totalItems,_that.totalWeightKg,_that.volumeUtilizationPct,_that.createdBy,_that.createdAt);case _:
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

@optionalTypeArgs TResult when<TResult extends Object?>(TResult Function( String id,  String code,  String title,  String status,  int totalItems,  double totalWeightKg,  double? volumeUtilizationPct,  String createdBy,  DateTime createdAt)  $default,) {final _that = this;
switch (_that) {
case _PlanModel():
return $default(_that.id,_that.code,_that.title,_that.status,_that.totalItems,_that.totalWeightKg,_that.volumeUtilizationPct,_that.createdBy,_that.createdAt);case _:
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

@optionalTypeArgs TResult? whenOrNull<TResult extends Object?>(TResult? Function( String id,  String code,  String title,  String status,  int totalItems,  double totalWeightKg,  double? volumeUtilizationPct,  String createdBy,  DateTime createdAt)?  $default,) {final _that = this;
switch (_that) {
case _PlanModel() when $default != null:
return $default(_that.id,_that.code,_that.title,_that.status,_that.totalItems,_that.totalWeightKg,_that.volumeUtilizationPct,_that.createdBy,_that.createdAt);case _:
  return null;

}
}

}

/// @nodoc


class _PlanModel implements PlanModel {
  const _PlanModel({required this.id, required this.code, required this.title, required this.status, required this.totalItems, required this.totalWeightKg, this.volumeUtilizationPct, required this.createdBy, required this.createdAt});
  

@override final  String id;
@override final  String code;
@override final  String title;
@override final  String status;
@override final  int totalItems;
@override final  double totalWeightKg;
@override final  double? volumeUtilizationPct;
@override final  String createdBy;
@override final  DateTime createdAt;

/// Create a copy of PlanModel
/// with the given fields replaced by the non-null parameter values.
@override @JsonKey(includeFromJson: false, includeToJson: false)
@pragma('vm:prefer-inline')
_$PlanModelCopyWith<_PlanModel> get copyWith => __$PlanModelCopyWithImpl<_PlanModel>(this, _$identity);



@override
bool operator ==(Object other) {
  return identical(this, other) || (other.runtimeType == runtimeType&&other is _PlanModel&&(identical(other.id, id) || other.id == id)&&(identical(other.code, code) || other.code == code)&&(identical(other.title, title) || other.title == title)&&(identical(other.status, status) || other.status == status)&&(identical(other.totalItems, totalItems) || other.totalItems == totalItems)&&(identical(other.totalWeightKg, totalWeightKg) || other.totalWeightKg == totalWeightKg)&&(identical(other.volumeUtilizationPct, volumeUtilizationPct) || other.volumeUtilizationPct == volumeUtilizationPct)&&(identical(other.createdBy, createdBy) || other.createdBy == createdBy)&&(identical(other.createdAt, createdAt) || other.createdAt == createdAt));
}


@override
int get hashCode => Object.hash(runtimeType,id,code,title,status,totalItems,totalWeightKg,volumeUtilizationPct,createdBy,createdAt);

@override
String toString() {
  return 'PlanModel(id: $id, code: $code, title: $title, status: $status, totalItems: $totalItems, totalWeightKg: $totalWeightKg, volumeUtilizationPct: $volumeUtilizationPct, createdBy: $createdBy, createdAt: $createdAt)';
}


}

/// @nodoc
abstract mixin class _$PlanModelCopyWith<$Res> implements $PlanModelCopyWith<$Res> {
  factory _$PlanModelCopyWith(_PlanModel value, $Res Function(_PlanModel) _then) = __$PlanModelCopyWithImpl;
@override @useResult
$Res call({
 String id, String code, String title, String status, int totalItems, double totalWeightKg, double? volumeUtilizationPct, String createdBy, DateTime createdAt
});




}
/// @nodoc
class __$PlanModelCopyWithImpl<$Res>
    implements _$PlanModelCopyWith<$Res> {
  __$PlanModelCopyWithImpl(this._self, this._then);

  final _PlanModel _self;
  final $Res Function(_PlanModel) _then;

/// Create a copy of PlanModel
/// with the given fields replaced by the non-null parameter values.
@override @pragma('vm:prefer-inline') $Res call({Object? id = null,Object? code = null,Object? title = null,Object? status = null,Object? totalItems = null,Object? totalWeightKg = null,Object? volumeUtilizationPct = freezed,Object? createdBy = null,Object? createdAt = null,}) {
  return _then(_PlanModel(
id: null == id ? _self.id : id // ignore: cast_nullable_to_non_nullable
as String,code: null == code ? _self.code : code // ignore: cast_nullable_to_non_nullable
as String,title: null == title ? _self.title : title // ignore: cast_nullable_to_non_nullable
as String,status: null == status ? _self.status : status // ignore: cast_nullable_to_non_nullable
as String,totalItems: null == totalItems ? _self.totalItems : totalItems // ignore: cast_nullable_to_non_nullable
as int,totalWeightKg: null == totalWeightKg ? _self.totalWeightKg : totalWeightKg // ignore: cast_nullable_to_non_nullable
as double,volumeUtilizationPct: freezed == volumeUtilizationPct ? _self.volumeUtilizationPct : volumeUtilizationPct // ignore: cast_nullable_to_non_nullable
as double?,createdBy: null == createdBy ? _self.createdBy : createdBy // ignore: cast_nullable_to_non_nullable
as String,createdAt: null == createdAt ? _self.createdAt : createdAt // ignore: cast_nullable_to_non_nullable
as DateTime,
  ));
}


}

// dart format on
