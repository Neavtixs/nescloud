"use client";

import { useAtomValue } from "jotai";
import { userAtom } from "@/lib/atoms/auth-atoms";
import { useRouter } from "next/navigation";
import SidebarNav from "./sidebar-nav";
import { authApi } from "@/lib/api/api-call";
import { invalidateAuth } from "@/components/provider/auth-provider";
import { HardDrive, LogOut } from "lucide-react";

export default function Sidebar() {
  const user = useAtomValue(userAtom);
  const router = useRouter();

  async function handleLogout() {
    try {
      await authApi.logout();
    } finally {
      invalidateAuth();
      router.replace("/login");
    }
  }

  return (
    <aside className="flex h-screen w-[250px] flex-col border-r border-gray-200 bg-white">
      <div className="flex h-14 items-center gap-2 border-b border-gray-200 px-4">
        <HardDrive size={22} className="text-blue-600" />
        <span className="text-lg font-semibold text-gray-900">NESCloud</span>
      </div>

      <div className="py-4">
        <SidebarNav />
      </div>

      <div className="mt-auto border-t border-gray-200 px-4 py-4">
        <div className="mb-3">
          <div className="mb-1 text-xs font-medium text-gray-500">
            Storage
          </div>
          <div className="flex items-center justify-between text-xs text-gray-600">
            <span>0 GB</span>
            <span>1 GB</span>
          </div>
          <div className="mt-1.5 h-1.5 w-full overflow-hidden rounded-full bg-gray-200">
            <div
              className="h-full rounded-full bg-blue-500 transition-all"
              style={{ width: "0%" }}
            />
          </div>
        </div>

        {user && (
          <div className="flex items-center gap-3 border-t border-gray-200 pt-3">
            <div className="flex h-8 w-8 items-center justify-center rounded-full bg-blue-100 text-sm font-medium text-blue-700">
              {user.name.charAt(0).toUpperCase()}
            </div>
            <div className="min-w-0 flex-1">
              <p className="truncate text-sm font-medium text-gray-900">
                {user.name}
              </p>
              <p className="truncate text-xs text-gray-500">{user.email}</p>
            </div>
            <button
              onClick={handleLogout}
              className="rounded-lg p-1.5 text-gray-400 hover:bg-gray-100 hover:text-gray-600"
              title="Logout"
            >
              <LogOut size={16} />
            </button>
          </div>
        )}
      </div>
    </aside>
  );
}
