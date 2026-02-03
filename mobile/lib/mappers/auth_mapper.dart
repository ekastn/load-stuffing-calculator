import '../dtos/auth_dto.dart';
import '../models/user_model.dart';

class AuthMapper {
  static UserModel toUserModel(UserSummaryDto dto) {
    return UserModel(
      id: dto.id,
      username: dto.username,
      role: dto.role,
    );
  }
}
