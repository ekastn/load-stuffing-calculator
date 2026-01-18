# Outline Bab 5: Visualisasi 3D Interaktif dengan Three.js

## Tujuan Bab

Mentransformasi hasil kalkulasi packing menjadi panduan operasional visual yang interaktif. Pembaca akan belajar merender objek 3D di browser dan mengimplementasikan urutan pemuatan kronologis.

**Referensi Outline Utama:**
> - **The Goal:** Mentransformasi hasil kalkulasi menjadi panduan operasional visual yang interaktif.
> - **Technical Challenges:** Rendering objek 3D secara real-time di browser dan implementasi urutan pemuatan kronologis.
> - **Addressing Challenges:** Penggunaan WebGL via Three.js, representasi BoxGeometry, dan mekanisme Step Playback berbasis step_number.

---

## Struktur Bab

### Introduction

- Menghubungkan dengan hasil dari Bab 3 (Backend API) dan Bab 4 (Packing Service)
- Mengapa visualisasi penting: panduan operasional untuk pekerja gudang
- WebGL dan Three.js sebagai pilihan rendering 3D di browser
- Preview hasil akhir: container 3D dengan items bertahap

**Technical Requirements:**
- Node.js 20+
- Browser modern dengan WebGL support

---

### 5.1 Setup Project Frontend

- Create Next.js app dengan pnpm: `pnpm create next-app@latest web --yes`
- Setup shadcn/ui: `pnpm dlx shadcn@latest init`
- Install UI components: `pnpm dlx shadcn@latest add button slider card tabs select input table badge`
- Install dependencies: `pnpm add three && pnpm add -D @types/three`
- Struktur direktori:
  ```
  web/
  ├── app/
  │   ├── page.tsx           # Demo page
  │   └── layout.tsx
  ├── components/
  │   ├── animation-controls.tsx
  │   ├── stuffing-viewer.tsx
  │   └── ui/                # shadcn components
  └── lib/
      ├── utils.ts           # shadcn utils
      └── StuffingVisualizer/
  ```

**Files:**
- `package.json`
- `tsconfig.json`
- `app/globals.css`

---

### 5.2 Arsitektur StuffingVisualizer

- Class-based design untuk Three.js logic
- Separation of concerns dengan Manager pattern:
  - `StuffingVisualizer` - Main orchestrator
  - `SceneManager` - Three.js scene setup
  - `CameraManager` - OrthographicCamera configuration
  - `RendererManager` - WebGLRenderer dengan resize handling
  - `LightManager` - Ambient dan directional lights
  - `ControlsManager` - OrbitControls untuk interaktivitas
  - `AnimationManager` - Step-by-step animation playback
- Event system: `onStepChange`, `onPlayStateChange`
- Lifecycle: `loadData → attach → build → animate → dispose`

**Files:**
- `lib/StuffingVisualizer/index.ts`
- `lib/StuffingVisualizer/stuffing-visualizer.ts`
- `lib/StuffingVisualizer/types.ts`

---

### 5.3 Type Definitions

- `ContainerData`: name, length_mm, width_mm, height_mm, max_weight_kg
- `ItemData`: item_id, label, dimensions, weight, quantity, color_hex
- `PlacementData`: placement_id, item_id, pos_x, pos_y, pos_z, rotation, step_number
- `StuffingPlanData`: plan_id, container, items, placements, stats
- `SceneConfig`: backgroundColor, stepDuration

**Files:**
- `lib/StuffingVisualizer/types.ts`

---

### 5.4 Scene dan Renderer Setup

- `SceneManager`: Creating Three.js Scene dengan background
- `RendererManager`: WebGLRenderer configuration:
  - Antialiasing untuk smooth edges
  - Pixel ratio untuk retina displays
  - ResizeObserver untuk responsive resize handling
- Animation loop dengan `requestAnimationFrame`
- Cleanup dan resource disposal

**Files:**
- `lib/StuffingVisualizer/core/scene-manager.ts`
- `lib/StuffingVisualizer/core/renderer-manager.ts`

---

### 5.5 Camera dan Controls

- `CameraManager` dengan OrthographicCamera:
  - Frustum size untuk consistent zoom
  - Isometric-like view position
  - Aspect ratio update on resize
- `ControlsManager` dengan OrbitControls:
  - Damping untuk smooth movement
  - Max polar angle constraint
  - Zoom dan pan limits

**Files:**
- `lib/StuffingVisualizer/core/camera-manager.ts`
- `lib/StuffingVisualizer/core/controls-manager.ts`

---

### 5.6 Lighting Setup

- `LightManager` mengatur pencahayaan scene:
  - AmbientLight untuk overall illumination
  - DirectionalLight untuk depth perception
- Light positioning untuk optimal shadow

**Files:**
- `lib/StuffingVisualizer/core/light-manager.ts`

---

### 5.7 Container Rendering

- `ContainerMesh` class untuk container wireframe:
  - BoxGeometry untuk dimensi container
  - EdgesGeometry untuk wireframe representation
  - LineBasicMaterial untuk wire color
  - Platform mesh untuk floor visual
- Coordinate system: center origin

**Files:**
- `lib/StuffingVisualizer/components/container-mesh.ts`

---

### 5.8 Item Rendering dengan BoxGeometry

- `ItemMesh` class untuk setiap item placement:
  - BoxGeometry untuk item dimensions
  - MeshStandardMaterial dengan color dari data
  - EdgesGeometry untuk white wireframe outline
- Rotation handling: 6 rotation codes (0-5)
- `getRotatedDims()`: menghitung dimensi setelah rotasi
- Coordinate transformation dari API ke Three.js:
  - API: X=length, Y=width, Z=height, origin=corner
  - Three.js: X=horizontal, Y=vertical, Z=depth, origin=center

**Files:**
- `lib/StuffingVisualizer/components/item-mesh.ts`
- `lib/StuffingVisualizer/utils/conversions.ts`

---

### 5.9 Coordinate Conversion Utilities

- `mmToMeters()`: Konversi millimeter ke meter untuk Three.js scale
- `containerToThreeCoords()`: Transform API coordinates ke Three.js space
  - API: origin di corner, Z=height
  - Three.js: origin di center, Y=height
- `getItemCenterOffset()`: Offset untuk center-based positioning

**Files:**
- `lib/StuffingVisualizer/utils/conversions.ts`

---

### 5.10 Step-by-Step Animation

- `AnimationManager` untuk animation playback:
  - `setMaxStep()`: total steps dari placements
  - `setCurrentStep()`: visibility control
  - `play()`: auto-advance dengan interval
  - `pause()`: stop auto-advance
  - `reset()`: kembali ke step 0
- Event listeners: `stepListeners`, `playStateListeners`
- `updateVisibleItems()`: show/hide berdasarkan step_number

**Files:**
- `lib/StuffingVisualizer/core/animation-manager.ts`

---

### 5.11 React Integration

- `StuffingViewer` React component:
  - `useRef` untuk canvas container
  - `useEffect` untuk Three.js lifecycle
  - State syncing: currentStep, maxStep, isPlaying
  - Dynamic import untuk SSR compatibility
- `AnimationControls` component:
  - shadcn Button untuk play/pause/reset
  - shadcn Slider untuk step selection
  - Step counter display
- Event subscription dan cleanup

**Files:**
- `components/stuffing-viewer.tsx`
- `components/animation-controls.tsx`

---

### 5.12 Demo Page

- Root page dengan Load & Stuffing Calculator layout:
  - Container settings form (length, width, height, max_weight)
  - Item input form (label, quantity, dimensions, weight, color)
  - Items table dengan delete action
  - Calculate Packing button
  - 3D Visualization panel
- Simple packing algorithm untuk demo:
  - Left-to-right, front-to-back, bottom-to-top
  - Step number assignment

**Files:**
- `app/page.tsx`

---

### Summary

- Three.js memungkinkan rendering 3D di browser tanpa plugin
- Manager pattern memisahkan concerns dan memudahkan testing
- StuffingVisualizer mengenkapsulasi kompleksitas Three.js
- Step-by-step animation memberikan panduan operasional
- Coordinate conversion menghubungkan API data dengan Three.js space
- Integration dengan React melalui effects, refs, dan events

### Further Reading

- Three.js Documentation: https://threejs.org/docs/
- Three.js Fundamentals: https://threejs.org/manual/
- WebGL Fundamentals: https://webglfundamentals.org/
- OrbitControls: https://threejs.org/docs/#examples/en/controls/OrbitControls

---

## Estimasi

- **Panjang**: ~5000-6000 kata
- **Code snippets**: 16 files
- **Implementasi reference**: `docs/buku/source/bab_5/web/`

---

## Source Code Files

```
docs/buku/source/bab_5/web/
├── app/
│   └── page.tsx                           # Demo page
├── components/
│   ├── animation-controls.tsx             # Play/pause/slider controls
│   ├── stuffing-viewer.tsx                # React wrapper
│   └── ui/                                # shadcn components
│       ├── button.tsx
│       ├── slider.tsx
│       ├── card.tsx
│       ├── input.tsx
│       ├── table.tsx
│       └── badge.tsx
├── lib/
│   ├── utils.ts                           # shadcn utils
│   └── StuffingVisualizer/
│       ├── index.ts                       # Exports
│       ├── stuffing-visualizer.ts         # Main orchestrator
│       ├── types.ts                       # Type definitions
│       ├── core/
│       │   ├── scene-manager.ts
│       │   ├── camera-manager.ts
│       │   ├── renderer-manager.ts
│       │   ├── light-manager.ts
│       │   ├── controls-manager.ts
│       │   └── animation-manager.ts
│       ├── components/
│       │   ├── container-mesh.ts
│       │   └── item-mesh.ts
│       └── utils/
│           └── conversions.ts
└── package.json
```

---

## Notes

### Koneksi dengan Bab Lain

- **Bab 4**: Data placements dari Packing Service (pos_x, pos_y, pos_z, rotation, step_number)
- **Bab 6**: Advanced features dan optimization

### Pendekatan Penulisan

- Fokus pada Three.js concepts, bukan React/Next.js details
- Diagram untuk menjelaskan coordinate system transformation
- Screenshots untuk expected visual output
- Step-by-step build: setup → types → managers → meshes → animation → react

### Fitur yang Diimplementasikan

- Three.js scene, camera, renderer, lights
- OrbitControls untuk interaktivitas
- Container wireframe dengan floor platform
- Item boxes dengan rotation support
- Step-by-step animation dengan controls
- React integration dengan state syncing
- Demo page dengan form inputs

### Fitur yang Tidak Dibahas

- Raycasting dan hover effects
- Tooltip untuk item info
- PDF report generation
- Advanced lighting dan shadows
- Performance optimization untuk 1000+ items
- Mobile touch controls

### Guidelines Penulisan

- Diagram-First: Architecture diagram untuk Manager pattern
- Modular Code Snippets: Setiap file satu section
- Architectural Justification: Mengapa Manager pattern, mengapa OrthographicCamera
- In-code Comments: Setiap API Three.js dijelaskan
