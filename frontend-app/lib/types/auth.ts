export interface LoginPayload {
  email: string
  password: string
}

export interface RegisterPayload {
  name: string
  email: string
  password: string
}

export interface AuthResponse {
  message: string
  data: {
    access_token: string
    expires_in?: number
  }
}
