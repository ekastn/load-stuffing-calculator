import 'package:flutter/material.dart';
import 'package:provider/provider.dart';
import '../../components/viewers/plan_visualizer_view.dart';
import '../../providers/loading_provider.dart';
import 'barcode_scanner_page.dart';

class LoadingSessionPage extends StatefulWidget {
  final String planId;

  const LoadingSessionPage({super.key, required this.planId});

  @override
  State<LoadingSessionPage> createState() => _LoadingSessionPageState();
}

class _LoadingSessionPageState extends State<LoadingSessionPage> {
  @override
  void initState() {
    super.initState();
    WidgetsBinding.instance.addPostFrameCallback((_) {
      _startSession();
    });
  }

  Future<void> _startSession() async {
    final provider = context.read<LoadingProvider>();
    await provider.startSession(widget.planId);
  }

  @override
  Widget build(BuildContext context) {
    return Consumer<LoadingProvider>(
      builder: (context, provider, child) {
        final session = provider.currentSession;
        final expectedItem = provider.getCurrentExpectedItem();

        if (provider.isLoading) {
          return const Scaffold(
            body: Center(child: CircularProgressIndicator()),
          );
        }

        if (provider.error != null) {
          return Scaffold(
            appBar: AppBar(title: const Text('Loading Error')),
            body: Center(
              child: Column(
                mainAxisAlignment: MainAxisAlignment.center,
                children: [
                  const Icon(Icons.error_outline, color: Colors.red, size: 64),
                  const SizedBox(height: 16),
                  Text(provider.error!),
                  const SizedBox(height: 24),
                  ElevatedButton(
                    onPressed: _startSession,
                    child: const Text('Retry'),
                  ),
                ],
              ),
            ),
          );
        }

        if (session == null) {
          return const Scaffold(
            body: Center(child: Text('Session not initialized')),
          );
        }

        return Scaffold(
          appBar: AppBar(
            title: Text(
              'Loading - Step ${session.currentStepIndex + 1}/${session.totalItems}',
            ),
          ),
          body: Column(
            children: [
              // Progress bar
              _buildProgressBar(provider),

              // Current item info
              _buildItemInfo(expectedItem),

              // 3D WebView Viewer
              Expanded(child: _build3DViewer(provider)),

              // Action buttons
              _buildActionButtons(provider, expectedItem),
            ],
          ),
        );
      },
    );
  }

  Widget _buildProgressBar(LoadingProvider provider) {
    final session = provider.currentSession!;
    final progress = session.totalItems == 0
        ? 0.0
        : session.validatedCount / session.totalItems;

    return Container(
      padding: const EdgeInsets.all(16),
      child: Column(
        children: [
          LinearProgressIndicator(
            value: progress,
            backgroundColor: Colors.grey[200],
            minHeight: 8,
          ),
          const SizedBox(height: 8),
          Row(
            mainAxisAlignment: MainAxisAlignment.spaceBetween,
            children: [
              Text(
                '${session.validatedCount}/${session.totalItems} items loaded',
                style: TextStyle(fontSize: 14, color: Colors.grey[600]),
              ),
              Text(
                '${(progress * 100).round()}%',
                style: const TextStyle(fontWeight: FontWeight.bold),
              ),
            ],
          ),
        ],
      ),
    );
  }

  Widget _buildItemInfo(ExpectedItem? item) {
    if (item == null) {
      return const Card(
        margin: EdgeInsets.all(16),
        child: Padding(
          padding: EdgeInsets.all(16),
          child: Center(
            child: Text(
              'All items loaded!',
              style: TextStyle(fontSize: 18, fontWeight: FontWeight.bold),
            ),
          ),
        ),
      );
    }

    return Card(
      margin: const EdgeInsets.symmetric(horizontal: 16),
      child: Padding(
        padding: const EdgeInsets.all(16),
        child: Column(
          crossAxisAlignment: CrossAxisAlignment.start,
          children: [
            Row(
              children: [
                const Icon(Icons.inventory, size: 16, color: Colors.blue),
                const SizedBox(width: 8),
                Text(
                  'Next Item:',
                  style: TextStyle(fontSize: 12, color: Colors.grey[600]),
                ),
              ],
            ),
            const SizedBox(height: 4),
            Text(
              item.itemLabel,
              style: const TextStyle(fontSize: 20, fontWeight: FontWeight.bold),
            ),
            const Divider(),
            Row(
              mainAxisAlignment: MainAxisAlignment.spaceBetween,
              children: [
                _buildInfoBit('Dimensions', item.dimensions),
                _buildInfoBit('Position', item.position),
              ],
            ),
          ],
        ),
      ),
    );
  }

  Widget _buildInfoBit(String label, String value) {
    return Column(
      crossAxisAlignment: CrossAxisAlignment.start,
      children: [
        Text(label, style: const TextStyle(fontSize: 11, color: Colors.grey)),
        Text(
          value,
          style: const TextStyle(fontSize: 13, fontWeight: FontWeight.w500),
        ),
      ],
    );
  }

  Widget _build3DViewer(LoadingProvider provider) {
    final session = provider.currentSession;
    if (session == null) return const SizedBox.shrink();

    // We reuse PlanVisualizerView but pass the current step to highlight
    return PlanVisualizerView(planId: session.planId);
  }

  Widget _buildActionButtons(LoadingProvider provider, ExpectedItem? item) {
    if (item == null) {
      return Padding(
        padding: const EdgeInsets.all(16),
        child: ElevatedButton(
          onPressed: () => _completeSession(provider),
          style: ElevatedButton.styleFrom(
            minimumSize: const Size(double.infinity, 54),
            backgroundColor: Colors.green,
            foregroundColor: Colors.white,
          ),
          child: const Text('Complete Loading Session'),
        ),
      );
    }

    return Padding(
      padding: const EdgeInsets.all(16),
      child: Column(
        children: [
          ElevatedButton.icon(
            onPressed: () => _openBarcodeScanner(provider),
            icon: const Icon(Icons.qr_code_scanner),
            label: const Text('Scan QR Code'),
            style: ElevatedButton.styleFrom(
              minimumSize: const Size(double.infinity, 54),
              backgroundColor: Theme.of(context).primaryColor,
              foregroundColor: Colors.white,
            ),
          ),
          const SizedBox(height: 8),
          Row(
            children: [
              Expanded(
                child: OutlinedButton(
                  onPressed: () => _manualConfirm(provider),
                  child: const Text('Manual OK'),
                ),
              ),
              const SizedBox(width: 8),
              Expanded(
                child: OutlinedButton(
                  onPressed: () => _skipItem(provider),
                  child: const Text('Skip'),
                ),
              ),
            ],
          ),
        ],
      ),
    );
  }

  Future<void> _openBarcodeScanner(LoadingProvider provider) async {
    final scannedBarcode = await Navigator.push<String>(
      context,
      MaterialPageRoute(builder: (_) => const BarcodeScannerPage()),
    );

    if (scannedBarcode != null) {
      final result = provider.validateBarcode(scannedBarcode);

      if (result.matched) {
        _showValidationFeedback(
          true,
          'Correct! ${result.expectedItem?.itemLabel}',
        );
      } else {
        _showErrorDialog(result);
      }
    }
  }

  void _showValidationFeedback(bool success, String message) {
    ScaffoldMessenger.of(context).showSnackBar(
      SnackBar(
        content: Text(message),
        backgroundColor: success ? Colors.green : Colors.red,
        duration: const Duration(seconds: 2),
      ),
    );
  }

  void _showErrorDialog(ValidationResult result) {
    showDialog(
      context: context,
      builder: (context) => AlertDialog(
        title: Row(
          children: [
            Icon(Icons.error, color: Colors.red[700]),
            const SizedBox(width: 8),
            const Text('Mismatch!'),
          ],
        ),
        content: Column(
          mainAxisSize: MainAxisSize.min,
          crossAxisAlignment: CrossAxisAlignment.start,
          children: [
            Text('Expected: ${result.expectedItem?.itemLabel}'),
            const SizedBox(height: 8),
            Text('Scanned barcode: ${result.validation?.scannedBarcode}'),
            const SizedBox(height: 16),
            Text(
              result.status == 'OUT_OF_SEQUENCE'
                  ? 'This item is out of sequence.'
                  : 'This barcode does not match the expected item.',
              style: const TextStyle(fontWeight: FontWeight.bold),
            ),
          ],
        ),
        actions: [
          TextButton(
            onPressed: () => Navigator.pop(context),
            child: const Text('Cancel'),
          ),
          if (result.status == 'OUT_OF_SEQUENCE')
            ElevatedButton(
              onPressed: () {
                context.read<LoadingProvider>().manualConfirm(
                  notes: 'Out of sequence but accepted',
                );
                Navigator.pop(context);
              },
              child: const Text('Load Anyway'),
            ),
        ],
      ),
    );
  }

  void _manualConfirm(LoadingProvider provider) {
    final controller = TextEditingController();
    showDialog(
      context: context,
      builder: (context) => AlertDialog(
        title: const Text('Manual Confirmation'),
        content: TextField(
          controller: controller,
          decoration: const InputDecoration(hintText: 'Notes (optional)'),
        ),
        actions: [
          TextButton(
            onPressed: () => Navigator.pop(context),
            child: const Text('Cancel'),
          ),
          ElevatedButton(
            onPressed: () {
              provider.manualConfirm(
                notes: controller.text.isEmpty ? null : controller.text,
              );
              Navigator.pop(context);
            },
            child: const Text('Confirm'),
          ),
        ],
      ),
    );
  }

  void _skipItem(LoadingProvider provider) {
    final controller = TextEditingController();
    showDialog(
      context: context,
      builder: (context) => AlertDialog(
        title: const Text('Skip Item'),
        content: TextField(
          controller: controller,
          decoration: const InputDecoration(hintText: 'Reason for skipping'),
        ),
        actions: [
          TextButton(
            onPressed: () => Navigator.pop(context),
            child: const Text('Cancel'),
          ),
          ElevatedButton(
            onPressed: () {
              provider.skipItem(
                reason: controller.text.isEmpty ? null : controller.text,
              );
              Navigator.pop(context);
            },
            child: const Text('Skip'),
          ),
        ],
      ),
    );
  }

  Future<void> _completeSession(LoadingProvider provider) async {
    final confirmed = await showDialog<bool>(
      context: context,
      builder: (context) => AlertDialog(
        title: const Text('Complete Session?'),
        content: const Text(
          'This will finish the loading session and sync data to the server.',
        ),
        actions: [
          TextButton(
            onPressed: () => Navigator.pop(context, false),
            child: const Text('Cancel'),
          ),
          ElevatedButton(
            onPressed: () => Navigator.pop(context, true),
            child: const Text('Complete'),
          ),
        ],
      ),
    );

    if (confirmed == true) {
      await provider.completeSession();
      if (mounted) Navigator.pop(context);
    }
  }
}
