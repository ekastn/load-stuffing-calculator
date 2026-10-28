"use client"

import { cn } from "@/lib/utils"

type AuthFieldProps = {
  id: string
  label: string
  value: string
  onChange: (value: string) => void
  type?: "text" | "email" | "password"
  autoComplete?: string
  help?: string
  error?: string
  disabled?: boolean
  required?: boolean
}

/** Labelled input for the auth screens; errors replace the border colour and add a message below. */
export function AuthField({
  id,
  label,
  value,
  onChange,
  type = "text",
  autoComplete,
  help,
  error,
  disabled = false,
  required = true,
}: AuthFieldProps) {
  const helpId = `${id}-help`
  const errorId = `${id}-error`
  const describedBy = [help ? helpId : null, error ? errorId : null].filter(Boolean).join(" ")

  return (
    <div className="flex flex-col gap-1.5">
      <label htmlFor={id} className="font-landing text-sm font-semibold text-ink">
        {label}
      </label>
      <input
        id={id}
        name={id}
        type={type}
        value={value}
        onChange={(event) => onChange(event.target.value)}
        autoComplete={autoComplete}
        required={required}
        disabled={disabled}
        aria-invalid={error ? true : undefined}
        aria-describedby={describedBy || undefined}
        className={cn(
          "h-10 w-full rounded-[4px] border bg-form px-3 font-landing text-base text-ink placeholder:text-steel/70",
          "focus:outline-2 focus:outline-offset-2 focus:outline-signal",
          "disabled:cursor-not-allowed disabled:bg-yard disabled:text-steel",
          error ? "border-danger" : "border-field",
        )}
      />
      {help && (
        <p id={helpId} className="text-sm text-steel">
          {help}
        </p>
      )}
      {error && (
        <p id={errorId} className="text-sm text-danger">
          {error}
        </p>
      )}
    </div>
  )
}
