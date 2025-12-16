export interface CreateContainerRequest {
  name: string
  inner_length_mm: number
  inner_width_mm: number
  inner_height_mm: number
  max_weight_kg: number
  description?: string
}

export interface UpdateContainerRequest {
  name: string
  inner_length_mm: number
  inner_width_mm: number
  inner_height_mm: number
  max_weight_kg: number
  description?: string
}

export interface ContainerResponse {
  id: string
  name: string
  inner_length_mm: number
  inner_width_mm: number
  inner_height_mm: number
  max_weight_kg: number
  description?: string
}
