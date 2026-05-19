/** Mirror of the currently-active session. `apiClient` reads this
 *  directly to attach the bearer header — keeping it as a single key
 *  avoids threading wallet context through every fetch. */
const ACTIVE_KEY = "privacy-proxy-ui:auth";
/** Per-EOA cache of every session the user has signed in with on this
 *  proxy origin. Restoring from this is what makes account switching
 *  feel sticky — flipping the wallet between two admins shouldn't
 *  cost a sign-in each way. */
const CACHE_KEY = "privacy-proxy-ui:auth:cache";

export interface Session {
  token: string;
  expiresAt: number; // unix seconds
  eoa: `0x${string}`;
}

type SessionCache = Record<string, Session>;

const nowSecs = () => Math.floor(Date.now() / 1000);

function isValidSession(v: unknown): v is Session {
  if (!v || typeof v !== "object") return false;
  const s = v as Partial<Session>;
  return (
    typeof s.token === "string" &&
    typeof s.expiresAt === "number" &&
    typeof s.eoa === "string"
  );
}

export function readSession(): Session | null {
  try {
    const raw = localStorage.getItem(ACTIVE_KEY);
    if (!raw) return null;
    const parsed = JSON.parse(raw);
    if (!isValidSession(parsed)) return null;
    if (parsed.expiresAt <= nowSecs()) {
      localStorage.removeItem(ACTIVE_KEY);
      return null;
    }
    return parsed;
  } catch {
    return null;
  }
}

function readCache(): SessionCache {
  try {
    const raw = localStorage.getItem(CACHE_KEY);
    if (!raw) return {};
    const parsed = JSON.parse(raw);
    if (!parsed || typeof parsed !== "object") return {};
    const cache: SessionCache = {};
    const now = nowSecs();
    for (const [k, v] of Object.entries(parsed)) {
      if (isValidSession(v) && v.expiresAt > now) cache[k.toLowerCase()] = v;
    }
    return cache;
  } catch {
    return {};
  }
}

function writeCache(cache: SessionCache): void {
  localStorage.setItem(CACHE_KEY, JSON.stringify(cache));
}

/** Returns a non-expired session for the given EOA from the cache, or
 *  null. Lowercases on lookup so address casing doesn't bite. */
export function getSessionFor(eoa: string): Session | null {
  const cache = readCache();
  const s = cache[eoa.toLowerCase()];
  return s && s.expiresAt > nowSecs() ? s : null;
}

/** Persist a session into BOTH the per-EOA cache and the active mirror
 *  (so apiClient picks it up on the very next request). */
export function writeSession(session: Session): void {
  const cache = readCache();
  cache[session.eoa.toLowerCase()] = session;
  writeCache(cache);
  localStorage.setItem(ACTIVE_KEY, JSON.stringify(session));
}

/** Mark a session already in the cache as the active one (the path
 *  the AuthContext takes when restoring on account switch). */
export function setActiveSession(session: Session | null): void {
  if (session) localStorage.setItem(ACTIVE_KEY, JSON.stringify(session));
  else localStorage.removeItem(ACTIVE_KEY);
}

/** Remove the active-slot mirror without touching the cache. Used when
 *  the wallet switches to an EOA we have no token for — the user may
 *  switch back, and we want that session waiting. */
export function clearActiveSession(): void {
  localStorage.removeItem(ACTIVE_KEY);
}

/** Drop a specific EOA's session (cache + active mirror if it was the
 *  active one). Used by explicit Sign-Out. */
export function clearSessionFor(eoa: string): void {
  const cache = readCache();
  const key = eoa.toLowerCase();
  if (cache[key]) {
    delete cache[key];
    writeCache(cache);
  }
  const active = readSession();
  if (active && active.eoa.toLowerCase() === key) {
    localStorage.removeItem(ACTIVE_KEY);
  }
}

/** Nuke every session. Kept for the 401-handler safety net. */
export function clearSession(): void {
  localStorage.removeItem(ACTIVE_KEY);
  localStorage.removeItem(CACHE_KEY);
}
