import { LandingLink } from "@/components/landing/landing-link"

type AuthFooterLinksProps = {
  href: string
  label: string
}

/** The link to the other auth screen, plus the way back to the landing page. */
export function AuthFooterLinks({ href, label }: AuthFooterLinksProps) {
  return (
    <div className="mt-6 flex flex-wrap items-center justify-between gap-x-6 gap-y-2 border-t border-ink/15 pt-4">
      <LandingLink href={href} variant="plain" className="font-semibold">
        {label}
      </LandingLink>
      <LandingLink href="/" variant="plain" className="text-steel">
        Go to home
      </LandingLink>
    </div>
  )
}
