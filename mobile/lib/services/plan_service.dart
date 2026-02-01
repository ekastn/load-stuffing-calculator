import '../dtos/api_response_dto.dart';
import '../dtos/plan_dto.dart';
import 'api_service.dart';

class PlanService {
  final ApiService _api;

  PlanService(this._api);

  Future<List<PlanListItemDto>> getPlans({int page = 1, int limit = 10}) async {
    final response = await _api.get(
      '/plans',
      queryParameters: {
        'page': page,
        'limit': limit,
      },
    );

    final apiResponse = ApiResponseDto<List<PlanListItemDto>>.fromJson(
      response.data,
      (json) => (json as List)
          .map((e) => PlanListItemDto.fromJson(e as Map<String, dynamic>))
          .toList(),
    );

    if (!apiResponse.success) {
      throw Exception(
          apiResponse.errors?.firstOrNull?.message ?? 'Failed to fetch plans');
    }

    return apiResponse.data ?? [];
  }
}
