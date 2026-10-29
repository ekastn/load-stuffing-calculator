"use client"

import { useEffect, useState, type ReactNode } from "react"
import Link from "next/link"

import { LogoMark } from "@/components/brand/logo-mark"
import { AuthSidePanel } from "./auth-side-panel"

const WIDE_QUERY = "(min-width: 1024px)"

/** Form sections sit directly on the page, divided by a hairline rule. */
export const AUTH_FORM_CLASS = "mt-6 flex flex-col gap-5 border-t border-ink/15 pt-6"

/**
 * Full-height split shared by /login and /register: the form column with the
 * logo on top, and the ink panel flush to the right edge. The panel mounts only
 * on wide screens, so phones never download the 3D scene or GSAP.
 */
export function AuthShell({ children }: { children: ReactNode }) {
  const [isWide, setIsWide] = useState(false)

  useEffect(() => {
    const media = window.matchMedia(WIDE_QUERY)
    const update = () => setIsWide(media.matches)
    update()
    media.addEventListener("change", update)
    return () => media.removeEventListener("change", update)
  }, [])

  return (
    <div className="grid min-h-svh grid-cols-1 lg:grid-cols-[5fr_7fr]">
      <div className="flex flex-col px-4 py-6 md:px-8 lg:py-8">
        <Link
          href="/"
          aria-label="LoadIQ home"
          className="flex w-fit items-center gap-2.5 focus-visible:outline-2 focus-visible:outline-offset-4 focus-visible:outline-signal"
        >
          <LogoMark className="size-9" />
          <span aria-hidden="true" className="font-stencil text-2xl font-extrabold tracking-wide text-ink">
            LoadIQ
          </span>
        </Link>
        <div className="w-full max-w-[28rem] pb-8 pt-14 lg:pt-[12vh]">{children}</div>
      </div>
      <div className="hidden lg:block">{isWide ? <AuthSidePanel /> : null}</div>
    </div>
  )
}
