// GENERATED CODE - DO NOT MODIFY BY HAND
// coverage:ignore-file
// ignore_for_file: type=lint
// ignore_for_file: unused_element, deprecated_member_use, deprecated_member_use_from_same_package, use_function_type_syntax_for_parameters, unnecessary_const, avoid_init_to_null, invalid_override_different_default_values_named, prefer_expression_function_bodies, annotate_overrides, invalid_annotation_target, unnecessary_question_mark

part of 'plan_detail_model.dart';

// **************************************************************************
// FreezedGenerator
// **************************************************************************

// dart format off
T _$identity<T>(T value) => value;
/// @nodoc
mixin _$PlanDetailModel {

 String get id; String get code; String get title; String? get notes; String get status; ContainerInfo get container; PlanStats get stats; List<PlanItem> get items; CalculationResult? get calculation; String get createdBy; DateTime get createdAt; DateTime get updatedAt; DateTime? get completedAt;
/// Create a copy of PlanDetailModel
/// with the given fields replaced by the non-null parameter values.
@JsonKey(includeFromJson: false, includeToJson: false)
@pragma('vm:prefer-inline')
$PlanDetailModelCopyWith<PlanDetailModel> get copyWith => _$PlanDetailModelCopyWithImpl<PlanDetailModel>(this as PlanDetailModel, _$identity);



@override
bool operator ==(Object other) {
  return identical(this, other) || (other.runtimeType == runtimeType&&other is PlanDetailModel&&(identical(other.id, id) || other.id == id)&&(identical(other.code, code) || other.code == code)&&(identical(other.title, title) || other.title == title)&&(identical(other.notes, notes) || other.notes == notes)&&(identical(other.status, status) || other.status == status)&&(identical(other.container, container) || other.container == container)&&(identical(other.stats, stats) || other.stats == stats)&&const DeepCollectionEquality().equals(other.items, items)&&(identical(other.calculation, calculation) || other.calculation == calculation)&&(identical(other.createdBy, createdBy) || other.createdBy == createdBy)&&(identical(other.createdAt, createdAt) || other.createdAt == createdAt)&&(identical(other.updatedAt, updatedAt) || other.updatedAt == updatedAt)&&(identical(other.completedAt, completedAt) || other.completedAt == completedAt));
}


@override
int get hashCode => Object.hash(runtimeType,id,code,title,notes,status,container,stats,const DeepCollectionEquality().hash(items),calculation,createdBy,createdAt,updatedAt,completedAt);

@override
String toString() {
  return 'PlanDetailModel(id: $id, code: $code, title: $title, notes: $notes, status: $status, container: $container, stats: $stats, items: $items, calculation: $calculation, createdBy: $createdBy, createdAt: $createdAt, updatedAt: $updatedAt, completedAt: $completedAt)';
}


}

/// @nodoc
abstract mixin class $PlanDetailModelCopyWith<$Res>  {
  factory $PlanDetailModelCopyWith(PlanDetailModel value, $Res Function(PlanDetailModel) _then) = _$PlanDetailModelCopyWithImpl;
@useResult
$Res call({
 String id, String code, String title, String? notes, String status, ContainerInfo container, PlanStats stats, List<PlanItem> items, CalculationResult? calculation, String createdBy, DateTime createdAt, DateTime updatedAt, DateTime? completedAt
});


$ContainerInfoCopyWith<$Res> get container;$PlanStatsCopyWith<$Res> get stats;$CalculationResultCopyWith<$Res>? get calculation;

}
/// @nodoc
class _$PlanDetailModelCopyWithImpl<$Res>
    implements $PlanDetailModelCopyWith<$Res> {
  _$PlanDetailModelCopyWithImpl(this._self, this._then);

  final PlanDetailModel _self;
  final $Res Function(PlanDetailModel) _then;

/// Create a copy of PlanDetailModel
/// with the given fields replaced by the non-null parameter values.
@pragma('vm:prefer-inline') @override $Res call({Object? id = null,Object? code = null,Object? title = null,Object? notes = freezed,Object? status = null,Object? container = null,Object? stats = null,Object? items = null,Object? calculation = freezed,Object? createdBy = null,Object? createdAt = null,Object? updatedAt = null,Object? completedAt = freezed,}) {
  return _then(_self.copyWith(
id: null == id ? _self.id : id // ignore: cast_nullable_to_non_nullable
as String,code: null == code ? _self.code : code // ignore: cast_nullable_to_non_nullable
as String,title: null == title ? _self.title : title // ignore: cast_nullable_to_non_nullable
as String,notes: freezed == notes ? _self.notes : notes // ignore: cast_nullable_to_non_nullable
as String?,status: null == status ? _self.status : status // ignore: cast_nullable_to_non_nullable
as String,container: null == container ? _self.container : container // ignore: cast_nullable_to_non_nullable
as ContainerInfo,stats: null == stats ? _self.stats : stats // ignore: cast_nullable_to_non_nullable
as PlanStats,items: null == items ? _self.items : items // ignore: cast_nullable_to_non_nullable
as List<PlanItem>,calculation: freezed == calculation ? _self.calculation : calculation // ignore: cast_nullable_to_non_nullable
as CalculationResult?,createdBy: null == createdBy ? _self.createdBy : createdBy // ignore: cast_nullable_to_non_nullable
as String,createdAt: null == createdAt ? _self.createdAt : createdAt // ignore: cast_nullable_to_non_nullable
as DateTime,updatedAt: null == updatedAt ? _self.updatedAt : updatedAt // ignore: cast_nullable_to_non_nullable
as DateTime,completedAt: freezed == completedAt ? _self.completedAt : completedAt // ignore: cast_nullable_to_non_nullable
as DateTime?,
  ));
}
/// Create a copy of PlanDetailModel
/// with the given fields replaced by the non-null parameter values.
@override
@pragma('vm:prefer-inline')
$ContainerInfoCopyWith<$Res> get container {
  
  return $ContainerInfoCopyWith<$Res>(_self.container, (value) {
    return _then(_self.copyWith(container: value));
  });
}/// Create a copy of PlanDetailModel
/// with the given fields replaced by the non-null parameter values.
@override
@pragma('vm:prefer-inline')
$PlanStatsCopyWith<$Res> get stats {
  
  return $PlanStatsCopyWith<$Res>(_self.stats, (value) {
    return _then(_self.copyWith(stats: value));
  });
}/// Create a copy of PlanDetailModel
/// with the given fields replaced by the non-null parameter values.
@override
@pragma('vm:prefer-inline')
$CalculationResultCopyWith<$Res>? get calculation {
    if (_self.calculation == null) {
    return null;
  }

  return $CalculationResultCopyWith<$Res>(_self.calculation!, (value) {
    return _then(_self.copyWith(calculation: value));
  });
}
}


/// Adds pattern-matching-related methods to [PlanDetailModel].
extension PlanDetailModelPatterns on PlanDetailModel {
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

@optionalTypeArgs TResult maybeMap<TResult extends Object?>(TResult Function( _PlanDetailModel value)?  $default,{required TResult orElse(),}){
final _that = this;
switch (_that) {
case _PlanDetailModel() when $default != null:
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

@optionalTypeArgs TResult map<TResult extends Object?>(TResult Function( _PlanDetailModel value)  $default,){
final _that = this;
switch (_that) {
case _PlanDetailModel():
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

@optionalTypeArgs TResult? mapOrNull<TResult extends Object?>(TResult? Function( _PlanDetailModel value)?  $default,){
final _that = this;
switch (_that) {
case _PlanDetailModel() when $default != null:
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

@optionalTypeArgs TResult maybeWhen<TResult extends Object?>(TResult Function( String id,  String code,  String title,  String? notes,  String status,  ContainerInfo container,  PlanStats stats,  List<PlanItem> items,  CalculationResult? calculation,  String createdBy,  DateTime createdAt,  DateTime updatedAt,  DateTime? completedAt)?  $default,{required TResult orElse(),}) {final _that = this;
switch (_that) {
case _PlanDetailModel() when $default != null:
return $default(_that.id,_that.code,_that.title,_that.notes,_that.status,_that.container,_that.stats,_that.items,_that.calculation,_that.createdBy,_that.createdAt,_that.updatedAt,_that.completedAt);case _:
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

@optionalTypeArgs TResult when<TResult extends Object?>(TResult Function( String id,  String code,  String title,  String? notes,  String status,  ContainerInfo container,  PlanStats stats,  List<PlanItem> items,  CalculationResult? calculation,  String createdBy,  DateTime createdAt,  DateTime updatedAt,  DateTime? completedAt)  $default,) {final _that = this;
switch (_that) {
case _PlanDetailModel():
return $default(_that.id,_that.code,_that.title,_that.notes,_that.status,_that.container,_that.stats,_that.items,_that.calculation,_that.createdBy,_that.createdAt,_that.updatedAt,_that.completedAt);case _:
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

@optionalTypeArgs TResult? whenOrNull<TResult extends Object?>(TResult? Function( String id,  String code,  String title,  String? notes,  String status,  ContainerInfo container,  PlanStats stats,  List<PlanItem> items,  CalculationResult? calculation,  String createdBy,  DateTime createdAt,  DateTime updatedAt,  DateTime? completedAt)?  $default,) {final _that = this;
switch (_that) {
case _PlanDetailModel() when $default != null:
return $default(_that.id,_that.code,_that.title,_that.notes,_that.status,_that.container,_that.stats,_that.items,_that.calculation,_that.createdBy,_that.createdAt,_that.updatedAt,_that.completedAt);case _:
  return null;

}
}

}

/// @nodoc


class _PlanDetailModel implements PlanDetailModel {
  const _PlanDetailModel({required this.id, required this.code, required this.title, this.notes, required this.status, required this.container, required this.stats, required final  List<PlanItem> items, this.calculation, required this.createdBy, required this.createdAt, required this.updatedAt, this.completedAt}): _items = items;
  

@override final  String id;
@override final  String code;
@override final  String title;
@override final  String? notes;
@override final  String status;
@override final  ContainerInfo container;
@override final  PlanStats stats;
 final  List<PlanItem> _items;
@override List<PlanItem> get items {
  if (_items is EqualUnmodifiableListView) return _items;
  // ignore: implicit_dynamic_type
  return EqualUnmodifiableListView(_items);
}

@override final  CalculationResult? calculation;
@override final  String createdBy;
@override final  DateTime createdAt;
@override final  DateTime updatedAt;
@override final  DateTime? completedAt;

/// Create a copy of PlanDetailModel
/// with the given fields replaced by the non-null parameter values.
@override @JsonKey(includeFromJson: false, includeToJson: false)
@pragma('vm:prefer-inline')
_$PlanDetailModelCopyWith<_PlanDetailModel> get copyWith => __$PlanDetailModelCopyWithImpl<_PlanDetailModel>(this, _$identity);



@override
bool operator ==(Object other) {
  return identical(this, other) || (other.runtimeType == runtimeType&&other is _PlanDetailModel&&(identical(other.id, id) || other.id == id)&&(identical(other.code, code) || other.code == code)&&(identical(other.title, title) || other.title == title)&&(identical(other.notes, notes) || other.notes == notes)&&(identical(other.status, status) || other.status == status)&&(identical(other.container, container) || other.container == container)&&(identical(other.stats, stats) || other.stats == stats)&&const DeepCollectionEquality().equals(other._items, _items)&&(identical(other.calculation, calculation) || other.calculation == calculation)&&(identical(other.createdBy, createdBy) || other.createdBy == createdBy)&&(identical(other.createdAt, createdAt) || other.createdAt == createdAt)&&(identical(other.updatedAt, updatedAt) || other.updatedAt == updatedAt)&&(identical(other.completedAt, completedAt) || other.completedAt == completedAt));
}


@override
int get hashCode => Object.hash(runtimeType,id,code,title,notes,status,container,stats,const DeepCollectionEquality().hash(_items),calculation,createdBy,createdAt,updatedAt,completedAt);

@override
String toString() {
  return 'PlanDetailModel(id: $id, code: $code, title: $title, notes: $notes, status: $status, container: $container, stats: $stats, items: $items, calculation: $calculation, createdBy: $createdBy, createdAt: $createdAt, updatedAt: $updatedAt, completedAt: $completedAt)';
}


}

/// @nodoc
abstract mixin class _$PlanDetailModelCopyWith<$Res> implements $PlanDetailModelCopyWith<$Res> {
  factory _$PlanDetailModelCopyWith(_PlanDetailModel value, $Res Function(_PlanDetailModel) _then) = __$PlanDetailModelCopyWithImpl;
@override @useResult
$Res call({
 String id, String code, String title, String? notes, String status, ContainerInfo container, PlanStats stats, List<PlanItem> items, CalculationResult? calculation, String createdBy, DateTime createdAt, DateTime updatedAt, DateTime? completedAt
});


@override $ContainerInfoCopyWith<$Res> get container;@override $PlanStatsCopyWith<$Res> get stats;@override $CalculationResultCopyWith<$Res>? get calculation;

}
/// @nodoc
class __$PlanDetailModelCopyWithImpl<$Res>
    implements _$PlanDetailModelCopyWith<$Res> {
  __$PlanDetailModelCopyWithImpl(this._self, this._then);

  final _PlanDetailModel _self;
  final $Res Function(_PlanDetailModel) _then;

/// Create a copy of PlanDetailModel
/// with the given fields replaced by the non-null parameter values.
@override @pragma('vm:prefer-inline') $Res call({Object? id = null,Object? code = null,Object? title = null,Object? notes = freezed,Object? status = null,Object? container = null,Object? stats = null,Object? items = null,Object? calculation = freezed,Object? createdBy = null,Object? createdAt = null,Object? updatedAt = null,Object? completedAt = freezed,}) {
  return _then(_PlanDetailModel(
id: null == id ? _self.id : id // ignore: cast_nullable_to_non_nullable
as String,code: null == code ? _self.code : code // ignore: cast_nullable_to_non_nullable
as String,title: null == title ? _self.title : title // ignore: cast_nullable_to_non_nullable
as String,notes: freezed == notes ? _self.notes : notes // ignore: cast_nullable_to_non_nullable
as String?,status: null == status ? _self.status : status // ignore: cast_nullable_to_non_nullable
as String,container: null == container ? _self.container : container // ignore: cast_nullable_to_non_nullable
as ContainerInfo,stats: null == stats ? _self.stats : stats // ignore: cast_nullable_to_non_nullable
as PlanStats,items: null == items ? _self._items : items // ignore: cast_nullable_to_non_nullable
as List<PlanItem>,calculation: freezed == calculation ? _self.calculation : calculation // ignore: cast_nullable_to_non_nullable
as CalculationResult?,createdBy: null == createdBy ? _self.createdBy : createdBy // ignore: cast_nullable_to_non_nullable
as String,createdAt: null == createdAt ? _self.createdAt : createdAt // ignore: cast_nullable_to_non_nullable
as DateTime,updatedAt: null == updatedAt ? _self.updatedAt : updatedAt // ignore: cast_nullable_to_non_nullable
as DateTime,completedAt: freezed == completedAt ? _self.completedAt : completedAt // ignore: cast_nullable_to_non_nullable
as DateTime?,
  ));
}

/// Create a copy of PlanDetailModel
/// with the given fields replaced by the non-null parameter values.
@override
@pragma('vm:prefer-inline')
$ContainerInfoCopyWith<$Res> get container {
  
  return $ContainerInfoCopyWith<$Res>(_self.container, (value) {
    return _then(_self.copyWith(container: value));
  });
}/// Create a copy of PlanDetailModel
/// with the given fields replaced by the non-null parameter values.
@override
@pragma('vm:prefer-inline')
$PlanStatsCopyWith<$Res> get stats {
  
  return $PlanStatsCopyWith<$Res>(_self.stats, (value) {
    return _then(_self.copyWith(stats: value));
  });
}/// Create a copy of PlanDetailModel
/// with the given fields replaced by the non-null parameter values.
@override
@pragma('vm:prefer-inline')
$CalculationResultCopyWith<$Res>? get calculation {
    if (_self.calculation == null) {
    return null;
  }

  return $CalculationResultCopyWith<$Res>(_self.calculation!, (value) {
    return _then(_self.copyWith(calculation: value));
  });
}
}

/// @nodoc
mixin _$ContainerInfo {

 String? get containerId; String? get name; double get lengthMm; double get widthMm; double get heightMm; double get maxWeightKg; double get volumeM3;
/// Create a copy of ContainerInfo
/// with the given fields replaced by the non-null parameter values.
@JsonKey(includeFromJson: false, includeToJson: false)
@pragma('vm:prefer-inline')
$ContainerInfoCopyWith<ContainerInfo> get copyWith => _$ContainerInfoCopyWithImpl<ContainerInfo>(this as ContainerInfo, _$identity);



@override
bool operator ==(Object other) {
  return identical(this, other) || (other.runtimeType == runtimeType&&other is ContainerInfo&&(identical(other.containerId, containerId) || other.containerId == containerId)&&(identical(other.name, name) || other.name == name)&&(identical(other.lengthMm, lengthMm) || other.lengthMm == lengthMm)&&(identical(other.widthMm, widthMm) || other.widthMm == widthMm)&&(identical(other.heightMm, heightMm) || other.heightMm == heightMm)&&(identical(other.maxWeightKg, maxWeightKg) || other.maxWeightKg == maxWeightKg)&&(identical(other.volumeM3, volumeM3) || other.volumeM3 == volumeM3));
}


@override
int get hashCode => Object.hash(runtimeType,containerId,name,lengthMm,widthMm,heightMm,maxWeightKg,volumeM3);

@override
String toString() {
  return 'ContainerInfo(containerId: $containerId, name: $name, lengthMm: $lengthMm, widthMm: $widthMm, heightMm: $heightMm, maxWeightKg: $maxWeightKg, volumeM3: $volumeM3)';
}


}

/// @nodoc
abstract mixin class $ContainerInfoCopyWith<$Res>  {
  factory $ContainerInfoCopyWith(ContainerInfo value, $Res Function(ContainerInfo) _then) = _$ContainerInfoCopyWithImpl;
@useResult
$Res call({
 String? containerId, String? name, double lengthMm, double widthMm, double heightMm, double maxWeightKg, double volumeM3
});




}
/// @nodoc
class _$ContainerInfoCopyWithImpl<$Res>
    implements $ContainerInfoCopyWith<$Res> {
  _$ContainerInfoCopyWithImpl(this._self, this._then);

  final ContainerInfo _self;
  final $Res Function(ContainerInfo) _then;

/// Create a copy of ContainerInfo
/// with the given fields replaced by the non-null parameter values.
@pragma('vm:prefer-inline') @override $Res call({Object? containerId = freezed,Object? name = freezed,Object? lengthMm = null,Object? widthMm = null,Object? heightMm = null,Object? maxWeightKg = null,Object? volumeM3 = null,}) {
  return _then(_self.copyWith(
containerId: freezed == containerId ? _self.containerId : containerId // ignore: cast_nullable_to_non_nullable
as String?,name: freezed == name ? _self.name : name // ignore: cast_nullable_to_non_nullable
as String?,lengthMm: null == lengthMm ? _self.lengthMm : lengthMm // ignore: cast_nullable_to_non_nullable
as double,widthMm: null == widthMm ? _self.widthMm : widthMm // ignore: cast_nullable_to_non_nullable
as double,heightMm: null == heightMm ? _self.heightMm : heightMm // ignore: cast_nullable_to_non_nullable
as double,maxWeightKg: null == maxWeightKg ? _self.maxWeightKg : maxWeightKg // ignore: cast_nullable_to_non_nullable
as double,volumeM3: null == volumeM3 ? _self.volumeM3 : volumeM3 // ignore: cast_nullable_to_non_nullable
as double,
  ));
}

}


/// Adds pattern-matching-related methods to [ContainerInfo].
extension ContainerInfoPatterns on ContainerInfo {
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

@optionalTypeArgs TResult maybeMap<TResult extends Object?>(TResult Function( _ContainerInfo value)?  $default,{required TResult orElse(),}){
final _that = this;
switch (_that) {
case _ContainerInfo() when $default != null:
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

@optionalTypeArgs TResult map<TResult extends Object?>(TResult Function( _ContainerInfo value)  $default,){
final _that = this;
switch (_that) {
case _ContainerInfo():
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

@optionalTypeArgs TResult? mapOrNull<TResult extends Object?>(TResult? Function( _ContainerInfo value)?  $default,){
final _that = this;
switch (_that) {
case _ContainerInfo() when $default != null:
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

@optionalTypeArgs TResult maybeWhen<TResult extends Object?>(TResult Function( String? containerId,  String? name,  double lengthMm,  double widthMm,  double heightMm,  double maxWeightKg,  double volumeM3)?  $default,{required TResult orElse(),}) {final _that = this;
switch (_that) {
case _ContainerInfo() when $default != null:
return $default(_that.containerId,_that.name,_that.lengthMm,_that.widthMm,_that.heightMm,_that.maxWeightKg,_that.volumeM3);case _:
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

@optionalTypeArgs TResult when<TResult extends Object?>(TResult Function( String? containerId,  String? name,  double lengthMm,  double widthMm,  double heightMm,  double maxWeightKg,  double volumeM3)  $default,) {final _that = this;
switch (_that) {
case _ContainerInfo():
return $default(_that.containerId,_that.name,_that.lengthMm,_that.widthMm,_that.heightMm,_that.maxWeightKg,_that.volumeM3);case _:
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

@optionalTypeArgs TResult? whenOrNull<TResult extends Object?>(TResult? Function( String? containerId,  String? name,  double lengthMm,  double widthMm,  double heightMm,  double maxWeightKg,  double volumeM3)?  $default,) {final _that = this;
switch (_that) {
case _ContainerInfo() when $default != null:
return $default(_that.containerId,_that.name,_that.lengthMm,_that.widthMm,_that.heightMm,_that.maxWeightKg,_that.volumeM3);case _:
  return null;

}
}

}

/// @nodoc


class _ContainerInfo implements ContainerInfo {
  const _ContainerInfo({this.containerId, this.name, required this.lengthMm, required this.widthMm, required this.heightMm, required this.maxWeightKg, required this.volumeM3});
  

@override final  String? containerId;
@override final  String? name;
@override final  double lengthMm;
@override final  double widthMm;
@override final  double heightMm;
@override final  double maxWeightKg;
@override final  double volumeM3;

/// Create a copy of ContainerInfo
/// with the given fields replaced by the non-null parameter values.
@override @JsonKey(includeFromJson: false, includeToJson: false)
@pragma('vm:prefer-inline')
_$ContainerInfoCopyWith<_ContainerInfo> get copyWith => __$ContainerInfoCopyWithImpl<_ContainerInfo>(this, _$identity);



@override
bool operator ==(Object other) {
  return identical(this, other) || (other.runtimeType == runtimeType&&other is _ContainerInfo&&(identical(other.containerId, containerId) || other.containerId == containerId)&&(identical(other.name, name) || other.name == name)&&(identical(other.lengthMm, lengthMm) || other.lengthMm == lengthMm)&&(identical(other.widthMm, widthMm) || other.widthMm == widthMm)&&(identical(other.heightMm, heightMm) || other.heightMm == heightMm)&&(identical(other.maxWeightKg, maxWeightKg) || other.maxWeightKg == maxWeightKg)&&(identical(other.volumeM3, volumeM3) || other.volumeM3 == volumeM3));
}


@override
int get hashCode => Object.hash(runtimeType,containerId,name,lengthMm,widthMm,heightMm,maxWeightKg,volumeM3);

@override
String toString() {
  return 'ContainerInfo(containerId: $containerId, name: $name, lengthMm: $lengthMm, widthMm: $widthMm, heightMm: $heightMm, maxWeightKg: $maxWeightKg, volumeM3: $volumeM3)';
}


}

/// @nodoc
abstract mixin class _$ContainerInfoCopyWith<$Res> implements $ContainerInfoCopyWith<$Res> {
  factory _$ContainerInfoCopyWith(_ContainerInfo value, $Res Function(_ContainerInfo) _then) = __$ContainerInfoCopyWithImpl;
@override @useResult
$Res call({
 String? containerId, String? name, double lengthMm, double widthMm, double heightMm, double maxWeightKg, double volumeM3
});




}
/// @nodoc
class __$ContainerInfoCopyWithImpl<$Res>
    implements _$ContainerInfoCopyWith<$Res> {
  __$ContainerInfoCopyWithImpl(this._self, this._then);

  final _ContainerInfo _self;
  final $Res Function(_ContainerInfo) _then;

/// Create a copy of ContainerInfo
/// with the given fields replaced by the non-null parameter values.
@override @pragma('vm:prefer-inline') $Res call({Object? containerId = freezed,Object? name = freezed,Object? lengthMm = null,Object? widthMm = null,Object? heightMm = null,Object? maxWeightKg = null,Object? volumeM3 = null,}) {
  return _then(_ContainerInfo(
containerId: freezed == containerId ? _self.containerId : containerId // ignore: cast_nullable_to_non_nullable
as String?,name: freezed == name ? _self.name : name // ignore: cast_nullable_to_non_nullable
as String?,lengthMm: null == lengthMm ? _self.lengthMm : lengthMm // ignore: cast_nullable_to_non_nullable
as double,widthMm: null == widthMm ? _self.widthMm : widthMm // ignore: cast_nullable_to_non_nullable
as double,heightMm: null == heightMm ? _self.heightMm : heightMm // ignore: cast_nullable_to_non_nullable
as double,maxWeightKg: null == maxWeightKg ? _self.maxWeightKg : maxWeightKg // ignore: cast_nullable_to_non_nullable
as double,volumeM3: null == volumeM3 ? _self.volumeM3 : volumeM3 // ignore: cast_nullable_to_non_nullable
as double,
  ));
}


}

/// @nodoc
mixin _$PlanStats {

 int get totalItems; double get totalWeightKg; double get totalVolumeM3; double get volumeUtilizationPct; double get weightUtilizationPct;
/// Create a copy of PlanStats
/// with the given fields replaced by the non-null parameter values.
@JsonKey(includeFromJson: false, includeToJson: false)
@pragma('vm:prefer-inline')
$PlanStatsCopyWith<PlanStats> get copyWith => _$PlanStatsCopyWithImpl<PlanStats>(this as PlanStats, _$identity);



@override
bool operator ==(Object other) {
  return identical(this, other) || (other.runtimeType == runtimeType&&other is PlanStats&&(identical(other.totalItems, totalItems) || other.totalItems == totalItems)&&(identical(other.totalWeightKg, totalWeightKg) || other.totalWeightKg == totalWeightKg)&&(identical(other.totalVolumeM3, totalVolumeM3) || other.totalVolumeM3 == totalVolumeM3)&&(identical(other.volumeUtilizationPct, volumeUtilizationPct) || other.volumeUtilizationPct == volumeUtilizationPct)&&(identical(other.weightUtilizationPct, weightUtilizationPct) || other.weightUtilizationPct == weightUtilizationPct));
}


@override
int get hashCode => Object.hash(runtimeType,totalItems,totalWeightKg,totalVolumeM3,volumeUtilizationPct,weightUtilizationPct);

@override
String toString() {
  return 'PlanStats(totalItems: $totalItems, totalWeightKg: $totalWeightKg, totalVolumeM3: $totalVolumeM3, volumeUtilizationPct: $volumeUtilizationPct, weightUtilizationPct: $weightUtilizationPct)';
}


}

/// @nodoc
abstract mixin class $PlanStatsCopyWith<$Res>  {
  factory $PlanStatsCopyWith(PlanStats value, $Res Function(PlanStats) _then) = _$PlanStatsCopyWithImpl;
@useResult
$Res call({
 int totalItems, double totalWeightKg, double totalVolumeM3, double volumeUtilizationPct, double weightUtilizationPct
});




}
/// @nodoc
class _$PlanStatsCopyWithImpl<$Res>
    implements $PlanStatsCopyWith<$Res> {
  _$PlanStatsCopyWithImpl(this._self, this._then);

  final PlanStats _self;
  final $Res Function(PlanStats) _then;

/// Create a copy of PlanStats
/// with the given fields replaced by the non-null parameter values.
@pragma('vm:prefer-inline') @override $Res call({Object? totalItems = null,Object? totalWeightKg = null,Object? totalVolumeM3 = null,Object? volumeUtilizationPct = null,Object? weightUtilizationPct = null,}) {
  return _then(_self.copyWith(
totalItems: null == totalItems ? _self.totalItems : totalItems // ignore: cast_nullable_to_non_nullable
as int,totalWeightKg: null == totalWeightKg ? _self.totalWeightKg : totalWeightKg // ignore: cast_nullable_to_non_nullable
as double,totalVolumeM3: null == totalVolumeM3 ? _self.totalVolumeM3 : totalVolumeM3 // ignore: cast_nullable_to_non_nullable
as double,volumeUtilizationPct: null == volumeUtilizationPct ? _self.volumeUtilizationPct : volumeUtilizationPct // ignore: cast_nullable_to_non_nullable
as double,weightUtilizationPct: null == weightUtilizationPct ? _self.weightUtilizationPct : weightUtilizationPct // ignore: cast_nullable_to_non_nullable
as double,
  ));
}

}


/// Adds pattern-matching-related methods to [PlanStats].
extension PlanStatsPatterns on PlanStats {
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

@optionalTypeArgs TResult maybeMap<TResult extends Object?>(TResult Function( _PlanStats value)?  $default,{required TResult orElse(),}){
final _that = this;
switch (_that) {
case _PlanStats() when $default != null:
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

@optionalTypeArgs TResult map<TResult extends Object?>(TResult Function( _PlanStats value)  $default,){
final _that = this;
switch (_that) {
case _PlanStats():
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

@optionalTypeArgs TResult? mapOrNull<TResult extends Object?>(TResult? Function( _PlanStats value)?  $default,){
final _that = this;
switch (_that) {
case _PlanStats() when $default != null:
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

@optionalTypeArgs TResult maybeWhen<TResult extends Object?>(TResult Function( int totalItems,  double totalWeightKg,  double totalVolumeM3,  double volumeUtilizationPct,  double weightUtilizationPct)?  $default,{required TResult orElse(),}) {final _that = this;
switch (_that) {
case _PlanStats() when $default != null:
return $default(_that.totalItems,_that.totalWeightKg,_that.totalVolumeM3,_that.volumeUtilizationPct,_that.weightUtilizationPct);case _:
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

@optionalTypeArgs TResult when<TResult extends Object?>(TResult Function( int totalItems,  double totalWeightKg,  double totalVolumeM3,  double volumeUtilizationPct,  double weightUtilizationPct)  $default,) {final _that = this;
switch (_that) {
case _PlanStats():
return $default(_that.totalItems,_that.totalWeightKg,_that.totalVolumeM3,_that.volumeUtilizationPct,_that.weightUtilizationPct);case _:
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

@optionalTypeArgs TResult? whenOrNull<TResult extends Object?>(TResult? Function( int totalItems,  double totalWeightKg,  double totalVolumeM3,  double volumeUtilizationPct,  double weightUtilizationPct)?  $default,) {final _that = this;
switch (_that) {
case _PlanStats() when $default != null:
return $default(_that.totalItems,_that.totalWeightKg,_that.totalVolumeM3,_that.volumeUtilizationPct,_that.weightUtilizationPct);case _:
  return null;

}
}

}

/// @nodoc


class _PlanStats implements PlanStats {
  const _PlanStats({required this.totalItems, required this.totalWeightKg, required this.totalVolumeM3, required this.volumeUtilizationPct, required this.weightUtilizationPct});
  

@override final  int totalItems;
@override final  double totalWeightKg;
@override final  double totalVolumeM3;
@override final  double volumeUtilizationPct;
@override final  double weightUtilizationPct;

/// Create a copy of PlanStats
/// with the given fields replaced by the non-null parameter values.
@override @JsonKey(includeFromJson: false, includeToJson: false)
@pragma('vm:prefer-inline')
_$PlanStatsCopyWith<_PlanStats> get copyWith => __$PlanStatsCopyWithImpl<_PlanStats>(this, _$identity);



@override
bool operator ==(Object other) {
  return identical(this, other) || (other.runtimeType == runtimeType&&other is _PlanStats&&(identical(other.totalItems, totalItems) || other.totalItems == totalItems)&&(identical(other.totalWeightKg, totalWeightKg) || other.totalWeightKg == totalWeightKg)&&(identical(other.totalVolumeM3, totalVolumeM3) || other.totalVolumeM3 == totalVolumeM3)&&(identical(other.volumeUtilizationPct, volumeUtilizationPct) || other.volumeUtilizationPct == volumeUtilizationPct)&&(identical(other.weightUtilizationPct, weightUtilizationPct) || other.weightUtilizationPct == weightUtilizationPct));
}


@override
int get hashCode => Object.hash(runtimeType,totalItems,totalWeightKg,totalVolumeM3,volumeUtilizationPct,weightUtilizationPct);

@override
String toString() {
  return 'PlanStats(totalItems: $totalItems, totalWeightKg: $totalWeightKg, totalVolumeM3: $totalVolumeM3, volumeUtilizationPct: $volumeUtilizationPct, weightUtilizationPct: $weightUtilizationPct)';
}


}

/// @nodoc
abstract mixin class _$PlanStatsCopyWith<$Res> implements $PlanStatsCopyWith<$Res> {
  factory _$PlanStatsCopyWith(_PlanStats value, $Res Function(_PlanStats) _then) = __$PlanStatsCopyWithImpl;
@override @useResult
$Res call({
 int totalItems, double totalWeightKg, double totalVolumeM3, double volumeUtilizationPct, double weightUtilizationPct
});




}
/// @nodoc
class __$PlanStatsCopyWithImpl<$Res>
    implements _$PlanStatsCopyWith<$Res> {
  __$PlanStatsCopyWithImpl(this._self, this._then);

  final _PlanStats _self;
  final $Res Function(_PlanStats) _then;

/// Create a copy of PlanStats
/// with the given fields replaced by the non-null parameter values.
@override @pragma('vm:prefer-inline') $Res call({Object? totalItems = null,Object? totalWeightKg = null,Object? totalVolumeM3 = null,Object? volumeUtilizationPct = null,Object? weightUtilizationPct = null,}) {
  return _then(_PlanStats(
totalItems: null == totalItems ? _self.totalItems : totalItems // ignore: cast_nullable_to_non_nullable
as int,totalWeightKg: null == totalWeightKg ? _self.totalWeightKg : totalWeightKg // ignore: cast_nullable_to_non_nullable
as double,totalVolumeM3: null == totalVolumeM3 ? _self.totalVolumeM3 : totalVolumeM3 // ignore: cast_nullable_to_non_nullable
as double,volumeUtilizationPct: null == volumeUtilizationPct ? _self.volumeUtilizationPct : volumeUtilizationPct // ignore: cast_nullable_to_non_nullable
as double,weightUtilizationPct: null == weightUtilizationPct ? _self.weightUtilizationPct : weightUtilizationPct // ignore: cast_nullable_to_non_nullable
as double,
  ));
}


}

/// @nodoc
mixin _$PlanItem {

 String get itemId; String? get productSku; String? get label; double get lengthMm; double get widthMm; double get heightMm; double get weightKg; int get quantity; double get totalWeightKg; double get totalVolumeM3; bool get allowRotation; int get stackingLimit; String? get colorHex; DateTime get createdAt;
/// Create a copy of PlanItem
/// with the given fields replaced by the non-null parameter values.
@JsonKey(includeFromJson: false, includeToJson: false)
@pragma('vm:prefer-inline')
$PlanItemCopyWith<PlanItem> get copyWith => _$PlanItemCopyWithImpl<PlanItem>(this as PlanItem, _$identity);



@override
bool operator ==(Object other) {
  return identical(this, other) || (other.runtimeType == runtimeType&&other is PlanItem&&(identical(other.itemId, itemId) || other.itemId == itemId)&&(identical(other.productSku, productSku) || other.productSku == productSku)&&(identical(other.label, label) || other.label == label)&&(identical(other.lengthMm, lengthMm) || other.lengthMm == lengthMm)&&(identical(other.widthMm, widthMm) || other.widthMm == widthMm)&&(identical(other.heightMm, heightMm) || other.heightMm == heightMm)&&(identical(other.weightKg, weightKg) || other.weightKg == weightKg)&&(identical(other.quantity, quantity) || other.quantity == quantity)&&(identical(other.totalWeightKg, totalWeightKg) || other.totalWeightKg == totalWeightKg)&&(identical(other.totalVolumeM3, totalVolumeM3) || other.totalVolumeM3 == totalVolumeM3)&&(identical(other.allowRotation, allowRotation) || other.allowRotation == allowRotation)&&(identical(other.stackingLimit, stackingLimit) || other.stackingLimit == stackingLimit)&&(identical(other.colorHex, colorHex) || other.colorHex == colorHex)&&(identical(other.createdAt, createdAt) || other.createdAt == createdAt));
}


@override
int get hashCode => Object.hash(runtimeType,itemId,productSku,label,lengthMm,widthMm,heightMm,weightKg,quantity,totalWeightKg,totalVolumeM3,allowRotation,stackingLimit,colorHex,createdAt);

@override
String toString() {
  return 'PlanItem(itemId: $itemId, productSku: $productSku, label: $label, lengthMm: $lengthMm, widthMm: $widthMm, heightMm: $heightMm, weightKg: $weightKg, quantity: $quantity, totalWeightKg: $totalWeightKg, totalVolumeM3: $totalVolumeM3, allowRotation: $allowRotation, stackingLimit: $stackingLimit, colorHex: $colorHex, createdAt: $createdAt)';
}


}

/// @nodoc
abstract mixin class $PlanItemCopyWith<$Res>  {
  factory $PlanItemCopyWith(PlanItem value, $Res Function(PlanItem) _then) = _$PlanItemCopyWithImpl;
@useResult
$Res call({
 String itemId, String? productSku, String? label, double lengthMm, double widthMm, double heightMm, double weightKg, int quantity, double totalWeightKg, double totalVolumeM3, bool allowRotation, int stackingLimit, String? colorHex, DateTime createdAt
});




}
/// @nodoc
class _$PlanItemCopyWithImpl<$Res>
    implements $PlanItemCopyWith<$Res> {
  _$PlanItemCopyWithImpl(this._self, this._then);

  final PlanItem _self;
  final $Res Function(PlanItem) _then;

/// Create a copy of PlanItem
/// with the given fields replaced by the non-null parameter values.
@pragma('vm:prefer-inline') @override $Res call({Object? itemId = null,Object? productSku = freezed,Object? label = freezed,Object? lengthMm = null,Object? widthMm = null,Object? heightMm = null,Object? weightKg = null,Object? quantity = null,Object? totalWeightKg = null,Object? totalVolumeM3 = null,Object? allowRotation = null,Object? stackingLimit = null,Object? colorHex = freezed,Object? createdAt = null,}) {
  return _then(_self.copyWith(
itemId: null == itemId ? _self.itemId : itemId // ignore: cast_nullable_to_non_nullable
as String,productSku: freezed == productSku ? _self.productSku : productSku // ignore: cast_nullable_to_non_nullable
as String?,label: freezed == label ? _self.label : label // ignore: cast_nullable_to_non_nullable
as String?,lengthMm: null == lengthMm ? _self.lengthMm : lengthMm // ignore: cast_nullable_to_non_nullable
as double,widthMm: null == widthMm ? _self.widthMm : widthMm // ignore: cast_nullable_to_non_nullable
as double,heightMm: null == heightMm ? _self.heightMm : heightMm // ignore: cast_nullable_to_non_nullable
as double,weightKg: null == weightKg ? _self.weightKg : weightKg // ignore: cast_nullable_to_non_nullable
as double,quantity: null == quantity ? _self.quantity : quantity // ignore: cast_nullable_to_non_nullable
as int,totalWeightKg: null == totalWeightKg ? _self.totalWeightKg : totalWeightKg // ignore: cast_nullable_to_non_nullable
as double,totalVolumeM3: null == totalVolumeM3 ? _self.totalVolumeM3 : totalVolumeM3 // ignore: cast_nullable_to_non_nullable
as double,allowRotation: null == allowRotation ? _self.allowRotation : allowRotation // ignore: cast_nullable_to_non_nullable
as bool,stackingLimit: null == stackingLimit ? _self.stackingLimit : stackingLimit // ignore: cast_nullable_to_non_nullable
as int,colorHex: freezed == colorHex ? _self.colorHex : colorHex // ignore: cast_nullable_to_non_nullable
as String?,createdAt: null == createdAt ? _self.createdAt : createdAt // ignore: cast_nullable_to_non_nullable
as DateTime,
  ));
}

}


/// Adds pattern-matching-related methods to [PlanItem].
extension PlanItemPatterns on PlanItem {
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

@optionalTypeArgs TResult maybeMap<TResult extends Object?>(TResult Function( _PlanItem value)?  $default,{required TResult orElse(),}){
final _that = this;
switch (_that) {
case _PlanItem() when $default != null:
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

@optionalTypeArgs TResult map<TResult extends Object?>(TResult Function( _PlanItem value)  $default,){
final _that = this;
switch (_that) {
case _PlanItem():
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

@optionalTypeArgs TResult? mapOrNull<TResult extends Object?>(TResult? Function( _PlanItem value)?  $default,){
final _that = this;
switch (_that) {
case _PlanItem() when $default != null:
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

@optionalTypeArgs TResult maybeWhen<TResult extends Object?>(TResult Function( String itemId,  String? productSku,  String? label,  double lengthMm,  double widthMm,  double heightMm,  double weightKg,  int quantity,  double totalWeightKg,  double totalVolumeM3,  bool allowRotation,  int stackingLimit,  String? colorHex,  DateTime createdAt)?  $default,{required TResult orElse(),}) {final _that = this;
switch (_that) {
case _PlanItem() when $default != null:
return $default(_that.itemId,_that.productSku,_that.label,_that.lengthMm,_that.widthMm,_that.heightMm,_that.weightKg,_that.quantity,_that.totalWeightKg,_that.totalVolumeM3,_that.allowRotation,_that.stackingLimit,_that.colorHex,_that.createdAt);case _:
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

@optionalTypeArgs TResult when<TResult extends Object?>(TResult Function( String itemId,  String? productSku,  String? label,  double lengthMm,  double widthMm,  double heightMm,  double weightKg,  int quantity,  double totalWeightKg,  double totalVolumeM3,  bool allowRotation,  int stackingLimit,  String? colorHex,  DateTime createdAt)  $default,) {final _that = this;
switch (_that) {
case _PlanItem():
return $default(_that.itemId,_that.productSku,_that.label,_that.lengthMm,_that.widthMm,_that.heightMm,_that.weightKg,_that.quantity,_that.totalWeightKg,_that.totalVolumeM3,_that.allowRotation,_that.stackingLimit,_that.colorHex,_that.createdAt);case _:
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

@optionalTypeArgs TResult? whenOrNull<TResult extends Object?>(TResult? Function( String itemId,  String? productSku,  String? label,  double lengthMm,  double widthMm,  double heightMm,  double weightKg,  int quantity,  double totalWeightKg,  double totalVolumeM3,  bool allowRotation,  int stackingLimit,  String? colorHex,  DateTime createdAt)?  $default,) {final _that = this;
switch (_that) {
case _PlanItem() when $default != null:
return $default(_that.itemId,_that.productSku,_that.label,_that.lengthMm,_that.widthMm,_that.heightMm,_that.weightKg,_that.quantity,_that.totalWeightKg,_that.totalVolumeM3,_that.allowRotation,_that.stackingLimit,_that.colorHex,_that.createdAt);case _:
  return null;

}
}

}

/// @nodoc


class _PlanItem implements PlanItem {
  const _PlanItem({required this.itemId, this.productSku, this.label, required this.lengthMm, required this.widthMm, required this.heightMm, required this.weightKg, required this.quantity, required this.totalWeightKg, required this.totalVolumeM3, required this.allowRotation, required this.stackingLimit, this.colorHex, required this.createdAt});
  

@override final  String itemId;
@override final  String? productSku;
@override final  String? label;
@override final  double lengthMm;
@override final  double widthMm;
@override final  double heightMm;
@override final  double weightKg;
@override final  int quantity;
@override final  double totalWeightKg;
@override final  double totalVolumeM3;
@override final  bool allowRotation;
@override final  int stackingLimit;
@override final  String? colorHex;
@override final  DateTime createdAt;

/// Create a copy of PlanItem
/// with the given fields replaced by the non-null parameter values.
@override @JsonKey(includeFromJson: false, includeToJson: false)
@pragma('vm:prefer-inline')
_$PlanItemCopyWith<_PlanItem> get copyWith => __$PlanItemCopyWithImpl<_PlanItem>(this, _$identity);



@override
bool operator ==(Object other) {
  return identical(this, other) || (other.runtimeType == runtimeType&&other is _PlanItem&&(identical(other.itemId, itemId) || other.itemId == itemId)&&(identical(other.productSku, productSku) || other.productSku == productSku)&&(identical(other.label, label) || other.label == label)&&(identical(other.lengthMm, lengthMm) || other.lengthMm == lengthMm)&&(identical(other.widthMm, widthMm) || other.widthMm == widthMm)&&(identical(other.heightMm, heightMm) || other.heightMm == heightMm)&&(identical(other.weightKg, weightKg) || other.weightKg == weightKg)&&(identical(other.quantity, quantity) || other.quantity == quantity)&&(identical(other.totalWeightKg, totalWeightKg) || other.totalWeightKg == totalWeightKg)&&(identical(other.totalVolumeM3, totalVolumeM3) || other.totalVolumeM3 == totalVolumeM3)&&(identical(other.allowRotation, allowRotation) || other.allowRotation == allowRotation)&&(identical(other.stackingLimit, stackingLimit) || other.stackingLimit == stackingLimit)&&(identical(other.colorHex, colorHex) || other.colorHex == colorHex)&&(identical(other.createdAt, createdAt) || other.createdAt == createdAt));
}


@override
int get hashCode => Object.hash(runtimeType,itemId,productSku,label,lengthMm,widthMm,heightMm,weightKg,quantity,totalWeightKg,totalVolumeM3,allowRotation,stackingLimit,colorHex,createdAt);

@override
String toString() {
  return 'PlanItem(itemId: $itemId, productSku: $productSku, label: $label, lengthMm: $lengthMm, widthMm: $widthMm, heightMm: $heightMm, weightKg: $weightKg, quantity: $quantity, totalWeightKg: $totalWeightKg, totalVolumeM3: $totalVolumeM3, allowRotation: $allowRotation, stackingLimit: $stackingLimit, colorHex: $colorHex, createdAt: $createdAt)';
}


}

/// @nodoc
abstract mixin class _$PlanItemCopyWith<$Res> implements $PlanItemCopyWith<$Res> {
  factory _$PlanItemCopyWith(_PlanItem value, $Res Function(_PlanItem) _then) = __$PlanItemCopyWithImpl;
@override @useResult
$Res call({
 String itemId, String? productSku, String? label, double lengthMm, double widthMm, double heightMm, double weightKg, int quantity, double totalWeightKg, double totalVolumeM3, bool allowRotation, int stackingLimit, String? colorHex, DateTime createdAt
});




}
/// @nodoc
class __$PlanItemCopyWithImpl<$Res>
    implements _$PlanItemCopyWith<$Res> {
  __$PlanItemCopyWithImpl(this._self, this._then);

  final _PlanItem _self;
  final $Res Function(_PlanItem) _then;

/// Create a copy of PlanItem
/// with the given fields replaced by the non-null parameter values.
@override @pragma('vm:prefer-inline') $Res call({Object? itemId = null,Object? productSku = freezed,Object? label = freezed,Object? lengthMm = null,Object? widthMm = null,Object? heightMm = null,Object? weightKg = null,Object? quantity = null,Object? totalWeightKg = null,Object? totalVolumeM3 = null,Object? allowRotation = null,Object? stackingLimit = null,Object? colorHex = freezed,Object? createdAt = null,}) {
  return _then(_PlanItem(
itemId: null == itemId ? _self.itemId : itemId // ignore: cast_nullable_to_non_nullable
as String,productSku: freezed == productSku ? _self.productSku : productSku // ignore: cast_nullable_to_non_nullable
as String?,label: freezed == label ? _self.label : label // ignore: cast_nullable_to_non_nullable
as String?,lengthMm: null == lengthMm ? _self.lengthMm : lengthMm // ignore: cast_nullable_to_non_nullable
as double,widthMm: null == widthMm ? _self.widthMm : widthMm // ignore: cast_nullable_to_non_nullable
as double,heightMm: null == heightMm ? _self.heightMm : heightMm // ignore: cast_nullable_to_non_nullable
as double,weightKg: null == weightKg ? _self.weightKg : weightKg // ignore: cast_nullable_to_non_nullable
as double,quantity: null == quantity ? _self.quantity : quantity // ignore: cast_nullable_to_non_nullable
as int,totalWeightKg: null == totalWeightKg ? _self.totalWeightKg : totalWeightKg // ignore: cast_nullable_to_non_nullable
as double,totalVolumeM3: null == totalVolumeM3 ? _self.totalVolumeM3 : totalVolumeM3 // ignore: cast_nullable_to_non_nullable
as double,allowRotation: null == allowRotation ? _self.allowRotation : allowRotation // ignore: cast_nullable_to_non_nullable
as bool,stackingLimit: null == stackingLimit ? _self.stackingLimit : stackingLimit // ignore: cast_nullable_to_non_nullable
as int,colorHex: freezed == colorHex ? _self.colorHex : colorHex // ignore: cast_nullable_to_non_nullable
as String?,createdAt: null == createdAt ? _self.createdAt : createdAt // ignore: cast_nullable_to_non_nullable
as DateTime,
  ));
}


}

/// @nodoc
mixin _$CalculationResult {

 String get jobId; String get status; String get algorithm; DateTime? get calculatedAt; int? get durationMs; double? get efficiencyScore; double? get volumeUtilizationPct; String get visualizationUrl; List<PlacementDetail>? get placements;
/// Create a copy of CalculationResult
/// with the given fields replaced by the non-null parameter values.
@JsonKey(includeFromJson: false, includeToJson: false)
@pragma('vm:prefer-inline')
$CalculationResultCopyWith<CalculationResult> get copyWith => _$CalculationResultCopyWithImpl<CalculationResult>(this as CalculationResult, _$identity);



@override
bool operator ==(Object other) {
  return identical(this, other) || (other.runtimeType == runtimeType&&other is CalculationResult&&(identical(other.jobId, jobId) || other.jobId == jobId)&&(identical(other.status, status) || other.status == status)&&(identical(other.algorithm, algorithm) || other.algorithm == algorithm)&&(identical(other.calculatedAt, calculatedAt) || other.calculatedAt == calculatedAt)&&(identical(other.durationMs, durationMs) || other.durationMs == durationMs)&&(identical(other.efficiencyScore, efficiencyScore) || other.efficiencyScore == efficiencyScore)&&(identical(other.volumeUtilizationPct, volumeUtilizationPct) || other.volumeUtilizationPct == volumeUtilizationPct)&&(identical(other.visualizationUrl, visualizationUrl) || other.visualizationUrl == visualizationUrl)&&const DeepCollectionEquality().equals(other.placements, placements));
}


@override
int get hashCode => Object.hash(runtimeType,jobId,status,algorithm,calculatedAt,durationMs,efficiencyScore,volumeUtilizationPct,visualizationUrl,const DeepCollectionEquality().hash(placements));

@override
String toString() {
  return 'CalculationResult(jobId: $jobId, status: $status, algorithm: $algorithm, calculatedAt: $calculatedAt, durationMs: $durationMs, efficiencyScore: $efficiencyScore, volumeUtilizationPct: $volumeUtilizationPct, visualizationUrl: $visualizationUrl, placements: $placements)';
}


}

/// @nodoc
abstract mixin class $CalculationResultCopyWith<$Res>  {
  factory $CalculationResultCopyWith(CalculationResult value, $Res Function(CalculationResult) _then) = _$CalculationResultCopyWithImpl;
@useResult
$Res call({
 String jobId, String status, String algorithm, DateTime? calculatedAt, int? durationMs, double? efficiencyScore, double? volumeUtilizationPct, String visualizationUrl, List<PlacementDetail>? placements
});




}
/// @nodoc
class _$CalculationResultCopyWithImpl<$Res>
    implements $CalculationResultCopyWith<$Res> {
  _$CalculationResultCopyWithImpl(this._self, this._then);

  final CalculationResult _self;
  final $Res Function(CalculationResult) _then;

/// Create a copy of CalculationResult
/// with the given fields replaced by the non-null parameter values.
@pragma('vm:prefer-inline') @override $Res call({Object? jobId = null,Object? status = null,Object? algorithm = null,Object? calculatedAt = freezed,Object? durationMs = freezed,Object? efficiencyScore = freezed,Object? volumeUtilizationPct = freezed,Object? visualizationUrl = null,Object? placements = freezed,}) {
  return _then(_self.copyWith(
jobId: null == jobId ? _self.jobId : jobId // ignore: cast_nullable_to_non_nullable
as String,status: null == status ? _self.status : status // ignore: cast_nullable_to_non_nullable
as String,algorithm: null == algorithm ? _self.algorithm : algorithm // ignore: cast_nullable_to_non_nullable
as String,calculatedAt: freezed == calculatedAt ? _self.calculatedAt : calculatedAt // ignore: cast_nullable_to_non_nullable
as DateTime?,durationMs: freezed == durationMs ? _self.durationMs : durationMs // ignore: cast_nullable_to_non_nullable
as int?,efficiencyScore: freezed == efficiencyScore ? _self.efficiencyScore : efficiencyScore // ignore: cast_nullable_to_non_nullable
as double?,volumeUtilizationPct: freezed == volumeUtilizationPct ? _self.volumeUtilizationPct : volumeUtilizationPct // ignore: cast_nullable_to_non_nullable
as double?,visualizationUrl: null == visualizationUrl ? _self.visualizationUrl : visualizationUrl // ignore: cast_nullable_to_non_nullable
as String,placements: freezed == placements ? _self.placements : placements // ignore: cast_nullable_to_non_nullable
as List<PlacementDetail>?,
  ));
}

}


/// Adds pattern-matching-related methods to [CalculationResult].
extension CalculationResultPatterns on CalculationResult {
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

@optionalTypeArgs TResult maybeMap<TResult extends Object?>(TResult Function( _CalculationResult value)?  $default,{required TResult orElse(),}){
final _that = this;
switch (_that) {
case _CalculationResult() when $default != null:
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

@optionalTypeArgs TResult map<TResult extends Object?>(TResult Function( _CalculationResult value)  $default,){
final _that = this;
switch (_that) {
case _CalculationResult():
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

@optionalTypeArgs TResult? mapOrNull<TResult extends Object?>(TResult? Function( _CalculationResult value)?  $default,){
final _that = this;
switch (_that) {
case _CalculationResult() when $default != null:
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

@optionalTypeArgs TResult maybeWhen<TResult extends Object?>(TResult Function( String jobId,  String status,  String algorithm,  DateTime? calculatedAt,  int? durationMs,  double? efficiencyScore,  double? volumeUtilizationPct,  String visualizationUrl,  List<PlacementDetail>? placements)?  $default,{required TResult orElse(),}) {final _that = this;
switch (_that) {
case _CalculationResult() when $default != null:
return $default(_that.jobId,_that.status,_that.algorithm,_that.calculatedAt,_that.durationMs,_that.efficiencyScore,_that.volumeUtilizationPct,_that.visualizationUrl,_that.placements);case _:
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

@optionalTypeArgs TResult when<TResult extends Object?>(TResult Function( String jobId,  String status,  String algorithm,  DateTime? calculatedAt,  int? durationMs,  double? efficiencyScore,  double? volumeUtilizationPct,  String visualizationUrl,  List<PlacementDetail>? placements)  $default,) {final _that = this;
switch (_that) {
case _CalculationResult():
return $default(_that.jobId,_that.status,_that.algorithm,_that.calculatedAt,_that.durationMs,_that.efficiencyScore,_that.volumeUtilizationPct,_that.visualizationUrl,_that.placements);case _:
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

@optionalTypeArgs TResult? whenOrNull<TResult extends Object?>(TResult? Function( String jobId,  String status,  String algorithm,  DateTime? calculatedAt,  int? durationMs,  double? efficiencyScore,  double? volumeUtilizationPct,  String visualizationUrl,  List<PlacementDetail>? placements)?  $default,) {final _that = this;
switch (_that) {
case _CalculationResult() when $default != null:
return $default(_that.jobId,_that.status,_that.algorithm,_that.calculatedAt,_that.durationMs,_that.efficiencyScore,_that.volumeUtilizationPct,_that.visualizationUrl,_that.placements);case _:
  return null;

}
}

}

/// @nodoc


class _CalculationResult implements CalculationResult {
  const _CalculationResult({required this.jobId, required this.status, required this.algorithm, this.calculatedAt, this.durationMs, this.efficiencyScore, this.volumeUtilizationPct, required this.visualizationUrl, final  List<PlacementDetail>? placements}): _placements = placements;
  

@override final  String jobId;
@override final  String status;
@override final  String algorithm;
@override final  DateTime? calculatedAt;
@override final  int? durationMs;
@override final  double? efficiencyScore;
@override final  double? volumeUtilizationPct;
@override final  String visualizationUrl;
 final  List<PlacementDetail>? _placements;
@override List<PlacementDetail>? get placements {
  final value = _placements;
  if (value == null) return null;
  if (_placements is EqualUnmodifiableListView) return _placements;
  // ignore: implicit_dynamic_type
  return EqualUnmodifiableListView(value);
}


/// Create a copy of CalculationResult
/// with the given fields replaced by the non-null parameter values.
@override @JsonKey(includeFromJson: false, includeToJson: false)
@pragma('vm:prefer-inline')
_$CalculationResultCopyWith<_CalculationResult> get copyWith => __$CalculationResultCopyWithImpl<_CalculationResult>(this, _$identity);



@override
bool operator ==(Object other) {
  return identical(this, other) || (other.runtimeType == runtimeType&&other is _CalculationResult&&(identical(other.jobId, jobId) || other.jobId == jobId)&&(identical(other.status, status) || other.status == status)&&(identical(other.algorithm, algorithm) || other.algorithm == algorithm)&&(identical(other.calculatedAt, calculatedAt) || other.calculatedAt == calculatedAt)&&(identical(other.durationMs, durationMs) || other.durationMs == durationMs)&&(identical(other.efficiencyScore, efficiencyScore) || other.efficiencyScore == efficiencyScore)&&(identical(other.volumeUtilizationPct, volumeUtilizationPct) || other.volumeUtilizationPct == volumeUtilizationPct)&&(identical(other.visualizationUrl, visualizationUrl) || other.visualizationUrl == visualizationUrl)&&const DeepCollectionEquality().equals(other._placements, _placements));
}


@override
int get hashCode => Object.hash(runtimeType,jobId,status,algorithm,calculatedAt,durationMs,efficiencyScore,volumeUtilizationPct,visualizationUrl,const DeepCollectionEquality().hash(_placements));

@override
String toString() {
  return 'CalculationResult(jobId: $jobId, status: $status, algorithm: $algorithm, calculatedAt: $calculatedAt, durationMs: $durationMs, efficiencyScore: $efficiencyScore, volumeUtilizationPct: $volumeUtilizationPct, visualizationUrl: $visualizationUrl, placements: $placements)';
}


}

/// @nodoc
abstract mixin class _$CalculationResultCopyWith<$Res> implements $CalculationResultCopyWith<$Res> {
  factory _$CalculationResultCopyWith(_CalculationResult value, $Res Function(_CalculationResult) _then) = __$CalculationResultCopyWithImpl;
@override @useResult
$Res call({
 String jobId, String status, String algorithm, DateTime? calculatedAt, int? durationMs, double? efficiencyScore, double? volumeUtilizationPct, String visualizationUrl, List<PlacementDetail>? placements
});




}
/// @nodoc
class __$CalculationResultCopyWithImpl<$Res>
    implements _$CalculationResultCopyWith<$Res> {
  __$CalculationResultCopyWithImpl(this._self, this._then);

  final _CalculationResult _self;
  final $Res Function(_CalculationResult) _then;

/// Create a copy of CalculationResult
/// with the given fields replaced by the non-null parameter values.
@override @pragma('vm:prefer-inline') $Res call({Object? jobId = null,Object? status = null,Object? algorithm = null,Object? calculatedAt = freezed,Object? durationMs = freezed,Object? efficiencyScore = freezed,Object? volumeUtilizationPct = freezed,Object? visualizationUrl = null,Object? placements = freezed,}) {
  return _then(_CalculationResult(
jobId: null == jobId ? _self.jobId : jobId // ignore: cast_nullable_to_non_nullable
as String,status: null == status ? _self.status : status // ignore: cast_nullable_to_non_nullable
as String,algorithm: null == algorithm ? _self.algorithm : algorithm // ignore: cast_nullable_to_non_nullable
as String,calculatedAt: freezed == calculatedAt ? _self.calculatedAt : calculatedAt // ignore: cast_nullable_to_non_nullable
as DateTime?,durationMs: freezed == durationMs ? _self.durationMs : durationMs // ignore: cast_nullable_to_non_nullable
as int?,efficiencyScore: freezed == efficiencyScore ? _self.efficiencyScore : efficiencyScore // ignore: cast_nullable_to_non_nullable
as double?,volumeUtilizationPct: freezed == volumeUtilizationPct ? _self.volumeUtilizationPct : volumeUtilizationPct // ignore: cast_nullable_to_non_nullable
as double?,visualizationUrl: null == visualizationUrl ? _self.visualizationUrl : visualizationUrl // ignore: cast_nullable_to_non_nullable
as String,placements: freezed == placements ? _self._placements : placements // ignore: cast_nullable_to_non_nullable
as List<PlacementDetail>?,
  ));
}


}

/// @nodoc
mixin _$PlacementDetail {

 String get placementId; String get itemId; double get posX; double get posY; double get posZ; int get rotation; int get stepNumber;
/// Create a copy of PlacementDetail
/// with the given fields replaced by the non-null parameter values.
@JsonKey(includeFromJson: false, includeToJson: false)
@pragma('vm:prefer-inline')
$PlacementDetailCopyWith<PlacementDetail> get copyWith => _$PlacementDetailCopyWithImpl<PlacementDetail>(this as PlacementDetail, _$identity);



@override
bool operator ==(Object other) {
  return identical(this, other) || (other.runtimeType == runtimeType&&other is PlacementDetail&&(identical(other.placementId, placementId) || other.placementId == placementId)&&(identical(other.itemId, itemId) || other.itemId == itemId)&&(identical(other.posX, posX) || other.posX == posX)&&(identical(other.posY, posY) || other.posY == posY)&&(identical(other.posZ, posZ) || other.posZ == posZ)&&(identical(other.rotation, rotation) || other.rotation == rotation)&&(identical(other.stepNumber, stepNumber) || other.stepNumber == stepNumber));
}


@override
int get hashCode => Object.hash(runtimeType,placementId,itemId,posX,posY,posZ,rotation,stepNumber);

@override
String toString() {
  return 'PlacementDetail(placementId: $placementId, itemId: $itemId, posX: $posX, posY: $posY, posZ: $posZ, rotation: $rotation, stepNumber: $stepNumber)';
}


}

/// @nodoc
abstract mixin class $PlacementDetailCopyWith<$Res>  {
  factory $PlacementDetailCopyWith(PlacementDetail value, $Res Function(PlacementDetail) _then) = _$PlacementDetailCopyWithImpl;
@useResult
$Res call({
 String placementId, String itemId, double posX, double posY, double posZ, int rotation, int stepNumber
});




}
/// @nodoc
class _$PlacementDetailCopyWithImpl<$Res>
    implements $PlacementDetailCopyWith<$Res> {
  _$PlacementDetailCopyWithImpl(this._self, this._then);

  final PlacementDetail _self;
  final $Res Function(PlacementDetail) _then;

/// Create a copy of PlacementDetail
/// with the given fields replaced by the non-null parameter values.
@pragma('vm:prefer-inline') @override $Res call({Object? placementId = null,Object? itemId = null,Object? posX = null,Object? posY = null,Object? posZ = null,Object? rotation = null,Object? stepNumber = null,}) {
  return _then(_self.copyWith(
placementId: null == placementId ? _self.placementId : placementId // ignore: cast_nullable_to_non_nullable
as String,itemId: null == itemId ? _self.itemId : itemId // ignore: cast_nullable_to_non_nullable
as String,posX: null == posX ? _self.posX : posX // ignore: cast_nullable_to_non_nullable
as double,posY: null == posY ? _self.posY : posY // ignore: cast_nullable_to_non_nullable
as double,posZ: null == posZ ? _self.posZ : posZ // ignore: cast_nullable_to_non_nullable
as double,rotation: null == rotation ? _self.rotation : rotation // ignore: cast_nullable_to_non_nullable
as int,stepNumber: null == stepNumber ? _self.stepNumber : stepNumber // ignore: cast_nullable_to_non_nullable
as int,
  ));
}

}


/// Adds pattern-matching-related methods to [PlacementDetail].
extension PlacementDetailPatterns on PlacementDetail {
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

@optionalTypeArgs TResult maybeMap<TResult extends Object?>(TResult Function( _PlacementDetail value)?  $default,{required TResult orElse(),}){
final _that = this;
switch (_that) {
case _PlacementDetail() when $default != null:
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

@optionalTypeArgs TResult map<TResult extends Object?>(TResult Function( _PlacementDetail value)  $default,){
final _that = this;
switch (_that) {
case _PlacementDetail():
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

@optionalTypeArgs TResult? mapOrNull<TResult extends Object?>(TResult? Function( _PlacementDetail value)?  $default,){
final _that = this;
switch (_that) {
case _PlacementDetail() when $default != null:
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

@optionalTypeArgs TResult maybeWhen<TResult extends Object?>(TResult Function( String placementId,  String itemId,  double posX,  double posY,  double posZ,  int rotation,  int stepNumber)?  $default,{required TResult orElse(),}) {final _that = this;
switch (_that) {
case _PlacementDetail() when $default != null:
return $default(_that.placementId,_that.itemId,_that.posX,_that.posY,_that.posZ,_that.rotation,_that.stepNumber);case _:
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

@optionalTypeArgs TResult when<TResult extends Object?>(TResult Function( String placementId,  String itemId,  double posX,  double posY,  double posZ,  int rotation,  int stepNumber)  $default,) {final _that = this;
switch (_that) {
case _PlacementDetail():
return $default(_that.placementId,_that.itemId,_that.posX,_that.posY,_that.posZ,_that.rotation,_that.stepNumber);case _:
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

@optionalTypeArgs TResult? whenOrNull<TResult extends Object?>(TResult? Function( String placementId,  String itemId,  double posX,  double posY,  double posZ,  int rotation,  int stepNumber)?  $default,) {final _that = this;
switch (_that) {
case _PlacementDetail() when $default != null:
return $default(_that.placementId,_that.itemId,_that.posX,_that.posY,_that.posZ,_that.rotation,_that.stepNumber);case _:
  return null;

}
}

}

/// @nodoc


class _PlacementDetail implements PlacementDetail {
  const _PlacementDetail({required this.placementId, required this.itemId, required this.posX, required this.posY, required this.posZ, required this.rotation, required this.stepNumber});
  

@override final  String placementId;
@override final  String itemId;
@override final  double posX;
@override final  double posY;
@override final  double posZ;
@override final  int rotation;
@override final  int stepNumber;

/// Create a copy of PlacementDetail
/// with the given fields replaced by the non-null parameter values.
@override @JsonKey(includeFromJson: false, includeToJson: false)
@pragma('vm:prefer-inline')
_$PlacementDetailCopyWith<_PlacementDetail> get copyWith => __$PlacementDetailCopyWithImpl<_PlacementDetail>(this, _$identity);



@override
bool operator ==(Object other) {
  return identical(this, other) || (other.runtimeType == runtimeType&&other is _PlacementDetail&&(identical(other.placementId, placementId) || other.placementId == placementId)&&(identical(other.itemId, itemId) || other.itemId == itemId)&&(identical(other.posX, posX) || other.posX == posX)&&(identical(other.posY, posY) || other.posY == posY)&&(identical(other.posZ, posZ) || other.posZ == posZ)&&(identical(other.rotation, rotation) || other.rotation == rotation)&&(identical(other.stepNumber, stepNumber) || other.stepNumber == stepNumber));
}


@override
int get hashCode => Object.hash(runtimeType,placementId,itemId,posX,posY,posZ,rotation,stepNumber);

@override
String toString() {
  return 'PlacementDetail(placementId: $placementId, itemId: $itemId, posX: $posX, posY: $posY, posZ: $posZ, rotation: $rotation, stepNumber: $stepNumber)';
}


}

/// @nodoc
abstract mixin class _$PlacementDetailCopyWith<$Res> implements $PlacementDetailCopyWith<$Res> {
  factory _$PlacementDetailCopyWith(_PlacementDetail value, $Res Function(_PlacementDetail) _then) = __$PlacementDetailCopyWithImpl;
@override @useResult
$Res call({
 String placementId, String itemId, double posX, double posY, double posZ, int rotation, int stepNumber
});




}
/// @nodoc
class __$PlacementDetailCopyWithImpl<$Res>
    implements _$PlacementDetailCopyWith<$Res> {
  __$PlacementDetailCopyWithImpl(this._self, this._then);

  final _PlacementDetail _self;
  final $Res Function(_PlacementDetail) _then;

/// Create a copy of PlacementDetail
/// with the given fields replaced by the non-null parameter values.
@override @pragma('vm:prefer-inline') $Res call({Object? placementId = null,Object? itemId = null,Object? posX = null,Object? posY = null,Object? posZ = null,Object? rotation = null,Object? stepNumber = null,}) {
  return _then(_PlacementDetail(
placementId: null == placementId ? _self.placementId : placementId // ignore: cast_nullable_to_non_nullable
as String,itemId: null == itemId ? _self.itemId : itemId // ignore: cast_nullable_to_non_nullable
as String,posX: null == posX ? _self.posX : posX // ignore: cast_nullable_to_non_nullable
as double,posY: null == posY ? _self.posY : posY // ignore: cast_nullable_to_non_nullable
as double,posZ: null == posZ ? _self.posZ : posZ // ignore: cast_nullable_to_non_nullable
as double,rotation: null == rotation ? _self.rotation : rotation // ignore: cast_nullable_to_non_nullable
as int,stepNumber: null == stepNumber ? _self.stepNumber : stepNumber // ignore: cast_nullable_to_non_nullable
as int,
  ));
}


}

// dart format on
