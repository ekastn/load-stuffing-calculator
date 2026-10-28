export function splitWords(text: string): string[] {
  const trimmed = text.trim()
  return trimmed ? trimmed.split(/\s+/) : []
}
