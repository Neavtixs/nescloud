import AuthProvider from "@/components/provider/auth-provider";

export default function MainLayout({
  children,
}: {
  children: React.ReactNode;
}) {
  return <AuthProvider>{children}</AuthProvider>;
}
