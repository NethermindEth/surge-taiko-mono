const KEY = "privacy-proxy-ui:auth";

export interface Session {
  token: string;
  expiresAt: number; // unix seconds
  eoa: `0x${string}`;
}

export function readSession(): Session | null {
  try {
    const raw = localStorage.getItem(KEY);
    if (!raw) return null;
    const parsed = JSON.parse(raw) as Partial<Session>;
    if (
      typeof parsed.token !== "string" ||
      typeof parsed.expiresAt !== "number" ||
      typeof parsed.eoa !== "string"
    ) {
      return null;
    }
    if (parsed.expiresAt <= Math.floor(Date.now() / 1000)) {
      // Expired — drop on read.
      localStorage.removeItem(KEY);
      return null;
    }
    return parsed as Session;
  } catch {
    return null;
  }
}

export function writeSession(session: Session): void {
  localStorage.setItem(KEY, JSON.stringify(session));
}

export function clearSession(): void {
  localStorage.removeItem(KEY);
}
