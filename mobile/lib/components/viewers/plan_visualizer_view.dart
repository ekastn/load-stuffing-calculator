import 'dart:io' show Platform;
import 'package:flutter/material.dart';
import 'package:webview_flutter/webview_flutter.dart';
import '../../config/constants.dart';
import '../../services/storage_service.dart';
import 'package:url_launcher/url_launcher.dart';

class PlanVisualizerView extends StatefulWidget {
  final String planId;

  const PlanVisualizerView({super.key, required this.planId});

  @override
  State<PlanVisualizerView> createState() => _PlanVisualizerViewState();
}

class _PlanVisualizerViewState extends State<PlanVisualizerView> {
  WebViewController? _controller;
  bool _isLoading = true;
  String? _error;
  bool _isDesktop = false;
  bool _isInitialized = false;

  @override
  void initState() {
    super.initState();

    // Check if running on desktop platform (WebView not fully supported)
    _isDesktop = Platform.isLinux || Platform.isWindows || Platform.isMacOS;

    if (!_isDesktop) {
      _initializeWebView();
    }
  }

  Future<void> _initializeWebView() async {
    try {
      // Get the access token from secure storage
      final storageService = StorageService();
      final token = await storageService.getAccessToken();

      if (!mounted) return;

      // Build URL with token parameter for authentication
      final baseUrl =
          '${Constants.webBaseUrl}/embed/shipments/${widget.planId}';
      final url = token != null ? '$baseUrl?token=$token' : baseUrl;

      _controller = WebViewController()
        ..setJavaScriptMode(JavaScriptMode.unrestricted)
        ..setBackgroundColor(Colors.white)
        ..setNavigationDelegate(
          NavigationDelegate(
            onPageStarted: (url) {
              if (mounted) {
                setState(() {
                  _isLoading = true;
                  _error = null;
                });
              }
            },
            onPageFinished: (url) {
              if (mounted) {
                setState(() {
                  _isLoading = false;
                });
              }
            },
            onWebResourceError: (error) {
              if (mounted) {
                setState(() {
                  _isLoading = false;
                  _error = 'Failed to load 3D viewer: ${error.description}';
                });
              }
              debugPrint(
                'WebView error: ${error.errorCode} - ${error.description}',
              );
            },
          ),
        )
        ..loadRequest(Uri.parse(url));

      if (mounted) {
        setState(() {
          _isInitialized = true;
        });
      }
    } catch (e) {
      if (mounted) {
        setState(() {
          _isLoading = false;
          _error = 'Failed to initialize viewer: $e';
        });
      }
      debugPrint('WebView initialization error: $e');
    }
  }

  Future<void> _openInBrowser() async {
    final url = Uri.parse(
      '${Constants.webBaseUrl}/embed/shipments/${widget.planId}',
    );
    if (await canLaunchUrl(url)) {
      await launchUrl(url, mode: LaunchMode.externalApplication);
    }
  }

  Future<void> _retry() async {
    setState(() {
      _error = null;
      _isLoading = true;
      _isInitialized = false;
    });
    await _initializeWebView();
  }

  @override
  Widget build(BuildContext context) {
    // Show fallback UI for desktop platforms
    if (_isDesktop) {
      return Container(
        color: Colors.white,
        child: Center(
          child: Column(
            mainAxisAlignment: MainAxisAlignment.center,
            children: [
              Icon(Icons.desktop_mac, size: 64, color: Colors.grey[400]),
              const SizedBox(height: 24),
              Text(
                '3D Viewer Not Available on Desktop',
                style: TextStyle(
                  fontSize: 18,
                  fontWeight: FontWeight.bold,
                  color: Colors.grey[700],
                ),
              ),
              const SizedBox(height: 8),
              Text(
                'WebView is not supported on Linux desktop',
                style: TextStyle(fontSize: 14, color: Colors.grey[600]),
              ),
              const SizedBox(height: 24),
              FilledButton.icon(
                onPressed: _openInBrowser,
                icon: const Icon(Icons.open_in_browser),
                label: const Text('Open in Browser'),
              ),
              const SizedBox(height: 16),
              Text(
                'The 3D viewer will open in your default web browser',
                style: TextStyle(fontSize: 12, color: Colors.grey[500]),
              ),
            ],
          ),
        ),
      );
    }

    // Mobile WebView UI
    return Stack(
      children: [
        if (_controller != null && _isInitialized)
          WebViewWidget(controller: _controller!),
        if (_isLoading)
          Container(
            color: Colors.white,
            child: const Center(
              child: Column(
                mainAxisAlignment: MainAxisAlignment.center,
                children: [
                  CircularProgressIndicator(),
                  SizedBox(height: 16),
                  Text(
                    'Loading 3D Viewer...',
                    style: TextStyle(fontSize: 14, color: Colors.grey),
                  ),
                ],
              ),
            ),
          ),
        if (_error != null)
          Container(
            color: Colors.white,
            child: Center(
              child: Column(
                mainAxisAlignment: MainAxisAlignment.center,
                children: [
                  const Icon(Icons.error_outline, size: 48, color: Colors.red),
                  const SizedBox(height: 16),
                  Padding(
                    padding: const EdgeInsets.symmetric(horizontal: 24),
                    child: Text(
                      _error!,
                      style: const TextStyle(
                        fontSize: 14,
                        fontWeight: FontWeight.bold,
                        color: Colors.red,
                      ),
                      textAlign: TextAlign.center,
                    ),
                  ),
                  const SizedBox(height: 16),
                  TextButton(onPressed: _retry, child: const Text('Retry')),
                ],
              ),
            ),
          ),
      ],
    );
  }
}
