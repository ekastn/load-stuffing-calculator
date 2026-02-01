import 'package:json_annotation/json_annotation.dart';

part 'api_response_dto.g.dart';

@JsonSerializable(genericArgumentFactories: true)
class ApiResponseDto<T> {
  final bool success;
  final T? data;
  final List<ErrorDetailDto>? errors;

  ApiResponseDto({
    required this.success,
    this.data,
    this.errors,
  });

  factory ApiResponseDto.fromJson(
    Map<String, dynamic> json,
    T Function(Object? json) fromJsonT,
  ) =>
      _$ApiResponseDtoFromJson(json, fromJsonT);
}

@JsonSerializable()
class ErrorDetailDto {
  final String message;
  final String? code;

  ErrorDetailDto({
    required this.message,
    this.code,
  });

  factory ErrorDetailDto.fromJson(Map<String, dynamic> json) =>
      _$ErrorDetailDtoFromJson(json);
}
