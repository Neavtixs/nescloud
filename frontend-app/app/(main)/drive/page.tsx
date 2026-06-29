import Link from "next/link";

export default function DrivePage() {
  return (
    <div>
      <h1>my drive</h1>
      <Link href="/home" replace>
        back
      </Link>
    </div>
  );
}
