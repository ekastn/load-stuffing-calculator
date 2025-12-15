export interface MockUser {
  id: string
  email: string
  name: string
  role: "admin" | "planner" | "operator"
  password: string
}

export const MOCK_USERS: MockUser[] = [
  {
    id: "user_admin_001",
    email: "admin@example.com",
    name: "Admin User",
    role: "admin",
    password: "admin123",
  },
  {
    id: "user_planner_001",
    email: "planner@example.com",
    name: "Planner User",
    role: "planner",
    password: "planner123",
  },
  {
    id: "user_operator_001",
    email: "operator@example.com",
    name: "Operator User",
    role: "operator",
    password: "operator123",
  },
]
