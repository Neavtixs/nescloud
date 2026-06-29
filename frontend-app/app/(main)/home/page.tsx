"use client";

import { useAtomValue } from "jotai";
import { userAtom } from "@/lib/atoms/auth-atoms";
import {
  Folder,
  File,
  HardDrive,
  Link,
  Download,
  MoreHorizontal,
  Upload,
  FolderPlus,
  FileText,
  ImageIcon,
  FileSpreadsheet,
} from "lucide-react";
import LinkNext from "next/link";

type MockFile = {
  id: string;
  name: string;
  mime_type: string;
  size_bytes: number;
  modified_at: string;
};

const mockRecentFiles: MockFile[] = [
  {
    id: "f1",
    name: "report-q2.pdf",
    mime_type: "application/pdf",
    size_bytes: 2516582,
    modified_at: "2026-06-29T08:30:00Z",
  },
  {
    id: "f2",
    name: "holiday-photo.png",
    mime_type: "image/png",
    size_bytes: 1153433,
    modified_at: "2026-06-28T16:45:00Z",
  },
  {
    id: "f3",
    name: "budget-2026.xlsx",
    mime_type:
      "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet",
    size_bytes: 876544,
    modified_at: "2026-06-27T10:15:00Z",
  },
  {
    id: "f4",
    name: "meeting-notes.txt",
    mime_type: "text/plain",
    size_bytes: 12288,
    modified_at: "2026-06-26T14:00:00Z",
  },
  {
    id: "f5",
    name: "presentation.pptx",
    mime_type:
      "application/vnd.openxmlformats-officedocument.presentationml.presentation",
    size_bytes: 5242880,
    modified_at: "2026-06-25T09:20:00Z",
  },
];

const stats = {
  folders: 12,
  files: 48,
  used_bytes: 335544320,
  total_bytes: 1073741824,
  shared_links: 3,
};

function formatSize(bytes: number): string {
  if (bytes >= 1073741824) return `${(bytes / 1073741824).toFixed(1)} GB`;
  if (bytes >= 1048576) return `${(bytes / 1048576).toFixed(1)} MB`;
  if (bytes >= 1024) return `${(bytes / 1024).toFixed(0)} KB`;
  return `${bytes} B`;
}

function formatRelativeTime(iso: string): string {
  const diff = Date.now() - new Date(iso).getTime();
  const mins = Math.floor(diff / 60000);
  if (mins < 1) return "Just now";
  if (mins < 60) return `${mins} min ago`;
  const hrs = Math.floor(mins / 60);
  if (hrs < 24) return `${hrs} hours ago`;
  const days = Math.floor(hrs / 24);
  if (days === 1) return "Yesterday";
  if (days < 30) return `${days} days ago`;
  return new Date(iso).toLocaleDateString("en-US", {
    year: "numeric",
    month: "short",
    day: "numeric",
  });
}

function FileIcon({ mime_type }: { mime_type: string }) {
  if (mime_type.startsWith("image/"))
    return <ImageIcon size={18} className="text-purple-500" />;
  if (mime_type.startsWith("text/"))
    return <FileText size={18} className="text-gray-500" />;
  if (
    mime_type.includes("spreadsheet") ||
    mime_type.includes("excel") ||
    mime_type.includes("csv")
  )
    return <FileSpreadsheet size={18} className="text-green-600" />;
  if (mime_type.includes("presentation") || mime_type.includes("powerpoint"))
    return <FileText size={18} className="text-orange-500" />;
  if (mime_type === "application/pdf")
    return <FileText size={18} className="text-red-500" />;
  return <File size={18} className="text-gray-400" />;
}

export default function HomePage() {
  const user = useAtomValue(userAtom);
  const used = formatSize(stats.used_bytes);
  const total = formatSize(stats.total_bytes);
  const usagePercent = Math.round((stats.used_bytes / stats.total_bytes) * 100);

  return (
    <div className="mx-auto max-w-5xl">
      <div className="mb-8">
        <h1 className="text-2xl font-semibold text-gray-900">
          Welcome back{user ? `, ${user.name.split(" ")[0]}` : ""}
        </h1>
        <p className="mt-1 text-sm text-gray-500">
          Overview of your cloud storage
        </p>
      </div>

      <div className="mb-8 grid grid-cols-2 gap-4 lg:grid-cols-4">
        <div className="rounded-xl border border-gray-200 bg-white p-4">
          <div className="flex items-center gap-2">
            <div className="rounded-lg bg-blue-50 p-2">
              <Folder size={18} className="text-blue-600" />
            </div>
            <span className="text-xs font-medium text-gray-500">Folders</span>
          </div>
          <p className="mt-2 text-2xl font-semibold text-gray-900">
            {stats.folders}
          </p>
        </div>

        <div className="rounded-xl border border-gray-200 bg-white p-4">
          <div className="flex items-center gap-2">
            <div className="rounded-lg bg-green-50 p-2">
              <File size={18} className="text-green-600" />
            </div>
            <span className="text-xs font-medium text-gray-500">Files</span>
          </div>
          <p className="mt-2 text-2xl font-semibold text-gray-900">
            {stats.files}
          </p>
        </div>

        <div className="rounded-xl border border-gray-200 bg-white p-4">
          <div className="flex items-center gap-2">
            <div className="rounded-lg bg-purple-50 p-2">
              <HardDrive size={18} className="text-purple-600" />
            </div>
            <span className="text-xs font-medium text-gray-500">Storage</span>
          </div>
          <p className="mt-2 text-2xl font-semibold text-gray-900">
            {used}
          </p>
          <div className="mt-2 h-1.5 w-full overflow-hidden rounded-full bg-gray-100">
            <div
              className="h-full rounded-full bg-purple-500"
              style={{ width: `${usagePercent}%` }}
            />
          </div>
          <p className="mt-1 text-xs text-gray-400">
            {usagePercent}% of {total}
          </p>
        </div>

        <div className="rounded-xl border border-gray-200 bg-white p-4">
          <div className="flex items-center gap-2">
            <div className="rounded-lg bg-amber-50 p-2">
              <Link size={18} className="text-amber-600" />
            </div>
            <span className="text-xs font-medium text-gray-500">Shared</span>
          </div>
          <p className="mt-2 text-2xl font-semibold text-gray-900">
            {stats.shared_links}
          </p>
        </div>
      </div>

      <div className="mb-6 rounded-xl border border-gray-200 bg-white">
        <div className="flex items-center justify-between border-b border-gray-200 px-5 py-4">
          <h2 className="text-sm font-semibold text-gray-900">
            Recent Files
          </h2>
          <LinkNext
            href="/drive"
            className="text-xs font-medium text-blue-600 hover:text-blue-700"
          >
            View all
          </LinkNext>
        </div>

        {mockRecentFiles.length === 0 ? (
          <div className="flex flex-col items-center gap-2 py-12 text-gray-400">
            <File size={32} />
            <p className="text-sm">No recent files</p>
          </div>
        ) : (
          <table className="w-full">
            <thead>
              <tr className="border-b border-gray-100 text-left text-xs font-medium text-gray-500">
                <th className="px-5 py-3 font-medium">Name</th>
                <th className="px-5 py-3 font-medium">Size</th>
                <th className="hidden px-5 py-3 font-medium sm:table-cell">
                  Modified
                </th>
                <th className="px-5 py-3 font-medium"></th>
              </tr>
            </thead>
            <tbody>
              {mockRecentFiles.map((file) => (
                <tr
                  key={file.id}
                  className="border-b border-gray-50 hover:bg-gray-50"
                >
                  <td className="px-5 py-3">
                    <div className="flex items-center gap-3">
                      <FileIcon mime_type={file.mime_type} />
                      <span className="text-sm font-medium text-gray-900">
                        {file.name}
                      </span>
                    </div>
                  </td>
                  <td className="px-5 py-3 text-sm text-gray-500">
                    {formatSize(file.size_bytes)}
                  </td>
                  <td className="hidden px-5 py-3 text-sm text-gray-500 sm:table-cell">
                    {formatRelativeTime(file.modified_at)}
                  </td>
                  <td className="px-5 py-3">
                    <div className="flex items-center gap-1">
                      <button
                        className="rounded p-1 text-gray-400 hover:bg-gray-100 hover:text-gray-600"
                        title="Download"
                        onClick={() => alert("Download coming soon")}
                      >
                        <Download size={15} />
                      </button>
                      <button
                        className="rounded p-1 text-gray-400 hover:bg-gray-100 hover:text-gray-600"
                        title="More"
                        onClick={() => alert("Actions coming soon")}
                      >
                        <MoreHorizontal size={15} />
                      </button>
                    </div>
                  </td>
                </tr>
              ))}
            </tbody>
          </table>
        )}
      </div>

      <div className="flex items-center gap-3">
        <button
          className="inline-flex items-center gap-2 rounded-lg bg-blue-600 px-4 py-2 text-sm font-medium text-white hover:bg-blue-700"
          onClick={() => alert("Upload coming soon")}
        >
          <Upload size={16} />
          Upload
        </button>
        <button
          className="inline-flex items-center gap-2 rounded-lg border border-gray-300 bg-white px-4 py-2 text-sm font-medium text-gray-700 hover:bg-gray-50"
          onClick={() => alert("New folder coming soon")}
        >
          <FolderPlus size={16} />
          New Folder
        </button>
      </div>
    </div>
  );
}
