# Portal Frontend Redesign – Design Specification

## 1. Problem Statement

The landing page (`web/app/(public)`) uses a warm industrial palette (`DESIGN.md`), while the portal (`web/app/(app)`) still uses default template blues (`oklch(0.3 0.15 260)`) and hardcoded Tailwind utilities (`bg-white`, `bg-slate-50/50`, `text-slate-600`, `text-blue-500`, `text-violet-500`).

However, redesigning every page's layout or forcing the landing page's stark stencil/box style into the portal is too heavy-handed. The existing portal layouts and `shadcn/ui` components already work well – they just need:
1. The **warm-stone background and sidebar palette** (`#f8f8f6` canvas, `#f2f3ef` sidebar, `#ffffff` cards) paired with LoadIQ's **Ink (`#16181b`)**, **Steel (`#525860`)**, and **Signal Orange (`#ff5b14`)** accents.
2. **Proper semantic CSS variables** in `web/app/globals.css` (`:root` and `.dark`) so changing `color-scheme` or tweaking colors updates the entire app from one file.
3. **Light, modern polish on `shadcn/ui` primitives** (`Card`, `Button`, `Badge`, `Input`, `Select`, `Tabs`, `Table`) and replacing hardcoded `slate`/`blue`/`white` classes across `web/app/(app)`.

---

## 2. Objectives

1. **Keep existing page layouts intact**: Preserve the current Dashboard (stat cards + quick actions + recent activity), Shipments (filter bar + shipment card grid), Shipment Detail (header + *3D Simulation / Cargo Breakdown / Summary* tabs), and `DataTable` screens.
2. **Apply the semantic color system (`globals.css`)**:
   - Use the warm-stone page background (`--background: #f8f8f6`), warm-stone sidebar (`--sidebar: #f2f3ef`), crisp white cards (`--card: #ffffff`), deep ink text/buttons (`--foreground: #16181b`), and safety orange brand accents (`--signal: #ff5b14`).
   - Provide matching `.dark` tokens so switching `color-scheme` works out of the box without touching individual components.
3. **Refine `shadcn/ui` components & remove hardcoded colors**:
   - Standardize `--radius: 0.5rem` (`8px` cards, `6px` buttons/inputs/badges) and subtle borders (`--border: #e3e5e0`).
   - Replace hardcoded `bg-white`, `bg-slate-*`, `text-slate-*`, and rainbow colors in `web/app/(app)` with semantic tokens (`bg-card`, `bg-muted`, `text-foreground`, `text-muted-foreground`, `bg-primary`, `bg-signal`).

---

## 3. Semantic Token System (`web/app/globals.css`)

| Token | Light (`:root`) | Dark (`.dark`) | Usage |
|---|---|---|---|
| `--background` | `#f8f8f6` | `#0f1114` | Main page background |
| `--sidebar` | `#f2f3ef` | `#121519` | Sidebar background |
| `--card` / `--popover` | `#ffffff` | `#16191e` | Cards, tables, dialogs, selects, inputs |
| `--foreground` | `#16181b` | `#f3f4f6` | Primary text & headings (`ink`) |
| `--muted` / `--secondary` | `#f1f2ee` | `#1c2026` | Muted backgrounds, tab strips, hover fills |
| `--muted-foreground` | `#525860` | `#949ca6` | Secondary text & labels (`steel`) |
| `--primary` | `#16181b` | `#f3f4f6` | Primary buttons at rest (`ink`) |
| `--primary-foreground` | `#f5f5f2` | `#16181b` | Text on `--primary` |
| `--signal` | `#ff5b14` | `#ff5b14` | Primary button hover fill, active sidebar icon, focus ring, progress bars |
| `--signal-foreground` | `#16181b` | `#16181b` | Dark ink text on `--signal` hover (5.72:1 contrast) |
| `--signal-soft` | `rgba(255, 91, 20, 0.10)` | `rgba(255, 91, 20, 0.15)` | Subtle orange tint for active nav / badges |
| `--signal-text` | `#c23d06` | `#ff854d` | Accessible orange text on light/dark |
| `--border` | `#e3e5e0` | `#262a32` | Card, sidebar, and table borders |
| `--input` | `#cfd3cd` | `#323741` | Form input and select borders |
| `--ring` | `#ff5b14` | `#ff5b14` | Focus ring |
| `--destructive` | `#d92d20` | `#f04438` | Destructive actions & failed status |
| `--status-ok-bg` / `fg` | `#ecfdf3` / `#067647` | `rgba(23,178,106,0.14)` / `#75e0a7` | `COMPLETED` badge & positive metrics |
| `--status-warn-bg` / `fg` | `#fffaeb` / `#b54708` | `rgba(247,144,9,0.14)` / `#fec84b` | `PARTIAL` badge & warnings |
| `--status-danger-bg` / `fg` | `#fef3f2` / `#b42318` | `rgba(240,68,56,0.14)` / `#fda29b` | `FAILED` badge & errors |
| `--radius` | `0.5rem` | `0.5rem` | `8px` cards (`rounded-lg`), `6px` controls (`rounded-md`) |

---

## 4. What Changes vs. What Stays the Same

| Area | Stays the Same | Light Polish & Token Update |
|---|---|---|
| **Sidebar (`dashboard-layout.tsx`)** | Same navigation groups, workspace switcher, mobile drawer, and user footer | Uses `--sidebar` (`#f2f3ef`), `LogoMark` + `LoadIQ` title, clean active state (`bg-card` with subtle border & orange icon accent) |
| **Dashboard (`dashboard/page.tsx`)** | Same 4-card stats grid, Quick Actions card, and Recent Activity card | Replaces rainbow icon colors (`text-blue-500`, `text-violet-500`) with semantic token accents; clean `bg-card` surfaces |
| **Shipments (`shipments/page.tsx`)** | Same Filter Bar card + 3-column Shipment Cards grid | Replaces hardcoded `bg-white` with `bg-card`, uses semantic status badges, and highlights hover state cleanly |
| **Shipment Detail (`shipments/[id]/page.tsx`)** | Same compact top header + 3 Tabs (*3D Simulation*, *Cargo Breakdown* donut + table, *Summary* 3 cards + table) | Replaces hardcoded `bg-white`, `bg-slate-50/50`, `text-slate-*`, `bg-blue-500` with semantic tokens (`bg-card`, `bg-muted/40`, `text-foreground`, `bg-signal`) |
| **Tables (`products`, `containers`, `members`, etc.)** | Same `DataTable` component, columns, toolbar, and dialogs | Removes redundant double borders, uses `bg-card` and semantic text/border tokens |
| **`shadcn/ui` (`components/ui/*`)** | Same Radix primitives, props, and variants | Tuned to `--radius: 0.5rem`, `shadow-xs`, semantic borders, plus optional `signal` variant on `Button` and status variants on `Badge` |
