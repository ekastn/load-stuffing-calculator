type LogoMarkProps = {
  className?: string
  /** Accessible name; omit when the mark sits next to the "LoadIQ" text. */
  title?: string
}

/**
 * LoadIQ's mark: a container end view with the first box seated in the
 * bottom-left corner, where every packing plan starts. The frame follows
 * the current text colour; the box is always signal orange.
 */
export function LogoMark({ className, title }: LogoMarkProps) {
  const labelling = title ? { role: "img" } : { "aria-hidden": true }

  return (
    <svg xmlns="http://www.w3.org/2000/svg" viewBox="0 0 32 32" className={className} {...labelling}>
      {title && <title>{title}</title>}
      <rect x="3.5" y="5.5" width="25" height="21" fill="none" stroke="currentColor" strokeWidth="3" />
      <rect x="1" y="3" width="6" height="6" fill="currentColor" />
      <rect x="25" y="3" width="6" height="6" fill="currentColor" />
      <rect x="1" y="23" width="6" height="6" fill="currentColor" />
      <rect x="25" y="23" width="6" height="6" fill="currentColor" />
      <rect x="5" y="14" width="11" height="11" fill="#ff5b14" />
    </svg>
  )
}
