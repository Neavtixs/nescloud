import { atom } from "jotai"

export type UserData = {
  id: string
  name: string
  email: string
  created_at: string
}

export const userAtom = atom<UserData | null>(null)
export const isAuthLoadingAtom = atom(true)
