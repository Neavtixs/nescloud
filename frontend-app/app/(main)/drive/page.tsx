"use client";

import { useState } from "react";
import {
  Folder,
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

const mockItems: DriveItem[] = [
  {
    id: "d1",
    type: "folder",
    name: "Documents",
    modified_at: "2026-06-28T14:30:00Z",
  },
  {
    id: "d2",
    type: "folder",
    name: "Photos",
    modified_at: "2026-06-27T11:15:00Z",
  },
  {
    id: "d3",
    type: "folder",
    name: "Work",
    modified_at: "2026-06-25T09:00:00Z",
  },
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

function FileIcon({ mime_type }: { mime_type: string }) {
  if (mime_type.startsWith("image/"))
    return <ImageIcon size={18} className="text-purple-500 dark:text-purple-400" />;
  if (mime_type.startsWith("text/"))
    return <FileText size={18} className="text-gray-500 dark:text-gray-400" />;
  if (
    mime_type.includes("spreadsheet") ||
    mime_type.includes("excel") ||
    mime_type.includes("csv")
  )
    return <FileSpreadsheet size={18} className="text-green-600 dark:text-green-400" />;
  if (mime_type.includes("presentation") || mime_type.includes("powerpoint"))
    return <FileText size={18} className="text-orange-500 dark:text-orange-400" />;
  if (mime_type === "application/pdf")
    return <FileText size={18} className="text-red-500 dark:text-red-400" />;
  return <File size={18} className="text-gray-400 dark:text-gray-500" />;
}

function ItemMenu({
  type,
  onClose,
}: {
  type: ItemType;
  onClose: () => void;
}) {
  return (
    <>
      <div className="fixed inset-0 z-10" onClick={onClose} />
      <div className="absolute right-0 top-8 z-20 w-44 rounded-lg border border-gray-200 bg-white py-1 shadow-lg dark:border-gray-800 dark:bg-gray-900">
        {type === "file" && (
          <button
            className="flex w-full items-center gap-2 px-3 py-2 text-sm text-gray-700 hover:bg-gray-50 dark:text-gray-300 dark:hover:bg-gray-800"
            onClick={() => {
              onClose();
              alert("Download coming soon");
            }}
          >
            <Download size={14} />
            Download
          </button>
        )}
        <button
          className="flex w-full items-center gap-2 px-3 py-2 text-sm text-gray-700 hover:bg-gray-50 dark:text-gray-300 dark:hover:bg-gray-800"
          onClick={() => {
            onClose();
            alert("Rename coming soon");
          }}
        >
          <Pencil size={14} />
          Rename
        </button>
        {type === "file" && (
          <button
            className="flex w-full items-center gap-2 px-3 py-2 text-sm text-gray-700 hover:bg-gray-50 dark:text-gray-300 dark:hover:bg-gray-800"
            onClick={() => {
              onClose();
              alert("Share coming soon");
            }}
          >
            <Link size={14} />
            Share
          </button>
        )}
        <button
          className="flex w-full items-center gap-2 px-3 py-2 text-sm text-red-600 hover:bg-red-50 dark:text-red-400 dark:hover:bg-red-900/20"
          onClick={() => {
            onClose();
            alert("Delete coming soon");
          }}
        >
          <Trash2 size={14} />
          Delete
        </button>
      </div>
    </>
  );
}

export default function DrivePage() {
  const [search, setSearch] = useState("");
  const [viewMode, setViewMode] = useState<"list" | "grid">("list");
  const [openMenuId, setOpenMenuId] = useState<string | null>(null);

  const breadcrumb = [{ label: "My Drive", href: "/drive" }];

  const filtered = search
    ? mockItems.filter((item) =>
        item.name.toLowerCase().includes(search.toLowerCase()),
      )
    : mockItems;

  const sorted = [...filtered].sort((a, b) => {
    if (a.type !== b.type) return a.type === "folder" ? -1 : 1;
    return a.name.localeCompare(b.name);
  });

  return (
    <div className="mx-auto max-w-5xl">
      <div className="mb-6 flex items-center gap-2 text-sm text-gray-500 dark:text-gray-400">
        {breadcrumb.map((crumb, i) => (
          <span key={crumb.href} className="flex items-center gap-2">
            {i > 0 && <ChevronRight size={14} />}
            {i === 0 && <Home size={14} />}
            <span
              className={
                i === breadcrumb.length - 1
                  ? "font-medium text-gray-900 dark:text-gray-100"
                  : ""
              }
            >
              {crumb.label}
            </span>
          </span>
        ))}
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
            onChange={(e) => setSearch(e.target.value)}
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
          <button
            className="inline-flex items-center gap-2 rounded-lg border border-gray-300 bg-white px-3 py-2 text-sm font-medium text-gray-700 hover:bg-gray-50 dark:border-gray-700 dark:bg-gray-800 dark:text-gray-300 dark:hover:bg-gray-700"
            onClick={() => alert("New folder coming soon")}
          >
            <FolderPlus size={16} />
            New Folder
          </button>
          <button
            className="inline-flex items-center gap-2 rounded-lg bg-blue-600 px-4 py-2 text-sm font-medium text-white hover:bg-blue-700 dark:bg-blue-700 dark:hover:bg-blue-600"
            onClick={() => alert("Upload coming soon")}
          >
            <Upload size={16} />
            Upload
          </button>
        </div>
      </div>

      {sorted.length === 0 ? (
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
              onClick={() => alert("Upload coming soon")}
            >
              <Upload size={16} />
              Upload Files
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
              {sorted.map((item) => (
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
                  <td className="px-5 py-3">
                    <div className="flex items-center gap-3">
                      {item.type === "folder" ? (
                        <Folder size={18} className="text-blue-500 dark:text-blue-400" />
                      ) : (
                        <FileIcon mime_type={item.mime_type!} />
                      )}
                      <span
                        className={`text-sm font-medium ${
                          item.type === "folder"
                            ? "cursor-pointer text-blue-600 hover:underline dark:text-blue-400 dark:hover:text-blue-300"
                            : "text-gray-900 dark:text-gray-100"
                        }`}
                        onClick={() => {
                          if (item.type === "folder")
                            alert(`Navigate to folder: ${item.name}`);
                        }}
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
                    <button
                      className="rounded p-1 text-gray-400 hover:bg-gray-100 hover:text-gray-600 dark:text-gray-500 dark:hover:bg-gray-800 dark:hover:text-gray-300"
                      onClick={() =>
                        setOpenMenuId(
                          openMenuId === item.id ? null : item.id,
                        )
                      }
                    >
                      <MoreHorizontal size={15} />
                    </button>

                    {openMenuId === item.id && (
                      <ItemMenu
                        type={item.type}
                        onClose={() => setOpenMenuId(null)}
                      />
                    )}
                  </td>
                </tr>
              ))}
            </tbody>
          </table>

          <div className="flex items-center justify-between border-t border-gray-200 px-5 py-3 dark:border-gray-800">
            <p className="text-xs text-gray-500 dark:text-gray-400">
              1–{sorted.length} of {mockItems.length} items
            </p>
            <div className="flex items-center gap-1">
              <button className="rounded px-2 py-1 text-xs text-gray-400 hover:bg-gray-100 hover:text-gray-600 dark:text-gray-500 dark:hover:bg-gray-800 dark:hover:text-gray-300">
                ←
              </button>
              <button className="rounded bg-blue-50 px-2 py-1 text-xs font-medium text-blue-600 dark:bg-blue-900/30 dark:text-blue-400">
                1
              </button>
              <button className="rounded px-2 py-1 text-xs text-gray-500 hover:bg-gray-100 dark:text-gray-400 dark:hover:bg-gray-800">
                2
              </button>
              <button className="rounded px-2 py-1 text-xs text-gray-500 hover:bg-gray-100 dark:text-gray-400 dark:hover:bg-gray-800">
                3
              </button>
              <span className="px-1 text-xs text-gray-400 dark:text-gray-500">...</span>
              <button className="rounded px-2 py-1 text-xs text-gray-500 hover:bg-gray-100 dark:text-gray-400 dark:hover:bg-gray-800">
                6
              </button>
              <button className="rounded px-2 py-1 text-xs text-gray-400 hover:bg-gray-100 hover:text-gray-600 dark:text-gray-500 dark:hover:bg-gray-800 dark:hover:text-gray-300">
                →
              </button>
            </div>
          </div>
        </div>
      ) : (
        <div className="grid grid-cols-2 gap-4 sm:grid-cols-3 lg:grid-cols-4">
          {sorted.map((item) => (
            <div
              key={item.id}
              className="group relative rounded-xl border border-gray-200 bg-white p-4 dark:border-gray-800 dark:bg-gray-900"
            >
              <div className="flex items-start justify-between">
                {item.type === "folder" ? (
                  <div className="rounded-lg bg-blue-50 p-2.5 dark:bg-blue-900/30">
                    <Folder size={22} className="text-blue-500 dark:text-blue-400" />
                  </div>
                ) : (
                  <FileIcon mime_type={item.mime_type!} />
                )}
                <button
                  className="rounded p-1 text-gray-400 opacity-0 transition-opacity hover:bg-gray-100 hover:text-gray-600 group-hover:opacity-100 dark:text-gray-500 dark:hover:bg-gray-800 dark:hover:text-gray-300"
                  onClick={() =>
                    setOpenMenuId(openMenuId === item.id ? null : item.id)
                  }
                >
                  <MoreHorizontal size={15} />
                </button>

                {openMenuId === item.id && (
                  <ItemMenu
                    type={item.type}
                    onClose={() => setOpenMenuId(null)}
                  />
                )}
              </div>

              <p
                className={`mt-3 truncate text-sm font-medium ${
                  item.type === "folder"
                    ? "cursor-pointer text-blue-600 hover:underline dark:text-blue-400"
                    : "text-gray-900 dark:text-gray-100"
                }`}
                onClick={() => {
                  if (item.type === "folder")
                    alert(`Navigate to folder: ${item.name}`);
                }}
              >
                {item.name}
              </p>

              <p className="mt-1 text-xs text-gray-400 dark:text-gray-500">
                {item.type === "folder"
                  ? "--"
                  : formatSize(item.size_bytes!)}
                {" · "}
                {formatDate(item.modified_at)}
              </p>
            </div>
          ))}
        </div>
      )}
    </div>
  );
}
