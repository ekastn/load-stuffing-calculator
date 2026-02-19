# Mobile Client Developer Handbook

This document serves as the single source of truth for developing the Load Stuffing Calculator mobile client. It defines the architecture, coding standards, testing strategies, and operational procedures in detail.

## 1. Core Architecture

### Philosophy
- **Personal User Focus:** This mobile client is exclusively for "Personal" type users. It ignores Organizational features (team management, multi-workspace switching, etc.).
- **Separation of Concerns:**
  - **Services (Business Logic):** Pure Dart classes for API, transformation, and storage. Stateless & Testable.
  - **Providers (State Management):** `ChangeNotifier` classes that hold UI state and consume Services.
  - **Pages (UI):** Widgets that consume Providers via the `provider` package.
- **Layer-First Structure:** Organized into `dtos/`, `models/`, `mappers/`, `services/`, `providers/`, `pages/`.
- **Web-Driven Intelligence:** 3D Visualization and Packing logic reside in the WebView; Hardware (Scanner) and Master Data reside in Native.

### Directory Structure
```
lib/
├── main.dart                 # Composition Root & App Entry
├── config/                   # Configuration & Constants
│   ├── assets.dart           # Asset paths (images, icons)
│   ├── constants.dart        # API Endpoints, Magic Numbers
│   ├── routes.dart           # GoRouter Config
│   └── theme.dart            # ThemeData & AppColors
├── components/               # Reusable UI Widgets
│   ├── buttons/              # PrimaryButton, AppButton
│   ├── dialogs/              # ErrorDialog, ConfirmDialog
│   ├── inputs/               # AppTextField, SearchBar
│   ├── overlays/             # ScannerOverlay, LoadingOverlay
│   └── misc/                 # PermissionGuard, Gap
├── models/                   # UI Entities (Freezed) - Used by UI & Providers
│   ├── user_model.dart
│   ├── token_model.dart
│   ├── product_model.dart
│   ├── container_model.dart
│   └── scan_event_model.dart
├── dtos/                     # API JSON Objects (JsonSerializable) - Used by Services
│   ├── auth_dto.dart         # LoginRequestDto, LoginResponseDto
│   ├── product_dto.dart      # ProductResponseDto, ProductListDto
│   ├── container_dto.dart    # ContainerResponseDto
│   └── error_dto.dart        # ApiErrorDto
├── mappers/                  # DTO <-> Model Transformers (Pure Logic)
│   ├── auth_mapper.dart
│   ├── product_mapper.dart
│   └── container_mapper.dart
├── pages/                    # Screens (Scaffolds)
│   ├── auth/                 # LoginPage, SplashPage
│   ├── home/                 # DashboardPage
│   ├── master/               # ProductsPage, ContainersPage, ProductDetailPage
│   └── shipment/             # ViewerPage (WebView Container)
├── services/                 # Pure Logic (Data Access)
│   ├── api_service.dart      # Dio Client (Interceptors, Retry)
│   ├── auth_service.dart     # Login/Logout logic (Stateless)
│   ├── bridge_service.dart   # Web Controller Logic
│   ├── master_service.dart   # Products Repositories
│   └── storage_service.dart  # Secure Storage Wrapper
├── providers/                # App State (ChangeNotifiers)
│   ├── auth_provider.dart    # User State, Auth Status
│   ├── master_provider.dart  # Product List, Search Filters
│   └── viewer_provider.dart  # WebView State, Bridge Events
└── utils/                    # Helpers & Extensions
    ├── extensions/           # ContextExt, StringExt
    ├── formatters/           # DateFormatter, CurrencyFormatter
    └── validators/           # EmailValidator, PasswordValidator
test/                         # Testing Strategy
├── unit/                     # Services, Mappers, Providers
├── widget/                   # Components, Pages
└── integration/              # Full app flows
```

---

## 2. Authentication & Authorization

### A. JWT-Enclosed Context
The backend Go API uses a **JWT-based Workspace Scoping** strategy. There is no need for custom headers like `X-Workspace-ID`.

1. **Login:** `AuthService.login()` calls `/auth/login`.
2. **Auto-Scoping:** The backend's `Login` service automatically resolves the user's `personal` workspace and embeds its `workspace_id` directly into the JWT Access Token claims.
3. **Usage:** The mobile app simply includes the `Authorization: Bearer <token>` header in every request. 
4. **Backend Enforcement:** The API's `JWT` and `Permission` middlewares automatically extract the workspace context from the token to filter all database queries.
5. **Persistence:** The `access_token` and `refresh_token` are stored in `FlutterSecureStorage`.

### B. Global Interceptor
A simple Dio interceptor handles the inclusion of the bearer token for all authenticated requests.

**File:** `services/api_service.dart`
```dart
class AuthInterceptor extends Interceptor {
  final StorageService _storage;
  
  AuthInterceptor(this._storage);

  @override
  void onRequest(RequestOptions options, RequestInterceptorHandler handler) async {
    final token = await _storage.getAccessToken();
    if (token != null) {
      options.headers['Authorization'] = 'Bearer $token';
    }
    return handler.next(options);
  }
  
  @override
  void onError(DioException err, ErrorInterceptorHandler handler) {
    if (err.response?.statusCode == 401) {
        // Trigger global logout or refresh flow
    }
    return handler.next(err);
  }
}
```

---

## 3. Data Strategy: DTO vs Model vs Mapper

We strictly separate API Data from App Data using Mappers to protect the UI from backend schema changes.

- **DTOs (`lib/dtos/`):** Flat files mirroring API JSON. No logic. Suffix with `Dto`. Use `json_serializable`.
- **Models (`lib/models/`):** Clean, type-safe data for UI. Suffix with `Model`. Use `freezed`.
- **Mappers (`lib/mappers/`):** Pure static functions for transformation (`toModel` and `toDto`).

**Example Data Flow:**
`ApiService` (JSON) -> `AuthDto` -> `AuthMapper.toUserModel()` -> `UserModel` -> `AuthProvider`.

---

## 4. Architecture Pattern: Service-Provider-UI

We use **Constructor Injection** to ensure testability. All dependencies are wired in the `main.dart` Composition Root.

### A. The Service
Pure Dart logic. Returns clean Models.
```dart
class AuthService {
  final ApiService _api;
  AuthService(this._api);

  Future<UserModel> login(String user, String pass) async {
    try {
      final response = await _api.post('/auth/login', {
        'username': user,
        'password': pass,
      });
      final dto = LoginResponseDto.fromJson(response.data);
      return AuthMapper.toUserModel(dto); // Map DTO to UI Model
    } catch (e) {
      throw AuthException.fromDio(e); // Custom Exception
    }
  }
}
```

### B. The Provider
App state manager. Extends `ChangeNotifier`.
```dart
class AuthProvider extends ChangeNotifier {
  final AuthService _authService;
  AuthProvider({required AuthService authService}) : _authService = authService;

  UserModel? _user;
  bool get isAuthenticated => _user != null;
  String? _error;

  Future<void> login(String u, String p) async {
    _error = null;
    notifyListeners();
    try {
      _user = await _authService.login(u, p);
    } catch (e) {
      _error = e.toString();
    } finally {
      notifyListeners();
    }
  }
}
```

---

## 5. Communication Bridge (WebView <-> Native)

The WebView handles 3D rendering. The Bridge uses a JSON protocol for bidirectional communication.

- **Native -> Web:** `window.stuffingBridge.onScan(barcode)` via `runJavaScript`.
- **Web -> Native:** `JavascriptChannel('FlutterBridge')` receiving JSON events like `SCAN_RESULT`.

**Error Handling in Bridge:**
- If Web logic fails (e.g., item not found), it sends `{ "type": "ERROR", "message": "Item not in plan" }`.
- Flutter receives this and shows a `SnackBar` (Red) + Haptic Feedback (Double Vibrate).

---

## 6. Implementation Roadmap

### Phase 1: Skeleton & Core
- [ ] Initialize project with `dio`, `provider`, `freezed`, `flutter_secure_storage`.
- [ ] Create `main.dart` Composition Root and Dependency Graph.
- [ ] Implement `ApiInterceptor` and `AppException` classes.
- [ ] Setup `build.yaml` for code generation.

### Phase 2: Personal Auth
- [ ] Implement `AuthService` and `AuthMapper`.
- [ ] Create `AuthProvider` to hold `UserModel`.
- [ ] Create `LoginPage` with validation.
- [ ] Implement `StorageService` for secure token handling.

### Phase 3: Hardware & 3D Viewer
- [ ] Create `ViewerPage` (WebView) loading the `/shipments/:id/embed` route.
- [ ] Implement `BridgeService` to abstract `WebViewController`.
- [ ] Integrate `mobile_scanner` overlay.
- [ ] Wire up Scan -> Bridge -> WebView -> Feedback loop.

### Phase 4: Master Data (Personal Context)
- [ ] Implement `ProductService` & `ContainerService` (API already scopes by JWT).
- [ ] Create `ProductListModel` and `ProductDetailModel`.
- [ ] Build Searchable List UI for Products.

### Phase 5: Polish & Release
- [ ] Add Loading States (Shimmer/Spinner).
- [ ] Implement Error Dialogs.
- [ ] App Icon & Splash Screen.
- [ ] CI/CD Setup (Fastlane).

---

## 7. Tech Stack & Versions

| Package | Version Constraint | Usage |
| :--- | :--- | :--- |
| `flutter` | `3.19.0+` | SDK |
| `provider` | `^6.1.1` | DI & State |
| `dio` | `^5.4.0` | HTTP Client |
| `freezed_annotation` | `^2.4.1` | Immutable Models |
| `json_annotation` | `^4.8.1` | JSON Serialization |
| `go_router` | `^13.0.0` | Navigation |
| `flutter_secure_storage` | `^9.0.0` | Token Storage |
| `webview_flutter` | `^4.4.2` | 3D Visualization |
| `mobile_scanner` | `^5.0.0` | Barcode Scanning |
| `mocktail` | `^1.0.0` | Unit Testing |

## 8. Testing Strategy

1.  **Unit Tests (`test/unit/`):**
    *   **Services:** Mock `Dio` adapter to test API handling and Error Mapping.
    *   **Mappers:** Verify JSON <-> Entity conversion logic.
    *   **Providers:** Mock Services to test State changes (`isLoading`, `error`, `data`).
2.  **Widget Tests (`test/widget/`):**
    *   Test standard Components (Buttons, Inputs).
    *   Test Pages (Login) using mocked Providers.
