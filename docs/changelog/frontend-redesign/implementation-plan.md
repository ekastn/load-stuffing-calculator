# Portal Frontend Redesign – Implementation Plan

## 1. Scope & Philosophy

We will **not** overhaul page layouts or replace standard `shadcn/ui` patterns. Instead, we will:
1. Update the **CSS token layer** (`web/app/globals.css`) with the warm-stone background (`#f8f8f6`), sidebar (`#f2f3ef`), crisp white cards (`#ffffff`), LoadIQ ink (`#16181b`), steel (`#525860`), signal orange (`#ff5b14`), and semantic status variables (with full `.dark` mode support).
2. Lightly tune **`web/components/ui/*`** (`button`, `badge`, `card`, `input`, `select`, `tabs`, `table`) so they consume the new tokens cleanly and look modern (`--radius: 0.5rem`, `shadow-xs`).
3. Replace **hardcoded color classes** (`bg-white`, `bg-slate-*`, `text-slate-*`, `text-blue-*`, `text-green-*`, `text-violet-*`) across `web/components/dashboard-layout.tsx` and `web/app/(app)/**` with semantic token classes so `color-scheme` changes work everywhere without touching components again.

---

## 2. Step-by-Step Checklist

### Step 1: Semantic Tokens (`web/app/globals.css`)
- [ ] Update `:root` (`color-scheme: light`) and `.dark` (`color-scheme: dark`):
  - `--background` (`#f8f8f6` light / `#0f1114` dark)
  - `--sidebar` (`#f2f3ef` light / `#121519` dark)
  - `--card` & `--popover` (`#ffffff` light / `#16191e` dark)
  - `--foreground` (`#16181b` light / `#f3f4f6` dark)
  - `--muted` & `--secondary` (`#f1f2ee` light / `#1c2026` dark)
  - `--muted-foreground` (`#525860` light / `#949ca6` dark)
  - `--primary` (`#16181b` light / `#f3f4f6` dark) and `--primary-foreground`
  - `--signal` (`#ff5b14`), `--signal-foreground` (`#ffffff`), `--signal-soft`, `--signal-text`
  - `--border` (`#e3e5e0` light / `#262a32` dark), `--input` (`#cfd3cd` light / `#323741` dark), `--ring` (`#ff5b14`)
  - Status tokens: `--status-ok-bg`, `--status-ok-fg`, `--status-warn-bg`, `--status-warn-fg`, `--status-danger-bg`, `--status-danger-fg`, `--status-progress-bg`, `--status-progress-fg`, `--status-draft-bg`, `--status-draft-fg`
  - `--radius: 0.5rem` (8px cards, 6px controls)
- [ ] Expose the new semantic tokens (`--color-signal`, `--color-signal-soft`, `--color-signal-text`, `--color-status-*`) in `@theme inline`.
- [ ] Keep `.landing` scoped to `web/app/(public)` so the landing and auth pages remain unchanged.

### Step 2: Shared `shadcn/ui` Components (`web/components/ui/*`)
- [ ] **`web/components/ui/button.tsx`**:
  - Keep existing sizes and `rounded-md` geometry.
  - Update `default` variant to rest in `ink` and transition to safety orange on hover: `bg-primary text-primary-foreground hover:bg-signal hover:text-signal-foreground shadow-2xs` (matching the "Start a load plan" button in the landing page header).
  - Clean up `outline` (`border-border bg-card hover:bg-accent text-foreground shadow-2xs`).
- [ ] **`web/components/ui/badge.tsx`**:
  - Add semantic status variants (`completed`, `in_progress`, `partial`, `draft`, `failed`) using the `--status-*` tokens and export a small `getStatusBadgeVariant(status)` helper.
- [ ] **`web/components/ui/card.tsx`**:
  - Use `rounded-lg border border-border bg-card text-card-foreground shadow-xs`.
- [ ] **`web/components/ui/input.tsx`, `select.tsx`, `tabs.tsx`, `table.tsx`**:
  - Ensure inputs and selects use `bg-card border-input` instead of `bg-transparent`.
  - Ensure `TabsList` and `TabsTrigger` use `bg-muted` and `data-[state=active]:bg-card`.

### Step 3: Sidebar & Shell (`web/components/dashboard-layout.tsx`)
- [ ] Use `bg-sidebar border-r border-border` on the sidebar container.
- [ ] Style active nav item with `bg-card text-foreground shadow-2xs border border-border/80 [&_svg]:text-signal` and inactive items with `text-muted-foreground hover:bg-accent/60 hover:text-foreground`.
- [ ] Style user avatar badge with `bg-primary text-primary-foreground` (or `bg-signal-soft text-signal-text`).

### Step 4: Portal Pages Token Cleanup (`web/app/(app)/**`)
- [ ] **`web/app/(app)/dashboard/page.tsx`**:
  - Keep the 4-card stats grid, Quick Actions card, and Recent Activity card.
  - Replace `border-border/50 bg-card/50 backdrop-blur-sm` with clean `Card` defaults (`bg-card border-border`).
  - Replace rainbow icon colors (`text-blue-500`, `text-violet-500`, `text-green-500`) with semantic token classes (`text-muted-foreground`, `text-signal`, `text-status-ok-fg`).
- [ ] **`web/app/(app)/shipments/page.tsx`**:
  - Keep the Filter Bar + 3-column Shipment Cards grid.
  - Replace hardcoded `bg-white` on the Filter Bar card with `bg-card`.
  - Replace ad-hoc `statusColors` (`bg-yellow-500/10`, `bg-green-500/10`) with the semantic status badge variants.
- [ ] **`web/app/(app)/shipments/[id]/page.tsx`, `web/components/stuffing-viewer.tsx`, `web/components/animation-controls.tsx`**:
  - Keep the compact top bar + 3 Tabs (*3D Simulation*, *Cargo Breakdown*, *Summary*) and the floating bottom-center `AnimationControls` bar (`Play/Pause`, `Reset`, `Download`, `Slider`, and `currentStep / maxStep` counter).
  - Update `AnimationControls` and `StuffingViewer` to use semantic tokens (`bg-card/95 border-border text-foreground shadow-md`) instead of hardcoded `bg-white border-zinc-200 text-zinc-700`, and `tabular-nums` on the step counter.
  - Replace all hardcoded `bg-white`, `bg-slate-50/50`, `border-slate-50`, `text-slate-900/700/600/500/400`, `bg-blue-500`, and `text-blue-600` with semantic tokens (`bg-card`, `bg-muted/40`, `border-border`, `text-foreground`, `text-muted-foreground`, `bg-signal`, `text-signal-text`).
- [ ] **`web/app/(app)/products/**`, `containers/**`, `settings/**`, `reports/**`, `(platform)/**`**:
  - Replace `border-border/50 bg-card/50` wrappers around `DataTable` with clean `bg-card border-border rounded-lg shadow-xs`.

---

## 3. Verification
- [ ] `pnpm -C web lint`
- [ ] `pnpm -C web test`
- [ ] `pnpm -C web build`
