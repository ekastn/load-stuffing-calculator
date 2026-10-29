import { RedirectIfAuthed } from "@/components/redirect-if-authed"
import { LoginForm } from "@/components/login-form"

export default function LoginPage() {
  return (
    <>
      <RedirectIfAuthed />
      <LoginForm />
    </>
  )
}
