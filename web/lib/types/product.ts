export interface CreateProductRequest {
  name: string
  length_mm: number
  width_mm: number
  height_mm: number
  weight_kg: number
  color_hex?: string
}

export interface UpdateProductRequest {
  name: string
  length_mm: number
  width_mm: number
  height_mm: number
  weight_kg: number
  color_hex?: string
}

export interface ProductResponse {
  id: string
  name: string
  length_mm: number
  width_mm: number
  height_mm: number
  weight_kg: number
  color_hex?: string
}
