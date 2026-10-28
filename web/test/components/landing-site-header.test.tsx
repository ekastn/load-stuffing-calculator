import { renderToStaticMarkup } from "react-dom/server"
import { describe, expect, it } from "vitest"

import { SiteHeader } from "@/components/landing/site-header"

describe("SiteHeader", () => {
  it("shows section links and both actions on the landing page", () => {
    const markup = renderToStaticMarkup(<SiteHeader />)

    expect(markup).toContain("How it works")
    expect(markup).toContain('href="/login"')
    expect(markup).toContain("Start a load plan")
    expect(markup).toContain('aria-label="Open menu"')
  })
})
