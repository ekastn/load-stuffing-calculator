import { RedirectIfAuthed } from "@/components/redirect-if-authed"
import { ClosingSection } from "@/components/landing/closing-section"
import { DemoSection } from "@/components/landing/demo-section"
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
      <RedirectIfAuthed />
      <SiteHeader />
      <div className="pt-16">
        <LoadSequence />
        <DemoSection />
        <ProductSection />
        <RolesSection />
        <ClosingSection />
      </div>
      <SiteFooter />
      <HeadingReveals />
    </SmoothScroll>
  )
}
