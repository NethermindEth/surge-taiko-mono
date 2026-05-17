/** Truncate an EOA / contract address to `0x1234…abcd`. */
export function shortenAddress(addr: string, head = 6, tail = 4): string {
  if (!addr) return "";
  if (addr.length <= head + tail + 2) return addr;
  return `${addr.slice(0, head)}…${addr.slice(-tail)}`;
}

const SECOND = 1;
const MINUTE = 60 * SECOND;
const HOUR = 60 * MINUTE;
const DAY = 24 * HOUR;

/** Human duration like "6d 4h" or "12m 33s" from now → unix seconds. */
export function durationUntil(unix: number): string {
  const diff = unix - Math.floor(Date.now() / 1000);
  if (diff <= 0) return "expired";
  const parts: string[] = [];
  let rem = diff;
  if (rem >= DAY) {
    parts.push(`${Math.floor(rem / DAY)}d`);
    rem %= DAY;
  }
  if (rem >= HOUR) {
    parts.push(`${Math.floor(rem / HOUR)}h`);
    rem %= HOUR;
  }
  if (parts.length === 0 && rem >= MINUTE) {
    parts.push(`${Math.floor(rem / MINUTE)}m`);
    rem %= MINUTE;
  }
  if (parts.length === 0) parts.push(`${rem}s`);
  return parts.slice(0, 2).join(" ");
}

/** Relative timestamp like "2 days ago" or "12 minutes ago". */
export function timeAgo(unix: number): string {
  const diff = Math.floor(Date.now() / 1000) - unix;
  if (diff < 60) return `${diff}s ago`;
  if (diff < HOUR) return `${Math.floor(diff / MINUTE)}m ago`;
  if (diff < DAY) return `${Math.floor(diff / HOUR)}h ago`;
  if (diff < 30 * DAY) return `${Math.floor(diff / DAY)}d ago`;
  const d = new Date(unix * 1000);
  return d.toLocaleDateString();
}

/** Validate a 20-byte hex address; tolerant of casing, requires 0x. */
export function isAddress(value: string): boolean {
  return /^0x[0-9a-fA-F]{40}$/.test(value);
}

/** Validate a 4-byte selector hex; tolerant of casing, requires 0x. */
export function isSelector(value: string): boolean {
  return /^0x[0-9a-fA-F]{8}$/.test(value);
}

export function normalizeAddress(value: string): string {
  return value.toLowerCase();
}
