import { getAuthToken } from "@/lib/auth";

// Server-only on purpose: every call here runs in a Server Component or
// Server Action (this file is never imported from a "use client" component),
// so it must be read at request time via a plain env var, not a NEXT_PUBLIC_
// one — those get inlined into the build at `npm run build` time, which
// means setting them as a container *runtime* env var in docker-compose.yml
// has no effect at all.
const API_BASE_URL = process.env.API_BASE_URL ?? "http://localhost:8080/api";

export class ApiError extends Error {
  status: number;
  constructor(status: number, message: string) {
    super(message);
    this.status = status;
  }
}

export interface Matter {
  id: string;
  title: string;
  description: string;
  client_id: string;
  court_id: string;
  judge_id: string;
  advocate_id: string;
  case_number: string;
  case_type: string;
  stage: string;
  status: string;
  priority: string;
  limitation_date: string;
  created_at: string;
  updated_at: string;
}

export interface Client {
  id: string;
  name: string;
  type: string;
  email: string;
  phone: string;
  address: string | null;
  pan: string | null;
  gstin: string | null;
  notes: string | null;
}

export interface Hearing {
  id: string;
  matter_id: string;
  date: string;
  notes: string;
  outcome: string;
  next_hearing_date: string | null;
}

export interface Order {
  id: string;
  matter_id: string;
  hearing_id: string | null;
  title: string;
  description: string;
  date: string;
  document_id: string | null;
}

export interface Document {
  id: string;
  matter_id: string;
  name: string;
  mime_type: string;
  size: number;
  ocr_status: string;
  created_at: string;
}

export interface TimelineEvent {
  type: "hearing" | "order";
  date: string;
  data: Hearing | Order;
}

export interface LoginResponse {
  token: string;
  refresh_token: string;
}

// login is unauthenticated by definition, so it bypasses apiFetch's token attachment.
export async function login(email: string, password: string): Promise<LoginResponse> {
  const res = await fetch(`${API_BASE_URL}/auth/login`, {
    method: "POST",
    headers: { "Content-Type": "application/json" },
    body: JSON.stringify({ email, password }),
    cache: "no-store",
  });

  if (!res.ok) {
    throw new ApiError(res.status, "Invalid email or password");
  }

  return res.json() as Promise<LoginResponse>;
}

// apiFetch is the authenticated fetch wrapper: it reads the caller's session
// token from the request cookies (server-side only — see lib/auth.ts) and
// attaches it as a Bearer token. A missing/expired token surfaces as a
// distinguishable ApiError(401) so pages can prompt re-authentication instead
// of showing a generic "failed to load" message.
async function apiFetch<T>(path: string, init: RequestInit = {}): Promise<T> {
  const token = await getAuthToken();

  const res = await fetch(`${API_BASE_URL}${path}`, {
    ...init,
    headers: {
      "Content-Type": "application/json",
      ...(token ? { Authorization: `Bearer ${token}` } : {}),
      ...init.headers,
    },
    cache: "no-store",
  });

  if (!res.ok) {
    throw new ApiError(res.status, `Request to ${path} failed with status ${res.status}`);
  }

  return res.json() as Promise<T>;
}

export function getMatters(limit = 20, offset = 0): Promise<Matter[]> {
  return apiFetch<Matter[]>(`/matters?limit=${limit}&offset=${offset}`);
}

export function getMatter(id: string): Promise<Matter> {
  return apiFetch<Matter>(`/matters/${id}`);
}

export function getClients(limit = 20, offset = 0): Promise<Client[]> {
  return apiFetch<Client[]>(`/clients?limit=${limit}&offset=${offset}`);
}

export function getMatterTimeline(matterID: string): Promise<TimelineEvent[]> {
  return apiFetch<TimelineEvent[]>(`/matters/${matterID}/timeline`);
}

export function getHearingsByMatter(matterID: string): Promise<Hearing[]> {
  return apiFetch<Hearing[]>(`/hearings/matter/${matterID}`);
}

export function getOrdersByMatter(matterID: string): Promise<Order[]> {
  return apiFetch<Order[]>(`/orders/matter/${matterID}`);
}

export function getDocumentsByMatter(matterID: string): Promise<Document[]> {
  return apiFetch<Document[]>(`/documents/matter/${matterID}?limit=50&offset=0`);
}
