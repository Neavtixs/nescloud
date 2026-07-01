"use client";

import { useEffect, useState } from "react";
import {
  Trash2,
  RotateCcw,
  AlertTriangle,
  Folder,
  Loader2,
} from "lucide-react";
import { api } from "@/lib/api/api-call";
import type { ApiResponse } from "@/lib/api/api-response";

type TrashItem = {
  id: string;
  name: string;
  parent_folder_id: string;
  deleted_at: string;
};

function formatDate(iso: string): string {
  return new Date(iso).toLocaleDateString("en-US", {
    year: "numeric",
    month: "short",
    day: "numeric",
  });
}

export default function TrashPage() {
  const [items, setItems] = useState<TrashItem[]>([]);
  const [selectedIds, setSelectedIds] = useState<Set<string>>(new Set());
  const [isLoading, setIsLoading] = useState(true);
  const [restoringIds, setRestoringIds] = useState<Set<string>>(new Set());
  const [deletingIds, setDeletingIds] = useState<Set<string>>(new Set());

  const allSelected = items.length > 0 && selectedIds.size === items.length;

  useEffect(() => {
    let cancelled = false;
    api
      .get<ApiResponse<TrashItem[]>>("/trash/folders")
      .then((res) => {
        if (!cancelled) setItems(res.data);
      })
      .catch(() => {})
      .finally(() => {
        if (!cancelled) setIsLoading(false);
      });
    return () => {
      cancelled = true;
    };
  }, []);

  function toggleSelect(id: string) {
    setSelectedIds((prev) => {
      const next = new Set(prev);
      if (next.has(id)) next.delete(id);
      else next.add(id);
      return next;
    });
  }

  function toggleSelectAll() {
    if (allSelected) {
      setSelectedIds(new Set());
    } else {
      setSelectedIds(new Set(items.map((i) => i.id)));
    }
  }

  async function handleRestore(id: string) {
    setRestoringIds((prev) => new Set(prev).add(id));
    try {
      await api.post<ApiResponse<null>>(`/trash/folders/${id}/restore`, {});
      setItems((prev) => prev.filter((i) => i.id !== id));
      setSelectedIds((prev) => {
        const next = new Set(prev);
        next.delete(id);
        return next;
      });
    } catch (err: unknown) {
      const msg = err instanceof Error ? err.message : "Failed to restore";
      alert(msg);
    } finally {
      setRestoringIds((prev) => {
        const next = new Set(prev);
        next.delete(id);
        return next;
      });
    }
  }

  async function handleDeleteSingle(id: string, name: string) {
    if (
      !confirm(
        `Permanently delete "${name}"? This cannot be undone.`,
      )
    )
      return;
    setDeletingIds((prev) => new Set(prev).add(id));
    try {
      await api.delete<ApiResponse<null>>(`/trash/folders/${id}`);
      setItems((prev) => prev.filter((i) => i.id !== id));
      setSelectedIds((prev) => {
        const next = new Set(prev);
        next.delete(id);
        return next;
      });
    } catch (err: unknown) {
      const msg = err instanceof Error ? err.message : "Failed to delete";
      alert(msg);
    } finally {
      setDeletingIds((prev) => {
        const next = new Set(prev);
        next.delete(id);
        return next;
      });
    }
  }

  async function handleDeleteSelected() {
    const count = selectedIds.size;
    if (
      !confirm(
        `Permanently delete ${count} selected ${count === 1 ? "item" : "items"}? This cannot be undone.`,
      )
    )
      return;
    try {
      await Promise.all(
        Array.from(selectedIds).map((id) =>
          api.delete<ApiResponse<null>>(`/trash/folders/${id}`),
        ),
      );
      setItems((prev) => prev.filter((i) => !selectedIds.has(i.id)));
      setSelectedIds(new Set());
    } catch (err: unknown) {
      const msg = err instanceof Error ? err.message : "Failed to delete items";
      alert(msg);
    }
  }

  async function handleEmptyTrash() {
    if (!confirm("Permanently delete all items in trash?")) return;
    try {
      await Promise.all(
        items.map((i) =>
          api.delete<ApiResponse<null>>(`/trash/folders/${i.id}`),
        ),
      );
      setItems([]);
      setSelectedIds(new Set());
    } catch (err: unknown) {
      const msg =
        err instanceof Error ? err.message : "Failed to empty trash";
      alert(msg);
    }
  }

  return (
    <div className="mx-auto max-w-5xl">
      <div className="mb-6 flex items-center justify-between">
        <h1 className="text-2xl font-semibold text-gray-900 dark:text-gray-100">
          Trash
        </h1>
        <div className="flex items-center gap-3">
          {selectedIds.size > 0 && (
            <span className="text-sm text-gray-500 dark:text-gray-400">
              {selectedIds.size} selected
            </span>
          )}
          {items.length > 0 && (
            <button
              className="inline-flex items-center gap-2 rounded-lg bg-red-600 px-4 py-2 text-sm font-medium text-white hover:bg-red-700 dark:bg-red-700 dark:hover:bg-red-600"
              onClick={selectedIds.size > 0 ? handleDeleteSelected : handleEmptyTrash}
            >
              <Trash2 size={16} />
              {selectedIds.size > 0 ? "Delete Selected" : "Empty Trash"}
            </button>
          )}
        </div>
      </div>

      <div className="mb-6 flex items-start gap-2 rounded-lg border border-amber-200 bg-amber-50 px-4 py-3 text-sm text-amber-800 dark:border-amber-900 dark:bg-amber-900/20 dark:text-amber-300">
        <AlertTriangle size={16} className="mt-0.5 shrink-0" />
        <span>
          Items in trash are permanently deleted after 24 hours. Restore items
          you want to keep.
        </span>
      </div>

      {isLoading ? (
        <div className="flex items-center justify-center py-20">
          <Loader2 size={32} className="animate-spin text-gray-400" />
        </div>
      ) : items.length === 0 ? (
        <div className="flex flex-col items-center gap-3 rounded-xl border border-dashed border-gray-300 bg-white py-16 dark:border-gray-700 dark:bg-gray-900">
          <Trash2 size={40} className="text-gray-300 dark:text-gray-600" />
          <div className="text-center">
            <p className="text-sm font-medium text-gray-900 dark:text-gray-100">
              Trash is empty
            </p>
            <p className="mt-1 text-xs text-gray-500 dark:text-gray-400">
              Deleted files and folders will appear here
            </p>
          </div>
        </div>
      ) : (
        <div className="overflow-hidden rounded-xl border border-gray-200 bg-white dark:border-gray-800 dark:bg-gray-900">
          <table className="w-full">
            <thead>
              <tr className="border-b border-gray-200 text-left text-xs font-medium text-gray-500 dark:border-gray-800 dark:text-gray-400">
                <th className="w-10 px-5 py-3">
                  <input
                    type="checkbox"
                    checked={allSelected}
                    onChange={toggleSelectAll}
                    className="rounded border-gray-300 dark:border-gray-600 dark:bg-gray-800"
                  />
                </th>
                <th className="px-5 py-3 font-medium">Name</th>
                <th className="px-5 py-3 font-medium">Deleted</th>
                <th className="px-5 py-3 font-medium">Actions</th>
              </tr>
            </thead>
            <tbody>
              {items.map((item) => (
                <tr
                  key={item.id}
                  className={`border-b border-gray-50 hover:bg-gray-50 dark:border-gray-800 dark:hover:bg-gray-800 last:border-b-0 ${
                    selectedIds.has(item.id) ? "bg-blue-50 dark:bg-blue-900/10" : ""
                  }`}
                >
                  <td className="px-5 py-3">
                    <input
                      type="checkbox"
                      checked={selectedIds.has(item.id)}
                      onChange={() => toggleSelect(item.id)}
                      className="rounded border-gray-300 dark:border-gray-600 dark:bg-gray-800"
                    />
                  </td>
                  <td className="px-5 py-3">
                    <div className="flex items-center gap-3">
                      <Folder size={18} className="text-blue-500 dark:text-blue-400" />
                      <span className="text-sm font-medium text-gray-900 dark:text-gray-100">
                        {item.name}
                      </span>
                    </div>
                  </td>
                  <td className="px-5 py-3 text-sm text-gray-500 dark:text-gray-400">
                    {formatDate(item.deleted_at)}
                  </td>
                  <td className="px-5 py-3">
                    <div className="flex items-center gap-1">
                      <button
                        className="inline-flex items-center gap-1 rounded px-2 py-1 text-xs font-medium text-blue-600 hover:bg-blue-50 dark:text-blue-400 dark:hover:bg-blue-900/30 disabled:opacity-50"
                        onClick={() => handleRestore(item.id)}
                        disabled={restoringIds.has(item.id)}
                      >
                        {restoringIds.has(item.id) ? (
                          <Loader2 size={13} className="animate-spin" />
                        ) : (
                          <RotateCcw size={13} />
                        )}
                        Restore
                      </button>
                      <button
                        className="inline-flex items-center gap-1 rounded px-2 py-1 text-xs font-medium text-red-600 hover:bg-red-50 dark:text-red-400 dark:hover:bg-red-900/20 disabled:opacity-50"
                        onClick={() => handleDeleteSingle(item.id, item.name)}
                        disabled={deletingIds.has(item.id)}
                      >
                        {deletingIds.has(item.id) ? (
                          <Loader2 size={13} className="animate-spin" />
                        ) : (
                          <Trash2 size={13} />
                        )}
                        Delete
                      </button>
                    </div>
                  </td>
                </tr>
              ))}
            </tbody>
          </table>
        </div>
      )}
    </div>
  );
}
