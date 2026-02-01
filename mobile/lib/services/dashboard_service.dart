import '../dtos/api_response_dto.dart';
import '../dtos/dashboard_dto.dart';
import 'api_service.dart';

class DashboardService {
  final ApiService _api;

  DashboardService(this._api);

  Future<DashboardStatsDto> getStats() async {
    final response = await _api.get('/dashboard');

    final apiResponse = ApiResponseDto<DashboardStatsDto>.fromJson(
      response.data,
      (json) => DashboardStatsDto.fromJson(json as Map<String, dynamic>),
    );

    if (!apiResponse.success || apiResponse.data == null) {
      throw Exception(
          apiResponse.errors?.firstOrNull?.message ?? 'Failed to fetch stats');
    }

    return apiResponse.data!;
  }
}
