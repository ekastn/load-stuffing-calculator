
import { LogoMark } from "@/components/brand/logo-mark"
import { LandingLink } from "./landing-link"
import { MEDIA_CREDITS } from "./media-credits"
import { SECTION_LINKS } from "./section-links"

const linkClass = "text-steel hover:text-ink"

export function SiteFooter() {
  const credits = Object.values(MEDIA_CREDITS)

  return (
    <footer className="border-t border-ink/15 py-14">
      <div className="mx-auto grid max-w-[90rem] gap-10 px-4 md:grid-cols-12 md:px-8">
        <div className="md:col-span-5">
          <div className="flex items-center gap-2.5">
            <LogoMark className="size-8" />
            <span className="font-stencil text-2xl font-extrabold tracking-wide text-ink">LoadIQ</span>
          </div>
          <p className="mt-4 max-w-[26rem] leading-relaxed text-steel">Container load planning for logistics teams.</p>
        </div>

        <nav aria-label="Footer" className="grid grid-cols-2 gap-8 md:col-span-7 md:grid-cols-3">
          <ul className="space-y-3">
            {SECTION_LINKS.map((link) => (
              <li key={link.href}>
                <LandingLink href={link.href} variant="plain" className={linkClass}>
                  {link.label}
                </LandingLink>
              </li>
            ))}
          </ul>
          <ul className="space-y-3">
            <li>
              <LandingLink href="/login" variant="plain" className={linkClass}>
                Sign in
              </LandingLink>
            </li>
            <li>
              <LandingLink href="/register" variant="plain" className={linkClass}>
                Create an account
              </LandingLink>
            </li>
          </ul>
        </nav>

        <div className="flex flex-col gap-2 border-t border-ink/15 pt-6 text-sm text-steel md:col-span-12 md:flex-row md:justify-between">
          <p>&copy; {new Date().getFullYear()} LoadIQ</p>
          <p>
            Photos by{" "}
            {credits.map((credit, index) => (
              <span key={credit.sourceUrl}>
                {index > 0 && " and "}
                <a href={credit.sourceUrl} className="underline underline-offset-2 hover:text-ink" rel="noopener" target="_blank">
                  {credit.photographer}
                </a>
              </span>
            ))}{" "}
            on Unsplash.
          </p>
        </div>
      </div>
    </footer>
  )
}
