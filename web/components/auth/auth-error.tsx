"use client"

/** Form-level failure notice: what happened, in the danger tone, never a decoration. */
export function AuthError({ message }: { message: string }) {
  return (
    <div
      role="alert"
      className="flex items-start gap-2 border border-danger/50 bg-form p-3 font-landing text-sm text-danger"
    >
      <span aria-hidden="true" className="mt-px icon-[solar--danger-triangle-linear] size-5 shrink-0" />
      <span>{message}</span>
    </div>
  )
}
