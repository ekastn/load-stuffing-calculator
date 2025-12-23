import Link from "next/link"

import { ArrowRight, Boxes, ClipboardList, Cuboid, LineChart, Package, PlayCircle } from "lucide-react"

import { ContainerIllustration } from "@/components/landing/container-illustration"
import { RedirectIfAuthed } from "@/components/redirect-if-authed"
import { Button } from "@/components/ui/button"
import { Card, CardContent, CardHeader, CardTitle } from "@/components/ui/card"

import type React from "react"

type Feature = {
  title: string
  description: string
  Icon: React.ComponentType<{ className?: string }>
}

function SectionTitle({ title, subtitle }: { title: string; subtitle: string }) {
  return (
    <div className="space-y-2">
      <h2 className="text-2xl font-semibold tracking-tight text-foreground sm:text-3xl">{title}</h2>
      <p className="max-w-2xl text-muted-foreground">{subtitle}</p>
    </div>
  )
}

function FeatureCard({ feature }: { feature: Feature }) {
  const Icon = feature.Icon

  return (
    <Card className="border-border/50 bg-card/50 shadow-sm backdrop-blur">
      <CardHeader className="space-y-3">
        <div className="flex h-10 w-10 items-center justify-center rounded-md bg-primary/10 text-primary">
          <Icon className="h-5 w-5" />
        </div>
        <CardTitle className="text-base">{feature.title}</CardTitle>
      </CardHeader>
      <CardContent>
        <p className="text-sm text-muted-foreground">{feature.description}</p>
      </CardContent>
    </Card>
  )
}

export default function LandingPage() {
  const features: Feature[] = [
    {
      title: "Build shipment plans",
      description: "Create a shipment, add items, and keep everything organized in one place.",
      Icon: ClipboardList,
    },
    {
      title: "Preview loading in 3D",
      description: "See how items fit inside the container before you start loading.",
      Icon: Cuboid,
    },
    {
      title: "Track space & weight",
      description: "Get a clear view of utilization so you can reduce wasted capacity.",
      Icon: LineChart,
    },
    {
      title: "Step-by-step loading guide",
      description: "Replay the loading sequence to help teams load faster and more consistently.",
      Icon: PlayCircle,
    },
    {
      title: "Manifests & reports",
      description: "Generate simple shipment documentation for operations and handoffs.",
      Icon: Package,
    },
    {
      title: "Works across teams",
      description: "Planning, warehouse, and supervisors can each focus on what they need.",
      Icon: Boxes,
    },
  ]

  return (
    <div className="min-h-screen bg-gradient-to-br from-primary/5 via-background to-background">
      <RedirectIfAuthed />

      <div className="mx-auto w-full max-w-6xl px-6">
        {/* Hero */}
        <header className="py-16 sm:py-20">
          <div className="grid items-center gap-10 lg:grid-cols-2">
            <div className="space-y-6">
              <p className="inline-flex w-fit items-center rounded-full border border-border/60 bg-background/60 px-3 py-1 text-xs font-medium text-muted-foreground">
                Load planning &amp; 3D stuffing preview
              </p>

              <h1 className="text-4xl font-bold tracking-tight text-foreground sm:text-5xl">
                Plan better loads. <span className="text-primary">Load faster.</span>
              </h1>

              <p className="max-w-xl text-lg text-muted-foreground">
                Create shipment plans, preview container loading in 3D, and follow a clear step-by-step loading
                sequence.
              </p>

              <div className="flex flex-wrap items-center gap-3">
                <Button asChild size="lg" className="gap-2">
                  <Link href="/login">
                    Sign in <ArrowRight className="h-4 w-4" />
                  </Link>
                </Button>

                <Button asChild size="lg" variant="outline">
                  <Link href="#how-it-works">How it works</Link>
                </Button>
              </div>

              <p className="text-xs text-muted-foreground">
                Designed for planning teams, warehouse operations, and supervisors.
              </p>
            </div>

            {/* Lightweight product preview (no WebGL) */}
            <div className="rounded-xl border border-border/60 bg-background/60 p-4 shadow-sm backdrop-blur">
              <div className="grid gap-4 sm:grid-cols-5">
                <div className="sm:col-span-3">
                  <div className="flex items-center justify-between">
                    <p className="text-xs font-medium text-muted-foreground">3D preview</p>
                    <p className="text-xs text-muted-foreground">Step 12 / 42</p>
                  </div>
                  <div className="mt-2 overflow-hidden rounded-lg bg-gradient-to-br from-slate-100 to-slate-50 ring-1 ring-border/60 dark:from-slate-950 dark:to-slate-900">
                    <div className="flex items-center justify-center px-3 py-6">
                      <ContainerIllustration className="h-auto w-full max-w-sm" />
                    </div>
                  </div>
                </div>

                <div className="sm:col-span-2">
                  <p className="text-xs font-medium text-muted-foreground">Utilization</p>

                  <div className="mt-2 space-y-3">
                    <div>
                      <div className="flex items-center justify-between text-xs text-muted-foreground">
                        <span>Space</span>
                        <span>78%</span>
                      </div>
                      <div className="mt-1 h-2 w-full rounded-full bg-slate-100 ring-1 ring-border/60">
                        <div className="h-2 w-[78%] rounded-full bg-primary" />
                      </div>
                    </div>

                    <div>
                      <div className="flex items-center justify-between text-xs text-muted-foreground">
                        <span>Weight</span>
                        <span>64%</span>
                      </div>
                      <div className="mt-1 h-2 w-full rounded-full bg-slate-100 ring-1 ring-border/60">
                        <div className="h-2 w-[64%] rounded-full bg-accent" />
                      </div>
                    </div>

                    <div className="rounded-lg border border-border/60 bg-background p-3">
                      <p className="text-xs font-medium text-foreground">Loading guidance</p>
                      <p className="mt-1 text-xs text-muted-foreground">Follow a clear sequence so loading stays consistent.</p>
                    </div>
                  </div>
                </div>
              </div>
            </div>
          </div>
        </header>

        {/* Features */}
        <section className="pb-16 sm:pb-20">
          <SectionTitle
            title="Everything you need to plan and load"
            subtitle="From planning to execution, keep shipments clear and consistent."
          />

          <div className="mt-8 grid gap-4 sm:grid-cols-2 lg:grid-cols-3">
            {features.map((feature) => (
              <FeatureCard key={feature.title} feature={feature} />
            ))}
          </div>
        </section>

        {/* How it works */}
        <section id="how-it-works" className="pb-16 sm:pb-20">
          <SectionTitle
            title="How it works"
            subtitle="A simple workflow that fits real shipping and warehouse operations."
          />

          <div className="mt-8 grid gap-4 md:grid-cols-2">
            <Card className="border-border/50 bg-card/50 shadow-sm backdrop-blur">
              <CardHeader>
                <CardTitle className="text-base">1) Create a shipment</CardTitle>
              </CardHeader>
              <CardContent className="text-sm text-muted-foreground">
                Start a new plan, select a container, and define what needs to be shipped.
              </CardContent>
            </Card>

            <Card className="border-border/50 bg-card/50 shadow-sm backdrop-blur">
              <CardHeader>
                <CardTitle className="text-base">2) Add items</CardTitle>
              </CardHeader>
              <CardContent className="text-sm text-muted-foreground">
                Add quantities and dimensions so the system can calculate a realistic load.
              </CardContent>
            </Card>

            <Card className="border-border/50 bg-card/50 shadow-sm backdrop-blur">
              <CardHeader>
                <CardTitle className="text-base">3) Generate a 3D load plan</CardTitle>
              </CardHeader>
              <CardContent className="text-sm text-muted-foreground">
                Run the calculation and review the 3D placement result along with utilization insights.
              </CardContent>
            </Card>

            <Card className="border-border/50 bg-card/50 shadow-sm backdrop-blur">
              <CardHeader>
                <CardTitle className="text-base">4) Load with guidance</CardTitle>
              </CardHeader>
              <CardContent className="text-sm text-muted-foreground">
                Use the step-by-step sequence to load consistently, reduce rework, and speed up handoffs.
              </CardContent>
            </Card>
          </div>
        </section>

        {/* Teams */}
        <section className="pb-16 sm:pb-20">
          <SectionTitle
            title="Made for everyday operations"
            subtitle="Different teams get the views they need without extra complexity."
          />

          <div className="mt-8 grid gap-4 md:grid-cols-3">
            <Card className="border-border/50 bg-card/50 shadow-sm backdrop-blur">
              <CardHeader>
                <CardTitle className="text-base">Planning teams</CardTitle>
              </CardHeader>
              <CardContent className="text-sm text-muted-foreground">
                Build shipment plans, optimize loads, and communicate a clear loading sequence.
              </CardContent>
            </Card>

            <Card className="border-border/50 bg-card/50 shadow-sm backdrop-blur">
              <CardHeader>
                <CardTitle className="text-base">Warehouse teams</CardTitle>
              </CardHeader>
              <CardContent className="text-sm text-muted-foreground">
                Follow step-by-step guidance while loading and reduce last-minute shuffling.
              </CardContent>
            </Card>

            <Card className="border-border/50 bg-card/50 shadow-sm backdrop-blur">
              <CardHeader>
                <CardTitle className="text-base">Supervisors</CardTitle>
              </CardHeader>
              <CardContent className="text-sm text-muted-foreground">
                Track shipment progress and reference simple reports for coordination.
              </CardContent>
            </Card>
          </div>
        </section>

        {/* Footer */}
        <footer className="border-t border-border/60 py-10">
          <div className="flex flex-col gap-4 sm:flex-row sm:items-center sm:justify-between">
            <div>
              <p className="text-sm font-medium text-foreground">Load &amp; Stuffing</p>
              <p className="text-xs text-muted-foreground">Container load planning with 3D preview and guided steps.</p>
            </div>

            <Button asChild className="w-fit">
              <Link href="/login">Sign in</Link>
            </Button>
          </div>
        </footer>
      </div>
    </div>
  )
}
