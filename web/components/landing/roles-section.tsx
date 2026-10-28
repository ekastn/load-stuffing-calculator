import Image from "next/image"

import { MEDIA_CREDITS } from "./media-credits"
import { SplitHeading } from "./motion/split-heading"

const ROLES = [
  {
    name: "Admin",
    body: "Sets up the workspace: who has access and in which role, the product catalog, and the container types you ship in.",
  },
  {
    name: "Planner",
    body: "Creates load plans, adds items, runs the calculation and adjusts the cargo until the load works.",
  },
  {
    name: "Operator",
    body: "Opens finished plans, replays them, and loads the container in the order the plan gives.",
  },
]

export function RolesSection() {
  const photo = MEDIA_CREDITS.warehouseForklift

  return (
    <section id="roles" tabIndex={-1} aria-labelledby="roles-heading" className="scroll-mt-16 border-t border-ink/15 py-24 outline-none lg:py-36">
      <div className="mx-auto grid max-w-[90rem] gap-12 px-4 md:px-8 lg:grid-cols-12 lg:gap-10">
        <div className="relative aspect-[4/3] overflow-hidden bg-steel lg:col-span-5 lg:aspect-auto lg:min-h-[36rem]">
          <Image
            src={photo.src}
            alt="A reach truck carrying a wrapped pallet past tall warehouse racks"
            fill
            loading="lazy"
            sizes="(min-width: 1024px) 40vw, 100vw"
            className="object-cover object-[45%_60%] grayscale-[35%]"
          />
        </div>

        <div className="lg:col-span-7 lg:pl-6">
          <SplitHeading
            id="roles-heading"
            reveal
            text="One plan, three jobs"
            className="font-stencil text-5xl font-bold leading-[0.95] text-ink md:text-7xl"
          />
          <p className="mt-6 max-w-[34rem] text-lg leading-relaxed text-steel">
            Everyone in the workspace works from the same plan. Each role sees what it needs and can change only what
            it should.
          </p>
          <dl className="mt-12">
            {ROLES.map((role) => (
              <div key={role.name} className="grid gap-2 border-t border-ink/15 py-7 md:grid-cols-[12rem_1fr] md:gap-8">
                <dt className="font-stencil text-3xl font-bold text-ink">{role.name}</dt>
                <dd className="max-w-[34rem] text-lg leading-relaxed text-steel">{role.body}</dd>
              </div>
            ))}
          </dl>
        </div>
      </div>
    </section>
  )
}
