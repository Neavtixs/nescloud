"use client";

import { Sun, Moon } from "lucide-react";
import { useTheme } from "next-themes";
import { useEffect, useState } from "react";

export default function ThemeSwitcher() {
  const { theme, setTheme } = useTheme();
  const [mounted, setMounted] = useState(false);

  useEffect(() => {
    setMounted(true);
  }, []);

  if (!mounted) {
    return (
      <button
        disabled
        className="fixed bottom-6 right-6 z-50 flex h-11 w-11 items-center justify-center rounded-full border border-gray-200 bg-white text-gray-600 shadow-lg transition-all hover:scale-105 hover:shadow-xl dark:border-gray-700 dark:bg-gray-800 dark:text-gray-300"
      >
        <div className="h-[18px] w-[18px]" />
      </button>
    );
  }

  return (
    <button
      onClick={() => setTheme(theme === "dark" ? "light" : "dark")}
      className="fixed bottom-6 right-6 z-50 flex h-11 w-11 items-center justify-center rounded-full border border-gray-200 bg-white text-gray-600 shadow-lg transition-all hover:scale-105 hover:shadow-xl dark:border-gray-700 dark:bg-gray-800 dark:text-gray-300"
    >
      {theme === "dark" ? <Sun size={18} /> : <Moon size={18} />}
    </button>
  );
}
