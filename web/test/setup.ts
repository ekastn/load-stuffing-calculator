import { afterEach } from "vitest"

// jsdom does not implement matchMedia, which the landing and auth screens use
// for reduced-motion and breakpoint decisions. Default every query to "no match".
if (typeof window !== "undefined" && typeof window.matchMedia !== "function") {
  window.matchMedia = (query: string) =>
    ({
      matches: false,
      media: query,
      onchange: null,
      addEventListener: () => {},
      removeEventListener: () => {},
      addListener: () => {},
      removeListener: () => {},
      dispatchEvent: () => false,
    }) as unknown as MediaQueryList
}

// Vitest has no global afterEach, so React Testing Library cannot clean up by itself.
if (typeof document !== "undefined") {
  const { cleanup } = await import("@testing-library/react")
  afterEach(cleanup)
}
