"use client"

import { Suspense } from "react"
import dynamic from "next/dynamic"

import {
  HERO_CONTAINER_NAME,
  HERO_LOAD,
  HERO_TOTAL_WEIGHT_KG,
  HERO_VOLUME_USED_PCT,
} from "@/components/landing/hero-scene/hero-load"
import { LoadReadout } from "@/components/landing/load-sequence/load-readout"
import { useLoadPlayback } from "./use-load-playback"

// The canvas, Three.js and the load scene stay out of the auth pages' first load.
const HeroCanvas = dynamic(
  () => import("@/components/landing/hero-scene/hero-canvas").then((module) => module.HeroCanvas),
  { ssr: false },
)

const TOTAL_BOXES = HERO_LOAD.placements.length

/** The ink side of the auth frame: the hero load plays once, then holds on the full container. */
export function AuthSidePanel() {
  const { readoutRef, handleSceneChange } = useLoadPlayback({ totalBoxes: TOTAL_BOXES })

  return (
    <aside className="flex h-full min-h-svh flex-col bg-ink px-8 py-8 text-form">
      <div className="flex items-baseline justify-between gap-4 border-b border-form/30 pb-3">
        <p className="font-figures text-sm">{HERO_CONTAINER_NAME}</p>
        <p className="font-figures text-sm text-form/75">{TOTAL_BOXES} boxes</p>
      </div>

      <div className="flex min-h-0 flex-1 items-center py-6">
        {/* The poster and the live canvas share the poster's yard background, so the swap is invisible. */}
        <div className="relative aspect-[1572/1100] w-full bg-yard">
          <Suspense fallback={null}>
            <HeroCanvas onSceneChange={handleSceneChange} className="h-full w-full" />
          </Suspense>
        </div>
      </div>

      <p className="sr-only">
        A load plan from LoadIQ&apos;s packing engine: {TOTAL_BOXES} boxes in a {HERO_CONTAINER_NAME} container,{" "}
        {HERO_VOLUME_USED_PCT}% of its volume used, {HERO_TOTAL_WEIGHT_KG.toLocaleString("en-US")} kg total.
      </p>

      <LoadReadout
        ref={readoutRef}
        tone="dark"
        totalBoxes={TOTAL_BOXES}
        containerName={HERO_CONTAINER_NAME}
        volumeUsedPct={HERO_VOLUME_USED_PCT}
        totalWeightKg={HERO_TOTAL_WEIGHT_KG}
        className="mt-2"
      />
    </aside>
  )
}
