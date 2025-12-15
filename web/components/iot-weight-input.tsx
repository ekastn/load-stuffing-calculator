"use client"

import type React from "react"

import { useState } from "react"
import { Button } from "@/components/ui/button"
import { Input } from "@/components/ui/input"
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from "@/components/ui/card"
import { useExecution } from "@/lib/execution-context"
import { AlertCircle, CheckCircle, Scale } from "lucide-react"

interface IoTWeightInputProps {
  sessionId: string
  itemId: string
  itemName: string
  expectedWeight: number
  tolerance: number
  onValidated: (passed: boolean) => void
}

export function IoTWeightInput({
  sessionId,
  itemId,
  itemName,
  expectedWeight,
  tolerance,
  onValidated,
}: IoTWeightInputProps) {
  const { recordValidation } = useExecution()
  const [weight, setWeight] = useState("")
  const [result, setResult] = useState<{ passed: boolean; message: string } | null>(null)
  const [isSimulating, setIsSimulating] = useState(false)

  const handleSimulateIoT = () => {
    setIsSimulating(true)
    // Simulate IoT sensor reading with slight variance
    const variance = (Math.random() - 0.5) * tolerance * 1.5
    const simulatedWeight = expectedWeight + variance
    setWeight(simulatedWeight.toFixed(2))
    setIsSimulating(false)
  }

  const handleValidate = (e: React.FormEvent) => {
    e.preventDefault()
    const actualWeight = Number.parseFloat(weight)

    if (isNaN(actualWeight)) {
      setResult({ passed: false, message: "Invalid weight value" })
      return
    }

    const passed = recordValidation(sessionId, itemId, itemName, expectedWeight, actualWeight, tolerance)

    const delta = Math.abs(actualWeight - expectedWeight)

    if (passed) {
      setResult({
        passed: true,
        message: `✓ Valid! Delta: ${delta.toFixed(2)} kg (within ±${tolerance} kg)`,
      })
      setWeight("")
      setTimeout(() => {
        onValidated(true)
        setResult(null)
      }, 1500)
    } else {
      setResult({
        passed: false,
        message: `✗ Invalid! Delta: ${delta.toFixed(2)} kg (exceeds ±${tolerance} kg)`,
      })
    }
  }

  return (
    <Card className="border-border/50 bg-card/50">
      <CardHeader>
        <CardTitle className="flex items-center gap-2">
          <Scale className="h-5 w-5" />
          Weight Validation
        </CardTitle>
        <CardDescription>Scanning: {itemName}</CardDescription>
      </CardHeader>
      <CardContent className="space-y-4">
        <div className="rounded-lg border border-border/50 bg-background/50 p-4 space-y-2">
          <p className="text-sm text-muted-foreground">Expected Weight</p>
          <p className="text-2xl font-bold text-foreground">{expectedWeight.toFixed(2)} kg</p>
          <p className="text-xs text-muted-foreground">Tolerance: ±{tolerance} kg</p>
        </div>

        <form onSubmit={handleValidate} className="space-y-3">
          <div className="space-y-2">
            <label className="text-sm font-medium">Actual Weight (kg)</label>
            <div className="flex gap-2">
              <Input
                type="number"
                step="0.01"
                value={weight}
                onChange={(e) => setWeight(e.target.value)}
                placeholder="0.00"
                disabled={isSimulating}
                required
              />
              <Button type="button" variant="outline" onClick={handleSimulateIoT} disabled={isSimulating}>
                {isSimulating ? "Reading..." : "Scan IoT"}
              </Button>
            </div>
          </div>

          <Button type="submit" className="w-full" disabled={!weight}>
            Validate & Next
          </Button>
        </form>

        {result && (
          <div
            className={`flex items-center gap-3 rounded-lg p-3 ${
              result.passed
                ? "bg-green-500/10 border border-green-500/20"
                : "bg-destructive/10 border border-destructive/20"
            }`}
          >
            {result.passed ? (
              <CheckCircle className="h-5 w-5 text-green-500 flex-shrink-0" />
            ) : (
              <AlertCircle className="h-5 w-5 text-destructive flex-shrink-0" />
            )}
            <span
              className={result.passed ? "text-sm font-medium text-green-500" : "text-sm font-medium text-destructive"}
            >
              {result.message}
            </span>
          </div>
        )}
      </CardContent>
    </Card>
  )
}
