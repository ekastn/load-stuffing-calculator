import { ViewTransition, type ReactNode } from "react"

import { landingFontVariables } from "@/components/landing/fonts"

/**
 * Shared shell for the landing page and the auth screens: fonts, background,
 * skip link and the main landmark. Each page brings its own top of page.
 */
export default function PublicLayout({ children }: { children: ReactNode }) {
  return (
    <div
      className={`landing ${landingFontVariables} min-h-screen bg-yard font-landing text-ink antialiased selection:bg-signal selection:text-ink`}
    >
      <a
        href="#main"
        className="sr-only z-[60] bg-signal px-4 py-3 font-semibold text-ink focus:not-sr-only focus:fixed focus:left-4 focus:top-4"
      >
        Skip to content
      </a>
      <main id="main" tabIndex={-1} className="outline-none">
        <ViewTransition default="page-fade">{children}</ViewTransition>
      </main>
    </div>
  )
}
