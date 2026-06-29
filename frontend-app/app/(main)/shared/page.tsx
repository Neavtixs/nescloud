"use client";

import { useState } from "react";
import * as Tooltip from "@radix-ui/react-tooltip";
import {
  Link,
  Copy,
  XCircle,
  FileText,
  ImageIcon,
  FileSpreadsheet,
  File,
  Check,
} from "lucide-react";

type SharedItem = {
  id: string;
  file_name: string;
  mime_type: string;
  size_bytes: number;
  public_url: string;
  created_at: string;
};

const mockSharedItems: SharedItem[] = [
  {
    id: "s1",
    file_name: "report-q2.pdf",
    mime_type: "application/pdf",
    size_bytes: 2516582,
    public_url: "https://nescloud.example.com/public/abc123xyz",
    created_at: "2026-06-28T14:30:00Z",
  },
  {
    id: "s2",
    file_name: "banner-promo.png",
    mime_type: "image/png",
    size_bytes: 1153433,
    public_url: "https://nescloud.example.com/public/def456uvw",
    created_at: "2026-06-27T10:15:00Z",
  },
  {
    id: "s3",
    file_name: "data-report.xlsx",
    mime_type:
      "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet",
    size_bytes: 876544,
    public_url: "https://nescloud.example.com/public/ghi789rst",
    created_at: "2026-06-25T09:00:00Z",
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
    return (
      <ImageIcon size={18} className="text-purple-500 dark:text-purple-400" />
    );
  if (mime_type.startsWith("text/"))
    return <FileText size={18} className="text-gray-500 dark:text-gray-400" />;
  if (
    mime_type.includes("spreadsheet") ||
    mime_type.includes("excel") ||
    mime_type.includes("csv")
  )
    return (
      <FileSpreadsheet
        size={18}
        className="text-green-600 dark:text-green-400"
      />
    );
  if (mime_type.includes("presentation") || mime_type.includes("powerpoint"))
    return (
      <FileText size={18} className="text-orange-500 dark:text-orange-400" />
    );
  if (mime_type === "application/pdf")
    return <FileText size={18} className="text-red-500 dark:text-red-400" />;
  return <File size={18} className="text-gray-400 dark:text-gray-500" />;
}

export default function SharedPage() {
  const [copiedId, setCopiedId] = useState<string | null>(null);

  function handleCopy(url: string, id: string) {
    navigator.clipboard.writeText(url).then(() => {
      setCopiedId(id);
      setTimeout(() => setCopiedId(null), 2000);
    });
  }

  return (
    <div className="mx-auto max-w-5xl">
      <div className="mb-6">
        <h1 className="text-2xl font-semibold text-gray-900 dark:text-gray-100">
          Shared Links
        </h1>
        <p className="mt-1 text-sm text-gray-500 dark:text-gray-400">
          Files shared via public link. Anyone with the link can view them.
        </p>
      </div>

      {mockSharedItems.length === 0 ? (
        <div className="flex flex-col items-center gap-3 rounded-xl border border-dashed border-gray-300 bg-white py-16 dark:border-gray-700 dark:bg-gray-900">
          <Link size={40} className="text-gray-300 dark:text-gray-600" />
          <div className="text-center">
            <p className="text-sm font-medium text-gray-900 dark:text-gray-100">
              No shared links
            </p>
            <p className="mt-1 text-xs text-gray-500 dark:text-gray-400">
              Create a public link for any file from My Drive to share it
            </p>
          </div>
        </div>
      ) : (
        <Tooltip.Provider delayDuration={0}>
          <div className="overflow-hidden rounded-xl border border-gray-200 bg-white dark:border-gray-800 dark:bg-gray-900">
            <table className="w-full">
              <thead>
                <tr className="border-b border-gray-200 text-left text-xs font-medium text-gray-500 dark:border-gray-800 dark:text-gray-400">
                  <th className="px-5 py-3 font-medium">File</th>
                  <th className="hidden px-5 py-3 font-medium md:table-cell">
                    Size
                  </th>
                  <th className="px-5 py-3 font-medium">Public Link</th>
                  <th className="hidden px-5 py-3 font-medium sm:table-cell">
                    Created
                  </th>
                  <th className="px-5 py-3 font-medium">Actions</th>
                </tr>
              </thead>
              <tbody>
                {mockSharedItems.map((item) => (
                  <tr
                    key={item.id}
                    className="border-b border-gray-50 hover:bg-gray-50 dark:border-gray-800 dark:hover:bg-gray-800 last:border-b-0"
                  >
                    <td className="px-5 py-3">
                      <div className="flex items-center gap-3">
                        <FileIcon mime_type={item.mime_type} />
                        <span className="text-sm font-medium text-gray-900 dark:text-gray-100">
                          {item.file_name}
                        </span>
                      </div>
                    </td>
                    <td className="hidden px-5 py-3 text-sm text-gray-500 dark:text-gray-400 md:table-cell">
                      {formatSize(item.size_bytes)}
                    </td>
                  <td className="max-w-[180px] px-5 py-3">
                    <Tooltip.Root>
                      <Tooltip.Trigger asChild>
                        <span className="block truncate text-xs text-gray-500 dark:text-gray-400">
                          {item.public_url}
                        </span>
                      </Tooltip.Trigger>
                      <Tooltip.Portal>
                        <Tooltip.Content
                          side="top"
                          sideOffset={4}
                          className="max-w-[320px] break-all rounded-md bg-white px-2.5 py-1 text-xs text-gray-900 shadow-lg dark:bg-gray-700 dark:text-gray-100"
                        >
                          {item.public_url}
                          <Tooltip.Arrow className="fill-white dark:fill-gray-700" />
                        </Tooltip.Content>
                      </Tooltip.Portal>
                    </Tooltip.Root>
                  </td>
                    <td className="hidden px-5 py-3 text-sm text-gray-500 dark:text-gray-400 sm:table-cell">
                      {formatDate(item.created_at)}
                    </td>
                    <td className="px-5 py-3">
                      <div className="flex items-center gap-0.5">
                        <Tooltip.Root>
                          <Tooltip.Trigger asChild>
                            <button
                              className="rounded p-1.5 text-gray-400 hover:bg-gray-100 hover:text-gray-600 dark:text-gray-500 dark:hover:bg-gray-800 dark:hover:text-gray-300"
                              onClick={() =>
                                handleCopy(item.public_url, item.id)
                              }
                            >
                              {copiedId === item.id ? (
                                <Check
                                  size={15}
                                  className="text-green-600 dark:text-green-400"
                                />
                              ) : (
                                <Copy size={15} />
                              )}
                            </button>
                          </Tooltip.Trigger>
                          <Tooltip.Portal>
                            <Tooltip.Content
                              side="top"
                              sideOffset={4}
                              className="rounded-md bg-white px-2.5 py-1 text-xs text-gray-900 shadow-lg dark:bg-gray-700 dark:text-gray-100"
                            >
                              {copiedId === item.id ? "Copied" : "Copy"}
                              <Tooltip.Arrow className="fill-white dark:fill-gray-700" />
                            </Tooltip.Content>
                          </Tooltip.Portal>
                        </Tooltip.Root>

                        <Tooltip.Root>
                          <Tooltip.Trigger asChild>
                            <button
                              className="rounded p-1.5 text-gray-400 hover:bg-gray-100 hover:text-red-600 dark:text-gray-500 dark:hover:bg-gray-800 dark:hover:text-red-400"
                              onClick={() => {
                                if (
                                  confirm(
                                    `Revoke public link for "${item.file_name}"?`,
                                  )
                                ) {
                                  alert("Revoke coming soon");
                                }
                              }}
                            >
                              <XCircle size={15} />
                            </button>
                          </Tooltip.Trigger>
                          <Tooltip.Portal>
                            <Tooltip.Content
                              side="top"
                              sideOffset={4}
                              className="rounded-md bg-white px-2.5 py-1 text-xs text-gray-900 shadow-lg dark:bg-gray-700 dark:text-gray-100"
                            >
                              Revoke
                              <Tooltip.Arrow className="fill-white dark:fill-gray-700" />
                            </Tooltip.Content>
                          </Tooltip.Portal>
                        </Tooltip.Root>
                      </div>
                    </td>
                  </tr>
                ))}
              </tbody>
            </table>
          </div>
        </Tooltip.Provider>
      )}
    </div>
  );
}
