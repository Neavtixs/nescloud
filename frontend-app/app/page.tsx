import Link from "next/link";
import {
  HardDrive,
  Shield,
  Zap,
  Share2,
  Cloud,
  CheckCircle,
  ArrowRight,
  Users,
} from "lucide-react";

const features = [
  {
    icon: Shield,
    title: "Secure & Encrypted",
    description:
      "Your files are encrypted and stored securely. Only you have full access to your data.",
    iconBg: "bg-blue-50 dark:bg-blue-900/30",
    iconColor: "text-blue-600 dark:text-blue-400",
  },
  {
    icon: Zap,
    title: "Fast Access",
    description:
      "Access your files from any device, anytime. Automatic sync keeps your files up to date.",
    iconBg: "bg-amber-50 dark:bg-amber-900/20",
    iconColor: "text-amber-600 dark:text-amber-400",
  },
  {
    icon: Share2,
    title: "Easy Sharing",
    description:
      "Create public links to share files with anyone. Set permissions and expiration as needed.",
    iconBg: "bg-green-50 dark:bg-green-900/20",
    iconColor: "text-green-600 dark:text-green-400",
  },
  {
    icon: Cloud,
    title: "Generous Storage",
    description:
      "Enjoy up to 1 GB of free storage. Upgrade your capacity anytime based on your needs.",
    iconBg: "bg-purple-50 dark:bg-purple-900/20",
    iconColor: "text-purple-600 dark:text-purple-400",
  },
];

const stats = [
  { icon: Users, value: "10,000+", label: "Active Users" },
  { icon: CheckCircle, value: "99.9%", label: "Uptime" },
  { icon: HardDrive, value: "1 GB", label: "Free Storage" },
];

export default function Home() {
  return (
    <div className="flex min-h-screen flex-col dark:bg-gray-950">
      <header className="sticky top-0 z-50 border-b border-gray-200 bg-white/20 backdrop-blur-md dark:border-gray-800 dark:bg-gray-950/20">
        <div className="mx-auto flex h-14 max-w-6xl items-center justify-between px-6">
          <Link href="/" className="flex items-center gap-2">
            <HardDrive size={22} className="text-blue-600 dark:text-blue-400" />
            <span className="text-lg font-semibold text-gray-900 dark:text-gray-100">
              NESCloud
            </span>
          </Link>

          <nav className="flex items-center gap-4">
            <Link
              href="/login"
              className="text-sm font-medium text-gray-600 hover:text-gray-900 dark:text-gray-400 dark:hover:text-gray-100"
            >
              Sign In
            </Link>
            <Link
              href="/register"
              className="rounded-lg bg-blue-600 px-4 py-2 text-sm font-medium text-white hover:bg-blue-700 dark:bg-blue-700 dark:hover:bg-blue-600"
            >
              Get Started
            </Link>
          </nav>
        </div>
      </header>

      <section className="flex-1 bg-gray-50 dark:bg-gray-950">
        <div className="mx-auto flex max-w-6xl flex-col items-center px-6 py-20 text-center lg:flex-row lg:text-left">
          <div className="flex-1">
            <h1 className="text-4xl font-bold leading-tight tracking-tight text-gray-900 dark:text-gray-100 sm:text-5xl">
              Store, Access, and Share{" "}
              <span className="text-blue-600 dark:text-blue-400">
                Your Files
              </span>{" "}
              Securely
            </h1>
            <p className="mx-auto mt-4 max-w-xl text-lg text-gray-500 dark:text-gray-400 lg:mx-0">
              NESCloud is your personal cloud storage. Upload, organize, and
              share files easily from anywhere, anytime.
            </p>

            <div className="mt-8 flex flex-col items-center gap-3 sm:flex-row lg:justify-start">
              <Link
                href="/register"
                className="inline-flex items-center gap-2 rounded-lg bg-blue-600 px-6 py-3 text-sm font-medium text-white hover:bg-blue-700 dark:bg-blue-700 dark:hover:bg-blue-600"
              >
                Get Started Free
                <ArrowRight size={16} />
              </Link>
              <Link
                href="/login"
                className="inline-flex items-center gap-2 rounded-lg border border-gray-300 bg-white px-6 py-3 text-sm font-medium text-gray-700 hover:bg-gray-50 dark:border-gray-700 dark:bg-gray-800 dark:text-gray-400 dark:hover:bg-gray-700"
              >
                Sign In
              </Link>
            </div>
          </div>

          <div className="mt-12 hidden flex-1 lg:mt-0 lg:flex lg:justify-center">
            <div className="relative">
              <div className="flex h-72 w-72 items-center justify-center rounded-3xl border border-gray-200 bg-white shadow-xl dark:border-gray-800 dark:bg-gray-900">
                <HardDrive
                  size={80}
                  className="text-blue-600/30 dark:text-blue-400/30"
                />
              </div>
              <div className="absolute -bottom-4 -right-4 rounded-2xl border border-gray-200 bg-white px-4 py-3 shadow-lg dark:border-gray-800 dark:bg-gray-900">
                <div className="flex items-center gap-2 text-sm text-gray-600 dark:text-gray-400">
                  <CheckCircle
                    size={16}
                    className="text-green-500 dark:text-green-400"
                  />
                  <span className="font-medium text-gray-900 dark:text-gray-100">
                    +500
                  </span>{" "}
                  files secured
                </div>
              </div>
              <div className="absolute -left-4 -top-4 rounded-2xl border border-gray-200 bg-white px-4 py-3 shadow-lg dark:border-gray-800 dark:bg-gray-900">
                <div className="flex items-center gap-2 text-sm text-gray-600 dark:text-gray-400">
                  <Shield
                    size={16}
                    className="text-blue-500 dark:text-blue-400"
                  />
                  <span className="font-medium text-gray-900 dark:text-gray-100">
                    End-to-end
                  </span>{" "}
                  encrypted
                </div>
              </div>
            </div>
          </div>
        </div>
      </section>

      <section className="border-t border-gray-200 bg-white py-16 dark:border-gray-800 dark:bg-gray-950">
        <div className="mx-auto max-w-6xl px-6">
          <div className="text-center">
            <h2 className="text-2xl font-semibold text-gray-900 dark:text-gray-100">
              Why NESCloud?
            </h2>
            <p className="mt-2 text-sm text-gray-500 dark:text-gray-400">
              A secure, fast, and easy-to-use cloud storage solution
            </p>
          </div>

          <div className="mt-10 grid gap-6 sm:grid-cols-2 lg:grid-cols-4">
            {features.map((feature) => {
              const Icon = feature.icon;
              return (
                <div
                  key={feature.title}
                  className="rounded-xl border border-gray-200 bg-white p-6 dark:border-gray-800 dark:bg-gray-900"
                >
                  <div
                    className={`inline-flex rounded-lg p-2.5 ${feature.iconBg}`}
                  >
                    <Icon size={20} className={feature.iconColor} />
                  </div>
                  <h3 className="mt-4 text-sm font-semibold text-gray-900 dark:text-gray-100">
                    {feature.title}
                  </h3>
                  <p className="mt-2 text-sm leading-relaxed text-gray-500 dark:text-gray-400">
                    {feature.description}
                  </p>
                </div>
              );
            })}
          </div>
        </div>
      </section>

      <section className="border-t border-gray-200 bg-gray-50 py-12 dark:border-gray-800 dark:bg-gray-950">
        <div className="mx-auto max-w-4xl px-6">
          <div className="grid grid-cols-1 gap-8 sm:grid-cols-3">
            {stats.map((stat) => {
              const Icon = stat.icon;
              return (
                <div key={stat.label} className="text-center">
                  <Icon
                    size={24}
                    className="mx-auto text-blue-600 dark:text-blue-400"
                  />
                  <p className="mt-2 text-3xl font-bold text-gray-900 dark:text-gray-100">
                    {stat.value}
                  </p>
                  <p className="mt-1 text-sm text-gray-500 dark:text-gray-400">
                    {stat.label}
                  </p>
                </div>
              );
            })}
          </div>
        </div>
      </section>

      <section className="border-t border-gray-200 bg-white py-16 dark:border-gray-800 dark:bg-gray-950">
        <div className="mx-auto max-w-3xl px-6 text-center">
          <h2 className="text-2xl font-semibold text-gray-900 dark:text-gray-100">
            Ready to Get Started?
          </h2>
          <p className="mt-2 text-sm text-gray-500 dark:text-gray-400">
            Create a free account and enjoy the ease of storing and sharing
            files in the cloud.
          </p>
          <Link
            href="/register"
            className="mt-8 inline-flex items-center gap-2 rounded-lg bg-blue-600 px-6 py-3 text-sm font-medium text-white hover:bg-blue-700 dark:bg-blue-700 dark:hover:bg-blue-600"
          >
            Create Free Account
            <ArrowRight size={16} />
          </Link>
        </div>
      </section>

      <footer className="border-t border-gray-200 bg-white py-6 dark:border-gray-800 dark:bg-gray-950">
        <div className="mx-auto px-6 text-center text-sm text-gray-400 dark:text-gray-500">
          &copy; 2026 NESCloud. All rights reserved.
        </div>
      </footer>
    </div>
  );
}
