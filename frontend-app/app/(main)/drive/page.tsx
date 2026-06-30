"use client";

import { useEffect, useMemo, useCallback, useState } from "react";
import { useAtom } from "jotai";
import { api } from "@/lib/api/api-call";
import type { ApiResponse, PaginatedResponse } from "@/lib/api/api-response";
import { foldersAtom, type FolderItem } from "@/lib/atoms/folder-atoms";
import * as Dialog from "@radix-ui/react-dialog";
import * as DropdownMenu from "@radix-ui/react-dropdown-menu";
import * as AlertDialog from "@radix-ui/react-alert-dialog";
import {
  Folder,
  FolderOpen,
  File,
  Upload,
  FolderPlus,
  Search,
  Download,
  Trash2,
  Pencil,
  Link,
  MoreHorizontal,
  ImageIcon,
  FileText,
  FileSpreadsheet,
  ChevronRight,
  Home,
  LayoutGrid,
  List,
  X,
  Loader2,
} from "lucide-react";

type ItemType = "folder" | "file";

type DriveItem = {
  id: string;
  type: ItemType;
  name: string;
  mime_type?: string;
  size_bytes?: number;
  modified_at: string;
};

const mockFiles: DriveItem[] = [
  {
    id: "f1",
    type: "file",
    name: "report-q2.pdf",
    mime_type: "application/pdf",
    size_bytes: 2516582,
    modified_at: "2026-06-29T08:30:00Z",
  },
  {
    id: "f2",
    type: "file",
    name: "holiday-photo.png",
    mime_type: "image/png",
    size_bytes: 1153433,
    modified_at: "2026-06-28T16:45:00Z",
  },
  {
    id: "f3",
    type: "file",
    name: "budget-2026.xlsx",
    mime_type:
      "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet",
    size_bytes: 876544,
    modified_at: "2026-06-27T10:15:00Z",
  },
  {
    id: "f4",
    type: "file",
    name: "meeting-notes.txt",
    mime_type: "text/plain",
    size_bytes: 12288,
    modified_at: "2026-06-26T14:00:00Z",
  },
  {
    id: "f5",
    type: "file",
    name: "presentation.pptx",
    mime_type:
      "application/vnd.openxmlformats-officedocument.presentationml.presentation",
    size_bytes: 5242880,
    modified_at: "2026-06-25T09:20:00Z",
  },
];

function formatSize(bytes: number): string {
  if (bytes >= 1073741824) return `${(bytes / 1073741824).toFixed(1)} GB`;
  if (bytes >= 1048576) return `${(bytes / 1048576).toFixed(1)} MB`;
  if (bytes >= 1024) return `${(bytes / 1024).toFixed(0)} KB`;
  return `${bytes} B`;
}

function formatDate(iso: string): string {
  return new Date(iso).toLocaleDateString("en-US", {
    year: "numeric",
    month: "short",
    day: "numeric",
  });
}

function FileIcon({ mime_type, size = 18 }: { mime_type: string; size?: number }) {
  if (mime_type.startsWith("image/"))
    return <ImageIcon size={size} className="text-purple-500 dark:text-purple-400" />;
  if (mime_type.startsWith("text/"))
    return <FileText size={size} className="text-gray-500 dark:text-gray-400" />;
  if (
    mime_type.includes("spreadsheet") ||
    mime_type.includes("excel") ||
    mime_type.includes("csv")
  )
    return <FileSpreadsheet size={size} className="text-green-600 dark:text-green-400" />;
  if (mime_type.includes("presentation") || mime_type.includes("powerpoint"))
    return <FileText size={size} className="text-orange-500 dark:text-orange-400" />;
  if (mime_type === "application/pdf")
    return <FileText size={size} className="text-red-500 dark:text-red-400" />;
  return <File size={size} className="text-gray-400 dark:text-gray-500" />;
}

const menuItem =
  "flex w-full items-center gap-2 px-3 py-2 text-sm text-gray-700 outline-none data-[highlighted]:bg-gray-50 dark:text-gray-300 dark:data-[highlighted]:bg-gray-800";
const menuItemDanger =
  "flex w-full items-center gap-2 px-3 py-2 text-sm text-red-600 outline-none data-[highlighted]:bg-red-50 dark:text-red-400 dark:data-[highlighted]:bg-red-900/20";

function ItemMenu({
  type,
  children,
  onRenameFolder,
  onOpenFolder,
  onDeleteFolder,
}: {
  type: ItemType;
  children: React.ReactNode;
  onRenameFolder?: () => void;
  onOpenFolder?: () => void;
  onDeleteFolder?: () => void;
}) {
  return (
    <DropdownMenu.Root>
      <DropdownMenu.Trigger asChild>
        {children}
      </DropdownMenu.Trigger>
      <DropdownMenu.Portal>
        <DropdownMenu.Content
          className="w-44 rounded-lg border border-gray-200 bg-white py-1 shadow-lg dark:border-gray-800 dark:bg-gray-900"
          sideOffset={8}
          align="end"
        >
          {type === "file" && (
            <DropdownMenu.Item
              onSelect={() => alert("Download coming soon")}
              className={menuItem}
            >
              <Download size={14} />
              Download
            </DropdownMenu.Item>
          )}
          <DropdownMenu.Item
            onSelect={() => {
              if (type === "folder" && onRenameFolder) {
                onRenameFolder();
              } else {
                alert("Rename coming soon");
              }
            }}
            className={menuItem}
          >
            <Pencil size={14} />
            Rename
          </DropdownMenu.Item>
          {type === "file" && (
            <DropdownMenu.Item
              onSelect={() => alert("Share coming soon")}
              className={menuItem}
            >
              <Link size={14} />
              Share
            </DropdownMenu.Item>
          )}
          {type === "folder" && (
            <>
              <DropdownMenu.Item
                onSelect={() => onOpenFolder?.()}
                className={menuItem}
              >
                <FolderOpen size={14} />
                Open
              </DropdownMenu.Item>
              <DropdownMenu.Separator className="my-1 border-t border-gray-200 dark:border-gray-800" />
              <DropdownMenu.Item
                onSelect={() => onDeleteFolder?.()}
                className={menuItemDanger}
              >
                <Trash2 size={14} />
                Delete
              </DropdownMenu.Item>
            </>
          )}
          {type === "file" && (
            <DropdownMenu.Item
              onSelect={() => alert("Delete coming soon")}
              className={menuItem}
            >
              <Trash2 size={14} />
              Delete
            </DropdownMenu.Item>
          )}
        </DropdownMenu.Content>
      </DropdownMenu.Portal>
    </DropdownMenu.Root>
  );
}

const FILE_LIMIT = 20;

export default function DrivePage() {
  const [search, setSearch] = useState("");
  const [viewMode, setViewMode] = useState<"list" | "grid">("list");
  const [folders, setFolders] = useAtom(foldersAtom);
  const [showCreateModal, setShowCreateModal] = useState(false);
  const [newFolderName, setNewFolderName] = useState("");
  const [isCreating, setIsCreating] = useState(false);
  const [createError, setCreateError] = useState("");
  const [filePage, setFilePage] = useState(1);

  const [showRenameModal, setShowRenameModal] = useState(false);
  const [renameTarget, setRenameTarget] = useState<{ id: string; name: string } | null>(null);
  const [renameValue, setRenameValue] = useState("");
  const [isRenaming, setIsRenaming] = useState(false);
  const [renameError, setRenameError] = useState("");

  const [deleteTarget, setDeleteTarget] = useState<{ id: string; name: string } | null>(null);
  const [isDeleting, setIsDeleting] = useState(false);

  const [currentFolderId, setCurrentFolderId] = useState("");
  const [folderPath, setFolderPath] = useState<{ id: string; name: string }[]>([]);
  const [isNavigating, setIsNavigating] = useState(false);

  useEffect(() => {
    setIsNavigating(true);
    const param = currentFolderId
      ? `parent_folder_id=${currentFolderId}&limit=999`
      : "parent_folder_id=&limit=999";
    api
      .get<PaginatedResponse<FolderItem>>(`/folders?${param}`)
      .then((res) => {
        setFolders(res.data);
      })
      .catch(() => {})
      .finally(() => setIsNavigating(false));
  }, [currentFolderId, setFolders]);

  const fileTotal = mockFiles.length;
  const fileTotalPages = Math.max(1, Math.ceil(fileTotal / FILE_LIMIT));

  const paginatedFiles = useMemo(() => {
    const offset = (filePage - 1) * FILE_LIMIT;
    return mockFiles.slice(offset, offset + FILE_LIMIT);
  }, [filePage]);

  const items = useMemo(() => {
    const folderItems: DriveItem[] = folders.map((f) => ({
      id: f.id,
      type: "folder" as const,
      name: f.name,
      modified_at: f.updated_at,
    }));

    const all = [...folderItems, ...paginatedFiles];

    const filtered = search
      ? all.filter((item) =>
          item.name.toLowerCase().includes(search.toLowerCase()),
        )
      : all;

    return filtered.sort((a, b) => {
      if (a.type !== b.type) return a.type === "folder" ? -1 : 1;
      return a.name.localeCompare(b.name);
    });
  }, [folders, paginatedFiles, search]);

  const handleSearch = useCallback((value: string) => {
    setSearch(value);
    setFilePage(1);
  }, []);

  const goToPage = useCallback(
    (page: number) => {
      if (page >= 1 && page <= fileTotalPages) {
        setFilePage(page);
      }
    },
    [fileTotalPages],
  );

  const pageNumbers = useMemo(() => {
    const pages: number[] = [];
    if (fileTotalPages <= 7) {
      for (let i = 1; i <= fileTotalPages; i++) pages.push(i);
    } else {
      pages.push(1);
      if (filePage > 3) pages.push(-1);
      const start = Math.max(2, filePage - 1);
      const end = Math.min(fileTotalPages - 1, filePage + 1);
      for (let i = start; i <= end; i++) pages.push(i);
      if (filePage < fileTotalPages - 2) pages.push(-1);
      pages.push(fileTotalPages);
    }
    return pages;
  }, [filePage, fileTotalPages]);

  const navigateToFolder = useCallback((folderId: string, folderName: string) => {
    setFolderPath((prev) => [...prev, { id: folderId, name: folderName }]);
    setCurrentFolderId(folderId);
  }, []);

  const handleCreateFolder = useCallback(async () => {
    setCreateError("");
    setIsCreating(true);

    try {
      const res = await api.post<ApiResponse<{ id: string }>>("/folders", {
        parent_folder_id: currentFolderId,
        name: newFolderName.trim(),
      });

      const newFolder: FolderItem = {
        id: res.data.id,
        name: newFolderName.trim(),
        parent_folder_id: currentFolderId,
        created_at: new Date().toISOString(),
        updated_at: new Date().toISOString(),
      };

      setFolders((prev) => [newFolder, ...prev]);
      setShowCreateModal(false);
      setNewFolderName("");
    } catch (err: unknown) {
      const message =
        err instanceof Error ? err.message : "Failed to create folder";
      setCreateError(message);
    } finally {
      setIsCreating(false);
    }
  }, [newFolderName, currentFolderId, setFolders]);

  const openRenameModal = useCallback((id: string, name: string) => {
    setRenameTarget({ id, name });
    setRenameValue(name);
    setShowRenameModal(true);
  }, []);

  const handleRenameFolder = useCallback(async () => {
    if (!renameTarget) return;
    setRenameError("");
    setIsRenaming(true);

    try {
      await api.patch<ApiResponse<{ id: string; name: string }>>(
        `/folders/${renameTarget.id}`,
        { name: renameValue.trim() },
      );

      setFolders((prev) =>
        prev.map((f) =>
          f.id === renameTarget.id
            ? { ...f, name: renameValue.trim(), updated_at: new Date().toISOString() }
            : f,
        ),
      );
      setShowRenameModal(false);
      setRenameTarget(null);
      setRenameValue("");
    } catch (err: unknown) {
      const message =
        err instanceof Error ? err.message : "Failed to rename folder";
      setRenameError(message);
    } finally {
      setIsRenaming(false);
    }
  }, [renameTarget, renameValue, setFolders]);

  const confirmDeleteFolder = useCallback(async () => {
    if (!deleteTarget) return;
    setIsDeleting(true);

    try {
      await api.delete<ApiResponse<null>>(`/folders/${deleteTarget.id}`);
      setFolders((prev) => prev.filter((f) => f.id !== deleteTarget.id));
      setDeleteTarget(null);
    } catch (err: unknown) {
      const message =
        err instanceof Error ? err.message : "Failed to delete folder";
      alert(message);
    } finally {
      setIsDeleting(false);
    }
  }, [deleteTarget, setFolders]);

  return (
    <>
      <div className="mx-auto max-w-5xl">
      <div className="mb-6 flex items-center gap-2 text-sm text-gray-500 dark:text-gray-400">
        {[{ id: "", label: "My Drive" }, ...folderPath.map((f) => ({ id: f.id, label: f.name }))].map(
          (crumb, i) => {
            const isCurrent = i === folderPath.length;
            return (
              <span key={crumb.id || "root"} className="flex items-center gap-2">
                {i > 0 && <ChevronRight size={14} />}
                {i === 0 && <Home size={14} />}
                {isCurrent ? (
                  <span className="font-medium text-gray-900 dark:text-gray-100">
                    {crumb.label}
                  </span>
                ) : (
                  <span
                    className="cursor-pointer hover:underline"
                    onClick={() => {
                      setFolderPath((prev) => prev.slice(0, i));
                      setCurrentFolderId(crumb.id);
                    }}
                  >
                    {crumb.label}
                  </span>
                )}
              </span>
            );
          },
        )}
      </div>

        <div className="mb-6 flex flex-col gap-3 sm:flex-row sm:items-center sm:justify-between">
          <div className="relative flex-1 max-w-sm">
            <Search
              size={16}
              className="absolute left-3 top-1/2 -translate-y-1/2 text-gray-400 dark:text-gray-500"
            />
            <input
              type="text"
              placeholder="Search files and folders..."
              value={search}
              onChange={(e) => handleSearch(e.target.value)}
              className="w-full rounded-lg border border-gray-300 bg-white py-2 pl-9 pr-3 text-sm outline-none focus:border-blue-400 focus:ring-2 focus:ring-blue-100 dark:border-gray-700 dark:bg-gray-800 dark:text-gray-100 dark:focus:border-blue-500 dark:focus:ring-blue-800"
            />
          </div>

          <div className="flex items-center gap-1 rounded-lg border border-gray-300 bg-white p-0.5 dark:border-gray-700 dark:bg-gray-800">
            <button
              onClick={() => setViewMode("list")}
              className={`rounded-md p-1.5 transition-colors ${
                viewMode === "list"
                  ? "bg-white text-gray-900 shadow-sm dark:bg-gray-700 dark:text-gray-100"
                  : "text-gray-400 hover:text-gray-600 dark:text-gray-500 dark:hover:text-gray-300"
              }`}
              title="List view"
            >
              <List size={16} />
            </button>
            <button
              onClick={() => setViewMode("grid")}
              className={`rounded-md p-1.5 transition-colors ${
                viewMode === "grid"
                  ? "bg-white text-gray-900 shadow-sm dark:bg-gray-700 dark:text-gray-100"
                  : "text-gray-400 hover:text-gray-600 dark:text-gray-500 dark:hover:text-gray-300"
              }`}
              title="Grid view"
            >
              <LayoutGrid size={16} />
            </button>
          </div>

          <div className="flex items-center gap-2">
            <Dialog.Root open={showCreateModal} onOpenChange={setShowCreateModal}>
              <Dialog.Trigger asChild>
                <button
                  className="inline-flex items-center gap-2 rounded-lg border border-gray-300 bg-white px-3 py-2 text-sm font-medium text-gray-700 hover:bg-gray-50 dark:border-gray-700 dark:bg-gray-800 dark:text-gray-300 dark:hover:bg-gray-700"
                >
                  <FolderPlus size={16} />
                  New Folder
                </button>
              </Dialog.Trigger>
              <Dialog.Portal>
                <Dialog.Overlay className="fixed inset-0 z-50 bg-black/50" />
                <Dialog.Content
                  className="fixed left-1/2 top-1/2 z-50 w-full max-w-md -translate-x-1/2 -translate-y-1/2 rounded-xl border border-gray-200 bg-white p-6 shadow-xl focus:outline-none dark:border-gray-700 dark:bg-gray-900"
                  onInteractOutside={(e) => e.preventDefault()}
                >
                  <div className="mb-4 flex items-center justify-between">
                    <Dialog.Title className="text-lg font-semibold text-gray-900 dark:text-gray-100">
                      Create Folder
                    </Dialog.Title>
                    <Dialog.Close asChild>
                      <button
                        className="rounded p-1 text-gray-400 hover:bg-gray-100 hover:text-gray-600 dark:text-gray-500 dark:hover:bg-gray-800 dark:hover:text-gray-300"
                      >
                        <X size={18} />
                      </button>
                    </Dialog.Close>
                  </div>

                  <div className="mb-4">
                    <label className="mb-1.5 block text-sm font-medium text-gray-700 dark:text-gray-300">
                      Folder Name
                    </label>
                    <input
                      type="text"
                      value={newFolderName}
                      onChange={(e) => setNewFolderName(e.target.value)}
                      placeholder="Enter folder name"
                      autoFocus
                      onKeyDown={(e) => {
                        if (e.key === "Enter" && newFolderName.trim() && !isCreating) {
                          handleCreateFolder();
                        }
                      }}
                      className="w-full rounded-lg border border-gray-300 bg-white px-3 py-2 text-sm outline-none focus:border-blue-400 focus:ring-2 focus:ring-blue-100 dark:border-gray-700 dark:bg-gray-800 dark:text-gray-100 dark:focus:border-blue-500 dark:focus:ring-blue-800"
                    />
                  </div>

                  {createError && (
                    <div className="mb-4 rounded-lg border border-red-200 bg-red-50 px-3 py-2 text-sm text-red-600 dark:border-red-800 dark:bg-red-900/20 dark:text-red-400">
                      {createError}
                    </div>
                  )}

                  <div className="flex justify-end gap-2">
                    <Dialog.Close asChild>
                      <button
                        className="rounded-lg border border-gray-300 bg-white px-4 py-2 text-sm font-medium text-gray-700 hover:bg-gray-50 dark:border-gray-700 dark:bg-gray-800 dark:text-gray-300 dark:hover:bg-gray-700"
                      >
                        Cancel
                      </button>
                    </Dialog.Close>
                    <button
                      onClick={handleCreateFolder}
                      disabled={!newFolderName.trim() || isCreating}
                      className="inline-flex items-center gap-2 rounded-lg bg-blue-600 px-4 py-2 text-sm font-medium text-white hover:bg-blue-700 disabled:cursor-not-allowed disabled:opacity-50 dark:bg-blue-700 dark:hover:bg-blue-600"
                    >
                      {isCreating ? (
                        <>
                          <Loader2 size={16} className="animate-spin" />
                          Creating...
                        </>
                      ) : (
                        <>Create</>
                      )}
                    </button>
                  </div>
                </Dialog.Content>
              </Dialog.Portal>
            </Dialog.Root>

            <button
              className="inline-flex items-center gap-2 rounded-lg bg-blue-600 px-4 py-2 text-sm font-medium text-white hover:bg-blue-700 dark:bg-blue-700 dark:hover:bg-blue-600"
              onClick={() => alert("Upload coming soon")}
            >
              <Upload size={16} />
              Upload
            </button>
          </div>
        </div>

        {isNavigating ? (
          <div className="flex items-center justify-center py-20">
            <Loader2 className="animate-spin text-gray-400" size={32} />
          </div>
        ) : items.length === 0 ? (
          <div className="flex flex-col items-center gap-3 rounded-xl border border-dashed border-gray-300 bg-white py-16 dark:border-gray-700 dark:bg-gray-900">
            <Folder size={40} className="text-gray-300 dark:text-gray-600" />
            <div className="text-center">
              <p className="text-sm font-medium text-gray-900 dark:text-gray-100">
                {search ? "No results found" : "This folder is empty"}
              </p>
              <p className="mt-1 text-xs text-gray-500 dark:text-gray-400">
                {search
                  ? "Try a different search term"
                  : "Upload files or create a new folder to get started"}
              </p>
            </div>
            {!search && (
              <button
                className="mt-2 inline-flex items-center gap-2 rounded-lg bg-blue-600 px-4 py-2 text-sm font-medium text-white hover:bg-blue-700 dark:bg-blue-700 dark:hover:bg-blue-600"
                onClick={() => setShowCreateModal(true)}
              >
                <FolderPlus size={16} />
                Create Folder
              </button>
            )}
          </div>
        ) : viewMode === "list" ? (
          <div className="rounded-xl border border-gray-200 bg-white dark:border-gray-800 dark:bg-gray-900">
            <table className="w-full">
              <thead>
                <tr className="border-b border-gray-200 text-left text-xs font-medium text-gray-500 dark:border-gray-800 dark:text-gray-400">
                  <th className="w-10 px-5 py-3">
                    <input
                      type="checkbox"
                      className="rounded border-gray-300 dark:border-gray-600 dark:bg-gray-800"
                    />
                  </th>
                  <th className="px-5 py-3 font-medium">Name</th>
                  <th className="hidden px-5 py-3 font-medium sm:table-cell">
                    Size
                  </th>
                  <th className="hidden px-5 py-3 font-medium md:table-cell">
                    Modified
                  </th>
                  <th className="px-5 py-3 font-medium"></th>
                </tr>
              </thead>
              <tbody>
                {items.map((item) => (
                  <tr
                    key={item.id}
                    className="border-b border-gray-50 hover:bg-gray-50 dark:border-gray-800 dark:hover:bg-gray-800"
                  >
                    <td className="px-5 py-3">
                      <input
                        type="checkbox"
                        className="rounded border-gray-300 dark:border-gray-600 dark:bg-gray-800"
                      />
                    </td>
                    <td
                      className={`px-5 py-3 ${item.type === "folder" ? "cursor-pointer" : ""}`}
                      onClick={() => {
                        if (item.type === "folder")
                          navigateToFolder(item.id, item.name);
                      }}
                    >
                      <div className="flex items-center gap-3">
                        {item.type === "folder" ? (
                          <Folder size={18} className="text-blue-500 dark:text-blue-400" />
                        ) : (
                          <FileIcon mime_type={item.mime_type!} />
                        )}
                        <span
                          className={`text-sm font-medium ${
                            item.type === "folder"
                              ? "text-blue-600 hover:underline dark:text-blue-400 dark:hover:text-blue-300"
                              : "text-gray-900 dark:text-gray-100"
                          }`}
                        >
                          {item.name}
                        </span>
                      </div>
                    </td>
                    <td className="hidden px-5 py-3 text-sm text-gray-500 dark:text-gray-400 sm:table-cell">
                      {item.type === "folder"
                        ? "--"
                        : formatSize(item.size_bytes!)}
                    </td>
                    <td className="hidden px-5 py-3 text-sm text-gray-500 dark:text-gray-400 md:table-cell">
                      {formatDate(item.modified_at)}
                    </td>
                    <td className="relative px-5 py-3">
                      <ItemMenu
                        type={item.type}
                        onRenameFolder={
                          item.type === "folder"
                            ? () => openRenameModal(item.id, item.name)
                            : undefined
                        }
                        onOpenFolder={
                          item.type === "folder"
                            ? () => navigateToFolder(item.id, item.name)
                            : undefined
                        }
                        onDeleteFolder={
                          item.type === "folder"
                            ? () => setDeleteTarget({ id: item.id, name: item.name })
                            : undefined
                        }
                      >
                        <button className="rounded p-1 text-gray-400 hover:bg-gray-100 hover:text-gray-600 dark:text-gray-500 dark:hover:bg-gray-800 dark:hover:text-gray-300">
                          <MoreHorizontal size={15} />
                        </button>
                      </ItemMenu>
                    </td>
                  </tr>
                ))}
              </tbody>
            </table>

            <div className="flex items-center justify-between border-t border-gray-200 px-5 py-3 dark:border-gray-800">
              <p className="text-xs text-gray-500 dark:text-gray-400">
                {fileTotal === 0
                  ? `${folders.length} folder${folders.length !== 1 ? "s" : ""}`
                  : `${(filePage - 1) * FILE_LIMIT + 1}–${Math.min(filePage * FILE_LIMIT, fileTotal)} of ${fileTotal} files`}
                {folders.length > 0 && fileTotal > 0 && " · "}
                {folders.length > 0 && `${folders.length} folder${folders.length !== 1 ? "s" : ""}`}
              </p>
              {fileTotalPages > 1 && (
                <div className="flex items-center gap-1">
                  <button
                    onClick={() => goToPage(filePage - 1)}
                    disabled={filePage <= 1}
                    className="rounded px-2 py-1 text-xs text-gray-400 hover:bg-gray-100 hover:text-gray-600 disabled:cursor-not-allowed disabled:opacity-30 dark:text-gray-500 dark:hover:bg-gray-800 dark:hover:text-gray-300"
                  >
                    ←
                  </button>
                  {pageNumbers.map((p, i) =>
                    p === -1 ? (
                      <span key={`ellipsis-${i}`} className="px-1 text-xs text-gray-400 dark:text-gray-500">
                        ...
                      </span>
                    ) : (
                      <button
                        key={p}
                        onClick={() => goToPage(p)}
                        className={`rounded px-2 py-1 text-xs ${
                          p === filePage
                            ? "bg-blue-50 font-medium text-blue-600 dark:bg-blue-900/30 dark:text-blue-400"
                            : "text-gray-500 hover:bg-gray-100 dark:text-gray-400 dark:hover:bg-gray-800"
                        }`}
                      >
                        {p}
                      </button>
                    ),
                  )}
                  <button
                    onClick={() => goToPage(filePage + 1)}
                    disabled={filePage >= fileTotalPages}
                    className="rounded px-2 py-1 text-xs text-gray-400 hover:bg-gray-100 hover:text-gray-600 disabled:cursor-not-allowed disabled:opacity-30 dark:text-gray-500 dark:hover:bg-gray-800 dark:hover:text-gray-300"
                  >
                    →
                  </button>
                </div>
              )}
            </div>
          </div>
        ) : (
          <div className="grid grid-cols-3 gap-3 sm:grid-cols-4 md:grid-cols-5 lg:grid-cols-6">
            {items.map((item) => (
              <div
                key={item.id}
                className={`group relative flex flex-col items-center justify-center rounded-xl border border-gray-200 bg-white p-3 aspect-square dark:border-gray-800 dark:bg-gray-900 ${item.type === "folder" ? "cursor-pointer" : ""}`}
                onClick={() => {
                  if (item.type === "folder")
                    navigateToFolder(item.id, item.name);
                }}
              >
                <div className="absolute right-1.5 top-1.5">
                  <ItemMenu
                    type={item.type}
                    onRenameFolder={
                      item.type === "folder"
                        ? () => openRenameModal(item.id, item.name)
                        : undefined
                    }
                    onOpenFolder={
                      item.type === "folder"
                        ? () => navigateToFolder(item.id, item.name)
                        : undefined
                    }
                    onDeleteFolder={
                      item.type === "folder"
                        ? () => setDeleteTarget({ id: item.id, name: item.name })
                        : undefined
                    }
                  >
                    <button
                      onClick={(e) => e.stopPropagation()}
                      className="rounded p-1 text-gray-400 opacity-0 transition-opacity hover:bg-gray-100 hover:text-gray-600 group-hover:opacity-100 dark:text-gray-500 dark:hover:bg-gray-800 dark:hover:text-gray-300"
                    >
                      <MoreHorizontal size={15} />
                    </button>
                  </ItemMenu>
                </div>

                <div className="flex flex-1 items-center justify-center">
                  {item.type === "folder" ? (
                    <Folder size={40} className="text-blue-500 dark:text-blue-400" />
                  ) : (
                    <FileIcon mime_type={item.mime_type!} size={36} />
                  )}
                </div>

                <p
                  className={`w-full truncate text-center text-xs font-medium ${
                    item.type === "folder"
                      ? "text-blue-600 hover:underline dark:text-blue-400"
                      : "text-gray-900 dark:text-gray-100"
                  }`}
                >
                  {item.name}
                </p>
              </div>
            ))}
          </div>
        )}
      </div>

      <Dialog.Root open={showRenameModal} onOpenChange={setShowRenameModal}>
        <Dialog.Portal>
          <Dialog.Overlay className="fixed inset-0 z-50 bg-black/50" />
          <Dialog.Content
            className="fixed left-1/2 top-1/2 z-50 w-full max-w-md -translate-x-1/2 -translate-y-1/2 rounded-xl border border-gray-200 bg-white p-6 shadow-xl focus:outline-none dark:border-gray-700 dark:bg-gray-900"
            onInteractOutside={(e) => e.preventDefault()}
          >
            <div className="mb-4 flex items-center justify-between">
              <Dialog.Title className="text-lg font-semibold text-gray-900 dark:text-gray-100">
                Rename Folder
              </Dialog.Title>
              <Dialog.Close asChild>
                <button className="rounded p-1 text-gray-400 hover:bg-gray-100 hover:text-gray-600 dark:text-gray-500 dark:hover:bg-gray-800 dark:hover:text-gray-300">
                  <X size={18} />
                </button>
              </Dialog.Close>
            </div>

            <div className="mb-4">
              <label className="mb-1.5 block text-sm font-medium text-gray-700 dark:text-gray-300">
                Folder Name
              </label>
              <input
                type="text"
                value={renameValue}
                onChange={(e) => setRenameValue(e.target.value)}
                placeholder="Enter folder name"
                autoFocus
                onKeyDown={(e) => {
                  if (e.key === "Enter" && renameValue.trim() && !isRenaming) {
                    handleRenameFolder();
                  }
                }}
                className="w-full rounded-lg border border-gray-300 bg-white px-3 py-2 text-sm outline-none focus:border-blue-400 focus:ring-2 focus:ring-blue-100 dark:border-gray-700 dark:bg-gray-800 dark:text-gray-100 dark:focus:border-blue-500 dark:focus:ring-blue-800"
              />
            </div>

            {renameError && (
              <div className="mb-4 rounded-lg border border-red-200 bg-red-50 px-3 py-2 text-sm text-red-600 dark:border-red-800 dark:bg-red-900/20 dark:text-red-400">
                {renameError}
              </div>
            )}

            <div className="flex justify-end gap-2">
              <Dialog.Close asChild>
                <button className="rounded-lg border border-gray-300 bg-white px-4 py-2 text-sm font-medium text-gray-700 hover:bg-gray-50 dark:border-gray-700 dark:bg-gray-800 dark:text-gray-300 dark:hover:bg-gray-700">
                  Cancel
                </button>
              </Dialog.Close>
              <button
                onClick={handleRenameFolder}
                disabled={!renameValue.trim() || isRenaming}
                className="inline-flex items-center gap-2 rounded-lg bg-blue-600 px-4 py-2 text-sm font-medium text-white hover:bg-blue-700 disabled:cursor-not-allowed disabled:opacity-50 dark:bg-blue-700 dark:hover:bg-blue-600"
              >
                {isRenaming ? (
                  <>
                    <Loader2 size={16} className="animate-spin" />
                    Renaming...
                  </>
                ) : (
                  <>Rename</>
                )}
              </button>
            </div>
          </Dialog.Content>
        </Dialog.Portal>
      </Dialog.Root>

      <AlertDialog.Root
        open={deleteTarget !== null}
        onOpenChange={(open) => { if (!open) setDeleteTarget(null); }}
      >
        <AlertDialog.Portal>
          <AlertDialog.Overlay className="fixed inset-0 z-50 bg-black/50" />
          <AlertDialog.Content className="fixed left-1/2 top-1/2 z-50 w-full max-w-md -translate-x-1/2 -translate-y-1/2 rounded-xl border border-gray-200 bg-white p-6 shadow-xl focus:outline-none dark:border-gray-700 dark:bg-gray-900">
            <AlertDialog.Title className="text-lg font-semibold text-gray-900 dark:text-gray-100">
              Move to trash?
            </AlertDialog.Title>
            <AlertDialog.Description className="mt-2 text-sm text-gray-500 dark:text-gray-400">
              Are you sure you want to move <span className="font-medium text-gray-700 dark:text-gray-300">&ldquo;{deleteTarget?.name}&rdquo;</span> to trash?
            </AlertDialog.Description>
            <div className="mt-6 flex justify-end gap-2">
              <AlertDialog.Cancel asChild>
                <button className="rounded-lg border border-gray-300 bg-white px-4 py-2 text-sm font-medium text-gray-700 hover:bg-gray-50 dark:border-gray-700 dark:bg-gray-800 dark:text-gray-300 dark:hover:bg-gray-700">
                  Cancel
                </button>
              </AlertDialog.Cancel>
              <AlertDialog.Action asChild>
                <button
                  onClick={confirmDeleteFolder}
                  disabled={isDeleting}
                  className="inline-flex items-center gap-2 rounded-lg bg-red-600 px-4 py-2 text-sm font-medium text-white hover:bg-red-700 disabled:cursor-not-allowed disabled:opacity-50 dark:bg-red-700 dark:hover:bg-red-600"
                >
                  {isDeleting ? (
                    <>
                      <Loader2 size={16} className="animate-spin" />
                      Deleting...
                    </>
                  ) : (
                    <>Delete</>
                  )}
                </button>
              </AlertDialog.Action>
            </div>
          </AlertDialog.Content>
        </AlertDialog.Portal>
      </AlertDialog.Root>
    </>
  );
}
