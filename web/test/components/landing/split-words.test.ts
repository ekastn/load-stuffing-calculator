import { describe, expect, it } from "vitest"

import { splitWords } from "@/components/landing/motion/split-words"

describe("splitWords", () => {
  const cases: Array<{ name: string; input: string; expected: string[] }> = [
    { name: "single word", input: "Load", expected: ["Load"] },
    { name: "keeps punctuation attached", input: "Plan it, then load.", expected: ["Plan", "it,", "then", "load."] },
    { name: "collapses repeated whitespace", input: "one   two\n three", expected: ["one", "two", "three"] },
    { name: "trims leading and trailing whitespace", input: "  edge case  ", expected: ["edge", "case"] },
    { name: "empty string", input: "", expected: [] },
    { name: "whitespace only", input: "   ", expected: [] },
  ]

  for (const { name, input, expected } of cases) {
    it(name, () => {
      expect(splitWords(input)).toEqual(expected)
    })
  }
})
