import 'package:freezed_annotation/freezed_annotation.dart';

part 'container_model.freezed.dart';

@freezed
abstract class ContainerModel with _$ContainerModel {
  const factory ContainerModel({
    required String id,
    required String name,
    required double innerLengthMm,
    required double innerWidthMm,
    required double innerHeightMm,
    required double maxWeightKg,
    String? description,
  }) = _ContainerModel;
}
