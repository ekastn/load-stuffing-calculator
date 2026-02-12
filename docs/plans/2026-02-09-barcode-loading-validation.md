# Barcode Loading Validation Feature

**Date:** 2026-02-09  
**Status:** Design Review  
**Author:** System Design  
**Type:** Feature Addition

## Overview

Add barcode scanning capability to the mobile client for warehouse operators to validate items as they load them into containers, ensuring correct items are loaded in the correct sequence according to the calculated 3D plan.

### Goals

- Enable warehouse operators to scan items during loading
- Validate scanned items match the expected loading sequence
- Show real-time 3D visualization of current loading position
- Provide offline-capable validation (works without network)
- Track loading progress and completion status
- Zero database migrations required (use on-the-fly barcode generation)

### Non-Goals

- Weight validation (future feature)
- Real-time IoT sensor integration (future feature)
- Multi-operator collaborative loading (future feature)
- Photo capture of loaded items (future feature)

## Problem Statement

**Current state:**
- Plans are calculated with optimal 3D placement and loading sequence
- No mechanism to verify correct items are loaded in correct order
- Operators manually follow printed instructions (error-prone)
- No tracking of loading progress or completion
- No audit trail of what was actually loaded

**Desired state:**
- Operators scan QR codes on items to validate against plan
- Mobile app shows 3D position of current item to load
- Real-time progress tracking (5/20 items loaded)
- System prevents loading wrong items or out-of-sequence
- Complete audit trail stored locally and optionally synced to backend

## Design Decisions

### 1. On-the-Fly Barcode Generation (No DB Changes)

**Decision:** Generate QR codes deterministically based on `plan_id + step_number + item_id` instead of storing barcodes in database.

**Rationale:**
- Zero database migrations required
- Barcodes are reproducible (same plan always generates same codes)
- Simpler implementation (no barcode management UI needed)
- Works with existing schema

**Barcode Format:**
```
PLAN-{plan_id_short}-STEP-{step_number}-{item_id_short}

Example: PLAN-a3f2e8b1-STEP-001-c4d9f2a3
         ↑            ↑        ↑
         First 8 chars  3-digit  First 8 chars
         of plan UUID   step #   of item UUID
```

**Tradeoffs:**
- ✅ No database changes
- ✅ No barcode data management
- ✅ Always consistent
- ❌ Cannot customize barcodes per deployment
- ❌ Longer QR codes (but QR handles this well)

### 2. Client-Side Validation (Offline-First)

**Decision:** Mobile app performs validation locally without requiring server connection.

**Rationale:**
- Warehouse WiFi often unreliable
- Validation logic is simple (string comparison)
- Better UX (instant feedback)
- Works completely offline

**Validation Flow:**
1. Parse scanned QR code
2. Extract `plan_id_short`, `step_number`, `item_id_short`
3. Compare with expected values (from cached plan data)
4. Return match/mismatch result
5. Optionally sync validation log to backend when online

**Tradeoffs:**
- ✅ Works offline
- ✅ Instant validation feedback
- ✅ No backend load
- ❌ Requires plan data cached locally
- ❌ Validation log might be lost if app uninstalled before sync

### 3. WebView 3D Viewer Integration

**Decision:** Embed existing web 3D viewer in mobile app, passing `step` parameter to highlight current item.

**Rationale:**
- Reuses existing StuffingVisualizer implementation
- Single source of truth for 3D rendering
- No need to port Three.js to Flutter
- Web viewer already supports step-by-step animation

**Implementation:**
```dart
final embedUrl = '$webBaseUrl/embed/shipments/${planId}?token=$token&step=$currentStep&highlight=true';
```

**Tradeoffs:**
- ✅ Reuses existing code
- ✅ Consistent 3D visualization
- ✅ Automatic updates when web improves
- ❌ Requires internet for 3D (but validation still works offline)
- ❌ Slower than native 3D

### 4. Minimal Backend API

**Decision:** Add two RESTful endpoints that extend the existing plans resource.

**Endpoints:**
- `GET /api/v1/plans/:id/barcodes` - Generate QR codes for every step (sub-resource of plans)
- `POST /api/v1/plans/:id/validations` - Validate a scanned barcode and emit a validation record

**Rationale:**
- Frontend already has `GET /plans/:id` for plan data
- Client still manages loading state locally
- Barcode validation is deterministic and can be done in two places (client and server)

**Tradeoffs:**
- ✅ RESTful, consistent with `/plans/:id/items`
- ✅ Minimal backend scope (1 handler + helper functions)
- ✅ Enables both QR sheet generation and optional audit logging
- ❌ No persistence of entire loading sessions (clients still responsible)

## Architecture

### System Components

```
┌─────────────────────────────────────────────────────────────────┐
│ Preparation Phase (Web Client)                                  │
├─────────────────────────────────────────────────────────────────┤
│                                                                  │
│  1. Planner creates plan → Calculate placements                 │
│         ↓                                                        │
│  2. Navigate to "Barcodes" page                                 │
│         ↓                                                        │
│  3. GET /plans/:id/barcodes                                      │
│     Backend generates QR codes for each placement               │
│         ↓                                                        │
│  4. Display QR codes in printable grid                          │
│     (Print labels OR display on tablet for operators)           │
│                                                                  │
└─────────────────────────────────────────────────────────────────┘

┌─────────────────────────────────────────────────────────────────┐
│ Loading Phase (Mobile Client)                                   │
├─────────────────────────────────────────────────────────────────┤
│                                                                  │
│  1. Operator selects plan from list                             │
│         ↓                                                        │
│  2. Start loading session (local state)                         │
│     - Fetch plan data (GET /plans/:id)                          │
│     - Cache placements sorted by step_number                    │
│     - Initialize progress tracker                               │
│         ↓                                                        │
│  3. Show loading screen:                                        │
│     ┌──────────────────────────────────┐                       │
│     │ Progress: 5/20 items (25%)       │                       │
│     │ Next: Box A (500×300×200mm)      │                       │
│     │ Position: (100, 200, 0)          │                       │
│     ├──────────────────────────────────┤                       │
│     │ [3D WebView - highlights item]   │                       │
│     ├──────────────────────────────────┤                       │
│     │ [Scan QR Code]  [Manual OK]      │                       │
│     └──────────────────────────────────┘                       │
│         ↓                                                        │
│  4. Operator scans QR code on physical item                     │
│         ↓                                                        │
│  5. Validate (100% client-side):                                │
│     - Parse QR: PLAN-a3f2-STEP-001-c4d9                        │
│     - Check plan_id matches ✓                                   │
│     - Check step_number = 1 (expected) ✓                       │
│     - Check item_id matches ✓                                   │
│         ↓                                                        │
│  6a. ✓ MATCH → Success feedback                                │
│      - Show green checkmark + item name                         │
│      - Move to next step (step 2)                               │
│      - Update 3D viewer (highlight next item)                   │
│         ↓                                                        │
│  6b. ✗ MISMATCH → Error dialog                                 │
│      - Show expected vs scanned                                 │
│      - Options: Re-scan / Skip / Load Anyway                    │
│         ↓                                                        │
│  7. Repeat until all items validated                            │
│         ↓                                                        │
│  8. Complete session                                            │
│     - Show completion summary                                   │
│     - Optionally sync to backend (PUT /plans/:id)               │
│                                                                  │
└─────────────────────────────────────────────────────────────────┘
```

### Data Flow

```
Mobile Client                 Backend API              Database
─────────────                 ───────────              ────────

[Start Session]
     │
     ├─GET /plans/:id────────►[PlanHandler]──────────►[load_plans]
     │                        [Get placements]        [plan_placements]
     │◄──plan + items──────────                       [load_items]
     │
     └─[Cache locally]

[Generate Barcodes - Web Only]
     │
     ├─GET /loading/barcodes/:id►[LoadingHandler]───►[load_plans]
     │                            [Generate QR codes] [plan_placements]
     │◄──barcode list──────────────

[Scan & Validate]
     │
     ├─[LOCAL VALIDATION]────────►[No server call]
     │  Parse QR code
     │  Compare with cached data
     │  Return match/mismatch
     │
     └─[Store in local state]

[Complete Session]
     │
     └─PUT /plans/:id────────────►[PlanHandler]──────►[load_plans]
       {status: "COMPLETED",      [Update status]      UPDATE status
        loading_notes: {...}}     [Optional]
```

## Components

### 1. Backend API

#### Updated Handler: PlanHandler

**File:** `internal/handler/plan_handler.go` (add methods to existing handler)

**New Methods:**
- `GetPlanBarcodes(c *gin.Context)` - Generate QR codes for plan
- `ValidatePlanBarcode(c *gin.Context)` - Server-side barcode validation

**Key Logic:**
```go
// GetPlanBarcodes returns generated barcodes for all placements in a plan
func (h *PlanHandler) GetPlanBarcodes(c *gin.Context) {
    planID := c.Param("id")
    
    // Get plan with placements (reuse existing method)
    plan, err := h.planService.GetPlanDetail(c.Request.Context(), planID)
    if err != nil {
        response.Error(c, http.StatusNotFound, "Plan not found")
        return
    }
    
    // Generate barcodes
    barcodes := make([]dto.BarcodeInfo, 0, len(plan.Placements))
    
    for _, placement := range plan.Placements {
        if placement.StepNumber == nil {
            continue // Skip placements without step numbers
        }
        
        item := findItemByID(plan.Items, placement.ItemID)
        if item == nil {
            continue
        }
        
        barcode := generateBarcode(plan.PlanID, *placement.StepNumber, placement.ItemID)
        
        barcodes = append(barcodes, dto.BarcodeInfo{
            StepNumber: *placement.StepNumber,
            ItemID:     placement.ItemID.String(),
            ItemLabel:  item.ItemLabel,
            Barcode:    barcode,
            Position: dto.Position{
                X: placement.PosX,
                Y: placement.PosY,
                Z: placement.PosZ,
            },
            Dimensions: dto.Dimensions{
                Length: item.LengthMM,
                Width:  item.WidthMM,
                Height: item.HeightMM,
            },
        })
    }
    
    // Sort by step number
    sort.Slice(barcodes, func(i, j int) bool {
        return barcodes[i].StepNumber < barcodes[j].StepNumber
    })
    
    response.Success(c, barcodes)
}

// ValidatePlanBarcode validates a scanned barcode against a plan
func (h *PlanHandler) ValidatePlanBarcode(c *gin.Context) {
    planID := c.Param("id")
    
    var req dto.ValidateBarcodeRequest
    if err := c.ShouldBindJSON(&req); err != nil {
        response.Error(c, http.StatusBadRequest, err.Error())
        return
    }
    
    // Parse scanned barcode
    parsed := parseBarcode(req.Barcode)
    if parsed == nil {
        response.Success(c, dto.ValidationResult{
            Valid:  false,
            Status: "INVALID_FORMAT",
            Error:  "Invalid barcode format",
        })
        return
    }
    
    // Verify barcode belongs to this plan
    planUUID, err := uuid.Parse(planID)
    if err != nil {
        response.Error(c, http.StatusBadRequest, "Invalid plan ID")
        return
    }
    
    planShort := planUUID.String()[:8]
    if parsed.PlanID != planShort {
        response.Success(c, dto.ValidationResult{
            Valid:      false,
            Status:     "WRONG_PLAN",
            Error:      "Barcode belongs to different plan",
            PlanID:     parsed.PlanID,
            StepNumber: parsed.StepNumber,
            ItemID:     parsed.ItemID,
        })
        return
    }
    
    // Check if step matches expected
    valid := parsed.StepNumber == req.ExpectedStep
    status := "MATCHED"
    if !valid {
        if req.ExpectedStep != nil {
            status = "OUT_OF_SEQUENCE"
        } else {
            status = "UNEXPECTED_STEP"
        }
    }
    
    response.Success(c, dto.ValidationResult{
        Valid:      valid,
        Status:     status,
        PlanID:     parsed.PlanID,
        StepNumber: parsed.StepNumber,
        ItemID:     parsed.ItemID,
        Barcode:    req.Barcode,
    })
}

// Helper: Generate deterministic barcode
func generateBarcode(planID uuid.UUID, stepNumber int, itemID uuid.UUID) string {
    planShort := planID.String()[:8]
    itemShort := itemID.String()[:8]
    return fmt.Sprintf("PLAN-%s-STEP-%03d-%s", planShort, stepNumber, itemShort)
}

// Helper: Parse scanned barcode
func parseBarcode(barcode string) *ParsedBarcode {
    parts := strings.Split(barcode, "-")
    if len(parts) != 5 || parts[0] != "PLAN" || parts[2] != "STEP" {
        return nil
    }
    
    stepNum, err := strconv.Atoi(parts[3])
    if err != nil {
        return nil
    }
    
    return &ParsedBarcode{
        PlanID:     parts[1],
        StepNumber: stepNum,
        ItemID:     parts[4],
    }
}

// Helper: Find item by ID
func findItemByID(items []dto.ItemDetail, itemID uuid.UUID) *dto.ItemDetail {
    for i := range items {
        if items[i].ItemID == itemID {
            return &items[i]
        }
    }
    return nil
}

// Helper struct
type ParsedBarcode struct {
    PlanID     string
    StepNumber int
    ItemID     string
}
```

**DTOs:**

**File:** `internal/dto/plan.go` (add to existing file)

```go
type BarcodeInfo struct {
    StepNumber int        `json:"step_number"`
    ItemID     string     `json:"item_id"`
    ItemLabel  string     `json:"item_label"`
    Barcode    string     `json:"barcode"`
    Position   Position   `json:"position"`
    Dimensions Dimensions `json:"dimensions"`
}

type Position struct {
    X float64 `json:"x"`
    Y float64 `json:"y"`
    Z float64 `json:"z"`
}

type Dimensions struct {
    Length int `json:"length"`
    Width  int `json:"width"`
    Height int `json:"height"`
}

type ValidateBarcodeRequest struct {
    Barcode      string `json:"barcode" binding:"required"`
    ExpectedStep *int   `json:"expected_step"` // Optional
}

type ValidationResult struct {
    Valid      bool   `json:"valid"`
    Status     string `json:"status"` // MATCHED, OUT_OF_SEQUENCE, WRONG_PLAN, INVALID_FORMAT
    PlanID     string `json:"plan_id,omitempty"`
    StepNumber int    `json:"step_number,omitempty"`
    ItemID     string `json:"item_id,omitempty"`
    Barcode    string `json:"barcode,omitempty"`
    Error      string `json:"error,omitempty"`
}
```

#### Updated Routes (RESTful Design)

**File:** `internal/api/routes.go`

```go
// Plans group (authenticated) - add to existing plans group
plans := api.Group("/plans")
plans.Use(middleware.JWTAuth())
{
    // Existing routes...
    plans.GET("", middleware.RequirePermission("plan:read"), handlers.Plan.ListPlans)
    plans.POST("", middleware.RequirePermission("plan:create"), handlers.Plan.CreatePlan)
    plans.GET("/:id", middleware.RequirePermission("plan:read"), handlers.Plan.GetPlan)
    plans.PUT("/:id", middleware.RequirePermission("plan:update"), handlers.Plan.UpdatePlan)
    plans.DELETE("/:id", middleware.RequirePermission("plan:delete"), handlers.Plan.DeletePlan)
    plans.POST("/:id/calculate", middleware.RequirePermission("plan:update"), handlers.Plan.CalculatePlan)
    
    // NEW: Barcode endpoints (sub-resources of plans)
    plans.GET("/:id/barcodes", middleware.RequirePermission("plan:read"), handlers.Plan.GetPlanBarcodes)
    plans.POST("/:id/validations", middleware.RequirePermission("plan:read"), handlers.Plan.ValidatePlanBarcode)
    
    // Existing items sub-resource...
    plans.POST("/:id/items", middleware.RequirePermission("plan:update"), handlers.Plan.AddItem)
    plans.PUT("/:id/items/:itemId", middleware.RequirePermission("plan:update"), handlers.Plan.UpdateItem)
    plans.DELETE("/:id/items/:itemId", middleware.RequirePermission("plan:update"), handlers.Plan.DeleteItem)
}
```

#### API Endpoints (RESTful)

**GET /api/v1/plans/:id/barcodes**

Generate QR codes for all placements in a plan.

Request:
```http
GET /api/v1/plans/a3f2e8b1-c4d9-4f2a-b1c8-3d9e8f7a6b5c/barcodes
Authorization: Bearer <token>
```

Response:
```json
{
  "data": [
    {
      "step_number": 1,
      "item_id": "c4d9f2a3-8b1c-4e5f-9a2b-6d7e8f9a0b1c",
      "item_label": "Box A",
      "barcode": "PLAN-a3f2e8b1-STEP-001-c4d9f2a3",
      "position": {
        "x": 100,
        "y": 200,
        "z": 0
      },
      "dimensions": {
        "length": 500,
        "width": 300,
        "height": 200
      }
    },
    {
      "step_number": 2,
      "item_id": "d5e0f3b4-9c2d-5f6a-0b3c-7e8f9a1b2c3d",
      "item_label": "Box B",
      "barcode": "PLAN-a3f2e8b1-STEP-002-d5e0f3b4",
      "position": {
        "x": 600,
        "y": 200,
        "z": 0
      },
      "dimensions": {
        "length": 400,
        "width": 400,
        "height": 300
      }
    }
  ]
}
```

**POST /api/v1/plans/:id/validations**

Validate a scanned barcode against the plan (creates validation result).

Request:
```http
POST /api/v1/plans/a3f2e8b1-c4d9-4f2a-b1c8-3d9e8f7a6b5c/validations
Authorization: Bearer <token>
Content-Type: application/json

{
  "barcode": "PLAN-a3f2e8b1-STEP-001-c4d9f2a3",
  "expected_step": 1
}
```

Response (Success - Matched):
```json
{
  "data": {
    "valid": true,
    "status": "MATCHED",
    "plan_id": "a3f2e8b1",
    "step_number": 1,
    "item_id": "c4d9f2a3",
    "barcode": "PLAN-a3f2e8b1-STEP-001-c4d9f2a3"
  }
}
```

Response (Out of Sequence):
```json
{
  "data": {
    "valid": false,
    "status": "OUT_OF_SEQUENCE",
    "plan_id": "a3f2e8b1",
    "step_number": 5,
    "item_id": "e6f1g4c5",
    "barcode": "PLAN-a3f2e8b1-STEP-005-e6f1g4c5"
  }
}
```

Response (Wrong Plan):
```json
{
  "data": {
    "valid": false,
    "status": "WRONG_PLAN",
    "plan_id": "b4g3h9c2",
    "step_number": 1,
    "item_id": "c4d9f2a3",
    "barcode": "PLAN-b4g3h9c2-STEP-001-c4d9f2a3",
    "error": "Barcode belongs to different plan"
  }
}
```

Response (Invalid Format):
```json
{
  "data": {
    "valid": false,
    "status": "INVALID_FORMAT",
    "error": "Invalid barcode format"
  }
}
```

### 2. Web Client

#### Barcode Display Page

**File:** `web/app/(app)/shipments/[id]/barcodes/page.tsx`

**Features:**
- Fetch barcodes from backend
- Generate QR code images using `qrcode` library
- Display in printable grid (3-4 columns)
- Print button for physical labels
- Show item details below each QR code

**Dependencies:**
- `npm install qrcode` - QR code generation
- `npm install @types/qrcode --save-dev`

**Layout:**
```
┌────────────────────────────────────────────────┐
│ Loading Barcodes           [Print Labels]      │
├────────────────────────────────────────────────┤
│  ┌──────────┐  ┌──────────┐  ┌──────────┐    │
│  │ Step 1   │  │ Step 2   │  │ Step 3   │    │
│  │ [QR Code]│  │ [QR Code]│  │ [QR Code]│    │
│  │ Box A    │  │ Box B    │  │ Box C    │    │
│  │ 500×300× │  │ 400×400× │  │ 300×300× │    │
│  │  200mm   │  │  300mm   │  │  400mm   │    │
│  └──────────┘  └──────────┘  └──────────┘    │
│                                                 │
│  ┌──────────┐  ┌──────────┐  ┌──────────┐    │
│  │ Step 4   │  │ Step 5   │  │ Step 6   │    │
│  │   ...    │  │   ...    │  │   ...    │    │
│  └──────────┘  └──────────┘  └──────────┘    │
└────────────────────────────────────────────────┘
```

#### Service Layer

**File:** `web/lib/services/plans.ts` (add to existing file)

```typescript
export interface BarcodeInfo {
  step_number: number;
  item_id: string;
  item_label: string;
  barcode: string;
  position: {
    x: number;
    y: number;
    z: number;
  };
  dimensions: {
    length: number;
    width: number;
    height: number;
  };
}

export interface ValidationResult {
  valid: boolean;
  status: string; // MATCHED, OUT_OF_SEQUENCE, WRONG_PLAN, INVALID_FORMAT
  plan_id?: string;
  step_number?: number;
  item_id?: string;
  barcode?: string;
  error?: string;
}

export interface ValidateBarcodeRequest {
  barcode: string;
  expected_step?: number;
}

export const getPlanBarcodes = async (planId: string): Promise<BarcodeInfo[]> => {
  const response = await apiClient.get(`/plans/${planId}/barcodes`);
  return response.data;
};

export const validatePlanBarcode = async (
  planId: string,
  request: ValidateBarcodeRequest
): Promise<ValidationResult> => {
  const response = await apiClient.post(`/plans/${planId}/validations`, request);
  return response.data;
};
```

#### Navigation Integration

Add "View Barcodes" button to shipment detail page:

**File:** `web/app/(app)/shipments/[id]/page.tsx`

```tsx
<Button onClick={() => router.push(`/shipments/${id}/barcodes`)}>
  <QrCode className="mr-2 h-4 w-4" />
  View Barcodes
</Button>
```

### 3. Mobile Client

#### Loading Session Page

**File:** `mobile/lib/pages/loading/loading_session_page.dart`

**Layout:**
```
┌──────────────────────────────────────┐
│ ← Loading - Step 5/20                │
├──────────────────────────────────────┤
│ ▓▓▓▓▓░░░░░░░░░░░░░░░  25%           │
│ 5/20 items loaded                    │
├──────────────────────────────────────┤
│ Next Item:                            │
│ Box A                                 │
│ Dimensions: 500 × 300 × 200 mm       │
│ Position: (100, 200, 0)              │
│ Expected QR: PLAN-a3f2-STEP-005-...  │
├──────────────────────────────────────┤
│                                       │
│    [3D WebView Viewer]               │
│    (shows container with             │
│     current item highlighted)        │
│                                       │
├──────────────────────────────────────┤
│ ┌───────────────────────────────────┐│
│ │  📷 Scan QR Code                  ││
│ └───────────────────────────────────┘│
│                                       │
│ ┌──────────────┐  ┌─────────────────┐│
│ │  Manual OK   │  │      Skip       ││
│ └──────────────┘  └─────────────────┘│
└──────────────────────────────────────┘
```

**Features:**
- Progress bar showing items loaded
- Current item details
- Embedded 3D viewer (WebView with step parameter)
- Scan QR button (opens camera scanner)
- Manual confirmation button (fallback)
- Skip button (with reason prompt)

#### Barcode Scanner Page

**File:** `mobile/lib/pages/loading/barcode_scanner_page.dart`

**Features:**
- Full-screen camera view using `mobile_scanner`
- Visual scan area indicator (square overlay)
- Auto-detect QR codes
- Haptic feedback on successful scan
- Manual entry fallback (if camera fails)

**Layout:**
```
┌──────────────────────────────────────┐
│ ← Scan Barcode                       │
├──────────────────────────────────────┤
│                                       │
│    ┌─────────────────────┐          │
│    │                     │          │
│    │   [Camera View]     │          │
│    │                     │          │
│    │     ┌─────────┐     │          │
│    │     │         │     │          │
│    │     │  Scan   │     │          │
│    │     │  Area   │     │          │
│    │     └─────────┘     │          │
│    │                     │          │
│    └─────────────────────┘          │
│                                       │
│  Point camera at QR code              │
│                                       │
│  [ Enter Code Manually ]             │
│                                       │
└──────────────────────────────────────┘
```

#### Loading Provider

**File:** `mobile/lib/providers/loading_provider.dart`

**Responsibilities:**
- Manage loading session state
- Fetch and cache plan data
- Generate barcodes (same logic as backend)
- Validate scanned barcodes (client-side)
- Track progress (validated, skipped counts)
- Persist session to local storage
- Optional: Sync completion to backend

**Key Methods:**
```dart
Future<void> startSession(String planId)
Future<void> resumeSession()
ExpectedItem? getCurrentExpectedItem()
ValidationResult validateBarcode(String scannedBarcode)
void manualConfirm({String? notes})
void skipItem({String? reason})
Future<void> completeSession({bool syncToBackend = true})
String generateBarcode(String planId, int stepNumber, String itemId)
ParsedBarcode? parseBarcode(String barcode)
```

**State:**
```dart
LoadingSession? _currentSession;
Plan? _currentPlan;
List<PlanPlacement>? _placements;  // Sorted by step_number
```

#### Models

**File:** `mobile/lib/models/loading_session.dart`

```dart
@freezed
class LoadingSession with _$LoadingSession {
  const factory LoadingSession({
    required String sessionId,        // Local UUID
    required String planId,
    required DateTime startedAt,
    DateTime? completedAt,
    
    required int totalItems,
    @Default(0) int validatedCount,
    @Default(0) int skippedCount,
    @Default(0) int currentStepIndex,
    
    @Default([]) List<ItemValidation> validations,
    @Default('IN_PROGRESS') String status,
  }) = _LoadingSession;
}

@freezed
class ItemValidation with _$ItemValidation {
  const factory ItemValidation({
    required String itemId,
    required String itemLabel,
    required int stepNumber,
    
    String? scannedBarcode,
    String? expectedBarcode,
    
    required String status,  // MATCHED, MISMATCHED, MANUAL_CONFIRMED, SKIPPED, OUT_OF_SEQUENCE
    required DateTime validatedAt,
    
    String? notes,
  }) = _ItemValidation;
}
```

## Implementation Plan

### Phase 1: Backend API (2-3 hours)

**Files to modify:**
- [x] `internal/handler/plan_handler.go` - Add two new methods
- [x] `internal/handler/plan_handler_test.go` - Add tests for new methods
- [x] `internal/dto/plan_dto.go` - Add new DTOs (implemented in `plan_dto.go`)
- [x] `internal/api/routes.go` - Add two new routes to plans group

**Tasks:**
1. [x] Add DTOs to `internal/dto/plan_dto.go`
2. [x] Add methods to `PlanHandler`
3. [x] Add helper functions to `plan_handler.go`
4. [x] Update routes in `internal/api/routes.go`
5. [x] Write handler tests (95%+ coverage)
6. [x] Test with Postman/curl

**Acceptance Criteria:**
- `GET /api/v1/plans/:id/barcodes` returns sorted barcode list
- `POST /api/v1/plans/:id/validations` validates barcodes correctly
- Barcodes are deterministic (same plan = same barcodes)
- Both endpoints require authentication and `plan:read` permission
- Workspace scoping works correctly (inherited from GetPlanDetail)
- 404 for non-existent plans
- Proper error responses for invalid barcodes
- Tests pass with 95%+ coverage

### Phase 2: Web Client - Barcode Display (2-3 hours)

**Files to create:**
- [x] `web/app/(app)/shipments/[id]/barcodes/page.tsx`

**Files to modify:**
- [x] `web/lib/services/plans.ts` - Add `getPlanBarcodes()` and `validatePlanBarcode()`
- [x] `web/lib/types/index.ts` - Add types for BarcodeInfo and ValidationResult
- [x] `web/app/(app)/shipments/[id]/page.tsx` - Add "View Barcodes" button
- [x] `web/package.json` - Add qrcode dependency

**Tasks:**
1. [x] Install `qrcode` library: `pnpm add qrcode @types/qrcode`
2. [x] Add types to `web/lib/types/index.ts`:
   - `BarcodeInfo` interface
   - `ValidationResult` interface
   - `ValidateBarcodeRequest` interface
3. [x] Add methods to `web/lib/services/plans.ts`:
   - `getPlanBarcodes()` - Fetch barcodes
   - `validatePlanBarcode()` - Server-side validation (optional use)
4. [x] Create barcode display page component
5. [x] Fetch barcodes from API
6. [x] Generate QR code images using qrcode library
7. [x] Display in printable grid (3 columns)
8. [x] Add print button with print CSS
9. [x] Add navigation from shipment detail page
10. [x] Manual testing in browser (see Testing Strategy section)

**Acceptance Criteria:**
- Barcode page loads and displays QR codes
- QR codes are scannable (test with phone camera)
- Print layout is clean and readable
- Each QR includes item label, step number, and dimensions
- Page handles loading and error states
- No console errors

### Phase 3: Mobile Client - Models & Provider (3-4 hours)

**Files to create:**
- [x] `mobile/lib/models/loading_session.dart`
- [x] `mobile/lib/models/loading_session.freezed.dart` (generated)
- [x] `mobile/lib/models/loading_session.g.dart` (generated)
- [x] `mobile/lib/providers/loading_provider.dart`

**Tasks:**
1. [x] Define `LoadingSession` model with Freezed
2. [x] Define `ItemValidation` model with Freezed
3. [x] Run code generation: `flutter pub run build_runner build`
4. [x] Create `LoadingProvider` class extending `ChangeNotifier`
5. [x] Implement `startSession()` - fetch plan and initialize
6. [x] Implement `resumeSession()` - restore from storage
7. [x] Implement `getCurrentExpectedItem()` - get next item
8. [x] Implement `generateBarcode()` - match backend logic
9. [x] Implement `parseBarcode()` - parse scanned QR
10. [x] Implement `validateBarcode()` - client-side validation
11. [x] Implement `manualConfirm()` - bypass scanning
12. [x] Implement `skipItem()` - skip with reason
13. [x] Implement `completeSession()` - finish and sync
14. [x] Add to provider setup in `main.dart`

**Acceptance Criteria:**
- Models generated successfully
- Provider compiles without errors
- Barcode generation matches backend format
- Barcode parsing handles invalid formats
- Validation logic correctly identifies matches/mismatches
- Session persists to storage
- Provider notifies listeners on state changes

### Phase 4: Mobile Client - Barcode Scanner (2-3 hours)

**Files to create:**
- [x] `mobile/lib/pages/loading/barcode_scanner_page.dart`

**Tasks:**
1. [x] Create `BarcodeScannerPage` stateful widget
2. [x] Initialize `MobileScannerController`
3. [x] Implement camera view with `MobileScanner` widget
4. [x] Add scan area overlay (visual guide)
5. [x] Implement `onDetect` callback
6. [x] Add haptic feedback on scan
7. [x] Return scanned barcode to previous page
8. [x] Add manual entry fallback (TextField)
9. [x] Handle camera permissions
10. [x] Handle "no camera" case
11. [x] Proper disposal of controller

**Acceptance Criteria:**
- Scanner opens camera successfully
- QR codes are detected and decoded
- Haptic feedback works on scan
- Scanner closes and returns barcode to parent
- Manual entry works as fallback
- Camera permission prompt appears if needed
- No memory leaks (controller disposed)

### Phase 5: Mobile Client - Loading Session UI (4-5 hours)

**Files to create:**
- [x] `mobile/lib/pages/loading/loading_session_page.dart`

**Files to modify:**
- [x] `mobile/lib/pages/plans/plan_detail_page.dart` - Add "Start Loading" button
- [x] `mobile/lib/config/routes.dart` - Add loading route

**Tasks:**
1. [x] Create `LoadingSessionPage` stateful widget
2. [x] Build progress bar widget
3. [x] Build current item info card
4. [x] Integrate 3D WebView with step parameter
5. [x] Add scan button (opens `BarcodeScannerPage`)
6. [x] Handle barcode validation result
7. [x] Show success dialog on match
8. [x] Show error dialog on mismatch
9. [x] Implement manual confirm dialog
10. [x] Implement skip dialog with reason input
11. [x] Add complete session button (when all items done)
12. [x] Wire up navigation from plan detail page
13. [x] Add loading session route to router

**Acceptance Criteria:**
- Loading session starts and shows first item
- Progress bar updates correctly
- 3D viewer loads and shows current step
- Scan button opens camera
- Validation success shows checkmark
- Validation failure shows error with options
- Manual confirm works
- Skip with reason works
- Session completes successfully
- Navigation works end-to-end

### Phase 6: Backend Testing & Documentation (2-3 hours)

**Backend API Tests:**
- [x] Handler tests for `GetPlanBarcodes` endpoint
- [x] Handler tests for `ValidatePlanBarcode` endpoint
- [x] Unit tests for `generateBarcode()` helper
- [x] Unit tests for `parseBarcode()` helper
- [x] Edge case: Plan with no placements
- [x] Edge case: Invalid plan ID
- [x] Edge case: Workspace scoping (user can't access other workspace's plans)
- [x] Edge case: Invalid barcode format
- [x] Edge case: Wrong plan ID in barcode
- [x] Edge case: Out of sequence step
- [x] Test with mocked SQLC store
- [x] Run `make test` to verify coverage
- [x] Aim for 95%+ coverage on new handler methods

**API Manual Testing (Postman/curl):**
- [x] GET `/plans/:id/barcodes` - Happy path
- [x] GET `/plans/:id/barcodes` - Plan not found (404)
- [x] GET `/plans/:id/barcodes` - Unauthorized (401)
- [x] POST `/plans/:id/validations` - Matched barcode
- [x] POST `/plans/:id/validations` - Mismatched barcode
- [x] POST `/plans/:id/validations` - Out of sequence
- [x] POST `/plans/:id/validations` - Wrong plan
- [x] POST `/plans/:id/validations` - Invalid format

**Web Client Manual Testing:**
- [x] See "Manual Testing Checklist" in Testing Strategy section below
- [x] Focus on QR code display, print functionality, browser compatibility

**Mobile Client Manual Testing:**
- [x] See "Manual Testing Checklist" in Testing Strategy section below
- [x] Full loading flow on Android and iOS
- [x] Offline mode, session persistence, edge cases

**Documentation:**
- [x] Update AGENTS.md with loading validation feature description
- [x] Update README.md with barcode scanning workflow
- [x] Add operator guide: "How to use loading validation"
- [x] Document API endpoints in Swagger (if using swag)

**Estimated time:** 2-3 hours (Backend tests: 1.5-2h, Manual testing: 30m, Docs: 30m)

## Testing Strategy

### Backend Tests

**PlanHandler Tests:**
```go
func TestGetPlanBarcodes(t *testing.T) {
    tests := []struct {
        name       string
        planID     string
        wantStatus int
        wantCount  int
    }{
        {"valid plan with 5 items", "valid-uuid", 200, 5},
        {"plan not found", "invalid-uuid", 404, 0},
        {"plan with no placements", "empty-uuid", 200, 0},
    }
    // ... test implementation
}

func TestValidatePlanBarcode(t *testing.T) {
    tests := []struct {
        name         string
        planID       string
        barcode      string
        expectedStep *int
        wantValid    bool
        wantStatus   string
    }{
        {
            "matched barcode",
            "a3f2e8b1-c4d9-4f2a-b1c8-3d9e8f7a6b5c",
            "PLAN-a3f2e8b1-STEP-001-c4d9f2a3",
            intPtr(1),
            true,
            "MATCHED",
        },
        {
            "out of sequence",
            "a3f2e8b1-c4d9-4f2a-b1c8-3d9e8f7a6b5c",
            "PLAN-a3f2e8b1-STEP-005-e6f1g4c5",
            intPtr(1),
            false,
            "OUT_OF_SEQUENCE",
        },
        {
            "wrong plan",
            "a3f2e8b1-c4d9-4f2a-b1c8-3d9e8f7a6b5c",
            "PLAN-b4g3h9c2-STEP-001-c4d9f2a3",
            intPtr(1),
            false,
            "WRONG_PLAN",
        },
        {
            "invalid format",
            "a3f2e8b1-c4d9-4f2a-b1c8-3d9e8f7a6b5c",
            "INVALID-BARCODE",
            intPtr(1),
            false,
            "INVALID_FORMAT",
        },
    }
    // ... test implementation
}

func TestGenerateBarcode(t *testing.T) {
    planID := uuid.MustParse("a3f2e8b1-c4d9-4f2a-b1c8-3d9e8f7a6b5c")
    itemID := uuid.MustParse("c4d9f2a3-8b1c-4e5f-9a2b-6d7e8f9a0b1c")
    
    barcode := generateBarcode(planID, 1, itemID)
    
    assert.Equal(t, "PLAN-a3f2e8b1-STEP-001-c4d9f2a3", barcode)
    
    // Test determinism
    barcode2 := generateBarcode(planID, 1, itemID)
    assert.Equal(t, barcode, barcode2)
}

func TestParseBarcode(t *testing.T) {
    tests := []struct {
        barcode  string
        wantNil  bool
        wantPlan string
        wantStep int
        wantItem string
    }{
        {
            "PLAN-a3f2e8b1-STEP-001-c4d9f2a3",
            false,
            "a3f2e8b1",
            1,
            "c4d9f2a3",
        },
        {
            "INVALID-FORMAT",
            true,
            "",
            0,
            "",
        },
        {
            "PLAN-a3f2-INVALID-001-c4d9",
            true,
            "",
            0,
            "",
        },
        {
            "PLAN-a3f2e8b1-STEP-ABC-c4d9f2a3",
            true,
            "",
            0,
            "",
        },
    }
    // ... test implementation
}
```

### Manual Testing Only

**Web Client - Manual Testing Checklist:**

1. **Barcode Display Page** (`/shipments/[id]/barcodes`)
   - [ ] Navigate to shipment detail → Click "View Barcodes" button
   - [ ] Verify all QR codes display correctly
   - [ ] Check step numbers are in correct order (1, 2, 3...)
   - [ ] Verify item labels match plan items
   - [ ] Test print functionality (print preview)
   - [ ] Test with small plan (3 items)
   - [ ] Test with large plan (50+ items)
   - [ ] Verify QR codes are scannable with phone camera

2. **Error States**
   - [ ] Plan with no calculated placements (empty barcode list)
   - [ ] Plan not found (404 error handling)
   - [ ] Network error during fetch (error message display)

3. **Browser Compatibility**
   - [ ] Chrome/Edge (desktop)
   - [ ] Safari (desktop)
   - [ ] Firefox (desktop)
   - [ ] Chrome (mobile)

**Mobile Client - Manual Testing Checklist:**

1. **Loading Session Flow**
   - [ ] Select plan from list
   - [ ] Start loading session
   - [ ] Verify progress bar shows 0/N items
   - [ ] Verify current item card displays correct information
   - [ ] Verify 3D WebView loads and shows container

2. **Barcode Scanner**
   - [ ] Open scanner (camera permission prompt)
   - [ ] Scan valid QR code → Success dialog → Next item
   - [ ] Scan invalid QR code → Error dialog
   - [ ] Scan out-of-sequence QR → Warning with "Load Anyway" option
   - [ ] Scanner auto-closes after successful scan
   - [ ] Test flashlight toggle (if available)
   - [ ] Test manual entry fallback

3. **Validation States**
   - [ ] MATCHED: Green success, auto-advance to next item
   - [ ] MISMATCHED: Red error, show expected vs scanned
   - [ ] OUT_OF_SEQUENCE: Yellow warning, option to proceed
   - [ ] WRONG_PLAN: Red error, cannot proceed
   - [ ] INVALID_FORMAT: Red error, invalid QR code

4. **Manual Actions**
   - [ ] Manual Confirm button → Add notes → Item validated
   - [ ] Skip button → Add reason → Item marked skipped
   - [ ] Progress updates correctly after each action

5. **3D Viewer Integration**
   - [ ] WebView loads embed URL with token
   - [ ] Viewer shows correct container
   - [ ] Step parameter highlights current item (if implemented in web)
   - [ ] Viewer remains visible during scanning
   - [ ] Viewer updates when moving to next step

6. **Session Management**
   - [ ] Complete session → Confirmation dialog → Return to plans
   - [ ] Close app mid-session → Reopen → Resume session
   - [ ] Verify session persists across app restarts
   - [ ] Multiple sessions (start second session while first is paused)

7. **Offline Mode**
   - [ ] Enable airplane mode
   - [ ] Start session (should fail - need plan data)
   - [ ] Load plan while online, then go offline
   - [ ] Scan items offline → Validations work
   - [ ] Complete session offline (if not syncing to backend)
   - [ ] Return online → Check if data syncs (if implemented)

8. **Platform Testing**
   - [ ] Android phone (real device)
   - [ ] Android tablet (real device)
   - [ ] iOS phone (real device)
   - [ ] iOS tablet (real device)
   - [ ] Test camera on different devices

9. **Edge Cases**
   - [ ] Plan with 1 item
   - [ ] Plan with 100+ items
   - [ ] Plan with items missing barcodes
   - [ ] All items skipped → Complete session
   - [ ] Mix of scanned, manual confirmed, and skipped items
   - [ ] Low battery mode (camera performance)
   - [ ] Poor lighting conditions (scan difficulty)

## Error Handling

### Backend Errors

| Error | Status | Response | Handling |
|-------|--------|----------|----------|
| Plan not found | 404 | `{"error": "Plan not found"}` | Show error page |
| Unauthorized | 401 | `{"error": "Unauthorized"}` | Redirect to login |
| No permission | 403 | `{"error": "Requires plan:read"}` | Show permission denied |
| Workspace mismatch | 403 | `{"error": "Plan not accessible"}` | Show error |

### Mobile Errors

| Error | Scenario | Handling |
|-------|----------|----------|
| Camera permission denied | User denies camera access | Show settings prompt + manual entry fallback |
| No camera available | Desktop/Linux | Show "Scanner not available" + manual entry |
| Invalid QR format | Scanned non-plan QR | Show "Invalid QR code" dialog |
| Wrong plan | Scanned QR from different plan | Show "QR belongs to different plan" error |
| Out of sequence | Scanned step 5 when expecting step 2 | Show warning + "Load anyway?" option |
| Network timeout | Offline when fetching plan | Show cached data or error |
| Session persistence failure | Storage write fails | Log error, continue (volatile session) |

### Offline Handling

**Mobile app works offline:**
1. Plan data cached after first fetch
2. Validation runs 100% client-side
3. Session state persists to local storage
4. Completion notes stored locally
5. Auto-sync when network restored

**Offline indicators:**
- Show "Offline Mode" banner
- Disable 3D viewer (or show cached version)
- Queue completion sync for later

## Security Considerations

### Backend

1. **Authentication:** All endpoints require JWT token
2. **Authorization:** Require `plan:read` permission
3. **Workspace scoping:** Users can only access plans in their workspace
4. **Input validation:** Validate plan ID format (UUID)
5. **Rate limiting:** Apply rate limits to prevent abuse

### Mobile

1. **Token storage:** Access token stored in `flutter_secure_storage`
2. **Session persistence:** Use secure storage for session state
3. **Barcode validation:** Parse and validate format before processing
4. **Plan verification:** Always verify scanned barcode belongs to current plan

## Performance Considerations

### Backend

- **Barcode generation:** O(n) where n = number of placements
- **Caching:** Consider caching barcode list (same plan = same barcodes)
- **Database:** Existing indexes on `plan_placements` sufficient

### Web

- **QR generation:** Generate images client-side (no server load)
- **Print optimization:** Use CSS print media queries
- **Large plans:** Consider pagination for 100+ items

### Mobile

- **Plan caching:** Cache full plan data locally
- **Validation speed:** Client-side validation is instant (string comparison)
- **WebView performance:** 3D viewer may be slow on low-end devices
- **Memory:** Dispose scanner controller properly

## Future Enhancements

### Phase 2 (Future)

1. **Weight validation:**
   - Integrate weight scale via Bluetooth
   - Compare actual vs expected weight
   - Alert if weight mismatch exceeds tolerance

2. **Photo capture:**
   - Option to capture photo of loaded item
   - Store photos in backend or cloud storage
   - Include in audit trail

3. **Multi-operator loading:**
   - Multiple operators scan different items in parallel
   - Real-time sync of loading progress
   - Prevent duplicate scans

4. **IoT sensor integration:**
   - RFID tag scanning (faster than QR)
   - Weight sensor automation
   - Position sensors for validation

5. **Advanced analytics:**
   - Loading time per item (efficiency metrics)
   - Error rate tracking (mismatches, skips)
   - Operator performance comparison

6. **Offline improvements:**
   - Offline 3D viewer (cache Three.js assets)
   - Background sync queue
   - Conflict resolution for multi-device

7. **Barcode customization:**
   - Support multiple barcode formats (EAN13, Code128, QR)
   - Custom barcode prefixes per workspace
   - Barcode printing integration (Zebra, Brother)

## Success Metrics

### Operational Metrics

- **Loading time reduction:** Measure time from plan creation to completion
- **Error reduction:** Track mismatch/skip rate before vs after
- **Operator satisfaction:** Survey scores

### Technical Metrics

- **Scan success rate:** % of successful scans on first attempt
- **Offline usage:** % of sessions completed offline
- **Performance:** Average validation response time < 100ms
- **Reliability:** Session persistence success rate > 99%

### Adoption Metrics

- **Feature usage:** % of plans that use loading validation
- **Mobile app downloads:** Increase in mobile installs
- **Web barcode page views:** Frequency of barcode printing

## Risks & Mitigations

| Risk | Impact | Probability | Mitigation |
|------|--------|-------------|------------|
| QR codes not scannable (print quality) | High | Medium | Provide high-res QR generation, print quality guidelines |
| Camera not working on some devices | High | Low | Mandatory manual entry fallback |
| Barcode format conflicts with existing systems | Medium | Low | Use unique prefix "PLAN-" to avoid conflicts |
| Users skip barcode and use manual confirm | Medium | High | Make manual confirm require reason/notes |
| 3D viewer slow on low-end phones | Medium | Medium | Provide option to disable 3D, show static image |
| Offline sync conflicts | Low | Low | Simple last-write-wins strategy (no concurrent edits) |

## Dependencies

### Backend
- None (uses existing Gin, SQLC, JWT infrastructure)

### Web
- `qrcode` - QR code image generation
- `@types/qrcode` - TypeScript types

### Mobile
- `mobile_scanner: ^7.1.4` - Already installed ✅
- `flutter_secure_storage: ^10.0.0` - Already installed ✅
- `webview_flutter: ^4.13.1` - Already installed ✅

## Rollout Plan

### Development Environment
1. Implement backend endpoints
2. Test with Postman
3. Implement web barcode page
4. Test QR code printing
5. Implement mobile app
6. Test end-to-end flow

### Staging Environment
1. Deploy backend changes
2. Deploy web changes
3. Deploy mobile app (TestFlight/Internal Testing)
4. Conduct UAT with test operators
5. Fix bugs and iterate

### Production Rollout
1. Deploy backend (backward compatible, no breaking changes)
2. Deploy web updates
3. Release mobile app update
4. Monitor error rates and performance
5. Gather user feedback
6. Iterate based on feedback

### Feature Flag (Optional)
- Add feature flag to enable/disable loading validation
- Gradual rollout to workspaces
- Easy rollback if issues arise

## Open Questions

1. **QR code size:** What physical size should QR codes be printed at? (Recommend: 2×2 inches minimum)

2. **Label material:** Should we recommend specific label printers or materials? (Suggest: waterproof labels for warehouse environment)

3. **Session timeout:** Should loading sessions auto-expire after X hours of inactivity? (Suggest: 24 hours)

4. **Multi-quantity items:** If an item has quantity=5, scan once or five times? (Recommend: Scan once, show "5x Box A" in UI)

5. **Sequence enforcement:** Strictly enforce sequence or allow out-of-order with warning? (Recommend: Warning only, allow override)

6. **Photo capture:** Should we include photo capture in Phase 1 or defer to Phase 2? (Recommend: Phase 2)

7. **Web barcode generation:** Should web client call backend or generate QR codes locally? (Current: Call backend, could optimize to local generation)

8. **Barcode uniqueness:** Should we include timestamp in barcode to prevent reuse across plan recalculations? (Current: No, barcodes are stable per plan version)

## Conclusion

This feature adds significant value to the load planning system by:

✅ **Reducing loading errors** - Operators scan to verify correct items  
✅ **Improving efficiency** - 3D viewer shows exactly where to place items  
✅ **Creating audit trail** - Track what was actually loaded vs planned  
✅ **Working offline** - No dependency on network during loading  
✅ **Minimal backend changes** - On-the-fly barcode generation, no DB migrations  
✅ **Reusing existing code** - WebView 3D viewer, existing plan endpoints  

**Estimated effort:** 18-22 hours total development + testing  
**Deployment risk:** Low (no database changes, backward compatible)  
**User impact:** High (warehouse operators' primary workflow)

**Recommended next steps:**
1. Review and approve this plan
2. Create implementation tasks in project management tool
3. Begin Phase 1 (Backend API) development
4. Iterate based on feedback

---

**Document Status:** COMPLETED  
**Last Updated:** 2026-02-09 (Implementation Finalized)  
**Next Review:** After implementation completion
