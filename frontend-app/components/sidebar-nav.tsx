"use client";

import Link from "next/link";
import { usePathname } from "next/navigation";
import { House, Folder, Trash2, Link as LinkIcon } from "lucide-react";

type NavItem = {
  href: string;
  label: string;
  icon: typeof House;
};

const navItems: NavItem[] = [
  { href: "/home", label: "Home", icon: House },
  { href: "/drive", label: "My Drive", icon: Folder },
  { href: "/trash", label: "Trash", icon: Trash2 },
  { href: "/shared", label: "Shared", icon: LinkIcon },
];

export default function SidebarNav() {
  const pathname = usePathname();

  return (
    <nav className="flex flex-col gap-1 px-3">
      {navItems.map((item) => {
        const Icon = item.icon;
        const isActive =
          pathname === item.href ||
          (item.href !== "/home" && pathname.startsWith(item.href));

        return (
          <Link
            key={item.href}
            href={item.href}
            className={`flex items-center gap-3 rounded-lg px-3 py-2 text-sm font-medium transition-colors ${
              isActive
                ? "bg-blue-50 text-blue-700 dark:bg-blue-900/30 dark:text-blue-400"
                : "text-gray-700 hover:bg-gray-100 dark:text-gray-300 dark:hover:bg-gray-800"
            }`}
          >
            <Icon size={18} />
            {item.label}
          </Link>
        );
      })}
    </nav>
  );
}
