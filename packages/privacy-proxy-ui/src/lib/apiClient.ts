import type { ApiErrorBody } from "../types/api";
import { readSession } from "./tokenStorage";

/** Same-origin base. Vite proxies /api/* through to the privacy-proxy. */
const BASE = "/api";

export class AdminApiError extends Error {
  constructor(
    public readonly status: number,
    public readonly code: string,
    message: string,
  ) {
    super(message);
    this.name = "AdminApiError";
  }
}

let onUnauthorized: (() => void) | null = null;

/** Wire the auth context up once at boot so 401s drop the session. */
export function setUnauthorizedHandler(handler: () => void): void {
  onUnauthorized = handler;
}

interface RequestOptions {
  method?: "GET" | "POST" | "PUT" | "DELETE";
  body?: unknown;
  /** Skip the bearer header (e.g. /auth/* endpoints). */
  anonymous?: boolean;
  /** AbortSignal for react-query cancellation. */
  signal?: AbortSignal;
}

export async function request<T>(
  path: string,
  opts: RequestOptions = {},
): Promise<T> {
  const headers: Record<string, string> = {
    Accept: "application/json",
  };
  if (opts.body !== undefined) {
    headers["Content-Type"] = "application/json";
  }
  if (!opts.anonymous) {
    // Read straight from localStorage so the very first request after a page
    // refresh has the bearer attached — by the time AuthProvider's effects
    // run, react-query has already dispatched its initial queries.
    const token = readSession()?.token ?? null;
    if (token) headers.Authorization = `Bearer ${token}`;
  }

  const res = await fetch(`${BASE}${path}`, {
    method: opts.method ?? "GET",
    headers,
    body: opts.body !== undefined ? JSON.stringify(opts.body) : undefined,
    signal: opts.signal,
  });

  if (res.status === 401 && !opts.anonymous) {
    onUnauthorized?.();
  }

  if (res.status === 204) {
    return undefined as T;
  }

  const text = await res.text();
  const data = text ? (JSON.parse(text) as unknown) : null;

  if (!res.ok) {
    const body = data as ApiErrorBody | null;
    const code = body?.error?.code ?? "unknown";
    const message = body?.error?.message ?? res.statusText ?? "request failed";
    throw new AdminApiError(res.status, code, message);
  }

  return data as T;
}
