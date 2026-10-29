import { SplitHeading } from "./motion/split-heading"
import { PlanViewIllustration } from "./plan-view-illustration"

const SCREEN_PARTS = [
  { title: "3D replay", body: "Play the load forward or drag the slider to any step." },
  { title: "Volume and weight", body: "How much of the container's space and weight limit the plan uses." },
  { title: "Cargo breakdown", body: "Package count and volume for each product in the load." },
]

export function ProductSection() {
  return (
    <section aria-labelledby="product-heading" className="py-24 lg:py-36">
      <div className="mx-auto max-w-[90rem] px-4 md:px-8">
        <SplitHeading
          id="product-heading"
          reveal
          text="The whole load on one screen"
          className="max-w-[16ch] font-stencil text-5xl font-bold leading-[0.95] text-ink md:text-7xl"
        />

        <figure className="mt-12 lg:mt-16">
          <PlanViewIllustration />
          <figcaption className="mt-8 grid gap-8 border-t border-ink/20 pt-8 md:grid-cols-3">
            {SCREEN_PARTS.map((part) => (
              <div key={part.title}>
                <p className="font-semibold text-ink">{part.title}</p>
                <p className="mt-1 max-w-[28rem] leading-relaxed text-steel">{part.body}</p>
              </div>
            ))}
          </figcaption>
        </figure>
      </div>
    </section>
  )
}
