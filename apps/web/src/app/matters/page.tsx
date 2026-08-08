import Link from "next/link";
import { ApiError, getMatters, type Matter } from "@/lib/api";

function statusPillClasses(status: string) {
  switch (status?.toLowerCase()) {
    case "closed":
      return "bg-zinc-100 text-zinc-600 dark:bg-zinc-800 dark:text-zinc-400";
    case "active":
    case "open":
      return "bg-emerald-100 text-emerald-700 dark:bg-emerald-950 dark:text-emerald-300";
    default:
      return "bg-blue-100 text-blue-700 dark:bg-blue-950 dark:text-blue-300";
  }
}

export default async function MattersListPage() {
  let matters: Matter[];
  try {
    matters = (await getMatters(50, 0)) ?? [];
  } catch (err) {
    if (err instanceof ApiError && (err.status === 401 || err.status === 403)) {
      return (
        <div className="mx-auto max-w-4xl px-6 py-12">
          <p className="text-sm text-amber-700 dark:text-amber-300">
            Your session has expired.{" "}
            <Link href="/login?next=/matters" className="font-medium underline">
              Sign in again
            </Link>
            .
          </p>
        </div>
      );
    }
    return (
      <div className="mx-auto max-w-4xl px-6 py-12">
        <p className="text-sm text-red-700 dark:text-red-400">
          Could not load matters. Is the backend API running and reachable?
        </p>
      </div>
    );
  }

  return (
    <div className="mx-auto flex max-w-4xl flex-col gap-6 px-6 py-12">
      <div className="flex items-center justify-between">
        <h1 className="text-2xl font-semibold">Matters</h1>
        <Link href="/clients" className="text-sm text-zinc-500 underline hover:text-zinc-700 dark:hover:text-zinc-300">
          View clients
        </Link>
      </div>

      {matters.length === 0 ? (
        <p className="text-sm text-zinc-500">No matters yet.</p>
      ) : (
        <ul className="divide-y divide-zinc-200 rounded-lg border border-zinc-200 dark:divide-zinc-800 dark:border-zinc-800">
          {matters.map((matter) => (
            <li key={matter.id}>
              <Link
                href={`/matters/${matter.id}`}
                className="flex items-center justify-between gap-4 px-5 py-4 hover:bg-zinc-50 dark:hover:bg-zinc-900"
              >
                <div>
                  <p className="font-medium">{matter.title || "Untitled matter"}</p>
                  <p className="mt-0.5 text-xs text-zinc-500">
                    {matter.case_number || "No case number"} · {matter.case_type || "—"}
                  </p>
                </div>
                <span className={`shrink-0 rounded-full px-2.5 py-1 text-xs font-medium ${statusPillClasses(matter.status)}`}>
                  {matter.status || "unknown"}
                </span>
              </Link>
            </li>
          ))}
        </ul>
      )}
    </div>
  );
}
