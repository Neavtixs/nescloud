import AuthProvider from "@/components/provider/auth-provider";
import Sidebar from "@/components/sidebar";
import ThemeSwitcher from "@/components/theme-switcher";

export default function MainLayout({
  children,
}: {
  children: React.ReactNode;
}) {
  return (
    <AuthProvider>
      <div className="flex h-screen">
        <Sidebar />
        <main className="flex-1 overflow-auto bg-gray-50 p-6 dark:bg-gray-950">
          {children}
        </main>
      </div>
    </AuthProvider>
  );
}
