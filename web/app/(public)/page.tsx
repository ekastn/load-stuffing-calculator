import { RedirectIfAuthed } from "@/components/redirect-if-authed"
import { ClosingSection } from "@/components/landing/closing-section"
import { DemoSection } from "@/components/landing/demo-section"
import { landingFontVariables } from "@/components/landing/fonts"
import { LoadSequence } from "@/components/landing/load-sequence/load-sequence"
import { HeadingReveals } from "@/components/landing/motion/heading-reveals"
import { SmoothScroll } from "@/components/landing/motion/smooth-scroll"
import { ProductSection } from "@/components/landing/product-section"
import { RolesSection } from "@/components/landing/roles-section"
import { SiteFooter } from "@/components/landing/site-footer"
import { SiteHeader } from "@/components/landing/site-header"

export default function LandingPage() {
  return (
    <SmoothScroll>
      <div
        className={`landing ${landingFontVariables} min-h-screen bg-yard font-landing text-ink antialiased selection:bg-signal selection:text-ink`}
      >
        <RedirectIfAuthed />
        <a
          href="#main"
          className="sr-only z-[60] bg-signal px-4 py-3 font-semibold text-ink focus:not-sr-only focus:fixed focus:left-4 focus:top-4"
        >
          Skip to content
        </a>
        <SiteHeader />
        <main id="main" tabIndex={-1} className="pt-16 outline-none">
          <LoadSequence />
          <DemoSection />
          <ProductSection />
          <RolesSection />
          <ClosingSection />
        </main>
        <SiteFooter />
        <HeadingReveals />
      </div>
    </SmoothScroll>
  )
}
