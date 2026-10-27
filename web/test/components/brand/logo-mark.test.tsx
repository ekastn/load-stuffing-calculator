import { renderToStaticMarkup } from "react-dom/server"
import { describe, expect, it } from "vitest"

import { LogoMark } from "@/components/brand/logo-mark"

describe("LogoMark", () => {
  it("is hidden from assistive technology without a title", () => {
    const markup = renderToStaticMarkup(<LogoMark />)
    expect(markup).toContain('aria-hidden="true"')
    expect(markup).not.toContain("<title>")
  })

  it("is announced as an image when given a title", () => {
    const markup = renderToStaticMarkup(<LogoMark title="LoadIQ" />)
    expect(markup).toContain('role="img"')
    expect(markup).toContain("<title>LoadIQ</title>")
    expect(markup).not.toContain("aria-hidden")
  })

  it("passes className through", () => {
    expect(renderToStaticMarkup(<LogoMark className="size-8" />)).toContain('class="size-8"')
  })
})
