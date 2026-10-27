import Image from "next/image"

import { LandingLink } from "./landing-link"
import { MEDIA_CREDITS } from "./media-credits"
import { SplitHeading } from "./motion/split-heading"

export function ClosingSection() {
  const photo = MEDIA_CREDITS.containerDoors

  return (
    <section aria-labelledby="closing-heading" className="relative isolate overflow-hidden bg-ink text-form">
      <Image
        src={photo.src}
        alt=""
        fill
        loading="lazy"
        sizes="100vw"
        className="-z-10 object-cover opacity-45 mix-blend-luminosity"
      />
      <div className="mx-auto flex min-h-[80svh] max-w-[90rem] flex-col justify-end gap-8 px-4 py-24 md:px-8 lg:py-32">
        <SplitHeading
          id="closing-heading"
          reveal
          text="Plan the next container"
          className="max-w-[12ch] font-stencil text-6xl font-extrabold uppercase leading-[0.88] md:text-8xl lg:text-[9rem]"
        />
        <p className="max-w-[32rem] text-lg leading-relaxed text-form/80">
          Create a workspace, add your products and containers, and run your first calculation.
        </p>
        <div className="flex flex-wrap gap-3">
          <LandingLink href="/register">Start a load plan</LandingLink>
          <LandingLink href="/login" variant="outline" className="border-form text-form hover:bg-form hover:text-ink">
            Sign in
          </LandingLink>
        </div>
      </div>
    </section>
  )
}
