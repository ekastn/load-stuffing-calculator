import Link from "next/link"
import Image from "next/image"
import { ArrowRight, Box, Check, CheckCircle2, Cuboid, LayoutDashboard, LineChart, Package, PlayCircle, Settings2, Smartphone, Truck } from "lucide-react"

import { ProductPreview } from "@/components/landing/product-preview"
import { RedirectIfAuthed } from "@/components/redirect-if-authed"
import { TrialLoadCalculator } from "@/components/trial-load-calculator"
import { Button } from "@/components/ui/button"
import { Badge } from "@/components/ui/badge"

export default function LandingPage() {
  return (
    <div className="flex min-h-screen flex-col bg-background selection:bg-primary/20">
      <RedirectIfAuthed />

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

              <div className="flex flex-col sm:flex-row gap-4 w-full sm:w-auto">
                <Button asChild size="lg" className="h-14 px-8 text-lg shadow-lg shadow-primary/20 hover:shadow-primary/30 transition-all hover:scale-105">
                  <Link href="/login">
                    Start Optimization <ArrowRight className="ml-2 h-5 w-5" />
                  </Link>
                </Button>
                <Button asChild variant="outline" size="lg" className="h-14 px-8 text-lg backdrop-blur-sm bg-background/50 hover:bg-background/80">
                  <Link href="#how-it-works">How it works</Link>
                </Button>
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
      <section className="py-32 bg-background border-y border-border/40 relative overflow-hidden">
         <div className="absolute inset-0 bg-[url('/grid.svg')] opacity-[0.25] pointer-events-none" />
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
      <section className="py-24 bg-background relative">
        {/* Dot Pattern Background */}
        <div className="absolute inset-0 bg-[image:radial-gradient(#e5e7eb_1px,transparent_1px)] [background-size:20px_20px] opacity-100 pointer-events-none" />
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
               <div className="mt-8 relative h-48 w-full rounded-xl border border-border/30 bg-background/50 overflow-hidden shadow-sm">
                  {/* Mock UI elements */}
                  <div className="absolute top-4 left-4 right-4 h-2 bg-muted rounded-full" />
                  <div className="absolute top-10 left-4 w-1/3 h-2 bg-muted rounded-full" />
                  <div className="absolute bottom-0 right-0 w-2/3 h-32 bg-primary/5 rounded-tl-xl border-t border-l border-border/30" />
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
               <div className="flex justify-center">
                  <div className="w-24 h-40 border-4 border-foreground/10 rounded-2xl bg-background" />
               </div>
            </div>

            {/* Feature 3 - Small */}
            <div className="md:col-span-1 rounded-3xl border border-border/50 bg-card p-8 hover:border-primary/50 transition-colors group relative overflow-hidden">
               <div className="absolute inset-0 bg-gradient-to-br from-primary/5 to-transparent opacity-0 group-hover:opacity-100 transition-opacity" />
               <div className="h-12 w-12 rounded-xl bg-primary/10 flex items-center justify-center text-primary mb-6">
                  <Cuboid className="h-6 w-6" />
               </div>
               <h3 className="text-xl font-bold mb-2">3D Visualization</h3>
               <p className="text-muted-foreground">Rotate, zoom, and inspect every layer of your cargo plan before loading.</p>
            </div>

            {/* Feature 4 - Wide */}
            <div className="md:col-span-2 rounded-3xl border border-border/50 bg-card p-8 flex flex-col md:flex-row items-center gap-8 hover:border-primary/50 transition-colors group relative overflow-hidden">
               <div className="flex-1 space-y-4 relative z-10">
                 <div className="h-12 w-12 rounded-xl bg-primary/10 flex items-center justify-center text-primary mb-4">
                    <LineChart className="h-6 w-6" />
                 </div>
                 <h3 className="text-2xl font-bold">Analytics & Reports</h3>
                 <p className="text-muted-foreground">Generate comprehensive manifests, weight reports, and utilization stats to identify cost-saving opportunities.</p>
               </div>
               <div className="flex-1 w-full">
                  <div className="grid grid-cols-2 gap-3">
                     <div className="bg-background/80 rounded-xl p-4 border border-border/50">
                        <div className="text-2xl font-bold text-foreground">98%</div>
                        <div className="text-xs text-muted-foreground">Volume Util.</div>
                     </div>
                     <div className="bg-background/80 rounded-xl p-4 border border-border/50">
                        <div className="text-2xl font-bold text-primary">-12%</div>
                        <div className="text-xs text-muted-foreground">Shipping Cost</div>
                     </div>
                  </div>
               </div>
            </div>

          </div>
        </div>
      </section>

      {/* Workflow Section - Timeline */}
      <section id="how-it-works" className="py-24 bg-background border-y border-border/40 relative overflow-hidden">
        {/* Subtle Noise Texture */}
        <div className="absolute inset-0 opacity-30 bg-[url('/noise.svg')] pointer-events-none" />
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
                <div className="sm:w-1/2 order-2">
                   {/* Placeholder for visual or extra text */}
                   <div className="p-4 bg-background rounded-xl border border-border/50 shadow-sm text-sm text-muted-foreground">
                      Input dimensions, weight limits, and constraints.
                   </div>
                </div>
              </div>

               {/* Step 2 */}
              <div className="relative flex flex-col sm:flex-row items-center gap-8 sm:gap-12">
                <div className="sm:w-1/2 order-2 sm:order-1 sm:text-right">
                   <div className="p-4 bg-background rounded-xl border border-border/50 shadow-sm text-sm text-muted-foreground">
                      Our engines test thousands of combinations.
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
                <div className="absolute left-0 sm:left-1/2 sm:-ml-4 flex h-8 w-8 items-center justify-center rounded-full border-4 border-background bg-orange-500 text-white shadow-lg z-10">
                   <Cuboid className="h-4 w-4" />
                </div>
                 <div className="sm:w-1/2 order-2">
                   <div className="p-4 bg-background rounded-xl border border-border/50 shadow-sm text-sm text-muted-foreground">
                      See step-by-step loading instructions.
                   </div>
                </div>
              </div>

               {/* Step 4 */}
              <div className="relative flex flex-col sm:flex-row items-center gap-8 sm:gap-12">
                 <div className="sm:w-1/2 order-2 sm:order-1 sm:text-right">
                   <div className="p-4 bg-background rounded-xl border border-border/50 shadow-sm text-sm text-muted-foreground">
                      Export to PDF, Excel, or share link.
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
                 <div className="relative h-8 w-8 overflow-hidden rounded-md border border-border/10 bg-white shadow-sm">
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
