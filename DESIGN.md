# LoadIQ Design System

This is the visual and interaction system for LoadIQ. The landing page (`web/app/(public)/page.tsx`) already uses it. The rest of the web app (dashboard, plans, master data, auth) still uses the original shadcn defaults and is meant to move to this system. Section 13 is the migration checklist.

Direction name: **industrial manifest**. The product plans how cargo goes into shipping containers, so the visual language comes from that world: container door stencils, cargo manifests, kraft boxes, safety markings, measured dimensions.

---

## 1. Principles

1. **Measure, don't decorate.** Every number shown is a real measurement with a unit. Structure (rules, numbering, labels) only appears when it carries information. Numbered markers are only for real sequences, such as plan steps.
2. **One loud thing per screen.** Safety orange marks the single most important thing: the primary action, the box being placed, the current step, the focused control. If two things are orange, one of them is wrong.
3. **Ink on concrete.** Light, cool grey surfaces with near-black ink. Flat color, hairline rules, no shadows, no gradients, no glass.
4. **Stencil for names, plain text for reading.** The stencil face is for headings and the wordmark. Body text and UI use a plain grotesk. Mono is only for measured values.
5. **Real data or nothing.** No fake testimonials, invented stats, logo walls, or placeholder charts. Illustrations of the product are built from real packing output.
6. **Motion shows what changed.** Animation answers an action or shows the load sequence. Nothing moves just to fill space, and everything works with reduced motion.

---

## 2. Color

### 2.1 Core palette

| Token | Hex | Role |
|---|---|---|
| `yard` | `#e3e4e0` | Page background. Cool concrete grey. |
| `form` | `#f5f5f2` | Raised surfaces: cards, panels, inputs, popovers. Light text on dark. |
| `ink` | `#16181b` | Primary text, headings, frames, dark sections, secondary buttons. |
| `steel` | `#4f585f` | Secondary text, captions, metadata. |
| `signal` | `#ff5b14` | Safety orange. Primary action fill, current item, focus ring, meter fill. Never text on light. |
| `kraft` | `#b8895a` | Cargo. Boxes in 3D views and the cargo swatches. Not for UI. |

Defined in `web/app/globals.css` under `@theme`, so Tailwind classes exist for each: `bg-yard`, `text-ink`, `border-ink/15`, `bg-signal`, and so on.

### 2.2 Supporting colors

| Token | Hex | Role |
|---|---|---|
| `signal-text` | `#a8380a` | Orange used as text or a small icon on `yard`/`form` (5.08:1 on yard). |
| `rule` | `ink` at 15% (`#c4c5c2` on yard) | Hairline dividers between rows and sections. Decorative only. |
| `rule-strong` | `ink` at 20–25% | Panel outlines, the manifest readout top rule. |
| `field` | `#7d8480` | Input, select and checkbox borders (3.51:1 on form, 3.0:1 on yard). |
| `danger` | `#b42318` | Errors, destructive actions, FAILED status. |
| `warn` | `#7a5000` | Warnings, PARTIAL status. |
| `ok` | `#2f6b3a` | Success, COMPLETED status. |

`signal-text` is in `globals.css` (`text-signal-text`). `field`, `danger`, `warn` and `ok` are not yet; add them under `@theme` when the app migrates.

### 2.3 Contrast (WCAG 2.2)

| Foreground on background | Ratio | Allowed use |
|---|---|---|
| ink on yard | 13.93 | All text |
| ink on form | 16.28 | All text |
| steel on yard | 5.68 | All text |
| steel on form | 6.64 | All text |
| form on ink | 16.28 | All text |
| form at 75% on ink | 9.53 | Body text on dark sections |
| ink on signal | 5.72 | Text on orange buttons |
| form on signal | **2.84** | **Never.** Orange buttons always use ink text. |
| signal on yard | **2.43** | **Never as text.** Fills, bars and large shapes only. |
| signal on ink | 5.72 | Text allowed on dark |
| signal-text on yard | 5.08 | Orange text on light surfaces |
| field on form | 3.51 | Form control borders |
| danger on form | 6.02 | Error text |
| ok on form | 5.85 | Success text |

### 2.4 Status colors (load plans)

Plan status is the main state in the app. Map it like this and always pair the color with the status word. Never show a color alone.

| Status | Treatment |
|---|---|
| DRAFT | `steel` text, `rule-strong` outline |
| IN_PROGRESS | `ink` text, `signal` 4px left bar or dot |
| COMPLETED | `form` text on `ink` fill (the "Completed" badge in the plan view illustration) |
| PARTIAL | `warn` text, `warn` outline |
| FAILED | `danger` text, `danger` outline |
| CANCELLED | `steel` text, no outline |

### 2.5 Dark mode (proposed for the app)

Dark mode keeps the same roles. Orange stays the same hex.

| Token | Light | Dark |
|---|---|---|
| background | `#e3e4e0` | `#16181b` |
| surface | `#f5f5f2` | `#202327` |
| text | `#16181b` | `#f5f5f2` (16.28 on bg, 14.44 on surface) |
| muted text | `#4f585f` | `#a9b0b5` (8.10 on bg, 7.19 on surface) |
| rule | ink 15% | `#3a3f44` |
| field | `#7d8480` | `#767d82` (3.78 on surface) |
| signal | `#ff5b14` | `#ff5b14` (5.72 on bg; text allowed) |
| danger | `#b42318` | `#f97066` (5.66 on surface) |
| warn | `#7a5000` | `#e0a43b` (7.17 on surface) |
| ok | `#2f6b3a` | `#6cc287` (7.29 on surface) |

The landing page is light only. Its dark sections (`#demo`, the closing section) use `bg-ink text-form`, which is the same pairing as dark mode.

---

## 3. Typography

### 3.1 Families

| Role | Family | Loaded as | Tailwind |
|---|---|---|---|
| Display: headings, wordmark, role names | Big Shoulders Stencil (variable, opsz + wght) | `--font-landing-stencil` | `font-stencil` |
| Text and UI | Archivo (variable, wdth + wght) | `--font-landing-sans` | `font-landing` |
| Measured values | IBM Plex Mono 400, 500 | `--font-landing-mono` | `font-figures` |

All three load through `next/font` in `web/components/landing/fonts.ts`. They are currently scoped to the landing page.

Rules:
- **Mono is only for numbers with units or counters**: `5898 × 2352 × 2393 mm`, `2,950 kg`, `Box 29 of 58`, `Step 2 of 4`, plan codes like `PLD-20261027-2259`. Never for labels or prose.
- **Uppercase only for the hero `h1` and the closing heading** on marketing pages, where it reads as door stencil. Everything else is sentence case. No all-caps labels or eyebrows.
- **Don't accent a single word** in a heading (no one orange or italic word).
- Body line length stays under about 70 characters (`max-w-[32rem]` to `max-w-[34rem]` at 18px).

### 3.2 Scale

| Style | Family | Size | Weight | Line height | Used for |
|---|---|---|---|---|---|
| Display XL | Stencil, uppercase | `clamp(3.25rem, 7.4vw, 7rem)` | 800 | 0.88 | Hero `h1` only |
| Display L | Stencil | 48 / 72px (`text-5xl md:text-7xl`) | 700 | 0.95 | Section headings |
| Display M | Stencil | 36 / 48px (`text-4xl md:text-5xl`) | 700 | 1.0 | Sub-sections, page titles in the app |
| Display S | Stencil | 30 / 36px (`text-3xl md:text-4xl`) | 700 | 1.1 | Step titles, role names, card titles |
| Wordmark | Stencil | 24px header, 26px+ elsewhere | 800 | 1 | "LoadIQ" next to the mark |
| Lead | Archivo | 18 / 20px (`text-lg md:text-xl`) | 400 | 1.6 | Intro paragraphs |
| Body | Archivo | 16–18px | 400 | 1.6 | Paragraphs, descriptions |
| UI | Archivo | 15px (`text-[0.95rem]`) | 600 | 1.2 | Buttons and links, with `word-spacing: 0.12em` |
| Label | Archivo | 14px | 600 | 1.3 | Form labels, table headers, list titles |
| Figure | Plex Mono | 14–16px | 400–500 | 1.4 | Measurements and counters |
| Figure small | Plex Mono | 12px | 400 | 1.4 | Secondary measurements, units under meters |

In the app, dense screens (tables, forms, sidebars) should stay at Body 14–16px and Display S/M for titles. Display L/XL are for marketing pages only.

---

## 4. Layout and spacing

- **Spacing unit:** 4px (Tailwind's default `--spacing: 0.25rem`).
- **Content width:** `max-w-[90rem]` (1440px), centred.
- **Gutters:** 16px on mobile (`px-4`), 32px from `md` (`px-8`).
- **Grid:** 12 columns from `lg`, `gap-10` (40px). Typical splits are 5/7 (copy/visual) and 7/5 (heading/intro).
- **Section rhythm (marketing):** `py-24` (96px) on mobile, `lg:py-36` (144px) on desktop.
- **Header height:** 64px (`h-16`), fixed. Content below it starts with `pt-16`, and anchored sections use `scroll-mt-16`.
- **Alignment:** left-aligned everywhere. Centred text is not used.
- **Rows instead of cards:** lists of things (steps, roles, cargo lines) are rows separated by hairline rules (`border-t border-ink/15`), not grids of identical cards.

App screens (proposed):
- Sidebar 256px (`w-64`), surface `form`, right rule.
- Page padding 24px mobile, 32px desktop. Page title is Display M, followed by a one-line steel description.
- Tables: rule between rows, no zebra striping, numbers right-aligned in mono.

---

## 5. Shape, borders, depth

- **Radius:** `--radius: 0.25rem` (4px) for form controls and small chips. Buttons, panels, images and the plan view are square (0). Don't put one radius on everything.
- **Borders:** 1px. `ink/15` for dividers, `ink/20` for panel outlines, `field` for form controls.
- **Shadows:** none. Depth comes from surface color (`yard` → `form` → `ink`). Remove shadcn's `shadow-*` usage when migrating.
- **Not used:** gradients, blurred blobs, glassmorphism (the header's `backdrop-blur-sm` over `bg-yard/90` is the only exception), glow, noise textures.

---

## 6. Components

These exist in `web/components/landing/` and should become shared components in `web/components/ui/` (or `components/brand/`) during the migration.

### 6.1 Button and link (`landing-link.tsx`)

One component, four variants. Minimum height 44px. Square corners. Archivo 600, 15px, `word-spacing: 0.12em`.

| Variant | Rest | Hover | Use |
|---|---|---|---|
| `signal` | `bg-signal text-ink` | `bg-ink text-form` | The one primary action per view ("Start a load plan", "Run calculation") |
| `ink` | `bg-ink text-form` | `bg-signal text-ink` | Primary action in the header or where orange is already used |
| `outline` | `border-ink text-ink` | `bg-ink text-form` | Secondary action ("Try the calculator", "Sign in") |
| `plain` | `text-ink`, no padding | underline | Nav and inline links |

- Focus: `outline-2 outline-offset-3 outline-signal`.
- Active: `translate-y-px`.
- Disabled (to add for the app): `ink/35` text, `ink/10` fill, no hover change, `cursor-not-allowed`.
- Labels say what happens: "Start a load plan", not "Get started". The same action keeps the same label everywhere.
- In-page `#hash` links scroll through Lenis.

### 6.2 Logo (`components/brand/logo-mark.tsx`)

A container end view (ink frame with four corner castings) with one orange box seated in the bottom-left corner, which is where every packing plan starts.

- The frame uses `currentColor`; the box is always `#ff5b14`.
- Pair it with the "LoadIQ" wordmark in the stencil face at about 70% of the mark's height. Gap 10px.
- Sizes: 16px favicon, 32px sidebar, 36px landing header, 80px auth pages.
- Pass `title` only when the mark appears without the wordmark.
- Browser icons: `web/app/icon.svg` (mark on a `yard` tile), `apple-icon.png`, `favicon.ico`.
- Don't recolor the box, add effects, put it on orange, or rotate it.

### 6.3 Manifest readout (`load-sequence/load-readout.tsx`)

A two-line data strip under a 3D view: top rule `ink/25`, counter + item name + dimensions on the first line, container summary on the second. Mono for numbers, Archivo 600 for the item name. Use this pattern wherever the app shows "what is happening to the load right now" (step playback, operator loading).

### 6.4 Meter (`plan-view-illustration.tsx`)

Label (Archivo 600) and percentage (mono) on one line, an 8px bar with `ink/10` track and `signal` fill, then `value of limit` in small mono steel. For volume and weight use. Square ends.

### 6.5 Data rows

Cargo lines, roles, steps: a list with `divide-y divide-ink/10` or `border-t border-ink/15` per row. Swatch 12px square (kraft tones), label in Archivo, counts and volumes right-aligned in mono. Column header row in Label style with units ("boxes, m³").

### 6.6 Step list (`load-sequence/load-sequence.tsx`)

Ordered list, each step with a mono counter ("Step 2 of 4"), a Display S title and a body paragraph. The active step is marked with `.is-active`, which turns its counter `text-signal-text` (`#a8380a`). Plain `text-signal` would be 2.43:1 on yard and fail as text.

### 6.7 Badge

Square, 12px Archivo 600, `px-2 py-0.5`. Variants follow the status table in 2.4. Sentence case ("Completed", not "COMPLETED").

### 6.8 Forms (proposed for the app)

- Input: `bg-form`, 1px `field` border, 4px radius, 40px height, Archivo 16px (prevents iOS zoom).
- Focus: 2px `signal` outline, offset 2px. Don't also change the border color.
- Label above the field in Label style. Help text below in 14px steel.
- Error: `danger` border and a message below that says what to fix ("Length must be between 1 and 20000 mm").
- Dimension inputs: three fields with a mono `×` between them and the unit as a suffix (`mm`, `kg`).

### 6.9 Header and navigation

Landing: fixed 64px bar, `bg-yard/90` with blur, bottom rule `ink/15`, wordmark left, section links centre, Sign in + ink button right. Below `md`, a menu button opens a full-width panel (Escape closes it).

App sidebar (proposed): `form` surface, group titles in Label style sentence case (replacing the current all-caps "OVERVIEW", "WORKSPACE"), active item with `ink` text and a 3px `signal` left bar instead of a filled pill.

---

## 7. Icons

- **Set:** Solar, through `@iconify/tailwind4` and `@iconify-json/solar`. Icons are compiled into CSS at build time, so there are no runtime requests.
- **Usage:** `<span aria-hidden="true" className="icon-[solar--hamburger-menu-linear] size-5" />`.
- **Style:** `linear` for UI, `bold` only for filled controls (play button).
- **Size:** 20px in UI (`size-5`), 16px inline with small text, 28px for the mobile menu button.
- Icon-only buttons need an `aria-label`.
- The app still uses `lucide-react`. Replace it screen by screen during the migration; don't mix both sets on one screen.

---

## 8. 3D views and data visualization

From `web/components/landing/hero-scene/scene-colors.ts`:

| Element | Color |
|---|---|
| Boxes by item type | `#b8895a`, `#8a6845`, `#cfa874`, `#dcc29a` (kraft tones, cycle in item order) |
| Box being placed / selected | `#ff5b14` |
| Container frame and box edges | `#3b4248` (edges at 45% opacity) |
| Container floor | `#9aa09c` |
| Far walls | frame color at 7% opacity |

- Camera: orthographic, three-quarter view from the door end, container centred.
- Light: hemisphere + one directional light, no shadows.
- Product colors from the database (`color_hex`, such as `#3357FF`) are saturated and clash with this palette. **Open decision:** either map product colors to the kraft set in the app's `StuffingVisualizer`, or keep user colors but desaturate them. Until decided, keep product colors in the app and kraft tones on marketing pages.
- Charts (reports): ink for the main series, steel for comparison, signal only for the highlighted value. Follow the same "one orange thing" rule.

---

## 9. Motion

- **Library:** GSAP + ScrollTrigger. **Smooth scroll:** Lenis, landing page only (`motion/smooth-scroll.tsx`). Don't add Lenis to the app; app screens use native scroll.
- **Headings:** marketing headings reveal word by word (`SplitHeading` with `reveal`, 0.8s, `power3.out`, 50ms stagger). Screen readers get the unsplit text.
- **App transitions:** CSS only, 150–200ms, `ease-out`, for hover, open/close and state changes. No entrance animations on page load.
- **3D playback:** boxes drop in with a cubic ease-out, about 1.1m drop height. The current box is orange until the next one starts.
- **Reduced motion:** no Lenis, no scrubbing, no reveals, final states shown immediately. The 3D hero is replaced by a static poster (`public/landing/hero-load-poster.png`).
- 3D renders on demand only (after progress or pointer changes), pauses offscreen, caps pixel ratio at 1.75, and disposes all GPU resources on unmount.

---

## 10. Imagery

- **Product visuals:** built from real packing output. The hero load (`hero-scene/hero-load.ts`) is a real API result: 58 boxes, 20ft Standard, 54.93% volume. Reuse it rather than inventing numbers.
- **Photos:** documentary industrial photos (warehouses, container surfaces). No faces as subjects, no visible shipping-line or brand logos. Credit every photo in `components/landing/media-credits.ts` and the footer.
- **Treatment:** slight desaturation (`grayscale-[35%]`) or luminosity blend on `ink` for dark sections.
- **Not used:** stock people, illustrated characters, fake UI screenshots with invented data, customer logos without permission.

---

## 11. Voice and copy

- Sentence case for everything except the two uppercase marketing headings.
- Plain verbs, active voice: "Run the calculation", "Add item", "Replay it box by box".
- Name things the way users do: "load plan", "container", "box", "step". Not "job", "entity", "bin".
- An action keeps its name through the flow: the "Calculate" button leads to a "Calculated" result, not "Processed".
- Errors say what happened and how to fix it. No apologies, no "Oops".
- Empty states invite the next action: "No load plans yet. Create one to see it packed in 3D."

**Number formats:**
- Dimensions: `5898 × 2352 × 2393 mm` (multiplication sign, spaces, unit once at the end).
- Weight: `2,950 kg`, with en-US grouping and a space before the unit.
- Volume: `18.23 m³`, two decimals.
- Percentages: one decimal in UI (`54.9%`), two only when quoting an exact result (`54.93%`).
- Counters: `Box 29 of 58`, `Step 2 of 4`.

---

## 12. Accessibility baseline

- Text contrast at least 4.5:1, large text and UI boundaries at least 3:1 (see 2.3).
- Visible focus on every interactive element: 2px `signal` outline with offset.
- Touch targets at least 44 × 44px.
- A skip link to `#main` on every layout.
- Canvases and decorative images are `aria-hidden` with a text equivalent nearby (`sr-only` summary of the load).
- Color is never the only signal: statuses carry words, the current box also appears in the readout.
- Everything works with reduced motion and with JavaScript off (the landing renders the poster and the full step list).

---

## 13. Migrating the app

Suggested order. Each step should leave the app working.

1. **Tokens.** Move the `.landing` overrides in `globals.css` into `:root`, add the supporting colors from 2.2, and add the dark values from 2.5 under `.dark`. Delete the old blue `--primary` and `--chart-*` values or remap them.
2. **Fonts.** Load the three fonts in `web/app/layout.tsx` instead of `components/landing/fonts.ts`, and remove the unused Geist imports.
3. **Shared components.** Move `LandingLink` → `components/ui/` as the app button, plus `LogoMark`, `Meter`, the manifest readout and the badge. Update `components/ui/button.tsx` variants to match 6.1 and remove shadows.
4. **Shell.** Restyle `components/dashboard-layout.tsx` (sidebar per 6.9, header, skip link).
5. **Forms.** Inputs, selects, dimension inputs per 6.8. Fix the field border contrast.
6. **Plan detail.** `web/app/shipments/[id]/page.tsx` and `StuffingVisualizer`: readout, meters, cargo rows, and the product color decision from 8.
7. **Icons.** Replace `lucide-react` with Solar screen by screen.
8. **Auth pages.** Login and register: remove the gradient background and card shadow, use the stencil wordmark.
9. **Reports and tables.** Mono right-aligned numbers, status badges, row rules.

---

## 14. Source files

| What | Where |
|---|---|
| Color and font tokens | `web/app/globals.css` (`@theme` blocks and `.landing`) |
| Fonts | `web/components/landing/fonts.ts` |
| Button / link | `web/components/landing/landing-link.tsx` |
| Logo mark | `web/components/brand/logo-mark.tsx`, `web/app/icon.svg` |
| Split headings and reveals | `web/components/landing/motion/` |
| Smooth scroll (Lenis + GSAP) | `web/components/landing/motion/smooth-scroll.tsx` |
| 3D scene, colors, data | `web/components/landing/hero-scene/` |
| Readout | `web/components/landing/load-sequence/load-readout.tsx` |
| Meter, cargo rows, badge | `web/components/landing/plan-view-illustration.tsx` |
| Photo credits | `web/components/landing/media-credits.ts` |
