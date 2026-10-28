import type { ElementType } from "react"

import { splitWords } from "./split-words"

type SplitHeadingProps = {
  as?: ElementType
  text: string
  className?: string
  id?: string
  /** Reveal word by word on scroll (see HeadingReveals). */
  reveal?: boolean
}

/**
 * Renders a heading whose words can be revealed one by one with GSAP
 * (target `[data-word]`). Screen readers get the unsplit text; the split
 * copy stays visible when JavaScript never runs.
 */
export function SplitHeading({ as: Tag = "h2", text, className, id, reveal = false }: SplitHeadingProps) {
  return (
    <Tag id={id} className={className} data-reveal={reveal || undefined}>
      <span className="sr-only">{text}</span>
      <span aria-hidden="true">
        {splitWords(text).map((word, index) => (
          <span key={index} className="inline-block overflow-hidden pb-[0.08em] align-bottom">
            <span data-word className="inline-block will-change-transform">
              {word}
            </span>
            {" "}
          </span>
        ))}
      </span>
    </Tag>
  )
}
