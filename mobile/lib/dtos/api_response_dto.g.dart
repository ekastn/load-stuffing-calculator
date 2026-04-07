// GENERATED CODE - DO NOT MODIFY BY HAND

part of 'api_response_dto.dart';

// **************************************************************************
// JsonSerializableGenerator
// **************************************************************************

ApiResponseDto<T> _$ApiResponseDtoFromJson<T>(
  Map<String, dynamic> json,
  T Function(Object? json) fromJsonT,
) => ApiResponseDto<T>(
  success: json['success'] as bool,
  data: _$nullableGenericFromJson(json['data'], fromJsonT),
  errors: (json['errors'] as List<dynamic>?)
      ?.map((e) => ErrorDetailDto.fromJson(e as Map<String, dynamic>))
      .toList(),
);

Map<String, dynamic> _$ApiResponseDtoToJson<T>(
  ApiResponseDto<T> instance,
  Object? Function(T value) toJsonT,
) => <String, dynamic>{
  'success': instance.success,
  'data': _$nullableGenericToJson(instance.data, toJsonT),
  'errors': instance.errors,
};

T? _$nullableGenericFromJson<T>(
  Object? input,
  T Function(Object? json) fromJson,
) => input == null ? null : fromJson(input);

Object? _$nullableGenericToJson<T>(
  T? input,
  Object? Function(T value) toJson,
) => input == null ? null : toJson(input);

ErrorDetailDto _$ErrorDetailDtoFromJson(Map<String, dynamic> json) =>
    ErrorDetailDto(
      message: json['message'] as String,
      code: json['code'] as String?,
    );

Map<String, dynamic> _$ErrorDetailDtoToJson(ErrorDetailDto instance) =>
    <String, dynamic>{'message': instance.message, 'code': instance.code};
