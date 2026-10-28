"use client"

import { useEffect, useId, useState } from "react"
import Link from "next/link"

import { LogoMark } from "@/components/brand/logo-mark"
import { cn } from "@/lib/utils"
import { LandingLink } from "./landing-link"
import { SECTION_LINKS } from "./section-links"

export function SiteHeader() {
  const [isMenuOpen, setIsMenuOpen] = useState(false)
  const menuId = useId()

  useEffect(() => {
    if (!isMenuOpen) return
    const closeOnEscape = (event: KeyboardEvent) => event.key === "Escape" && setIsMenuOpen(false)
    window.addEventListener("keydown", closeOnEscape)
    return () => window.removeEventListener("keydown", closeOnEscape)
  }, [isMenuOpen])

  const closeMenu = () => setIsMenuOpen(false)

  return (
    <header className="fixed inset-x-0 top-[var(--api-banner-offset,0px)] z-50 border-b border-ink/15 bg-yard/90 backdrop-blur-sm">
      <div className="mx-auto flex h-16 max-w-[90rem] items-center justify-between gap-6 px-4 md:px-8">
        <Link
          href="/"
          className="flex items-center gap-2.5 focus-visible:outline-2 focus-visible:outline-offset-4 focus-visible:outline-signal"
        >
          <LogoMark className="size-9" />
          <span className="font-stencil text-2xl font-extrabold tracking-wide text-ink">LoadIQ</span>
        </Link>

        <nav aria-label="Sections" className="hidden items-center gap-8 md:flex">
          {SECTION_LINKS.map((link) => (
            <LandingLink key={link.href} href={link.href} variant="plain" className="font-medium">
              {link.label}
            </LandingLink>
          ))}
        </nav>

        <div className="flex items-center gap-2">
          <LandingLink href="/login" variant="plain" className="hidden px-3 sm:inline-flex">
            Sign in
          </LandingLink>
          <LandingLink href="/register" variant="ink" className="hidden min-h-10 sm:inline-flex">
            Start a load plan
          </LandingLink>
          <button
            type="button"
            className="inline-flex size-11 items-center justify-center text-ink focus-visible:outline-2 focus-visible:outline-signal md:hidden"
            aria-expanded={isMenuOpen}
            aria-controls={menuId}
            aria-label={isMenuOpen ? "Close menu" : "Open menu"}
            onClick={() => setIsMenuOpen((open) => !open)}
          >
            <span aria-hidden="true" className={cn("size-7", isMenuOpen ? "icon-[solar--close-square-linear]" : "icon-[solar--hamburger-menu-linear]")} />
          </button>
        </div>
      </div>

      <div
        id={menuId}
        hidden={!isMenuOpen}
        className="border-t border-ink/15 bg-yard px-4 pb-6 pt-2 md:hidden"
      >
        <nav aria-label="Sections" className="flex flex-col">
          {SECTION_LINKS.map((link) => (
            <LandingLink
              key={link.href}
              href={link.href}
              variant="plain"
              className="min-h-12 justify-start border-b border-ink/10 text-lg"
              onClick={closeMenu}
            >
              {link.label}
            </LandingLink>
          ))}
        </nav>
        <div className="mt-5 grid grid-cols-2 gap-2">
          <LandingLink href="/login" variant="outline">
            Sign in
          </LandingLink>
          <LandingLink href="/register" variant="ink">
            Start a load plan
          </LandingLink>
        </div>
      </div>
    </header>
  )
}
