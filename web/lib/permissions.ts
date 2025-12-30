export type Permission = string

export function permissionMatches(granted: Permission, required: Permission) {
  if (granted === "*") return true
  if (granted === required) return true
  if (granted.endsWith(":*")) {
    const prefix = granted.slice(0, -1)
    return required.startsWith(prefix)
  }
  return false
}

export function hasAnyPermission(granted: Permission[], required: Permission[]) {
  if (required.length === 0) return true
  for (const requiredPerm of required) {
    for (const grantedPerm of granted) {
      if (permissionMatches(grantedPerm, requiredPerm)) return true
    }
  }
  return false
}
