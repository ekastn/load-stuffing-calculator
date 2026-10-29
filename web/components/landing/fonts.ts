import { Archivo, Big_Shoulders_Stencil, IBM_Plex_Mono } from "next/font/google"

/** Stencil display face, after the lettering painted on container doors. */
export const stencilFont = Big_Shoulders_Stencil({
  subsets: ["latin"],
  variable: "--font-landing-stencil",
  axes: ["opsz"],
  display: "swap",
})

export const archivoFont = Archivo({
  subsets: ["latin"],
  variable: "--font-landing-sans",
  axes: ["wdth"],
  display: "swap",
})

/** Reserved for measured values: dimensions, weights, step counts. */
export const plexMonoFont = IBM_Plex_Mono({
  subsets: ["latin"],
  variable: "--font-landing-mono",
  weight: ["400", "500"],
  display: "swap",
})

export const landingFontVariables = [stencilFont.variable, archivoFont.variable, plexMonoFont.variable].join(" ")
