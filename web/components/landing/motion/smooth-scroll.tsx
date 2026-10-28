"use client"

import { createContext, useCallback, useContext, useEffect, useRef, type ReactNode } from "react"
import { gsap } from "gsap"
import { ScrollTrigger } from "gsap/ScrollTrigger"
import Lenis from "lenis"
import "lenis/dist/lenis.css"

gsap.registerPlugin(ScrollTrigger)

const ScrollToContext = createContext<((hash: string) => void) | null>(null)

/** Scrolls to an in-page anchor through Lenis when it is running. Null outside the provider. */
export function useScrollTo() {
  return useContext(ScrollToContext)
}

/**
 * Anchor scrolling without Lenis: used by the provider when the viewer prefers
 * reduced motion, and by hash links rendered outside the provider (the shared
 * header sits above it in the layout).
 */
export function scrollToHashNatively(hash: string) {
  const target = document.querySelector<HTMLElement>(hash)
  if (!target) return
  const prefersReducedMotion = window.matchMedia("(prefers-reduced-motion: reduce)").matches
  target.scrollIntoView({ behavior: prefersReducedMotion ? "auto" : "smooth" })
  window.history.pushState(null, "", hash)
  target.focus({ preventScroll: true })
}

/**
 * Runs Lenis as the page's only smooth-scroll engine, driven by the GSAP
 * ticker so ScrollTrigger reads the same scroll position. Skipped entirely
 * when the viewer prefers reduced motion.
 */
export function SmoothScroll({ children }: { children: ReactNode }) {
  const lenisRef = useRef<Lenis | null>(null)

  useEffect(() => {
    const refresh = () => ScrollTrigger.refresh()
    document.fonts?.ready.then(refresh)
    window.addEventListener("load", refresh)

    if (window.matchMedia("(prefers-reduced-motion: reduce)").matches) {
      return () => window.removeEventListener("load", refresh)
    }

    const lenis = new Lenis({ autoRaf: false })
    lenisRef.current = lenis
    lenis.on("scroll", ScrollTrigger.update)
    const tick = (time: number) => lenis.raf(time * 1000)
    gsap.ticker.add(tick)
    gsap.ticker.lagSmoothing(0)

    return () => {
      window.removeEventListener("load", refresh)
      gsap.ticker.remove(tick)
      lenis.destroy()
      lenisRef.current = null
    }
  }, [])

  const scrollTo = useCallback((hash: string) => {
    if (lenisRef.current) {
      const target = document.querySelector<HTMLElement>(hash)
      if (!target) return
      lenisRef.current.scrollTo(target)
      window.history.pushState(null, "", hash)
      target.focus({ preventScroll: true })
      return
    }
    scrollToHashNatively(hash)
  }, [])

  return <ScrollToContext.Provider value={scrollTo}>{children}</ScrollToContext.Provider>
}
