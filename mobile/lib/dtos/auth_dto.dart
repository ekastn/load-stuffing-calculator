import 'package:json_annotation/json_annotation.dart';

part 'auth_dto.g.dart';

@JsonSerializable()
class LoginRequestDto {
  final String username;
  final String password;
  @JsonKey(name: 'guest_token')
  final String? guestToken;

  LoginRequestDto({
    required this.username,
    required this.password,
    this.guestToken,
  });

  Map<String, dynamic> toJson() => _$LoginRequestDtoToJson(this);
}

@JsonSerializable()
class LoginResponseDto {
  @JsonKey(name: 'access_token')
  final String accessToken;
  @JsonKey(name: 'refresh_token')
  final String refreshToken;
  @JsonKey(name: 'active_workspace_id')
  final String? activeWorkspaceId;
  final UserSummaryDto user;

  LoginResponseDto({
    required this.accessToken,
    required this.refreshToken,
    this.activeWorkspaceId,
    required this.user,
  });

  factory LoginResponseDto.fromJson(Map<String, dynamic> json) => _$LoginResponseDtoFromJson(json);
}

@JsonSerializable()
class UserSummaryDto {
  final String id;
  final String username;
  final String role;

  UserSummaryDto({
    required this.id,
    required this.username,
    required this.role,
  });

  factory UserSummaryDto.fromJson(Map<String, dynamic> json) => _$UserSummaryDtoFromJson(json);
}

@JsonSerializable()
class AuthMeResponseDto {
  final UserSummaryDto user;
  @JsonKey(name: 'active_workspace_id')
  final String? activeWorkspaceId;
  final List<String> permissions;
  @JsonKey(name: 'is_platform_member')
  final bool isPlatformMember;

  AuthMeResponseDto({
    required this.user,
    this.activeWorkspaceId,
    required this.permissions,
    required this.isPlatformMember,
  });

  factory AuthMeResponseDto.fromJson(Map<String, dynamic> json) => _$AuthMeResponseDtoFromJson(json);
}
