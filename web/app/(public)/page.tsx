"use client"

import Link from "next/link"
import Image from "next/image"
import { ArrowRight, Box, Check, CheckCircle2, Cuboid, LayoutDashboard, LineChart, Package, PlayCircle, Settings2, Smartphone, Truck } from "lucide-react"

import { ProductPreview } from "@/components/landing/product-preview"
import { RedirectIfAuthed } from "@/components/redirect-if-authed"
import { TrialLoadCalculator } from "@/components/trial-load-calculator"
import { Button } from "@/components/ui/button"
import { Badge } from "@/components/ui/badge"

export default function LandingPage() {
  const handleScroll = (e: React.MouseEvent<HTMLAnchorElement>, id: string) => {
    e.preventDefault()
    const element = document.getElementById(id)
    if (element) {
      element.scrollIntoView({ behavior: "smooth", block: "start" })
      window.history.pushState(null, "", `#${id}`)
    }
  }

  return (
    <div className="flex min-h-screen flex-col bg-background selection:bg-primary/20">
      <RedirectIfAuthed />

      {/* Navigation Bar */}
      <header className="sticky top-0 z-50 w-full border-b border-border/40 bg-background/80 backdrop-blur-md">
        <div className="container mx-auto flex h-16 items-center justify-between px-4 md:px-6">
          <Link href="/" className="flex items-center gap-2">
            <div className="relative h-12 w-12 overflow-hidden rounded-md border border-border/10 bg-white shadow-xs">
              <Image src="/logo.png" alt="LoadIQ" fill className="object-contain p-0.5" />
            </div>
            <span className="font-sans text-lg font-bold tracking-tight text-foreground">LoadIQ</span>
          </Link>
          
          <nav className="hidden md:flex gap-6 text-sm font-medium text-muted-foreground">
            <a href="#features" onClick={(e) => handleScroll(e, "features")} className="transition-colors hover:text-foreground">Features</a>
            <a href="#how-it-works" onClick={(e) => handleScroll(e, "how-it-works")} className="transition-colors hover:text-foreground">How it works</a>
            <a href="#demo" onClick={(e) => handleScroll(e, "demo")} className="transition-colors hover:text-foreground">Interactive Demo</a>
          </nav>
          
          <div className="flex items-center gap-3">
            <Button asChild variant="ghost" size="sm" className="h-9 px-4 text-sm">
              <Link href="/login">Sign In</Link>
            </Button>
            <Button asChild size="sm" className="h-9 px-4 text-sm shadow-sm shadow-primary/10 hover:shadow-primary/20">
              <Link href="/login">Get Started</Link>
            </Button>
          </div>
        </div>
      </header>

      {/* Hero Section */}
      <section className="relative overflow-hidden pt-20 md:pt-32 lg:pt-40 pb-32">
        {/* Background Gradients & Grid */}
        <div className="absolute top-0 -left-4 w-72 h-72 bg-primary/20 rounded-full blur-[100px] -z-10 animate-pulse" />
        <div className="absolute top-20 right-0 w-96 h-96 bg-purple-500/20 rounded-full blur-[100px] -z-10" />
        <div className="absolute bottom-0 inset-x-0 h-64 bg-[linear-gradient(to_bottom,transparent,rgba(255,255,255,0.8)_50%,#fff)] pointer-events-none z-0" />
        <div className="absolute inset-0 bg-[url('/grid.svg')] opacity-[0.2] pointer-events-none" />

        <div className="container px-4 md:px-6 mx-auto">
          <div className="grid gap-12 lg:grid-cols-2 lg:gap-16 items-center">
            
            {/* Hero Content */}
            <div className="flex flex-col gap-8 items-start">
              <div className="space-y-4">
                <h1 className="text-4xl font-bold tracking-tight text-foreground sm:text-5xl lg:text-6xl">
                  Optimize every <br className="hidden lg:block" />
                  <span className="text-transparent bg-clip-text bg-gradient-to-r from-primary to-purple-600">
                     cubic meter.
                  </span>
                </h1>
                <p className="max-w-[500px] text-lg text-muted-foreground leading-relaxed">
                  The intelligent 3D container load planning platform for modern logistics.
                  Visualize, pack, and ship with enterprise-grade precision.
                </p>
              </div>

              <div className="flex flex-col sm:flex-row gap-4 w-full sm:w-auto relative">
                <Button asChild size="lg" className="h-14 px-8 text-lg shadow-lg shadow-primary/20 hover:shadow-primary/30 transition-all hover:scale-105">
                  <Link href="/login">
                    Start Optimization <ArrowRight className="ml-2 h-5 w-5" />
                  </Link>
                </Button>
                <Button asChild variant="outline" size="lg" className="h-14 px-8 text-lg backdrop-blur-sm bg-background/50 hover:bg-background/80">
                  <a href="#how-it-works" onClick={(e) => handleScroll(e, "how-it-works")}>How it works</a>
                </Button>

                {/* Hero Background 3D Container Bounds Accent */}
                <div className="absolute left-[110%] top-[-80%] opacity-[0.03] dark:opacity-[0.08] pointer-events-none -z-10 text-foreground hidden xl:block select-none">
                  <div className="border-2 border-dashed border-current p-6 rounded-xl flex flex-col gap-2 w-72">
                    <span className="font-mono text-xs font-bold tracking-widest text-primary">CONTAINER_VOLUME: 76.2m³</span>
                    <span className="font-mono text-[10px] text-muted-foreground">MAX_STACK_WT: 26,000kg</span>
                    <span className="font-mono text-[10px] text-muted-foreground">STRATEGY: Height_Prioritized</span>
                  </div>
                </div>
              </div>
            </div>

            {/* Hero Visual - Product Preview */}
            <div className="relative mx-auto w-full max-w-[600px] lg:max-w-none">
               <ProductPreview />
            </div>
            
          </div>
        </div>
      </section>

      {/* Trial Calculator Section */}
      <section id="demo" className="py-32 bg-muted/30 dark:bg-zinc-900/10 border-y border-border/40 relative overflow-hidden isolate">
         {/* Background Gradients & Grid */}
         <div className="absolute top-1/2 left-1/4 w-[500px] h-[500px] bg-primary/5 rounded-full blur-[120px] -z-10" />
         <div className="absolute top-1/3 right-1/4 w-[400px] h-[400px] bg-purple-500/5 rounded-full blur-[120px] -z-10" />
         <div className="absolute inset-0 bg-[url('/grid.svg')] opacity-[0.25] pointer-events-none -z-10" />
         
         {/* Technical Axle Coordinate helper graphic */}
         <svg className="absolute right-10 top-10 opacity-[0.03] dark:opacity-[0.07] pointer-events-none text-foreground -z-10 hidden lg:block" width="200" height="200" viewBox="0 0 100 100">
           <path d="M50 50 L80 35 M50 50 L20 35 M50 50 L50 90" stroke="currentColor" strokeWidth="1.5" strokeDasharray="3,3" />
           <text x="82" y="32" fontSize="5" fill="currentColor" fontWeight="bold">X</text>
           <text x="14" y="32" fontSize="5" fill="currentColor" fontWeight="bold">Y</text>
           <text x="48" y="96" fontSize="5" fill="currentColor" fontWeight="bold">Z</text>
         </svg>

         {/* Floating decorative isometric package box */}
         <div className="absolute left-[6%] top-[20%] opacity-[0.08] dark:opacity-[0.15] pointer-events-none -z-10 animate-bounce hidden md:block" style={{ animationDuration: '8s' }}>
            <svg width="60" height="60" viewBox="0 0 40 40" fill="none" className="text-primary">
              <path d="M20 5 L35 13 L35 29 L20 37 L5 29 L5 13 Z" stroke="currentColor" strokeWidth="1.5" strokeLinejoin="round" />
              <path d="M20 5 L20 37 M5 13 L20 21 L35 13" stroke="currentColor" strokeWidth="1" strokeLinejoin="round" />
            </svg>
         </div>

         {/* Blueprint specs info tag */}
         <div className="absolute right-[5%] bottom-[10%] opacity-[0.08] dark:opacity-[0.15] pointer-events-none -z-10 font-mono text-[9px] select-none text-foreground hidden xl:block">
            <div className="border border-dashed border-border p-3 rounded-lg flex flex-col gap-1">
              <div>DIM_BOUNDS: [12000, 2350, 2690]</div>
              <div className="flex gap-2">
                <span className="text-primary font-bold">X-AXIS</span>
                <span className="text-purple-500 font-bold">Y-AXIS</span>
                <span className="text-emerald-500 font-bold">Z-AXIS</span>
              </div>
            </div>
         </div>

         <div className="container px-4 md:px-6 mx-auto relative z-10">
           <div className="text-center mb-12">
              <h2 className="text-3xl font-bold tracking-tight sm:text-4xl">Try our packing engine instantly.</h2>
              <p className="mt-4 text-muted-foreground text-lg">Interactive Demo • No account required</p>
           </div>
           
           <div className="max-w-6xl mx-auto transform transition-all hover:scale-[1.01] duration-500">
             <TrialLoadCalculator />
           </div>
         </div>
      </section>

      {/* Features Section - Bento Grid */}
      <section id="features" className="py-24 bg-background relative overflow-hidden isolate">
        {/* Ambient Glows */}
        <div className="absolute top-10 left-10 w-[500px] h-[500px] bg-purple-500/5 dark:bg-purple-500/10 rounded-full blur-[120px] -z-10" />
        <div className="absolute bottom-10 right-10 w-[500px] h-[500px] bg-emerald-500/5 dark:bg-emerald-500/10 rounded-full blur-[120px] -z-10" />
        
        {/* Dot Pattern Background */}
        <div className="absolute inset-0 bg-[image:radial-gradient(#e5e7eb_1px,transparent_1px)] [background-size:20px_20px] opacity-100 pointer-events-none -z-10" />

        {/* Blueprint Axes Graphic */}
        <svg className="absolute left-6 top-1/2 opacity-[0.02] dark:opacity-[0.05] pointer-events-none text-foreground -z-10 hidden xl:block" width="160" height="160" viewBox="0 0 100 100">
          <path d="M50 50 L80 35 M50 50 L20 35 M50 50 L50 90" stroke="currentColor" strokeWidth="1.5" strokeDasharray="3,3" />
        </svg>

        <div className="container px-4 md:px-6 mx-auto">
          <div className="text-center mb-16 space-y-4">
             <h2 className="text-3xl font-bold tracking-tight sm:text-5xl">Built for scale.</h2>
             <p className="text-muted-foreground text-lg max-w-2xl mx-auto">
               Everything you need to manage complex shipping operations in one unified platform.
             </p>
          </div>

          <div className="grid grid-cols-1 md:grid-cols-3 gap-6 max-w-6xl mx-auto">
            
            {/* Feature 1 - Large */}
            <div className="md:col-span-2 rounded-3xl border border-border/50 bg-card p-8 flex flex-col justify-between hover:border-primary/50 transition-colors group relative overflow-hidden">
               <div className="absolute inset-0 bg-gradient-to-br from-primary/5 to-transparent opacity-0 group-hover:opacity-100 transition-opacity" />
               <div className="relative z-10 space-y-4">
                 <div className="h-12 w-12 rounded-xl bg-primary/10 flex items-center justify-center text-primary mb-4">
                    <LayoutDashboard className="h-6 w-6" />
                 </div>
                 <h3 className="text-2xl font-bold">Intelligent Workspace</h3>
                 <p className="text-muted-foreground max-w-md">Manage multiple shipments, users, and cargo profiles from a single command center. Real-time syncing keeps your whole team aligned.</p>
               </div>
               <div className="mt-6 w-full rounded-2xl border border-border/40 bg-background/60 p-4 font-sans text-xs shadow-inner overflow-hidden select-none">
                  {/* Dashboard header mockup */}
                  <div className="flex items-center justify-between pb-3 border-b border-border/40 mb-3">
                    <div className="flex items-center gap-2">
                      <span className="font-semibold text-foreground">Active Load Plans</span>
                      <span className="px-1.5 py-0.5 rounded bg-primary/10 text-primary text-[10px] font-medium">Live</span>
                    </div>
                    <div className="flex gap-1.5">
                      <div className="w-2 h-2 rounded-full bg-border" />
                      <div className="w-2 h-2 rounded-full bg-border" />
                      <div className="w-2 h-2 rounded-full bg-border" />
                    </div>
                  </div>
                  {/* Shipment items list mockup */}
                  <div className="space-y-2">
                    {/* Header Row */}
                    <div className="grid grid-cols-12 gap-2 px-3 py-1.5 text-[9px] font-bold text-muted-foreground uppercase tracking-wider">
                      <div className="col-span-5">Shipment</div>
                      <div className="col-span-3">Container</div>
                      <div className="col-span-2">Volume</div>
                      <div className="col-span-2 text-right">Status</div>
                    </div>
                    
                    {/* Row 1 */}
                    <div className="grid grid-cols-12 gap-2 items-center p-3 rounded-xl bg-card border border-border/30 hover:border-primary/20 hover:shadow-2xs transition-all">
                      <div className="col-span-5 flex items-center gap-2.5 min-w-0">
                        <div className="p-1.5 rounded-lg bg-emerald-500/10 text-emerald-500 shrink-0">
                          <Truck className="h-3.5 w-3.5" />
                        </div>
                        <div className="min-w-0">
                          <div className="font-semibold text-foreground truncate">Hamburg #LPC-09</div>
                          <div className="text-[10px] text-muted-foreground">Dest: DEHAM</div>
                        </div>
                      </div>
                      <div className="col-span-3 text-muted-foreground">40ft HC</div>
                      <div className="col-span-2">
                        <div className="space-y-1">
                          <div className="font-bold text-foreground">94.2%</div>
                          <div className="h-1.5 w-full bg-muted rounded-full overflow-hidden">
                            <div className="h-full bg-emerald-500 rounded-full" style={{ width: "94.2%" }} />
                          </div>
                        </div>
                      </div>
                      <div className="col-span-2 text-right">
                        <span className="inline-flex px-2 py-0.5 rounded-full text-[9px] font-semibold bg-emerald-500/10 text-emerald-600 dark:text-emerald-400 border border-emerald-500/25">
                          Packed
                        </span>
                      </div>
                    </div>

                    {/* Row 2 */}
                    <div className="grid grid-cols-12 gap-2 items-center p-3 rounded-xl bg-card border border-border/30 hover:border-primary/20 hover:shadow-2xs transition-all">
                      <div className="col-span-5 flex items-center gap-2.5 min-w-0">
                        <div className="p-1.5 rounded-lg bg-blue-500/10 text-blue-500 shrink-0">
                          <Truck className="h-3.5 w-3.5" />
                        </div>
                        <div className="min-w-0">
                          <div className="font-semibold text-foreground truncate">Tokyo Express #LPC-10</div>
                          <div className="text-[10px] text-muted-foreground">Dest: JPTYO</div>
                        </div>
                      </div>
                      <div className="col-span-3 text-muted-foreground">20ft Standard</div>
                      <div className="col-span-2">
                        <div className="space-y-1">
                          <div className="font-bold text-foreground">78.5%</div>
                          <div className="h-1.5 w-full bg-muted rounded-full overflow-hidden">
                            <div className="h-full bg-blue-500 rounded-full" style={{ width: "78.5%" }} />
                          </div>
                        </div>
                      </div>
                      <div className="col-span-2 text-right">
                        <span className="inline-flex px-2 py-0.5 rounded-full text-[9px] font-semibold bg-blue-500/10 text-blue-600 dark:text-blue-400 border border-blue-500/25">
                          Draft
                        </span>
                      </div>
                    </div>

                    {/* Row 3 */}
                    <div className="grid grid-cols-12 gap-2 items-center p-3 rounded-xl bg-card border border-border/30 hover:border-primary/20 hover:shadow-2xs transition-all">
                      <div className="col-span-5 flex items-center gap-2.5 min-w-0">
                        <div className="p-1.5 rounded-lg bg-purple-500/10 text-purple-500 shrink-0">
                          <Truck className="h-3.5 w-3.5 animate-pulse" />
                        </div>
                        <div className="min-w-0">
                          <div className="font-semibold text-foreground truncate">Rotterdam #LPC-11</div>
                          <div className="text-[10px] text-muted-foreground">Dest: NLRTM</div>
                        </div>
                      </div>
                      <div className="col-span-3 text-muted-foreground">40ft HC</div>
                      <div className="col-span-2">
                        <div className="space-y-1">
                          <div className="font-bold text-purple-500 animate-pulse">Running</div>
                          <div className="h-1.5 w-full bg-muted rounded-full overflow-hidden">
                            <div className="h-full bg-gradient-to-r from-purple-500 to-indigo-500 rounded-full animate-pulse" style={{ width: "55%" }} />
                          </div>
                        </div>
                      </div>
                      <div className="col-span-2 text-right">
                        <span className="inline-flex px-2 py-0.5 rounded-full text-[9px] font-semibold bg-purple-500/10 text-purple-600 dark:text-purple-400 border border-purple-500/25 animate-pulse">
                          Packing
                        </span>
                      </div>
                    </div>

                    {/* Row 4 */}
                    <div className="grid grid-cols-12 gap-2 items-center p-3 rounded-xl bg-card border border-border/30 hover:border-primary/20 hover:shadow-2xs transition-all">
                      <div className="col-span-5 flex items-center gap-2.5 min-w-0">
                        <div className="p-1.5 rounded-lg bg-emerald-500/10 text-emerald-500 shrink-0">
                          <Truck className="h-3.5 w-3.5" />
                        </div>
                        <div className="min-w-0">
                          <div className="font-semibold text-foreground truncate">Singapore #LPC-12</div>
                          <div className="text-[10px] text-muted-foreground">Dest: SGSIN</div>
                        </div>
                      </div>
                      <div className="col-span-3 text-muted-foreground">40ft HC</div>
                      <div className="col-span-2">
                        <div className="space-y-1">
                          <div className="font-bold text-foreground">91.8%</div>
                          <div className="h-1.5 w-full bg-muted rounded-full overflow-hidden">
                            <div className="h-full bg-emerald-500 rounded-full" style={{ width: "91.8%" }} />
                          </div>
                        </div>
                      </div>
                      <div className="col-span-2 text-right">
                        <span className="inline-flex px-2 py-0.5 rounded-full text-[9px] font-semibold bg-emerald-500/10 text-emerald-600 dark:text-emerald-400 border border-emerald-500/25">
                          Packed
                        </span>
                      </div>
                    </div>

                    {/* Row 5 */}
                    <div className="grid grid-cols-12 gap-2 items-center p-3 rounded-xl bg-card border border-border/30 hover:border-primary/20 hover:shadow-2xs transition-all">
                      <div className="col-span-5 flex items-center gap-2.5 min-w-0">
                        <div className="p-1.5 rounded-lg bg-blue-500/10 text-blue-500 shrink-0">
                          <Truck className="h-3.5 w-3.5" />
                        </div>
                        <div className="min-w-0">
                          <div className="font-semibold text-foreground truncate">New York #LPC-13</div>
                          <div className="text-[10px] text-muted-foreground">Dest: USNYC</div>
                        </div>
                      </div>
                      <div className="col-span-3 text-muted-foreground">40ft Standard</div>
                      <div className="col-span-2">
                        <div className="space-y-1">
                          <div className="font-bold text-foreground">85.3%</div>
                          <div className="h-1.5 w-full bg-muted rounded-full overflow-hidden">
                            <div className="h-full bg-blue-500 rounded-full" style={{ width: "85.3%" }} />
                          </div>
                        </div>
                      </div>
                      <div className="col-span-2 text-right">
                        <span className="inline-flex px-2 py-0.5 rounded-full text-[9px] font-semibold bg-blue-500/10 text-blue-600 dark:text-blue-400 border border-blue-500/25">
                          Draft
                        </span>
                      </div>
                    </div>
                  </div>
               </div>
            </div>

            {/* Feature 2 - Tall/Small */}
            <div className="md:col-span-1 rounded-3xl border border-border/50 bg-card p-8 hover:border-primary/50 transition-colors group relative overflow-hidden">
               <div className="absolute inset-0 bg-gradient-to-br from-primary/5 to-transparent opacity-0 group-hover:opacity-100 transition-opacity" />
               <div className="h-12 w-12 rounded-xl bg-primary/10 flex items-center justify-center text-primary mb-6">
                  <Smartphone className="h-6 w-6" />
               </div>
               <h3 className="text-xl font-bold mb-2">Mobile Ready</h3>
               <p className="text-muted-foreground mb-8">Access plans and loading guides from any device on the warehouse floor.</p>
               <div className="flex justify-center mt-6">
                  {/* Phone frame mockup */}
                  <div className="relative w-56 h-[400px] border-[6px] border-zinc-800 dark:border-zinc-700 rounded-[32px] bg-background shadow-2xl overflow-hidden flex flex-col font-sans select-none">
                    {/* Speaker notch */}
                    <div className="absolute top-0 inset-x-0 h-4 flex justify-center z-20">
                      <div className="w-16 h-3 bg-zinc-800 dark:bg-zinc-700 rounded-b-xl" />
                    </div>
                    {/* Status bar */}
                    <div className="h-7 pt-1 px-4 flex justify-between items-center text-[8px] text-muted-foreground z-10 border-b border-border/30">
                      <span>09:41</span>
                      <div className="flex gap-1">
                        <div className="w-2 h-1.5 bg-muted-foreground rounded-xs" />
                        <div className="w-2 h-2.5 bg-muted-foreground rounded-xs" />
                      </div>
                    </div>
                    {/* Phone Body */}
                    <div className="flex-1 p-3 pb-1 flex flex-col justify-between">
                      {/* Step Indicator Header */}
                      <div className="flex items-center justify-between">
                        <span className="text-[10px] font-bold text-foreground">Step 8 of 24</span>
                        <Badge variant="outline" className="text-[8px] px-1.5 py-0 border-emerald-500/30 text-emerald-600 bg-emerald-500/5 font-semibold">
                          92% Stacked
                        </Badge>
                      </div>
                      
                      {/* Placement Instruction Card */}
                      <div className="my-1.5 p-2.5 rounded-xl border border-border/40 bg-muted/30 space-y-1.5">
                        <div className="flex items-center gap-1.5 text-foreground text-[10px] font-bold">
                          <Package className="h-3 w-3 text-primary" />
                          <span>Ergonomic Desk Frame</span>
                        </div>
                        <div className="grid grid-cols-2 gap-x-2 gap-y-1 text-[8px] text-muted-foreground">
                          <div>
                            <span className="block text-zinc-400">Dimensions</span>
                            <span className="font-semibold text-foreground">1.2 × 0.8 × 0.7m</span>
                          </div>
                          <div>
                            <span className="block text-zinc-400">Weight</span>
                            <span className="font-semibold text-foreground">24.5 kg</span>
                          </div>
                          <div className="col-span-2 border-t border-border/30 pt-1 mt-1">
                            <span className="block text-zinc-400">Position Placement</span>
                            <span className="font-bold text-primary">X: 2.4m, Y: 0.8m, Z: 0.0m</span>
                          </div>
                        </div>
                      </div>

                      {/* Direction Arrow visualization */}
                      <div className="flex items-center justify-center">
                        <div className="flex items-center gap-1 text-[8px] text-emerald-600 dark:text-emerald-400 bg-emerald-500/10 px-2.5 py-1 rounded-full font-semibold">
                          <CheckCircle2 className="h-3 w-3" />
                          <span>Load Back-Left Corner</span>
                        </div>
                      </div>

                      {/* Phone Confirm Button */}
                      <button type="button" className="w-full py-2 bg-primary hover:bg-primary/95 text-primary-foreground text-[10px] font-bold rounded-lg shadow-sm flex items-center justify-center gap-1 cursor-pointer">
                        <Check className="h-3 w-3 stroke-[3]" />
                        Confirm Loaded
                      </button>
                      
                      {/* Home Indicator */}
                      <div className="h-2 flex items-center justify-center pt-1 mt-1">
                        <div className="w-16 h-1 bg-zinc-800 dark:bg-zinc-700 rounded-full" />
                      </div>
                    </div>
                  </div>
               </div>
            </div>

            {/* Feature 3 - Small */}
            <div className="md:col-span-1 rounded-3xl border border-border/50 bg-card p-8 hover:border-primary/50 transition-colors group relative overflow-hidden animate-all duration-350 flex flex-col justify-between">
               <div className="absolute inset-0 bg-gradient-to-br from-primary/5 to-transparent opacity-0 group-hover:opacity-100 transition-opacity" />
               <div className="relative z-10 space-y-4">
                  <div className="h-12 w-12 rounded-xl bg-primary/10 flex items-center justify-center text-primary mb-6">
                     <Cuboid className="h-6 w-6" />
                  </div>
                  <h3 className="text-xl font-bold mb-2">3D Visualization</h3>
                  <p className="text-muted-foreground">Rotate, zoom, and inspect every layer of your cargo plan before loading.</p>
                  
                  {/* Decorative Camera Path Axis */}
                  <svg className="absolute -right-8 -top-8 opacity-[0.03] dark:opacity-[0.08] pointer-events-none text-foreground -z-10" width="120" height="120" viewBox="0 0 100 100">
                    <circle cx="50" cy="50" r="40" stroke="currentColor" strokeWidth="1" strokeDasharray="3,3" />
                    <circle cx="50" cy="50" r="20" stroke="currentColor" strokeWidth="1" />
                    <path d="M50 10 L50 90 M10 50 L90 50" stroke="currentColor" strokeWidth="0.5" />
                  </svg>
               </div>
               <div className="mt-6 relative w-full h-36 border border-border/30 rounded-2xl bg-background/50 overflow-hidden shadow-inner flex items-center justify-center">
                  <Image 
                     src="/visual_3d.png" 
                     alt="3D Stuffing Plan Viewer" 
                     fill 
                     className="object-cover rounded-xl transition-transform duration-500 group-hover:scale-105" 
                  />
               </div>
            </div>

            {/* Feature 4 - Wide */}
            <div className="md:col-span-2 rounded-3xl border border-border/50 bg-card p-8 flex flex-col md:flex-row items-center gap-8 hover:border-primary/50 transition-colors group relative overflow-hidden">
               <div className="flex-1 space-y-4 relative z-10">
                 <div className="h-12 w-12 rounded-xl bg-primary/10 flex items-center justify-center text-primary mb-4">
                    <LineChart className="h-6 w-6" />
                 </div>
                 <h3 className="text-2xl font-bold">Analytics & Reports</h3>
                 <p className="text-muted-foreground">Generate comprehensive manifests, weight reports, and utilization stats to identify cost-saving opportunities.</p>
                 
                 {/* Decorative chart vector */}
                 <svg className="absolute -left-6 -bottom-6 opacity-[0.03] dark:opacity-[0.08] pointer-events-none text-foreground -z-10" width="160" height="120" viewBox="0 0 100 80">
                   <path d="M10 70 L30 40 L50 60 L70 20 L90 30" fill="none" stroke="currentColor" strokeWidth="2" />
                   <line x1="10" y1="70" x2="90" y2="70" stroke="currentColor" strokeWidth="1" />
                 </svg>
               </div>
               <div className="flex-1 w-full mt-4 md:mt-0 grid grid-cols-2 gap-4 h-48">
                  <div className="relative h-full border border-border/30 rounded-2xl overflow-hidden shadow-inner bg-background/50">
                     <Image 
                        src="/visual_summary.png" 
                        alt="Loading Metrics Summary" 
                        fill 
                        className="object-contain p-2 rounded-xl transition-transform duration-500 group-hover:scale-105" 
                     />
                  </div>
                  <div className="relative h-full border border-border/30 rounded-2xl overflow-hidden shadow-inner bg-background/50">
                     <Image 
                        src="/visual_brakedown.png" 
                        alt="Cargo Breakdown Details" 
                        fill 
                        className="object-contain p-2 rounded-xl transition-transform duration-500 group-hover:scale-105" 
                     />
                  </div>
               </div>
            </div>

          </div>
        </div>
      </section>

      {/* Workflow Section - Timeline */}
      <section id="how-it-works" className="py-24 bg-muted/40 dark:bg-zinc-900/20 border-y border-border/40 relative overflow-hidden isolate">
        {/* Ambient Glows */}
        <div className="absolute -top-10 right-20 w-[450px] h-[450px] bg-primary/5 rounded-full blur-[120px] -z-10" />
        <div className="absolute -bottom-10 left-20 w-[450px] h-[450px] bg-purple-500/5 rounded-full blur-[120px] -z-10" />
        
        {/* Subtle Noise Texture */}
        <div className="absolute inset-0 opacity-30 bg-[url('/noise.svg')] pointer-events-none -z-10" />
        
        {/* Technical floating package vector */}
        <div className="absolute right-[8%] top-[10%] opacity-[0.05] dark:opacity-[0.10] pointer-events-none -z-10 animate-bounce hidden md:block" style={{ animationDuration: '9s' }}>
           <svg width="45" height="45" viewBox="0 0 40 40" fill="none" className="text-purple-500">
             <path d="M20 5 L35 13 L35 29 L20 37 L5 29 L5 13 Z" stroke="currentColor" strokeWidth="1.5" strokeLinejoin="round" />
             <path d="M20 5 L20 37 M5 13 L20 21 L35 13" stroke="currentColor" strokeWidth="1" strokeLinejoin="round" />
           </svg>
        </div>

        {/* Blueprint specs info tag */}
        <div className="absolute left-[4%] bottom-[15%] opacity-[0.06] dark:opacity-[0.12] pointer-events-none -z-10 font-mono text-[9px] select-none text-foreground hidden xl:block">
           <div className="border border-dashed border-border p-2.5 rounded-lg flex flex-col gap-0.5">
             <div>CONTAINER: 20ft_GP / 40ft_HC</div>
             <div>STACK_LIMIT: Stacked_Gravity</div>
           </div>
        </div>

        <div className="container px-4 md:px-6 mx-auto">
          <div className="text-center mb-16">
            <h2 className="text-3xl font-bold tracking-tight sm:text-4xl">How it works</h2>
            <p className="mt-4 text-muted-foreground text-lg">From chaotic spreadsheet to optimized 3D plan in minutes.</p>
          </div>

          <div className="max-w-4xl mx-auto">
            <div className="relative space-y-12 pl-8 sm:pl-0 sm:before:absolute sm:before:inset-0 sm:before:ml-auto sm:before:mr-auto sm:before:h-full sm:before:w-0.5 sm:before:bg-gradient-to-b sm:before:from-transparent sm:before:via-border sm:before:to-transparent">
              
              {/* Step 1 */}
              <div className="relative flex flex-col sm:flex-row items-center gap-8 sm:gap-12">
                <div className="sm:w-1/2 sm:text-right order-2 sm:order-1">
                  <h3 className="text-xl font-bold">1. Define Your Load</h3>
                  <p className="text-muted-foreground mt-2">Select your container type (20ft, 40ft, etc.) and input your cargo items from our catalog or manually.</p>
                </div>
                <div className="absolute left-0 sm:left-1/2 sm:-ml-4 flex h-8 w-8 items-center justify-center rounded-full border-4 border-background bg-primary text-primary-foreground shadow-lg z-10">
                   <Box className="h-4 w-4" />
                </div>
                <div className="sm:w-1/2 order-2 w-full">
                   <div className="p-4 bg-card rounded-2xl border border-border/50 shadow-sm space-y-2 select-none">
                      <div className="flex justify-between items-center text-xs font-semibold text-foreground">
                        <span>Pallet Cargo #A</span>
                        <span className="text-[10px] text-muted-foreground">Qty: 24</span>
                      </div>
                      <div className="flex justify-between text-[10px] text-muted-foreground font-sans">
                        <span>1,200 × 1,000 × 1,600 mm</span>
                        <span>Max Stacking: 3</span>
                      </div>
                      <div className="h-1.5 w-full bg-muted rounded-full">
                        <div className="h-full bg-primary rounded-full" style={{ width: "80%" }} />
                      </div>
                   </div>
                </div>
              </div>

               {/* Step 2 */}
              <div className="relative flex flex-col sm:flex-row items-center gap-8 sm:gap-12">
                <div className="sm:w-1/2 order-2 sm:order-1 sm:text-right w-full">
                   <div className="p-4 bg-card rounded-2xl border border-border/50 shadow-sm space-y-2 select-none text-left">
                      <div className="flex items-center gap-2 text-xs font-semibold text-foreground">
                        <div className="h-2 w-2 bg-purple-500 rounded-full animate-ping" />
                        <span>Heuristic Pack Optimizer</span>
                      </div>
                      <div className="text-[10px] font-mono text-muted-foreground space-y-1">
                        <div>&gt; Evaluating 10,240 packing permutations...</div>
                        <div>&gt; Applied rules: Gravity + Load Balance</div>
                        <div className="text-emerald-500 font-medium">&gt; Solution found: 94.2% yield in 1.2s</div>
                      </div>
                   </div>
                </div>
                <div className="absolute left-0 sm:left-1/2 sm:-ml-4 flex h-8 w-8 items-center justify-center rounded-full border-4 border-background bg-purple-500 text-white shadow-lg z-10">
                   <Settings2 className="h-4 w-4" />
                </div>
                <div className="sm:w-1/2 order-2">
                   <h3 className="text-xl font-bold">2. Auto-Pack</h3>
                   <p className="text-muted-foreground mt-2">Our advanced packing algorithm optimizes for volume, weight distribution, and stacking rules instantly.</p>
                </div>
              </div>

              {/* Step 3 */}
              <div className="relative flex flex-col sm:flex-row items-center gap-8 sm:gap-12">
                <div className="sm:w-1/2 sm:text-right order-2 sm:order-1">
                  <h3 className="text-xl font-bold">3. Visualize & Adjust</h3>
                  <p className="text-muted-foreground mt-2">Interact with the 3D plan. Rotate, zoom, and verify item placement. Make manual adjustments if needed.</p>
                </div>
                <div className="absolute left-0 sm:left-1/2 sm:-ml-4 flex h-8 w-8 items-center justify-center rounded-full border-4 border-background bg-primary text-white shadow-lg z-10">
                   <Cuboid className="h-4 w-4" />
                </div>
                 <div className="sm:w-1/2 order-2 w-full">
                   <div className="p-4 bg-card rounded-2xl border border-border/50 shadow-sm space-y-2 select-none">
                      <div className="flex justify-between items-center text-xs font-semibold text-foreground">
                        <span>Stack Safety Checklist</span>
                      </div>
                      <div className="space-y-1.5 text-[10px] text-left">
                        <div className="flex items-center gap-1.5 text-emerald-600 dark:text-emerald-400">
                          <Check className="h-3 w-3 stroke-[3]" />
                          <span>Heavier items placed first (Gravity lock)</span>
                        </div>
                        <div className="flex items-center gap-1.5 text-emerald-600 dark:text-emerald-400">
                          <Check className="h-3 w-3 stroke-[3]" />
                          <span>Stack weight limits verified</span>
                        </div>
                        <div className="flex items-center gap-1.5 text-emerald-600 dark:text-emerald-400">
                          <Check className="h-3 w-3 stroke-[3]" />
                          <span>Stacking orientation restrictions checked</span>
                        </div>
                      </div>
                   </div>
                </div>
              </div>

               {/* Step 4 */}
              <div className="relative flex flex-col sm:flex-row items-center gap-8 sm:gap-12">
                 <div className="sm:w-1/2 order-2 sm:order-1 sm:text-right w-full">
                   <div className="p-4 bg-card rounded-2xl border border-border/50 shadow-sm space-y-2.5 select-none text-left">
                      <div className="flex justify-between items-center text-xs font-semibold text-foreground">
                        <span>Dispatch & Share</span>
                      </div>
                      <div className="flex gap-2">
                        <div className="flex-1 p-1.5 border border-dashed border-border rounded-lg text-center text-[10px] text-muted-foreground bg-muted/20">
                          Download PDF Guide
                        </div>
                        <div className="flex-1 p-1.5 border border-dashed border-border rounded-lg text-center text-[10px] text-muted-foreground bg-muted/20 font-sans">
                          Web share link
                        </div>
                      </div>
                   </div>
                </div>
                <div className="absolute left-0 sm:left-1/2 sm:-ml-4 flex h-8 w-8 items-center justify-center rounded-full border-4 border-background bg-emerald-500 text-white shadow-lg z-10">
                   <Truck className="h-4 w-4" />
                </div>
                <div className="sm:w-1/2 order-2">
                   <h3 className="text-xl font-bold">4. Ship Confidently</h3>
                   <p className="text-muted-foreground mt-2">Generate manifests and share loading guides with your warehouse team to ensure perfect execution.</p>
                </div>
              </div>

            </div>
          </div>
        </div>
      </section>

      {/* CTA Section */}
      <section className="py-24 bg-zinc-950 text-white relative overflow-hidden">
        <div className="absolute inset-0 bg-[linear-gradient(to_right,#80808012_1px,transparent_1px),linear-gradient(to_bottom,#80808012_1px,transparent_1px)] bg-[size:24px_24px] opacity-20" />
        <div className="absolute top-0 right-0 w-96 h-96 bg-primary/20 rounded-full blur-[128px] pointer-events-none" />
        <div className="absolute bottom-0 left-0 w-64 h-64 bg-purple-500/20 rounded-full blur-[128px] pointer-events-none" />
        
        <div className="container px-4 md:px-6 mx-auto relative z-10 text-center">
           <h2 className="text-3xl font-bold tracking-tight sm:text-5xl mb-6">Ready to maximize your loads?</h2>
           <p className="text-zinc-400 text-xl max-w-2xl mx-auto mb-10">
             Join forward-thinking logistics teams using LoadIQ to save time and shipping costs today.
           </p>
           <Button asChild size="lg" className="h-14 px-8 text-lg bg-white text-zinc-950 hover:bg-zinc-200 shadow-xl hover:scale-105 transition-transform">
             <Link href="/login">Get Started for Free</Link>
           </Button>
           <p className="mt-6 text-sm text-zinc-500">No credit card required • Cancel anytime</p>
        </div>
      </section>

      {/* Footer */}
      <footer className="border-t border-border bg-card py-16">
        <div className="container px-4 md:px-6 mx-auto">
          <div className="grid grid-cols-1 md:grid-cols-4 gap-12">
             <div className="space-y-4">
               <div className="flex items-center gap-2">
                 <div className="relative h-12 w-12 overflow-hidden rounded-md border border-border/10 bg-white shadow-sm">
                    <Image src="/logo.png" alt="LoadIQ" fill className="object-contain p-0.5" />
                  </div>
                 <span className="text-xl font-bold">LoadIQ</span>
               </div>
               <p className="text-sm text-muted-foreground">
                 The intelligent container load planning system for modern logistics operations.
               </p>
             </div>

             <div>
               <h4 className="font-semibold mb-4">Product</h4>
               <ul className="space-y-2 text-sm text-muted-foreground">
                 <li><Link href="#" className="hover:text-foreground">Features</Link></li>
                 <li><Link href="#" className="hover:text-foreground">Pricing</Link></li>
                 <li><Link href="#" className="hover:text-foreground">API</Link></li>
                 <li><Link href="#" className="hover:text-foreground">Changelog</Link></li>
               </ul>
             </div>

             <div>
               <h4 className="font-semibold mb-4">Resources</h4>
               <ul className="space-y-2 text-sm text-muted-foreground">
                 <li><Link href="#" className="hover:text-foreground">Documentation</Link></li>
                 <li><Link href="#" className="hover:text-foreground">Guides</Link></li>
                 <li><Link href="#" className="hover:text-foreground">Support</Link></li>
                 <li><Link href="#" className="hover:text-foreground">Status</Link></li>
               </ul>
             </div>

             <div>
               <h4 className="font-semibold mb-4">Company</h4>
               <ul className="space-y-2 text-sm text-muted-foreground">
                 <li><Link href="#" className="hover:text-foreground">About</Link></li>
                 <li><Link href="#" className="hover:text-foreground">Blog</Link></li>
                 <li><Link href="#" className="hover:text-foreground">Careers</Link></li>
                 <li><Link href="#" className="hover:text-foreground">Legal</Link></li>
               </ul>
             </div>
          </div>

          <div className="mt-16 pt-8 border-t border-border flex flex-col md:flex-row items-center justify-between gap-4 text-xs text-muted-foreground">
             <p>&copy; {new Date().getFullYear()} LoadIQ. All rights reserved.</p>
             <div className="flex gap-6">
                <Link href="#" className="hover:text-foreground">Privacy Policy</Link>
                <Link href="#" className="hover:text-foreground">Terms of Service</Link>
             </div>
          </div>
        </div>
      </footer>
    </div>
  )
}
