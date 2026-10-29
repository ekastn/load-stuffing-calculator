"use client"

import React, { useState, useEffect } from "react"
import { NumericInput } from "@/components/numeric-input"
import { Select, SelectContent, SelectItem, SelectTrigger, SelectValue } from "@/components/ui/select"

export type WeightUnit = "g" | "kg" | "t"

interface WeightInputGroupProps {
  weight_kg: number
  onChange: (val: { weight_kg: number }) => void
  disabled?: boolean
  required?: boolean
  className?: string
  defaultUnit?: WeightUnit
  label?: string
  maxWeight_kg?: number
}

const UNIT_TO_KG = {
  g: 0.001,
  kg: 1,
  t: 1000,
}

export function WeightInputGroup({
  weight_kg,
  onChange,
  disabled = false,
  required = false,
  className = "",
  defaultUnit = "kg",
  label = "Weight",
  maxWeight_kg,
}: WeightInputGroupProps) {
  const [unit, setUnit] = useState<WeightUnit>(defaultUnit)
  const [displayValue, setDisplayValue] = useState<number>(weight_kg / UNIT_TO_KG[unit])

  useEffect(() => {
    setDisplayValue(weight_kg / UNIT_TO_KG[unit])
  }, [weight_kg, unit])

  const handleUnitChange = (newUnit: WeightUnit) => {
    setUnit(newUnit)
  }

  const handleValueChange = (val: number | undefined) => {
    const v = val || 0
    setDisplayValue(v)
    const vInKg = v * UNIT_TO_KG[unit]
    onChange({ weight_kg: vInKg })
  }

  const maxDisplay = maxWeight_kg ? maxWeight_kg / UNIT_TO_KG[unit] : undefined

  return (
    <div className={className}>
      <div className="space-y-2 flex-1">
        <label className="text-sm font-medium">{label}</label>
        <NumericInput
          value={displayValue || ""}
          onChange={handleValueChange}
          disabled={disabled}
          required={required}
          max={maxDisplay}
        />
      </div>
      <div className="space-y-2">
        <label className="text-sm font-medium invisible">Unit</label>
        <Select value={unit} onValueChange={(val) => handleUnitChange(val as WeightUnit)} disabled={disabled}>
          <SelectTrigger className="bg-input/50 w-[80px]">
            <SelectValue />
          </SelectTrigger>
          <SelectContent>
            <SelectItem value="g">g</SelectItem>
            <SelectItem value="kg">kg</SelectItem>
            <SelectItem value="t">t</SelectItem>
          </SelectContent>
        </Select>
      </div>
    </div>
  )
}
