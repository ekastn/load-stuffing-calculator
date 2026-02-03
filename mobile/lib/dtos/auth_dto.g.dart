// GENERATED CODE - DO NOT MODIFY BY HAND

part of 'auth_dto.dart';

// **************************************************************************
// JsonSerializableGenerator
// **************************************************************************

LoginRequestDto _$LoginRequestDtoFromJson(Map<String, dynamic> json) =>
    LoginRequestDto(
      username: json['username'] as String,
      password: json['password'] as String,
      guestToken: json['guest_token'] as String?,
    );

Map<String, dynamic> _$LoginRequestDtoToJson(LoginRequestDto instance) =>
    <String, dynamic>{
      'username': instance.username,
      'password': instance.password,
      'guest_token': instance.guestToken,
    };

LoginResponseDto _$LoginResponseDtoFromJson(Map<String, dynamic> json) =>
    LoginResponseDto(
      accessToken: json['access_token'] as String,
      refreshToken: json['refresh_token'] as String,
      activeWorkspaceId: json['active_workspace_id'] as String?,
      user: UserSummaryDto.fromJson(json['user'] as Map<String, dynamic>),
    );

Map<String, dynamic> _$LoginResponseDtoToJson(LoginResponseDto instance) =>
    <String, dynamic>{
      'access_token': instance.accessToken,
      'refresh_token': instance.refreshToken,
      'active_workspace_id': instance.activeWorkspaceId,
      'user': instance.user,
    };

UserSummaryDto _$UserSummaryDtoFromJson(Map<String, dynamic> json) =>
    UserSummaryDto(
      id: json['id'] as String,
      username: json['username'] as String,
      role: json['role'] as String,
    );

Map<String, dynamic> _$UserSummaryDtoToJson(UserSummaryDto instance) =>
    <String, dynamic>{
      'id': instance.id,
      'username': instance.username,
      'role': instance.role,
    };
