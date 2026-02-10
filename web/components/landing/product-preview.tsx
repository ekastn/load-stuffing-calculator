"use client"

import type React from "react"
import Image from "next/image" // Import needed
import { Box, Layers, Move3d, ZoomIn } from "lucide-react"

export function ProductPreview() {
  return (
    <div className="relative aspect-square w-full max-w-[500px] mx-auto group perspective-[1200px]">
      {/* Glassmorphic Container */}
      <div className="relative h-full w-full overflow-hidden rounded-3xl border border-white/20 bg-gradient-to-br from-white/10 to-white/5 p-8 shadow-2xl backdrop-blur-xl transition-all duration-500 hover:scale-[1.02] hover:shadow-primary/20">
        
        {/* Mock Window Controls */}
        <div className="flex items-center justify-between border-b border-white/10 pb-4 mb-4">
          <div className="flex gap-2">
            <div className="h-3 w-3 rounded-full bg-red-400/80" />
            <div className="h-3 w-3 rounded-full bg-yellow-400/80" />
            <div className="h-3 w-3 rounded-full bg-green-400/80" />
          </div>
          <div className="flex items-center gap-3 rounded-full bg-white/5 pr-4 pl-1.5 py-1.5 text-xs font-medium text-white/90">
             <div className="bg-white rounded-full p-1 flex items-center justify-center h-6 w-6 shadow-sm">
               <Image src="/logo.png" alt="LoadIQ" width={20} height={20} className="object-contain" />
             </div>
          </div>
        </div>

        {/* 3D Scene Mock */}
        <div className="relative flex h-full w-full items-center justify-center [transform-style:preserve-3d]">
          
          {/* Main 3D Container (Rotated Plane) */}
          {/* We rotate the entire scene to give isometric view */}
          <div className="relative w-64 h-64 transition-transform duration-700 ease-out [transform-style:preserve-3d] [transform:rotateX(60deg)_rotateZ(45deg)] group-hover:[transform:rotateX(55deg)_rotateZ(40deg)]">
             
             {/* Floor Grid */}
             <div className="absolute inset-0 grid grid-cols-4 grid-rows-4 border-2 border-white/20 bg-white/5 shadow-2xl [transform:translateZ(0px)]">
                {[...Array(16)].map((_, i) => (
                  <div key={i} className="border border-white/5" />
                ))}
             </div>

             {/* 3D Boxes - Layer 1 */}
             {/* Box 1: Bottom Left */}
             <div className="absolute bottom-0 left-0 w-1/4 h-1/4 p-1">
                <Cube color="bg-primary" height="h-12" />
             </div>
             {/* Box 2: Bottom Right */}
             <div className="absolute bottom-0 right-0 w-1/2 h-1/4 p-1">
                <Cube color="bg-orange-400" height="h-8" />
             </div>
             
             {/* Box 3: Center */}
             <div className="absolute top-1/4 left-1/4 w-1/2 h-1/2 p-1">
                <Cube color="bg-blue-500" height="h-4" />
             </div>

             {/* Box 4: Floating/Stacked */}
             <div className="absolute top-1/4 left-1/4 w-1/4 h-1/4 p-1 [transform:translateZ(40px)]">
                 <Cube color="bg-emerald-400" height="h-8" />
             </div>

          </div>

          {/* Overlay UI (Flat) */}
          <div className="absolute left-6 bottom-20 space-y-3 pointer-events-none">
             <StatCard label="Volume" value="78%" color="text-primary" />
             <StatCard label="Weight" value="64%" color="text-orange-400" />
          </div>

          <div className="absolute right-6 top-20 flex flex-col gap-2">
            <ToolButton icon={Move3d} />
            <ToolButton icon={ZoomIn} />
            <ToolButton icon={Layers} />
          </div>

        </div>
      </div>

      {/* Background Glows */}
      <div className="absolute -left-12 -top-12 h-64 w-64 rounded-full bg-primary/20 blur-[80px] -z-10" />
      <div className="absolute -right-12 -bottom-12 h-64 w-64 rounded-full bg-purple-500/20 blur-[80px] -z-10" />
    </div>
  )
}

function Cube({ color, height }: { color: string; height: string }) {
  // A pseudo-3D cube using CSS transforms
  // height prop controls the visual 'tallness' (Z-axis extraction)
  return (
    <div className={`relative w-full h-full group/cube`}>
       {/* Top Face */}
       <div className={`absolute inset-0 ${color} opacity-90 border border-white/20 transition-transform duration-300 group-hover/cube:-translate-y-2 shadow-lg [transform:translateZ(20px)]`} />
       
       {/* Side Faces (Simulated with borders/shadows for now to keep DOM light, or real faces if needed) */}
       <div className="absolute inset-x-0 -bottom-2 h-2 bg-black/20 blur-sm [transform:translateZ(0px)]" />
    </div>
  )
}

function StatCard({ label, value, color }: { label: string; value: string; color: string }) {
  return (
    <div className="flex items-center gap-3 rounded-lg border border-white/10 bg-black/50 px-4 py-2 backdrop-blur-md shadow-lg">
       <div className={`flex h-8 w-8 items-center justify-center rounded-full border-2 border-current ${color} text-[10px] font-bold`}>
         {value}
       </div>
       <span className="text-xs font-medium text-white/90">{label}</span>
    </div>
  )
}

function ToolButton({ icon: Icon }: { icon: React.ElementType }) {
  return (
    <div className="flex h-10 w-10 items-center justify-center rounded-xl border border-white/10 bg-black/30 text-white/50 backdrop-blur-md transition-all hover:bg-white/20 hover:text-white hover:scale-110 cursor-pointer shadow-sm">
      <Icon className="h-5 w-5" />
    </div>
  )
}
