export type ApiResponse<T> = {
  message: string
  data: T
}

export type PaginationMeta = {
  page: number
  limit: number
  total: number
  total_pages: number
}

export type PaginatedResponse<T> = ApiResponse<T[]> & {
  pagination: PaginationMeta
}

export type ApiErrorResponse = {
  message: string
  errors?: Record<string, string>
}
