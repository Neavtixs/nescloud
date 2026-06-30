import { atom } from "jotai";

export type FolderItem = {
  id: string;
  name: string;
  parent_folder_id: string;
  created_at: string;
  updated_at: string;
};

export const foldersAtom = atom<FolderItem[]>([]);
