import '../dtos/api_response_dto.dart';
import '../dtos/plan_dto.dart';
import '../dtos/plan_detail_dto.dart';
import '../models/plan_model.dart';
import '../models/plan_detail_model.dart';
import '../mappers/plan_mapper.dart';
import 'api_service.dart';

class PlanService {
  final ApiService _api;

  PlanService(this._api);

  Future<List<PlanModel>> getPlans({int page = 1, int limit = 10}) async {
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

    return apiResponse.data?.map(PlanMapper.fromDto).toList() ?? [];
  }

  Future<PlanModel> createPlan(CreatePlanRequestDto request) async {
    final response = await _api.post('/plans', data: request.toJson());

    final apiResponse = ApiResponseDto<CreatePlanResponseDto>.fromJson(
      response.data,
      (json) => CreatePlanResponseDto.fromJson(json as Map<String, dynamic>),
    );

    if (!apiResponse.success) {
      throw Exception(
          apiResponse.errors?.firstOrNull?.message ?? 'Failed to create plan');
    }

    final dto = apiResponse.data!;
    // Map CreatePlanResponseDto to PlanModel
    return PlanModel(
      id: dto.planId,
      code: dto.planCode,
      title: dto.title,
      status: dto.status,
      totalItems: dto.totalItems,
      totalWeightKg: dto.totalWeightKg,
      volumeUtilizationPct: null, // Not provided in create response
      createdBy: 'You', // Default value since not in response
      createdAt: DateTime.parse(dto.createdAt),
    );
  }

  Future<void> deletePlan(String planId) async {
    final response = await _api.delete('/plans/$planId');

    final apiResponse = ApiResponseDto<Map<String, dynamic>>.fromJson(
      response.data,
      (json) => json as Map<String, dynamic>,
    );

    if (!apiResponse.success) {
      throw Exception(
          apiResponse.errors?.firstOrNull?.message ?? 'Failed to delete plan');
    }
  }

  Future<PlanDetailModel> getPlanDetail(String planId) async {
    final response = await _api.get('/plans/$planId');

    final apiResponse = ApiResponseDto<PlanDetailDto>.fromJson(
      response.data,
      (json) => PlanDetailDto.fromJson(json as Map<String, dynamic>),
    );

    if (!apiResponse.success) {
      throw Exception(
          apiResponse.errors?.firstOrNull?.message ?? 'Failed to fetch plan detail');
    }

    return PlanMapper.fromDetailDto(apiResponse.data!);
  }

  Future<PlanDetailModel> recalculate(
    String planId, {
    String strategy = 'bestfitdecreasing',
    String? goal,
    bool gravity = true,
  }) async {
    final Map<String, dynamic> requestData = {
      'strategy': strategy,
      'gravity': gravity,
    };

    if (goal != null) {
      requestData['goal'] = goal;
    }

    final response = await _api.post(
      '/plans/$planId/calculate',
      data: requestData,
    );

    final apiResponse = ApiResponseDto<PlanDetailDto>.fromJson(
      response.data,
      (json) => PlanDetailDto.fromJson(json as Map<String, dynamic>),
    );

    if (!apiResponse.success) {
      throw Exception(
          apiResponse.errors?.firstOrNull?.message ?? 'Failed to recalculate plan');
    }

    return PlanMapper.fromDetailDto(apiResponse.data!);
  }
}
