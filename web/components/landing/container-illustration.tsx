import type React from "react"

type Props = {
  className?: string
  title?: string
}

export function ContainerIllustration({
  className,
  title = "Shipment loading preview illustration",
}: Props) {
  return (
    <svg
      viewBox="0 0 200 120"
      role="img"
      aria-label={title}
      className={className}
      fill="none"
      xmlns="http://www.w3.org/2000/svg"
    >
      <title>{title}</title>

      {/* Base floor */}
      <path
        d="M10 90 L140 120 L190 100 L60 70 Z"
        className="fill-slate-300 dark:fill-slate-700"
      />

      {/* Stacked boxes - back row */}
      <path d="M60 70 L110 82 L110 52 L60 40 Z" className="fill-primary" />
      <path d="M60 40 L110 52 L150 42 L100 30 Z" className="fill-primary/70" />
      <path d="M110 82 L150 72 L150 42 L110 52 Z" className="fill-primary/90" />

      {/* Stacked boxes - middle row */}
      <path d="M35 65 L85 77 L85 47 L35 35 Z" className="fill-chart-3" />
      <path d="M35 35 L85 47 L125 37 L75 25 Z" className="fill-chart-3/70" />
      <path d="M85 77 L125 67 L125 37 L85 47 Z" className="fill-chart-3/90" />

      {/* Stacked boxes - front row */}
      <path d="M10 90 L60 102 L60 72 L10 60 Z" className="fill-accent" />
      <path d="M10 60 L60 72 L100 62 L50 50 Z" className="fill-accent/70" />
      <path d="M60 102 L100 92 L100 62 L60 72 Z" className="fill-accent/90" />

      {/* Container frame overlay */}
      <path d="M10 30 L10 90" strokeWidth="2" className="stroke-slate-600 dark:stroke-slate-400" />
      <path d="M10 90 L140 120" strokeWidth="2" className="stroke-slate-600 dark:stroke-slate-400" />
      <path d="M140 120 L140 60" strokeWidth="2" className="stroke-slate-600 dark:stroke-slate-400" />

      {/* Subtle top edge */}
      <path
        d="M10 30 L60 10 L190 40 L140 60 Z"
        strokeWidth="2"
        className="stroke-slate-500/70 dark:stroke-slate-300/70"
      />
    </svg>
  )
}
