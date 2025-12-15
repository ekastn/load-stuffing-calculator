import type React from "react"
import type { Metadata } from "next"
import { Geist, Geist_Mono } from "next/font/google"
import { Analytics } from "@vercel/analytics/next"
import { AuthProvider } from "@/lib/auth-context"
import { StorageProvider } from "@/lib/storage-context"
import { PlanningProvider } from "@/lib/planning-context"
import { ExecutionProvider } from "@/lib/execution-context"
import { AuditProvider } from "@/lib/audit-context"
import "./globals.css"

const _geist = Geist({ subsets: ["latin"] })
const _geistMono = Geist_Mono({ subsets: ["latin"] })

export const metadata: Metadata = {
  title: "Load & Stuffing Calculator | Container Optimization System",
  description: "IoT-integrated platform for optimizing container loads with real-time weight validation",
  generator: "v0.app",
  icons: {
    icon: [
      {
        url: "/icon-light-32x32.png",
        media: "(prefers-color-scheme: light)",
      },
      {
        url: "/icon-dark-32x32.png",
        media: "(prefers-color-scheme: dark)",
      },
      {
        url: "/icon.svg",
        type: "image/svg+xml",
      },
    ],
    apple: "/apple-icon.png",
  },
}

export default function RootLayout({
  children,
}: Readonly<{
  children: React.ReactNode
}>) {
  return (
    <html lang="en" suppressHydrationWarning>
      <body className={`font-sans antialiased`}>
        <AuthProvider>
          <AuditProvider>
            <StorageProvider>
              <PlanningProvider>
                <ExecutionProvider>
                  {children}
                  <Analytics />
                </ExecutionProvider>
              </PlanningProvider>
            </StorageProvider>
          </AuditProvider>
        </AuthProvider>
      </body>
    </html>
  )
}
