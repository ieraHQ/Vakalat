import Link from "next/link";
import { redirect } from "next/navigation";
import { getAuthToken } from "@/lib/auth";
import { logoutAction } from "@/app/login/actions";

export default async function Home() {
  const token = await getAuthToken();
  if (!token) {
    redirect("/login");
  }

  return (
    <div className="flex flex-1 flex-col items-center justify-center gap-6 bg-zinc-50 px-6 text-center dark:bg-black">
      <div>
        <h1 className="text-2xl font-semibold">You&apos;re signed in</h1>
        <p className="mt-2 max-w-md text-sm text-zinc-500">
          Browse your matters and clients, or sign out below.
        </p>
      </div>
      <div className="flex gap-3">
        <Link
          href="/matters"
          className="rounded-md bg-zinc-900 px-4 py-2 text-sm font-medium text-white hover:bg-zinc-700 dark:bg-zinc-50 dark:text-zinc-900 dark:hover:bg-zinc-200"
        >
          View matters
        </Link>
        <Link
          href="/clients"
          className="rounded-md border border-zinc-300 px-4 py-2 text-sm font-medium hover:bg-zinc-100 dark:border-zinc-700 dark:hover:bg-zinc-900"
        >
          View clients
        </Link>
      </div>
      <form action={logoutAction}>
        <button
          type="submit"
          className="text-sm text-zinc-500 underline hover:text-zinc-700 dark:hover:text-zinc-300"
        >
          Sign out
        </button>
      </form>
    </div>
  );
}
