import { describe, expect, it } from "vitest"

import { isUuidV4 } from "@/lib/utils"

describe("isUuidV4", () => {
  it("returns true for a valid uuid v4", () => {
    expect(isUuidV4("550e8400-e29b-41d4-a716-446655440000")).toBe(true)
  })

  it("returns false for non-v4 uuid", () => {
    expect(isUuidV4("550e8400-e29b-11d4-a716-446655440000")).toBe(false)
  })

  it("returns false for invalid strings", () => {
    expect(isUuidV4("not-a-uuid")).toBe(false)
    expect(isUuidV4("")).toBe(false)
  })

  it("trims whitespace", () => {
    expect(isUuidV4("  550e8400-e29b-41d4-a716-446655440000  ")).toBe(true)
  })
})
