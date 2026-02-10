"use client"

import { useEffect, useState } from "react"
import { useParams, useRouter } from "next/navigation"
import QRCode from "qrcode"
import { jsPDF } from "jspdf"
import { Button } from "@/components/ui/button"
import { Card, CardContent, CardHeader, CardTitle } from "@/components/ui/card"
import { Badge } from "@/components/ui/badge"
import { BarcodeInfo } from "@/lib/types/plan"
import { PlanService } from "@/lib/services/plans"
import { ArrowLeft, Printer, Box, Loader2, QrCode } from "lucide-react"

export default function PlanBarcodesPage() {
  const params = useParams()
  const router = useRouter()
  const id = params.id as string
  const [barcodes, setBarcodes] = useState<BarcodeInfo[]>([])
  const [qrImages, setQrImages] = useState<Record<string, string>>({})
  const [loading, setLoading] = useState(true)
  const [error, setError] = useState<string | null>(null)

  useEffect(() => {
    loadBarcodes()
  }, [id])

  const loadBarcodes = async () => {
    try {
      setLoading(true)
      const data = await PlanService.getPlanBarcodes(id)
      setBarcodes(data)
      
      // Generate high-contrast QR code images for scan reliability
      const images: Record<string, string> = {}
      for (const barcode of data) {
        images[barcode.barcode] = await QRCode.toDataURL(barcode.barcode, {
          width: 400, // Higher resolution for print
          margin: 0,
          color: {
            dark: "#000000",
            light: "#ffffff",
          }
        })
      }
      setQrImages(images)
    } catch (err: any) {
      console.error("Failed to load barcodes:", err)
      setError(err.message || "Failed to load barcodes")
    } finally {
      setLoading(false)
    }
  }

  const handlePrint = () => {
    const doc = new jsPDF({
      orientation: "portrait",
      unit: "mm",
      format: "a4",
      putOnlyUsedFonts: true,
      floatPrecision: 16
    })

    const labelWidth = 70
    const labelHeight = 37.125
    const cols = 3
    const rows = 8
    const labelsPerPage = cols * rows

    barcodes.forEach((barcode, index) => {
      // Manage pages
      if (index > 0 && index % labelsPerPage === 0) {
        doc.addPage()
      }

      const pageIndex = index % labelsPerPage
      const col = pageIndex % cols
      const row = Math.floor(pageIndex / cols)
      const x = col * labelWidth
      const y = row * labelHeight

      // Label separator lines (light grey)
      doc.setDrawColor(230, 230, 230)
      doc.setLineWidth(0.1)
      doc.rect(x, y, labelWidth, labelHeight)

      const padding = 4
      const innerX = x + padding
      const innerY = y + padding

      // 1. Step Number Area
      doc.setTextColor(160, 160, 160)
      doc.setFont("helvetica", "bold")
      doc.setFontSize(7)
      doc.text("STEP", innerX + 4, innerY + 2, { align: "center" })

      doc.setTextColor(0, 0, 0)
      doc.setFontSize(28)
      doc.text(barcode.step_number.toString(), innerX + 4, innerY + 12, { align: "center" })

      // Vertical separator
      doc.setDrawColor(220, 220, 220)
      doc.line(innerX + 10, innerY, innerX + 10, innerY + labelHeight - (padding * 2))

      // 2. QR Code (Middle)
      const qrSize = 28
      const qrX = innerX + 12
      const qrY = innerY - 1
      if (qrImages[barcode.barcode]) {
        // We use the high-res QR generated in loadBarcodes
        doc.addImage(qrImages[barcode.barcode], "PNG", qrX, qrY, qrSize, qrSize, undefined, 'FAST')
      }

      // 3. Info Column (Right)
      const infoX = innerX + 42
      const infoY = innerY + 1
      
      // Item Label
      doc.setFont("helvetica", "bold")
      doc.setFontSize(8)
      doc.setTextColor(0, 0, 0)
      const splitLabel = doc.splitTextToSize(barcode.item_label, 22)
      doc.text(splitLabel, infoX, infoY)

      // Dimensions
      const dimY = innerY + 15
      doc.setFont("helvetica", "bold")
      doc.setFontSize(6)
      doc.setTextColor(160, 160, 160)
      doc.text("DIMENSIONS (MM)", infoX, dimY)

      doc.setFont("courier", "bold")
      doc.setFontSize(9)
      doc.setTextColor(60, 60, 60)
      doc.text(`${barcode.dimensions.length}x${barcode.dimensions.width}x${barcode.dimensions.height}`, infoX, dimY + 4)

      // 4. Footer Metadata
      const footerY = y + labelHeight - padding
      doc.setDrawColor(230, 230, 230)
      doc.setLineDashPattern([1, 1], 0)
      doc.line(infoX, footerY - 3, x + labelWidth - padding, footerY - 3)
      doc.setLineDashPattern([], 0) // reset

      doc.setFont("helvetica", "normal")
      doc.setFontSize(5)
      doc.setTextColor(180, 180, 180)
      doc.text("LOAD-STUFFING", infoX, footerY - 0.5)
      
      const shortId = barcode.item_id.split("-").pop() || ""
      doc.setFont("courier", "normal")
      doc.text(shortId, infoX, footerY + 1.5)

      doc.setFont("helvetica", "bold")
      doc.setFontSize(7)
      doc.setTextColor(200, 200, 200)
      doc.text(`#${barcode.step_number}`, x + labelWidth - padding, footerY + 1, { align: "right" })
    })

    doc.save(`loading-labels-${id}.pdf`)
  }

  if (loading) {
    return (
      <div className="flex flex-col items-center justify-center min-h-[400px] gap-4">
        <Loader2 className="w-8 h-8 animate-spin text-primary" />
        <p className="text-muted-foreground font-medium">Loading barcodes...</p>
      </div>
    )
  }

  if (error) {
    return (
      <div className="flex flex-col items-center justify-center min-h-[400px] gap-4 text-center">
        <Box className="w-12 h-12 text-destructive opacity-50" />
        <h2 className="text-xl font-semibold">Error Loading Barcodes</h2>
        <p className="text-muted-foreground max-w-md">{error}</p>
        <Button onClick={loadBarcodes} variant="outline" className="mt-2">
          Retry
        </Button>
      </div>
    )
  }

  if (barcodes.length === 0) {
    return (
      <div className="flex flex-col items-center justify-center min-h-[400px] gap-4 text-center">
        <QrCode className="w-12 h-12 text-muted-foreground opacity-50" />
        <h2 className="text-xl font-semibold">No Barcodes Available</h2>
        <p className="text-muted-foreground max-w-md">
          This shipment doesn't have any calculated placements yet. 
          Please run the optimization before printing loading labels.
        </p>
        <Button onClick={() => router.back()} variant="outline" className="mt-2">
          Back to Shipment
        </Button>
      </div>
    )
  }

  return (
    <div className="container mx-auto py-6 space-y-6">
      <div className="flex flex-col md:flex-row md:items-center justify-between gap-4">
        <div className="flex items-center gap-2">
          <Button variant="ghost" size="icon" onClick={() => router.back()}>
            <ArrowLeft className="h-4 w-4" />
          </Button>
          <div>
            <h1 className="text-2xl font-bold tracking-tight">Loading Labels</h1>
            <p className="text-sm text-muted-foreground">
              QR codes for step-by-step loading validation
            </p>
          </div>
        </div>
        <Button onClick={handlePrint} className="gap-2">
          <Printer className="h-4 w-4" />
          Print Labels (PDF)
        </Button>
      </div>

      <div className="grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-3 xl:grid-cols-4 gap-4">
        {barcodes.map((barcode) => (
          <Card key={barcode.barcode} className="overflow-hidden">
            <CardHeader className="pb-2 bg-muted/50">
              <div className="flex items-center justify-between">
                <CardTitle className="text-lg">Step {barcode.step_number}</CardTitle>
                <Badge variant="secondary" className="font-mono text-[10px]">
                  {barcode.item_id.split("-").pop()}
                </Badge>
              </div>
            </CardHeader>
            <CardContent className="pt-4 space-y-4">
              <div className="flex justify-center bg-white p-2 rounded-lg border">
                {qrImages[barcode.barcode] ? (
                  <img
                    src={qrImages[barcode.barcode]}
                    alt={`Step ${barcode.step_number}`}
                    className="w-32 h-32"
                    style={{ imageRendering: 'pixelated' }}
                  />
                ) : (
                  <div className="w-32 h-32 bg-muted animate-pulse rounded" />
                )}
              </div>
              <div className="space-y-1">
                <p className="text-sm font-semibold truncate" title={barcode.item_label}>
                  {barcode.item_label}
                </p>
                <p className="text-xs text-muted-foreground font-mono">
                  {barcode.dimensions.length}×{barcode.dimensions.width}×{barcode.dimensions.height} mm
                </p>
              </div>
            </CardContent>
          </Card>
        ))}
      </div>
      
      <p className="text-xs text-center text-muted-foreground pt-4">
        Generated for shipment ID: {id}
      </p>
    </div>
  )
}
