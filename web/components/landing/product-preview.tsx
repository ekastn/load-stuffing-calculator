"use client"

import React from "react"
import Image from "next/image"

export function ProductPreview() {
  return (
    <div className="relative w-full max-w-[800px] mx-auto group">
      {/* Minimalist Image Container */}
      <div className="relative rounded-2xl border border-border/80 bg-muted/20 p-1 shadow-2xl transition-all duration-500 hover:scale-[1.01] hover:shadow-primary/5">
        <div className="relative aspect-[3/2] w-full overflow-hidden rounded-xl border border-border bg-card">
          <Image 
            src="/antarmuka.png" 
            alt="LoadIQ Container Stuffing Planner Dashboard" 
            fill 
            priority
            className="object-cover transition-transform duration-700 ease-out group-hover:scale-[1.015]" 
          />
        </div>
      </div>

      {/* Behind background glows */}
      <div className="absolute -left-12 -top-12 h-64 w-64 rounded-full bg-primary/10 blur-[80px] -z-10 animate-pulse" />
      <div className="absolute -right-12 -bottom-12 h-64 w-64 rounded-full bg-purple-500/10 blur-[80px] -z-10" />
    </div>
  )
}
