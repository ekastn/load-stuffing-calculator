import * as React from "react"
import { Input } from "@/components/ui/input"
import { cn } from "@/lib/utils"

export interface NumericInputProps extends Omit<React.InputHTMLAttributes<HTMLInputElement>, 'value' | 'onChange'> {
  value?: number | string;
  onChange?: (value: number | undefined) => void;
  allowDecimals?: boolean;
}

const NumericInput = React.forwardRef<HTMLInputElement, NumericInputProps>(
  ({ className, value, onChange, allowDecimals = true, ...props }, ref) => {
    // Keep local string state for what the user is currently typing
    const [displayValue, setDisplayValue] = React.useState("")

    React.useEffect(() => {
      // Format incoming numeric value to string with thousand separators
      if (value === undefined || value === null || value === "") {
        setDisplayValue("")
        return
      }
      
      const numValue = typeof value === 'string' ? parseFloat(value) : value
      if (isNaN(numValue)) {
        setDisplayValue("")
        return
      }

      // If user is currently typing and the string ends with decimal or similar, don't reformat to avoid jumping cursor
      // but simple enough: just reformat when value prop changes from outside
      setDisplayValue(formatNumber(numValue.toString(), allowDecimals))
    }, [value, allowDecimals])

    const formatNumber = (val: string, allowDec: boolean) => {
      // Remove any non-digit character except decimal point (and minus sign if you want negatives)
      let cleaned = val.replace(/[^\d.-]/g, "")
      
      const parts = cleaned.split(".")
      let integerPart = parts[0]
      let decimalPart = parts[1]

      // apply thousand separator (e.g. dot) to integer part
      // Note: user requested "titik setiap 3 angka di belakangnya". So 1.000.000
      integerPart = integerPart.replace(/\B(?=(\d{3})+(?!\d))/g, ".")

      if (allowDec && parts.length > 1) {
        return `${integerPart},${decimalPart}`
      }
      return integerPart
    }

    const handleChange = (e: React.ChangeEvent<HTMLInputElement>) => {
      let rawVal = e.target.value
      
      // Let user type negative numbers, decimal separators
      // Here we parse knowing that '.' is thousand separator and ',' is decimal (Indonesian locale)
      // Since standard input assumes '.' is decimal, we need to convert for parsing
      let cleaned = rawVal.replace(/\./g, "").replace(/,/g, ".")
      
      if (cleaned === "" || cleaned === "-") {
        setDisplayValue(rawVal)
        onChange?.(undefined)
        return
      }

      // Only allow valid numbers
      if (!/^-?\d*\.?\d*$/.test(cleaned)) {
        return // ignore invalid typing
      }

      const parsed = parseFloat(cleaned)
      if (isNaN(parsed)) {
        setDisplayValue(rawVal)
        onChange?.(undefined)
        return
      }

      // Clamp to max if provided
      const maxLimit = typeof props.max === 'number' ? props.max : (typeof props.max === 'string' ? parseFloat(props.max) : undefined)
      const clamped = maxLimit !== undefined && !isNaN(maxLimit) ? Math.min(parsed, maxLimit) : parsed

      if (clamped !== parsed) {
        // Value was clamped — reformat the display
        setDisplayValue(formatNumber(clamped.toString(), allowDecimals))
        onChange?.(clamped)
        return
      }

      setDisplayValue(rawVal) // keep their raw typing to avoid cursor jumps
      onChange?.(parsed)
    }

    const handleBlur = () => {
      // On blur, reformat nicely
      if (value !== undefined && !isNaN(Number(value))) {
        const maxLimit = typeof props.max === 'number' ? props.max : (typeof props.max === 'string' ? parseFloat(props.max) : undefined)
        const numVal = Number(value)
        const clamped = maxLimit !== undefined && !isNaN(maxLimit) ? Math.min(numVal, maxLimit) : numVal
        setDisplayValue(formatNumber(String(clamped), allowDecimals))
      }
    }

    return (
      <Input
        type="text"
        className={cn("font-mono", className)}
        value={displayValue}
        onChange={handleChange}
        onBlur={handleBlur}
        ref={ref}
        inputMode="decimal"
        {...props}
      />
    )
  }
)
NumericInput.displayName = "NumericInput"

export { NumericInput }
