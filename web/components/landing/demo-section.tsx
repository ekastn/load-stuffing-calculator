import { TrialLoadCalculator } from "@/components/trial-load-calculator"
import { SplitHeading } from "./motion/split-heading"

export function DemoSection() {
  return (
    <section id="demo" tabIndex={-1} aria-labelledby="demo-heading" className="scroll-mt-16 bg-ink py-24 text-form outline-none lg:py-36">
      <div className="mx-auto max-w-[90rem] px-4 md:px-8">
        <div className="grid gap-6 lg:grid-cols-12 lg:items-end">
          <SplitHeading
            id="demo-heading"
            reveal
            text="Try it with your own numbers"
            className="font-stencil text-5xl font-bold leading-[0.95] md:text-7xl lg:col-span-7"
          />
          <p className="max-w-[30rem] text-lg leading-relaxed text-form/75 lg:col-span-5">
            No account needed. Pick a container, add a few items, and run the same packing engine your planners would
            use.
          </p>
        </div>

        <div className="mt-12 text-ink md:bg-form md:p-6 lg:mt-16">
          <TrialLoadCalculator />
        </div>
      </div>
    </section>
  )
}
