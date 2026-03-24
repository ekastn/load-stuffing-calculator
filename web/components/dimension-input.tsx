"use client"

import React, { useState, useEffect } from "react"
import { NumericInput } from "@/components/numeric-input"
import { Select, SelectContent, SelectItem, SelectTrigger, SelectValue } from "@/components/ui/select"
import { MAX_DIM_MM } from "@/lib/constants"

export type DimensionUnit = "mm" | "cm" | "m"

interface DimensionInputGroupProps {
  length_mm: number
  width_mm: number
  height_mm: number
  onChange: (dims: { length_mm: number; width_mm: number; height_mm: number }) => void
  disabled?: boolean
  required?: boolean
  className?: string
  defaultUnit?: DimensionUnit
  maxLength_mm?: number
  maxWidth_mm?: number
  maxHeight_mm?: number
}

const UNIT_MULTIPLIERS = {
  mm: 1,
  cm: 10,
  m: 1000,
}

export function DimensionInputGroup({
  length_mm,
  width_mm,
  height_mm,
  onChange,
  disabled = false,
  required = false,
  className = "grid gap-4 md:grid-cols-4",
  defaultUnit = "mm",
  maxLength_mm,
  maxWidth_mm,
  maxHeight_mm,
}: DimensionInputGroupProps) {
  const [unit, setUnit] = useState<DimensionUnit>(defaultUnit)

  // Local state for the displayed values based on the current unit
  const [displayLength, setDisplayLength] = useState<number>(length_mm / UNIT_MULTIPLIERS[unit])
  const [displayWidth, setDisplayWidth] = useState<number>(width_mm / UNIT_MULTIPLIERS[unit])
  const [displayHeight, setDisplayHeight] = useState<number>(height_mm / UNIT_MULTIPLIERS[unit])

  // Sync props to local state if props change externally (e.g. form reset, parent override)
  useEffect(() => {
    const divider = UNIT_MULTIPLIERS[unit]
    setDisplayLength(length_mm / divider)
    setDisplayWidth(width_mm / divider)
    setDisplayHeight(height_mm / divider)
  }, [length_mm, width_mm, height_mm, unit])

  const handleUnitChange = (newUnit: DimensionUnit) => {
     setUnit(newUnit)
  }

  const handleLocalChange = (field: "length" | "width" | "height", val: number) => {
     const multiplier = UNIT_MULTIPLIERS[unit]
     const valInMm = val * multiplier

     // Update local state for immediate feedback
     if (field === "length") setDisplayLength(val)
     if (field === "width") setDisplayWidth(val)
     if (field === "height") setDisplayHeight(val)

     // Send the millimeter value strictly to parent
     onChange({
       length_mm: field === "length" ? valInMm : length_mm,
       width_mm: field === "width" ? valInMm : width_mm,
       height_mm: field === "height" ? valInMm : height_mm
     })
  }

  const maxLenDisplay = (maxLength_mm ?? MAX_DIM_MM) / UNIT_MULTIPLIERS[unit]
  const maxWidDisplay = (maxWidth_mm ?? MAX_DIM_MM) / UNIT_MULTIPLIERS[unit]
  const maxHeiDisplay = (maxHeight_mm ?? MAX_DIM_MM) / UNIT_MULTIPLIERS[unit]

  return (
    <div className={className}>
      <div className="space-y-2">
        <label className="text-sm font-medium">Length</label>
        <NumericInput
           value={displayLength || ""}
           onChange={(val) => handleLocalChange("length", val || 0)}
           disabled={disabled}
           required={required}
           max={maxLenDisplay}
        />
      </div>
      <div className="space-y-2">
         <label className="text-sm font-medium">Width</label>
         <NumericInput
            value={displayWidth || ""}
            onChange={(val) => handleLocalChange("width", val || 0)}
            disabled={disabled}
            required={required}
            max={maxWidDisplay}
         />
      </div>
      <div className="space-y-2">
         <label className="text-sm font-medium">Height</label>
         <NumericInput
            value={displayHeight || ""}
            onChange={(val) => handleLocalChange("height", val || 0)}
            disabled={disabled}
            required={required}
            max={maxHeiDisplay}
         />
      </div>
      <div className="space-y-2">
         <label className="text-sm font-medium">Unit</label>
         <Select value={unit} onValueChange={(val) => handleUnitChange(val as DimensionUnit)} disabled={disabled}>
            <SelectTrigger className="bg-input/50">
               <SelectValue placeholder="Unit" />
            </SelectTrigger>
            <SelectContent>
               <SelectItem value="mm">mm</SelectItem>
               <SelectItem value="cm">cm</SelectItem>
               <SelectItem value="m">m</SelectItem>
            </SelectContent>
         </Select>
      </div>
    </div>
  )
}
