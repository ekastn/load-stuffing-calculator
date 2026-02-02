import 'package:flutter/material.dart';
import '../dtos/container_dto.dart';
import '../models/container_model.dart';
import '../services/container_service.dart';

class ContainerProvider with ChangeNotifier {
  final ContainerService _service;

  ContainerProvider(this._service);
  
  List<ContainerModel> _containers = [];
  bool _isLoading = false;
  String? _error;

  List<ContainerModel> get containers => _containers;
  bool get isLoading => _isLoading;
  String? get error => _error;

  Future<void> fetchContainers() async {
    _isLoading = true;
    _error = null;
    notifyListeners();

    try {
      _containers = await _service.getContainers();
    } catch (e) {
      _error = e.toString();
    } finally {
      _isLoading = false;
      notifyListeners();
    }
  }

  Future<void> createContainer(CreateContainerRequestDto request) async {
    _isLoading = true;
    _error = null;
    notifyListeners();

    try {
      final newContainer = await _service.createContainer(request);
      _containers.add(newContainer);
    } catch (e) {
      _error = e.toString();
      rethrow;
    } finally {
      _isLoading = false;
      notifyListeners();
    }
  }

  Future<void> updateContainer(String id, UpdateContainerRequestDto request) async {
    _isLoading = true;
    _error = null;
    notifyListeners();

    try {
      final updatedContainer = await _service.updateContainer(id, request);
      final index = _containers.indexWhere((c) => c.id == id);
      if (index != -1) {
        _containers[index] = updatedContainer;
      }
    } catch (e) {
      _error = e.toString();
      rethrow;
    } finally {
      _isLoading = false;
      notifyListeners();
    }
  }

  Future<void> deleteContainer(String id) async {
    _isLoading = true;
    _error = null;
    notifyListeners();

    try {
      await _service.deleteContainer(id);
      _containers.removeWhere((c) => c.id == id);
    } catch (e) {
      _error = e.toString();
      rethrow;
    } finally {
      _isLoading = false;
      notifyListeners();
    }
  }
}
