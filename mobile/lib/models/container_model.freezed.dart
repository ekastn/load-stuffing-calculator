// GENERATED CODE - DO NOT MODIFY BY HAND
// coverage:ignore-file
// ignore_for_file: type=lint
// ignore_for_file: unused_element, deprecated_member_use, deprecated_member_use_from_same_package, use_function_type_syntax_for_parameters, unnecessary_const, avoid_init_to_null, invalid_override_different_default_values_named, prefer_expression_function_bodies, annotate_overrides, invalid_annotation_target, unnecessary_question_mark

part of 'container_model.dart';

// **************************************************************************
// FreezedGenerator
// **************************************************************************

// dart format off
T _$identity<T>(T value) => value;
/// @nodoc
mixin _$ContainerModel {

 String get id; String get name; double get innerLengthMm; double get innerWidthMm; double get innerHeightMm; double get maxWeightKg; String? get description;
/// Create a copy of ContainerModel
/// with the given fields replaced by the non-null parameter values.
@JsonKey(includeFromJson: false, includeToJson: false)
@pragma('vm:prefer-inline')
$ContainerModelCopyWith<ContainerModel> get copyWith => _$ContainerModelCopyWithImpl<ContainerModel>(this as ContainerModel, _$identity);



@override
bool operator ==(Object other) {
  return identical(this, other) || (other.runtimeType == runtimeType&&other is ContainerModel&&(identical(other.id, id) || other.id == id)&&(identical(other.name, name) || other.name == name)&&(identical(other.innerLengthMm, innerLengthMm) || other.innerLengthMm == innerLengthMm)&&(identical(other.innerWidthMm, innerWidthMm) || other.innerWidthMm == innerWidthMm)&&(identical(other.innerHeightMm, innerHeightMm) || other.innerHeightMm == innerHeightMm)&&(identical(other.maxWeightKg, maxWeightKg) || other.maxWeightKg == maxWeightKg)&&(identical(other.description, description) || other.description == description));
}


@override
int get hashCode => Object.hash(runtimeType,id,name,innerLengthMm,innerWidthMm,innerHeightMm,maxWeightKg,description);

@override
String toString() {
  return 'ContainerModel(id: $id, name: $name, innerLengthMm: $innerLengthMm, innerWidthMm: $innerWidthMm, innerHeightMm: $innerHeightMm, maxWeightKg: $maxWeightKg, description: $description)';
}


}

/// @nodoc
abstract mixin class $ContainerModelCopyWith<$Res>  {
  factory $ContainerModelCopyWith(ContainerModel value, $Res Function(ContainerModel) _then) = _$ContainerModelCopyWithImpl;
@useResult
$Res call({
 String id, String name, double innerLengthMm, double innerWidthMm, double innerHeightMm, double maxWeightKg, String? description
});




}
/// @nodoc
class _$ContainerModelCopyWithImpl<$Res>
    implements $ContainerModelCopyWith<$Res> {
  _$ContainerModelCopyWithImpl(this._self, this._then);

  final ContainerModel _self;
  final $Res Function(ContainerModel) _then;

/// Create a copy of ContainerModel
/// with the given fields replaced by the non-null parameter values.
@pragma('vm:prefer-inline') @override $Res call({Object? id = null,Object? name = null,Object? innerLengthMm = null,Object? innerWidthMm = null,Object? innerHeightMm = null,Object? maxWeightKg = null,Object? description = freezed,}) {
  return _then(_self.copyWith(
id: null == id ? _self.id : id // ignore: cast_nullable_to_non_nullable
as String,name: null == name ? _self.name : name // ignore: cast_nullable_to_non_nullable
as String,innerLengthMm: null == innerLengthMm ? _self.innerLengthMm : innerLengthMm // ignore: cast_nullable_to_non_nullable
as double,innerWidthMm: null == innerWidthMm ? _self.innerWidthMm : innerWidthMm // ignore: cast_nullable_to_non_nullable
as double,innerHeightMm: null == innerHeightMm ? _self.innerHeightMm : innerHeightMm // ignore: cast_nullable_to_non_nullable
as double,maxWeightKg: null == maxWeightKg ? _self.maxWeightKg : maxWeightKg // ignore: cast_nullable_to_non_nullable
as double,description: freezed == description ? _self.description : description // ignore: cast_nullable_to_non_nullable
as String?,
  ));
}

}


/// Adds pattern-matching-related methods to [ContainerModel].
extension ContainerModelPatterns on ContainerModel {
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

@optionalTypeArgs TResult maybeMap<TResult extends Object?>(TResult Function( _ContainerModel value)?  $default,{required TResult orElse(),}){
final _that = this;
switch (_that) {
case _ContainerModel() when $default != null:
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

@optionalTypeArgs TResult map<TResult extends Object?>(TResult Function( _ContainerModel value)  $default,){
final _that = this;
switch (_that) {
case _ContainerModel():
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

@optionalTypeArgs TResult? mapOrNull<TResult extends Object?>(TResult? Function( _ContainerModel value)?  $default,){
final _that = this;
switch (_that) {
case _ContainerModel() when $default != null:
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

@optionalTypeArgs TResult maybeWhen<TResult extends Object?>(TResult Function( String id,  String name,  double innerLengthMm,  double innerWidthMm,  double innerHeightMm,  double maxWeightKg,  String? description)?  $default,{required TResult orElse(),}) {final _that = this;
switch (_that) {
case _ContainerModel() when $default != null:
return $default(_that.id,_that.name,_that.innerLengthMm,_that.innerWidthMm,_that.innerHeightMm,_that.maxWeightKg,_that.description);case _:
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

@optionalTypeArgs TResult when<TResult extends Object?>(TResult Function( String id,  String name,  double innerLengthMm,  double innerWidthMm,  double innerHeightMm,  double maxWeightKg,  String? description)  $default,) {final _that = this;
switch (_that) {
case _ContainerModel():
return $default(_that.id,_that.name,_that.innerLengthMm,_that.innerWidthMm,_that.innerHeightMm,_that.maxWeightKg,_that.description);case _:
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

@optionalTypeArgs TResult? whenOrNull<TResult extends Object?>(TResult? Function( String id,  String name,  double innerLengthMm,  double innerWidthMm,  double innerHeightMm,  double maxWeightKg,  String? description)?  $default,) {final _that = this;
switch (_that) {
case _ContainerModel() when $default != null:
return $default(_that.id,_that.name,_that.innerLengthMm,_that.innerWidthMm,_that.innerHeightMm,_that.maxWeightKg,_that.description);case _:
  return null;

}
}

}

/// @nodoc


class _ContainerModel implements ContainerModel {
  const _ContainerModel({required this.id, required this.name, required this.innerLengthMm, required this.innerWidthMm, required this.innerHeightMm, required this.maxWeightKg, this.description});
  

@override final  String id;
@override final  String name;
@override final  double innerLengthMm;
@override final  double innerWidthMm;
@override final  double innerHeightMm;
@override final  double maxWeightKg;
@override final  String? description;

/// Create a copy of ContainerModel
/// with the given fields replaced by the non-null parameter values.
@override @JsonKey(includeFromJson: false, includeToJson: false)
@pragma('vm:prefer-inline')
_$ContainerModelCopyWith<_ContainerModel> get copyWith => __$ContainerModelCopyWithImpl<_ContainerModel>(this, _$identity);



@override
bool operator ==(Object other) {
  return identical(this, other) || (other.runtimeType == runtimeType&&other is _ContainerModel&&(identical(other.id, id) || other.id == id)&&(identical(other.name, name) || other.name == name)&&(identical(other.innerLengthMm, innerLengthMm) || other.innerLengthMm == innerLengthMm)&&(identical(other.innerWidthMm, innerWidthMm) || other.innerWidthMm == innerWidthMm)&&(identical(other.innerHeightMm, innerHeightMm) || other.innerHeightMm == innerHeightMm)&&(identical(other.maxWeightKg, maxWeightKg) || other.maxWeightKg == maxWeightKg)&&(identical(other.description, description) || other.description == description));
}


@override
int get hashCode => Object.hash(runtimeType,id,name,innerLengthMm,innerWidthMm,innerHeightMm,maxWeightKg,description);

@override
String toString() {
  return 'ContainerModel(id: $id, name: $name, innerLengthMm: $innerLengthMm, innerWidthMm: $innerWidthMm, innerHeightMm: $innerHeightMm, maxWeightKg: $maxWeightKg, description: $description)';
}


}

/// @nodoc
abstract mixin class _$ContainerModelCopyWith<$Res> implements $ContainerModelCopyWith<$Res> {
  factory _$ContainerModelCopyWith(_ContainerModel value, $Res Function(_ContainerModel) _then) = __$ContainerModelCopyWithImpl;
@override @useResult
$Res call({
 String id, String name, double innerLengthMm, double innerWidthMm, double innerHeightMm, double maxWeightKg, String? description
});




}
/// @nodoc
class __$ContainerModelCopyWithImpl<$Res>
    implements _$ContainerModelCopyWith<$Res> {
  __$ContainerModelCopyWithImpl(this._self, this._then);

  final _ContainerModel _self;
  final $Res Function(_ContainerModel) _then;

/// Create a copy of ContainerModel
/// with the given fields replaced by the non-null parameter values.
@override @pragma('vm:prefer-inline') $Res call({Object? id = null,Object? name = null,Object? innerLengthMm = null,Object? innerWidthMm = null,Object? innerHeightMm = null,Object? maxWeightKg = null,Object? description = freezed,}) {
  return _then(_ContainerModel(
id: null == id ? _self.id : id // ignore: cast_nullable_to_non_nullable
as String,name: null == name ? _self.name : name // ignore: cast_nullable_to_non_nullable
as String,innerLengthMm: null == innerLengthMm ? _self.innerLengthMm : innerLengthMm // ignore: cast_nullable_to_non_nullable
as double,innerWidthMm: null == innerWidthMm ? _self.innerWidthMm : innerWidthMm // ignore: cast_nullable_to_non_nullable
as double,innerHeightMm: null == innerHeightMm ? _self.innerHeightMm : innerHeightMm // ignore: cast_nullable_to_non_nullable
as double,maxWeightKg: null == maxWeightKg ? _self.maxWeightKg : maxWeightKg // ignore: cast_nullable_to_non_nullable
as double,description: freezed == description ? _self.description : description // ignore: cast_nullable_to_non_nullable
as String?,
  ));
}


}

// dart format on
