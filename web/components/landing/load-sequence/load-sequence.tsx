"use client"

import { useCallback, useEffect, useRef } from "react"
import { gsap } from "gsap"
import { ScrollTrigger } from "gsap/ScrollTrigger"

import { HeroCanvas } from "../hero-scene/hero-canvas"
import { HERO_CONTAINER_NAME, HERO_LOAD, HERO_TOTAL_WEIGHT_KG, HERO_VOLUME_USED_PCT } from "../hero-scene/hero-load"
import type { LoadScene } from "../hero-scene/load-scene"
import { LandingLink } from "../landing-link"
import { SplitHeading } from "../motion/split-heading"
import { LOAD_STEPS } from "./load-steps"
import { LoadReadout, type LoadReadoutHandle } from "./load-readout"

gsap.registerPlugin(ScrollTrigger)

const TOTAL_BOXES = HERO_LOAD.placements.length
/** Boxes the intro drops in on page load; scrolling through the steps loads the rest. */
const INTRO_BOXES = 9

/**
 * The hero and the four-step walkthrough share one sticky 3D container:
 * the intro loads the first boxes, then scrolling the steps fills it.
 */
export function LoadSequence() {
  const rootRef = useRef<HTMLElement>(null)
  const stepsRef = useRef<HTMLOListElement>(null)
  const sceneRef = useRef<LoadScene | null>(null)
  const readoutRef = useRef<LoadReadoutHandle>(null)
  const progress = useRef({ intro: 0, scroll: 0 })

  const applyProgress = useCallback(() => {
    const scene = sceneRef.current
    if (!scene) return
    const step = Math.max(progress.current.intro, progress.current.scroll)
    scene.setProgress(step)
    readoutRef.current?.show(Math.ceil(step))
  }, [])

  const handleSceneChange = useCallback(
    (scene: LoadScene | null) => {
      sceneRef.current = scene
      if (scene) applyProgress()
      else readoutRef.current?.show(TOTAL_BOXES)
    },
    [applyProgress],
  )

  useEffect(() => {
    const root = rootRef.current
    const steps = stepsRef.current
    if (!root || !steps) return

    const ctx = gsap.context(() => {
      const media = gsap.matchMedia()
      media.add("(prefers-reduced-motion: no-preference)", () => {
        gsap.from(root.querySelectorAll("#hero-heading [data-word]"), {
          yPercent: 110,
          duration: 0.9,
          ease: "power3.out",
          stagger: 0.06,
        })
        gsap.to(progress.current, {
          intro: INTRO_BOXES,
          duration: 2.6,
          delay: 0.5,
          ease: "power1.inOut",
          onUpdate: applyProgress,
        })

        // Finish on the last step while the sticky canvas is still fully in view;
        // the list's bottom padding keeps it pinned past this point.
        ScrollTrigger.create({
          trigger: steps,
          start: "top 70%",
          endTrigger: steps.lastElementChild ?? steps,
          end: "bottom 70%",
          onUpdate: (self) => {
            progress.current.scroll = INTRO_BOXES + self.progress * (TOTAL_BOXES - INTRO_BOXES)
            applyProgress()
          },
        })
      })

      steps.querySelectorAll("li").forEach((step) => {
        ScrollTrigger.create({
          trigger: step,
          start: "top 60%",
          end: "bottom 60%",
          toggleClass: { targets: step, className: "is-active" },
        })
      })
    }, root)

    return () => ctx.revert()
  }, [applyProgress])

  return (
    <section ref={rootRef} aria-labelledby="hero-heading" className="relative">
      <div className="mx-auto grid max-w-[90rem] grid-cols-1 px-4 md:px-8 lg:grid-cols-12 lg:gap-10">
        <div className="sticky top-16 z-10 -mx-4 h-[44svh] border-b border-ink/15 bg-yard md:-mx-8 lg:order-last lg:col-span-7 lg:mx-0 lg:h-[calc(100svh-4rem)] lg:border-b-0">
          <HeroCanvas onSceneChange={handleSceneChange} className="h-full w-full" />
          <p className="sr-only">
            The animation replays a real result from LoadIQ&apos;s packing engine: {TOTAL_BOXES} boxes in a{" "}
            {HERO_CONTAINER_NAME} container, {HERO_VOLUME_USED_PCT}% of its volume used.
          </p>
          <LoadReadout
            ref={readoutRef}
            totalBoxes={TOTAL_BOXES}
            containerName={HERO_CONTAINER_NAME}
            volumeUsedPct={HERO_VOLUME_USED_PCT}
            totalWeightKg={HERO_TOTAL_WEIGHT_KG}
            className="absolute inset-x-4 bottom-3 md:inset-x-8 lg:inset-x-0 lg:bottom-8"
          />
        </div>

        <div className="lg:col-span-5">
          <div className="flex flex-col justify-center gap-8 pb-20 pt-10 lg:min-h-[calc(100svh-4rem)] lg:pb-16 lg:pt-16">
            <SplitHeading
              as="h1"
              id="hero-heading"
              text="Plan every box before it leaves the dock."
              className="font-stencil text-[clamp(3.25rem,7.4vw,7rem)] font-extrabold uppercase leading-[0.88] tracking-[0.01em] text-ink"
            />
            <p className="max-w-[34rem] text-lg leading-relaxed text-steel md:text-xl">
              LoadIQ turns an item list into a step-by-step 3D loading plan. Planners see what fits before anything
              ships, and the loading crew knows which box goes in next.
            </p>
            <div className="flex flex-wrap gap-3">
              <LandingLink href="/register">Start a load plan</LandingLink>
              <LandingLink href="#demo" variant="outline">
                Try the calculator
              </LandingLink>
            </div>
          </div>

          <SplitHeading
            id="how-it-works"
            reveal
            text="How a load plan comes together"
            className="scroll-mt-[calc(44svh+5rem)] font-stencil text-4xl font-bold text-ink md:text-5xl lg:scroll-mt-28"
          />
          <ol ref={stepsRef} className="pb-[35svh] pt-6">
            {LOAD_STEPS.map((step, index) => (
              <li
                key={step.title}
                className="group flex min-h-[55svh] flex-col justify-center gap-3 border-t border-ink/15 py-12 lg:min-h-[70svh]"
              >
                <span className="font-figures text-sm text-steel transition-colors duration-300 group-[.is-active]:text-signal-text">
                  Step {index + 1} of {LOAD_STEPS.length}
                </span>
                <h3 className="font-stencil text-3xl font-bold text-ink md:text-4xl">{step.title}</h3>
                <p className="max-w-[32rem] text-lg leading-relaxed text-steel">{step.body}</p>
              </li>
            ))}
          </ol>
        </div>
      </div>
    </section>
  )
}
