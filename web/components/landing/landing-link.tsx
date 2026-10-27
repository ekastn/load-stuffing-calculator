"use client"

import Link from "next/link"
import type { ComponentProps, MouseEvent } from "react"

import { cn } from "@/lib/utils"
import { useScrollTo } from "./motion/smooth-scroll"

const variants = {
  signal: "bg-signal text-ink hover:bg-ink hover:text-form",
  ink: "bg-ink text-form hover:bg-signal hover:text-ink",
  outline: "border border-ink text-ink hover:bg-ink hover:text-form",
  plain: "text-ink underline-offset-4 hover:underline",
}

export type LandingLinkVariant = keyof typeof variants

type LandingLinkProps = Omit<ComponentProps<typeof Link>, "href"> & {
  href: string
  variant?: LandingLinkVariant
}

/**
 * The landing page's single link/button style. In-page `#hash` links scroll
 * through Lenis; everything else is a normal Next.js link.
 */
export function LandingLink({ href, variant = "signal", className, onClick, ...props }: LandingLinkProps) {
  const scrollTo = useScrollTo()
  const isHash = href.startsWith("#")

  const handleClick = (event: MouseEvent<HTMLAnchorElement>) => {
    onClick?.(event)
    if (!isHash || event.defaultPrevented || event.metaKey || event.ctrlKey) return
    event.preventDefault()
    scrollTo(href)
  }

  return (
    <Link
      href={href}
      onClick={handleClick}
      className={cn(
        "inline-flex min-h-11 items-center justify-center gap-2 px-5 font-landing text-[0.95rem] font-semibold [word-spacing:0.12em] transition-colors duration-200",
        "focus-visible:outline-2 focus-visible:outline-offset-3 focus-visible:outline-signal active:translate-y-px",
        variant === "plain" && "min-h-0 px-0",
        variants[variant],
        className,
      )}
      {...props}
    />
  )
}
