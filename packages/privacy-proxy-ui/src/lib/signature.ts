import { parseAbiItem, toFunctionSelector } from "viem";

export interface ParsedParam {
  name?: string;
  type: string;
}

export interface ParsedSignature {
  name: string;
  params: ParsedParam[];
  selector: string;
}

export interface SignatureParseError {
  error: string;
}

export type ParseResult = ParsedSignature | SignatureParseError;

export function isParseError(r: ParseResult): r is SignatureParseError {
  return (r as SignatureParseError).error !== undefined;
}

export function parseSignature(input: string): ParseResult {
  const trimmed = input.trim();
  if (!trimmed) return { error: "Signature is empty." };
  const candidate = trimmed.startsWith("function ")
    ? trimmed
    : `function ${trimmed}`;
  try {
    const item = parseAbiItem(candidate);
    if (!item || item.type !== "function") {
      return { error: "Not a function signature." };
    }
    const params: ParsedParam[] = item.inputs.map((p) => ({
      name: p.name || undefined,
      type: p.type,
    }));
    const selector = toFunctionSelector(item);
    return { name: item.name, params, selector };
  } catch (e) {
    return { error: (e as Error).message };
  }
}

export function hasArrayParam(params: ParsedParam[]): boolean {
  return params.some((p) => p.type.includes("["));
}

export function signatureToSelector(input: string): string | null {
  const r = parseSignature(input);
  if (isParseError(r)) return null;
  return r.selector;
}
