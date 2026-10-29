import Image from "next/image"

import { cn } from "@/lib/utils"
import {
  HERO_CONTAINER_NAME,
  HERO_LOAD,
  HERO_MAX_WEIGHT_KG,
  HERO_TOTAL_WEIGHT_KG,
} from "./hero-scene/hero-load"
import { ITEM_COLORS } from "./hero-scene/scene-colors"
import { summarizeLoad } from "./hero-scene/summarize-load"

const numberFormat = new Intl.NumberFormat("en-US", { maximumFractionDigits: 2, minimumFractionDigits: 2 })
const integerFormat = new Intl.NumberFormat("en-US")

function Meter({ label, value, limit, pct }: { label: string; value: string; limit: string; pct: number }) {
  return (
    <div>
      <div className="flex items-baseline justify-between gap-4">
        <span className="font-semibold">{label}</span>
        <span className="font-figures text-sm">{pct.toFixed(1)}%</span>
      </div>
      <div className="mt-2 h-2 bg-ink/10">
        <div className="h-full bg-signal" style={{ width: `${pct}%` }} />
      </div>
      <p className="mt-1.5 font-figures text-xs text-steel">
        {value} of {limit}
      </p>
    </div>
  )
}

/**
 * A drawn version of LoadIQ's plan view, filled with the same real packing
 * result as the hero. It is an illustration of the screen, not a screenshot.
 */
export function PlanViewIllustration({ className }: { className?: string }) {
  const summary = summarizeLoad(HERO_LOAD)
  const totalBoxes = HERO_LOAD.placements.length
  const weightPct = (HERO_TOTAL_WEIGHT_KG / HERO_MAX_WEIGHT_KG) * 100
  const { container } = HERO_LOAD

  return (
    <div
      role="img"
      aria-label={`Illustration of the LoadIQ plan view: a ${HERO_CONTAINER_NAME} container loaded with ${totalBoxes} boxes in 3D, a step slider, volume and weight used, and a cargo breakdown by product.`}
      className={cn("border border-ink/20 bg-form text-ink", className)}
    >
      <div aria-hidden="true">
        <div className="flex flex-wrap items-center justify-between gap-x-6 gap-y-2 border-b border-ink/15 px-5 py-4 md:px-6">
          <div className="flex items-center gap-3">
            <span className="font-figures text-sm font-medium">PLD-20261027-2259</span>
            <span className="bg-ink px-2 py-0.5 text-xs font-semibold text-form">Completed</span>
          </div>
          <span className="font-figures text-xs text-steel md:text-sm">
            {HERO_CONTAINER_NAME}, {container.length} × {container.width} × {container.height} mm
          </span>
        </div>

        <div className="grid lg:grid-cols-[1fr_22rem]">
          <div className="flex flex-col border-b border-ink/15 lg:border-b-0 lg:border-r">
            <div className="relative aspect-[1572/1100] w-full bg-yard">
              <Image
                src="/landing/hero-load-poster.png"
                alt=""
                fill
                loading="lazy"
                sizes="(min-width: 1440px) 1000px, (min-width: 1024px) 65vw, 100vw"
                className="object-contain"
              />
            </div>
            <div className="mt-auto flex items-center gap-4 border-t border-ink/15 px-5 py-4 md:px-6">
              <span className="inline-flex size-9 shrink-0 items-center justify-center bg-ink text-form">
                <span className="icon-[solar--play-bold] size-4" />
              </span>
              <div className="relative h-1 flex-1 bg-ink/15">
                <div className="absolute inset-y-0 left-0 w-full bg-ink" />
                <span className="absolute right-0 top-1/2 size-3.5 -translate-y-1/2 translate-x-1/2 border-2 border-ink bg-signal" />
              </div>
              <span className="shrink-0 font-figures text-sm">
                Step {totalBoxes} of {totalBoxes}
              </span>
            </div>
          </div>

          <div className="flex flex-col gap-7 px-5 py-6 md:px-6">
            <Meter
              label="Volume used"
              pct={summary.usedVolumePct}
              value={`${numberFormat.format(summary.usedVolumeM3)} m³`}
              limit={`${numberFormat.format(summary.containerVolumeM3)} m³`}
            />
            <Meter
              label="Weight"
              pct={weightPct}
              value={`${integerFormat.format(HERO_TOTAL_WEIGHT_KG)} kg`}
              limit={`${integerFormat.format(HERO_MAX_WEIGHT_KG)} kg`}
            />
            <div>
              <div className="flex items-baseline justify-between border-b border-ink/15 pb-2">
                <span className="font-semibold">Cargo</span>
                <span className="font-figures text-xs text-steel">boxes, m³</span>
              </div>
              <ul className="divide-y divide-ink/10">
                {summary.items.map((item, index) => (
                  <li key={item.label} className="flex items-center gap-3 py-2.5 text-sm">
                    <span className="size-3 shrink-0" style={{ backgroundColor: ITEM_COLORS[index % ITEM_COLORS.length] }} />
                    <span className="flex-1">{item.label}</span>
                    <span className="w-8 text-right font-figures">{item.count}</span>
                    <span className="w-14 text-right font-figures text-steel">{numberFormat.format(item.volumeM3)}</span>
                  </li>
                ))}
              </ul>
            </div>
          </div>
        </div>
      </div>
    </div>
  )
}
