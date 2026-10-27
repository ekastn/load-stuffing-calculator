"use client"

import { useImperativeHandle, useRef, type Ref } from "react"

import { cn } from "@/lib/utils"
import { HERO_LOAD } from "../hero-scene/hero-load"

export type LoadReadoutHandle = {
  /** Shows the box at this 1-based step; updates the DOM directly to avoid per-frame renders. */
  show: (step: number) => void
}

type LoadReadoutProps = {
  ref?: Ref<LoadReadoutHandle>
  totalBoxes: number
  containerName: string
  volumeUsedPct: number
  totalWeightKg: number
  className?: string
}

function describeBox(step: number) {
  const placement = HERO_LOAD.placements.find((candidate) => candidate.step === step)
  if (!placement) return { label: "Empty container", size: "" }
  const item = HERO_LOAD.items[placement.itemIndex]
  return { label: item.label, size: `${item.length} × ${item.width} × ${item.height} mm` }
}

/** The manifest line under the 3D container, tracking the box being loaded. */
export function LoadReadout({ ref, totalBoxes, containerName, volumeUsedPct, totalWeightKg, className }: LoadReadoutProps) {
  const counterRef = useRef<HTMLSpanElement>(null)
  const labelRef = useRef<HTMLSpanElement>(null)
  const sizeRef = useRef<HTMLSpanElement>(null)
  const shownStep = useRef(totalBoxes)
  const finalBox = describeBox(totalBoxes)

  useImperativeHandle(ref, () => ({
    show(step) {
      const clamped = Math.min(totalBoxes, Math.max(0, step))
      if (clamped === shownStep.current) return
      shownStep.current = clamped
      const box = describeBox(clamped)
      if (counterRef.current) counterRef.current.textContent = String(clamped)
      if (labelRef.current) labelRef.current.textContent = box.label
      if (sizeRef.current) sizeRef.current.textContent = box.size
    },
  }))

  const { container } = HERO_LOAD

  return (
    <div
      aria-hidden="true"
      className={cn(
        "grid grid-cols-[auto_1fr] items-baseline gap-x-6 gap-y-1 border-t border-ink/25 pt-3 text-sm text-ink",
        className,
      )}
    >
      <span className="font-figures text-base">
        Box <span ref={counterRef} className="inline-block min-w-[2ch] text-right">{totalBoxes}</span> of {totalBoxes}
      </span>
      <span className="truncate">
        <span ref={labelRef} className="font-semibold">{finalBox.label}</span>{" "}
        <span ref={sizeRef} className="font-figures text-steel">{finalBox.size}</span>
      </span>
      <span className="col-span-2 truncate font-figures text-xs text-steel md:text-sm">
        {containerName}, {container.length} × {container.width} × {container.height} mm
        <span className="hidden sm:inline">, {totalWeightKg.toLocaleString("en-US")} kg, {volumeUsedPct}% of volume used</span>
      </span>
    </div>
  )
}
