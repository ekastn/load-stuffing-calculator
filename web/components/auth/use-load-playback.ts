"use client"

import { useCallback, useEffect, useRef } from "react"
import { gsap } from "gsap"

import type { LoadScene } from "@/components/landing/hero-scene/load-scene"
import type { LoadReadoutHandle } from "@/components/landing/load-sequence/load-readout"

/** Long enough to read as a load sequence, short enough to finish while the visitor types. */
const PLAYBACK_SECONDS = 6

/**
 * Drives the auth panel's 3D container through the whole load once, then holds
 * on the finished container. The scene is created asynchronously by
 * `HeroCanvas`, and never at all without WebGL or with reduced motion.
 */
export function useLoadPlayback({ totalBoxes }: { totalBoxes: number }) {
  const sceneRef = useRef<LoadScene | null>(null)
  const readoutRef = useRef<LoadReadoutHandle>(null)
  const timelineRef = useRef<gsap.core.Timeline | null>(null)
  const hasPlayed = useRef(false)
  const progress = useRef({ step: 0 })

  const applyProgress = useCallback(() => {
    sceneRef.current?.setProgress(progress.current.step)
    readoutRef.current?.show(Math.ceil(progress.current.step))
  }, [])

  const handleSceneChange = useCallback(
    (scene: LoadScene | null) => {
      timelineRef.current?.kill()
      timelineRef.current = null
      sceneRef.current = scene
      if (!scene) return

      // A live scene only arrives without reduced motion, so this is the playable path.
      if (!hasPlayed.current) {
        hasPlayed.current = true
        progress.current.step = 0
        applyProgress()
        const timeline = gsap.timeline()
        timeline.to(progress.current, {
          step: totalBoxes,
          duration: PLAYBACK_SECONDS,
          ease: "none",
          onUpdate: applyProgress,
        })
        timelineRef.current = timeline
        return
      }

      // Context loss hands us a fresh scene; show the finished container instead of replaying.
      progress.current.step = totalBoxes
      applyProgress()
    },
    [applyProgress, totalBoxes],
  )

  useEffect(
    () => () => {
      timelineRef.current?.kill()
      timelineRef.current = null
    },
    [],
  )

  return { readoutRef, handleSceneChange }
}
