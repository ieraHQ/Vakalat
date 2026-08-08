import { Suspense } from "react";
import Link from "next/link";
import {
  ApiError,
  getDocumentsByMatter,
  getMatter,
  getMatterTimeline,
  type Document,
  type Matter,
  type TimelineEvent,
} from "@/lib/api";
import { logoutAction } from "@/app/login/actions";

function SectionSkeleton({ label }: { label: string }) {
  return (
    <div className="animate-pulse rounded-lg border border-zinc-200 p-6 text-sm text-zinc-400 dark:border-zinc-800">
      Loading {label}...
    </div>
  );
}

function ErrorBox({ error, matterId }: { error: unknown; matterId: string }) {
  if (error instanceof ApiError && (error.status === 401 || error.status === 403)) {
    return (
      <div className="rounded-lg border border-amber-200 bg-amber-50 p-6 text-sm text-amber-800 dark:border-amber-900 dark:bg-amber-950 dark:text-amber-300">
        Your session has expired or you don&apos;t have access to this matter.{" "}
        <Link href={`/login?next=/matters/${matterId}`} className="font-medium underline">
          Sign in again
        </Link>
        .
      </div>
    );
  }

  return (
    <div className="rounded-lg border border-red-200 bg-red-50 p-6 text-sm text-red-700 dark:border-red-900 dark:bg-red-950 dark:text-red-300">
      Could not load this data. Is the backend API running and reachable?
    </div>
  );
}

async function MatterOverview({ matterId }: { matterId: string }) {
  let matter: Matter;
  try {
    matter = await getMatter(matterId);
  } catch (err) {
    return <ErrorBox error={err} matterId={matterId} />;
  }

  return (
    <div className="rounded-lg border border-zinc-200 p-6 dark:border-zinc-800">
      <h1 className="text-2xl font-semibold">{matter.title}</h1>
      <p className="mt-1 text-sm text-zinc-500">{matter.case_number || "No case number"}</p>
      <p className="mt-4 text-zinc-700 dark:text-zinc-300">{matter.description}</p>
      <dl className="mt-4 grid grid-cols-2 gap-4 text-sm sm:grid-cols-4">
        <div>
          <dt className="text-zinc-500">Stage</dt>
          <dd className="font-medium">{matter.stage || "—"}</dd>
        </div>
        <div>
          <dt className="text-zinc-500">Status</dt>
          <dd className="font-medium">{matter.status || "—"}</dd>
        </div>
        <div>
          <dt className="text-zinc-500">Priority</dt>
          <dd className="font-medium capitalize">{matter.priority || "—"}</dd>
        </div>
        <div>
          <dt className="text-zinc-500">Limitation Date</dt>
          <dd className="font-medium">{matter.limitation_date || "—"}</dd>
        </div>
      </dl>
    </div>
  );
}

async function MatterTimeline({ matterId }: { matterId: string }) {
  let events: TimelineEvent[];
  try {
    events = (await getMatterTimeline(matterId)) ?? [];
  } catch (err) {
    return <ErrorBox error={err} matterId={matterId} />;
  }

  if (events.length === 0) {
    return <p className="text-sm text-zinc-500">No hearings or orders recorded yet.</p>;
  }

  return (
    <ol className="space-y-3">
      {events.map((event, idx) => {
        const isHearing = event.type === "hearing";
        const title = isHearing
          ? (event.data as { notes?: string }).notes || "Hearing"
          : (event.data as { title?: string }).title || "Order";
        return (
          <li
            key={idx}
            className="flex gap-4 rounded-lg border border-zinc-200 p-4 dark:border-zinc-800"
          >
            <span
              className={`mt-0.5 inline-block h-2 w-2 shrink-0 rounded-full ${
                isHearing ? "bg-blue-500" : "bg-amber-500"
              }`}
            />
            <div>
              <p className="text-xs uppercase tracking-wide text-zinc-500">
                {event.type} · {event.date}
              </p>
              <p className="mt-1 text-sm">{title}</p>
            </div>
          </li>
        );
      })}
    </ol>
  );
}

async function MatterDocuments({ matterId }: { matterId: string }) {
  let documents: Document[];
  try {
    documents = (await getDocumentsByMatter(matterId)) ?? [];
  } catch (err) {
    return <ErrorBox error={err} matterId={matterId} />;
  }

  if (documents.length === 0) {
    return <p className="text-sm text-zinc-500">No documents uploaded yet.</p>;
  }

  return (
    <ul className="divide-y divide-zinc-200 dark:divide-zinc-800">
      {documents.map((doc) => (
        <li key={doc.id} className="flex items-center justify-between py-3 text-sm">
          <span>{doc.name}</span>
          <span className="text-xs text-zinc-500">{doc.ocr_status}</span>
        </li>
      ))}
    </ul>
  );
}

export default async function MatterPage({
  params,
}: {
  params: Promise<{ id: string }>;
}) {
  const { id } = await params;

  return (
    <div className="mx-auto flex max-w-4xl flex-col gap-8 px-6 py-12">
      <form action={logoutAction} className="flex justify-end">
        <button type="submit" className="text-sm text-zinc-500 underline hover:text-zinc-700 dark:hover:text-zinc-300">
          Sign out
        </button>
      </form>

      <Suspense fallback={<SectionSkeleton label="matter" />}>
        <MatterOverview matterId={id} />
      </Suspense>

      <section>
        <h2 className="mb-3 text-lg font-semibold">Timeline</h2>
        <Suspense fallback={<SectionSkeleton label="timeline" />}>
          <MatterTimeline matterId={id} />
        </Suspense>
      </section>

      <section>
        <h2 className="mb-3 text-lg font-semibold">Documents</h2>
        <Suspense fallback={<SectionSkeleton label="documents" />}>
          <MatterDocuments matterId={id} />
        </Suspense>
      </section>
    </div>
  );
}
