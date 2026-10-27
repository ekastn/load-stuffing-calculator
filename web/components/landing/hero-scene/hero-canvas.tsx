"use client"

import { useEffect, useRef, useState } from "react"
import Image from "next/image"

import { cn } from "@/lib/utils"
import { HERO_LOAD } from "./hero-load"
import type { LoadScene } from "./load-scene"
import { FLOOR_COLOR, FRAME_COLOR, HIGHLIGHT_COLOR, ITEM_COLORS } from "./scene-colors"
import { toSceneBoxes } from "./to-scene-boxes"

type HeroCanvasProps = {
  /** Receives the scene once it is live, and null when it goes away. */
  onSceneChange: (scene: LoadScene | null) => void
  className?: string
}

/**
 * Shows a static poster of the finished load, then swaps in the live
 * Three.js scene when WebGL is available and the viewer allows motion.
 */
export function HeroCanvas({ onSceneChange, className }: HeroCanvasProps) {
  const wrapperRef = useRef<HTMLDivElement>(null)
  const canvasRef = useRef<HTMLCanvasElement>(null)
  const [isLive, setIsLive] = useState(false)

  useEffect(() => {
    const wrapper = wrapperRef.current
    const canvas = canvasRef.current
    if (!wrapper || !canvas) return
    if (window.matchMedia("(prefers-reduced-motion: reduce)").matches || !supportsWebGL()) return

    let scene: LoadScene | null = null
    let isCancelled = false
    const cleanups: Array<() => void> = []

    const teardown = () => {
      cleanups.splice(0).forEach((cleanup) => cleanup())
      scene?.dispose()
      scene = null
      onSceneChange(null)
      setIsLive(false)
    }

    import("./load-scene").then(({ createLoadScene }) => {
      if (isCancelled) return
      scene = createLoadScene({
        canvas,
        container: HERO_LOAD.container,
        boxes: toSceneBoxes(HERO_LOAD),
        itemColors: ITEM_COLORS,
        highlightColor: HIGHLIGHT_COLOR,
        frameColor: FRAME_COLOR,
        floorColor: FLOOR_COLOR,
        onContextLost: teardown,
      })
      const liveScene = scene

      const resizeObserver = new ResizeObserver(([entry]) => {
        liveScene.resize(entry.contentRect.width, entry.contentRect.height)
      })
      resizeObserver.observe(wrapper)
      cleanups.push(() => resizeObserver.disconnect())

      const visibilityObserver = new IntersectionObserver(([entry]) => liveScene.setVisible(entry.isIntersecting))
      visibilityObserver.observe(wrapper)
      cleanups.push(() => visibilityObserver.disconnect())

      cleanups.push(trackPointer(liveScene))

      onSceneChange(liveScene)
      setIsLive(true)
    })

    return () => {
      isCancelled = true
      teardown()
    }
  }, [onSceneChange])

  return (
    <div ref={wrapperRef} className={cn("relative", className)}>
      <Image
        src="/landing/hero-load-poster.png"
        alt="A 20ft container cut away to show 58 crates, pallets and boxes stacked by the packing engine"
        fill
        priority
        sizes="(min-width: 1024px) 60vw, 100vw"
        className={cn("object-contain transition-opacity duration-500", isLive && "opacity-0")}
      />
      <canvas ref={canvasRef} aria-hidden="true" className={cn("absolute inset-0 h-full w-full", !isLive && "invisible")} />
    </div>
  )
}

function supportsWebGL() {
  try {
    const probe = document.createElement("canvas")
    return Boolean(probe.getContext("webgl2") ?? probe.getContext("webgl"))
  } catch {
    return false
  }
}

/** Tilts the camera with a fine pointer; re-centres on blur, hide, or leave. */
function trackPointer(scene: LoadScene) {
  if (!window.matchMedia("(pointer: fine)").matches) return () => {}

  let frameId = 0
  let latest = { x: 0, y: 0 }

  const handleMove = (event: PointerEvent) => {
    latest = { x: (event.clientX / window.innerWidth) * 2 - 1, y: (event.clientY / window.innerHeight) * 2 - 1 }
    if (frameId) return
    frameId = requestAnimationFrame(() => {
      frameId = 0
      scene.setPointer(latest.x, latest.y)
    })
  }
  const recentre = () => scene.setPointer(0, 0)
  const handleVisibility = () => document.hidden && recentre()

  window.addEventListener("pointermove", handleMove, { passive: true })
  window.addEventListener("blur", recentre)
  document.addEventListener("pointerleave", recentre)
  document.addEventListener("visibilitychange", handleVisibility)

  return () => {
    cancelAnimationFrame(frameId)
    window.removeEventListener("pointermove", handleMove)
    window.removeEventListener("blur", recentre)
    document.removeEventListener("pointerleave", recentre)
    document.removeEventListener("visibilitychange", handleVisibility)
  }
}
