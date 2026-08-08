import Link from "next/link";
import { ApiError, getClients, type Client } from "@/lib/api";

export default async function ClientsListPage() {
  let clients: Client[];
  try {
    clients = (await getClients(50, 0)) ?? [];
  } catch (err) {
    if (err instanceof ApiError && (err.status === 401 || err.status === 403)) {
      return (
        <div className="mx-auto max-w-4xl px-6 py-12">
          <p className="text-sm text-amber-700 dark:text-amber-300">
            Your session has expired.{" "}
            <Link href="/login?next=/clients" className="font-medium underline">
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
          Could not load clients. Is the backend API running and reachable?
        </p>
      </div>
    );
  }

  return (
    <div className="mx-auto flex max-w-4xl flex-col gap-6 px-6 py-12">
      <div className="flex items-center justify-between">
        <h1 className="text-2xl font-semibold">Clients</h1>
        <Link href="/matters" className="text-sm text-zinc-500 underline hover:text-zinc-700 dark:hover:text-zinc-300">
          View matters
        </Link>
      </div>

      {clients.length === 0 ? (
        <p className="text-sm text-zinc-500">No clients yet.</p>
      ) : (
        <ul className="divide-y divide-zinc-200 rounded-lg border border-zinc-200 dark:divide-zinc-800 dark:border-zinc-800">
          {clients.map((client) => (
            <li key={client.id} className="flex items-center justify-between gap-4 px-5 py-4">
              <div>
                <p className="font-medium">{client.name}</p>
                <p className="mt-0.5 text-xs text-zinc-500">
                  {client.email || "No email"} · {client.phone || "No phone"}
                </p>
              </div>
              <span className="shrink-0 rounded-full bg-zinc-100 px-2.5 py-1 text-xs font-medium capitalize text-zinc-600 dark:bg-zinc-800 dark:text-zinc-400">
                {client.type}
              </span>
            </li>
          ))}
        </ul>
      )}
    </div>
  );
}
