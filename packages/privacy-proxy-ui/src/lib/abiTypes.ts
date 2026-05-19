export type AbiCategory =
  | { kind: "bool" }
  | { kind: "address" }
  | { kind: "uint"; bits: number }
  | { kind: "int"; bits: number }
  | { kind: "bytes"; size: number }
  | { kind: "unknown"; raw: string };

const UINT_RE = /^uint(\d+)?$/;
const INT_RE = /^int(\d+)?$/;
const BYTES_RE = /^bytes(\d+)$/;

export function categorize(type: string | null | undefined): AbiCategory {
  const t = (type ?? "").trim();
  if (!t) return { kind: "unknown", raw: "" };
  if (t === "bool") return { kind: "bool" };
  if (t === "address") return { kind: "address" };
  let m = UINT_RE.exec(t);
  if (m) {
    const bits = m[1] ? Number(m[1]) : 256;
    if (Number.isFinite(bits) && bits >= 8 && bits <= 256 && bits % 8 === 0) {
      return { kind: "uint", bits };
    }
  }
  m = INT_RE.exec(t);
  if (m) {
    const bits = m[1] ? Number(m[1]) : 256;
    if (Number.isFinite(bits) && bits >= 8 && bits <= 256 && bits % 8 === 0) {
      return { kind: "int", bits };
    }
  }
  m = BYTES_RE.exec(t);
  if (m) {
    const size = Number(m[1]);
    if (Number.isFinite(size) && size >= 1 && size <= 32) {
      return { kind: "bytes", size };
    }
  }
  return { kind: "unknown", raw: t };
}

export interface EncodeOk {
  ok: true;
  hex: string;
}
export interface EncodeErr {
  ok: false;
  error: string;
}
export type EncodeResult = EncodeOk | EncodeErr;

function toHex32FromBigInt(value: bigint, signed: boolean, bits: number): EncodeResult {
  if (signed) {
    const lo = -(1n << BigInt(bits - 1));
    const hi = (1n << BigInt(bits - 1)) - 1n;
    if (value < lo || value > hi) {
      return { ok: false, error: `Value out of range for int${bits}` };
    }
    const mask = (1n << 256n) - 1n;
    const u = value < 0n ? (value + (1n << 256n)) & mask : value;
    return { ok: true, hex: `0x${u.toString(16).padStart(64, "0")}` };
  }
  if (value < 0n) {
    return { ok: false, error: "Unsigned value cannot be negative" };
  }
  const max = (1n << BigInt(bits)) - 1n;
  if (value > max) {
    return { ok: false, error: `Value exceeds uint${bits} range` };
  }
  return { ok: true, hex: `0x${value.toString(16).padStart(64, "0")}` };
}

export function encodeLiteral(value: string, type: string | null | undefined): EncodeResult {
  const v = (value ?? "").trim();
  if (!v) return { ok: false, error: "Value is empty" };
  const cat = categorize(type);

  switch (cat.kind) {
    case "bool": {
      if (v === "true") return { ok: true, hex: `0x${"0".repeat(63)}1` };
      if (v === "false") return { ok: true, hex: `0x${"0".repeat(64)}` };
      return { ok: false, error: "Boolean must be true or false" };
    }
    case "address": {
      if (!/^0x[0-9a-fA-F]{40}$/.test(v)) {
        return { ok: false, error: "Address must be 0x + 40 hex chars" };
      }
      return { ok: true, hex: `0x${"0".repeat(24)}${v.slice(2).toLowerCase()}` };
    }
    case "uint":
    case "int": {
      let big: bigint;
      try {
        if (/^-?0x[0-9a-fA-F]+$/.test(v)) {
          big = BigInt(v);
        } else if (/^-?[0-9]+$/.test(v)) {
          big = BigInt(v);
        } else {
          return { ok: false, error: "Enter a decimal or 0x-prefixed integer" };
        }
      } catch {
        return { ok: false, error: "Invalid integer" };
      }
      return toHex32FromBigInt(big, cat.kind === "int", cat.bits);
    }
    case "bytes": {
      const expected = cat.size * 2;
      const re = new RegExp(`^0x[0-9a-fA-F]{${expected}}$`);
      if (!re.test(v)) {
        return {
          ok: false,
          error: `bytes${cat.size} expects 0x + ${expected} hex chars`,
        };
      }
      const padded = `${v.slice(2).toLowerCase()}${"00".repeat(32 - cat.size)}`;
      return { ok: true, hex: `0x${padded}` };
    }
    case "unknown": {
      if (/^0x[0-9a-fA-F]{64}$/.test(v)) {
        return { ok: true, hex: v.toLowerCase() };
      }
      if (/^0x[0-9a-fA-F]{40}$/.test(v)) {
        return { ok: true, hex: `0x${"0".repeat(24)}${v.slice(2).toLowerCase()}` };
      }
      if (/^[0-9]+$/.test(v)) {
        try {
          const hex = BigInt(v).toString(16);
          if (hex.length > 64) return { ok: false, error: "Value too large" };
          return { ok: true, hex: `0x${hex.padStart(64, "0")}` };
        } catch {
          return { ok: false, error: "Invalid integer" };
        }
      }
      return {
        ok: false,
        error: "Enter a 32-byte hex, a 20-byte address, or a decimal integer",
      };
    }
  }
}

export function decodeLiteral(hex: string | null | undefined, type: string | null | undefined): string {
  const h = (hex ?? "").trim();
  if (!h) return "—";
  if (!/^0x[0-9a-fA-F]{64}$/.test(h)) return h;
  const cat = categorize(type);
  const bare = h.slice(2).toLowerCase();
  switch (cat.kind) {
    case "bool": {
      if (/^0+$/.test(bare)) return "false";
      if (/^0{63}1$/.test(bare)) return "true";
      return h;
    }
    case "address": {
      if (!/^0{24}/.test(bare)) return h;
      return `0x${bare.slice(24)}`;
    }
    case "uint": {
      try {
        return BigInt(h).toString(10);
      } catch {
        return h;
      }
    }
    case "int": {
      try {
        const u = BigInt(h);
        const top = 1n << BigInt(cat.bits - 1);
        const mod = 1n << BigInt(cat.bits);
        const lo = u & (mod - 1n);
        const signed = lo >= top ? lo - mod : lo;
        return signed.toString(10);
      } catch {
        return h;
      }
    }
    case "bytes": {
      if (cat.size === 32) return h;
      const head = bare.slice(0, cat.size * 2);
      const tail = bare.slice(cat.size * 2);
      if (!/^0+$/.test(tail)) return h;
      return `0x${head}`;
    }
    case "unknown":
      return h;
  }
}

export function typeLabel(type: string | null | undefined): string {
  const cat = categorize(type);
  switch (cat.kind) {
    case "bool":
      return "bool";
    case "address":
      return "address";
    case "uint":
      return `uint${cat.bits}`;
    case "int":
      return `int${cat.bits}`;
    case "bytes":
      return `bytes${cat.size}`;
    case "unknown":
      return cat.raw || "—";
  }
}
