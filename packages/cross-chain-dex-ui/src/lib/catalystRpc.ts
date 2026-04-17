import { Hex } from 'viem';

/// Shape of the `message` field returned by `surge_simulateReturnMessage`.
/// All numeric fields arrive as decimal strings from the Rust side.
export interface SimulatedMessage {
  id: string;
  fee: string;
  gasLimit: string;
  from: Hex;
  srcChainId: string;
  srcOwner: Hex;
  destChainId: string;
  destOwner: Hex;
  to: Hex;
  value: string;
  data: Hex;
}

export interface SimulateReturnMessageResponse {
  message: SimulatedMessage;
  signalSlot: Hex;
}

/// Narrow the strings to bigints in the `Message` tuple the vault expects.
export interface ResolvedMessage {
  id: bigint;
  fee: bigint;
  gasLimit: number;
  from: Hex;
  srcChainId: bigint;
  srcOwner: Hex;
  destChainId: bigint;
  destOwner: Hex;
  to: Hex;
  value: bigint;
  data: Hex;
}

/// Call Catalyst's `surge_simulateReturnMessage` JSON-RPC. Given a placeholder L2 tx
/// that contains a zero-filled return message, Catalyst traces the bridge call, simulates
/// the L1 callback, and returns the real return message (including the L1-computed
/// token amounts) that the caller should splice back into the final L2 tx.
export async function simulateReturnMessage(
  from: Hex,
  to: Hex,
  data: Hex,
  value?: bigint
): Promise<SimulateReturnMessageResponse> {
  // Use the Vite dev proxy (/api/builder → Catalyst RPC) to avoid CORS issues.
  // In production, VITE_CATALYST_RPC_URL would be used directly.
  const url = import.meta.env.DEV ? '/api/builder' : (import.meta.env.VITE_CATALYST_RPC_URL as string);
  if (!url) {
    throw new Error('VITE_CATALYST_RPC_URL is not configured');
  }

  const body = JSON.stringify({
    jsonrpc: '2.0',
    method: 'surge_simulateReturnMessage',
    params: [{ from, to, data, ...(value != null ? { value: value.toString() } : {}) }],
    id: 1,
  });

  const resp = await fetch(url, {
    method: 'POST',
    headers: { 'Content-Type': 'application/json' },
    body,
  });

  if (!resp.ok) {
    throw new Error(`Catalyst RPC HTTP ${resp.status}`);
  }

  const json = await resp.json();
  if (json.error) {
    throw new Error(
      `surge_simulateReturnMessage: ${json.error.message || JSON.stringify(json.error)}`
    );
  }
  return json.result as SimulateReturnMessageResponse;
}

export function resolveSimulatedMessage(m: SimulatedMessage): ResolvedMessage {
  return {
    id: BigInt(m.id),
    fee: BigInt(m.fee),
    gasLimit: Number(m.gasLimit),
    from: m.from,
    srcChainId: BigInt(m.srcChainId),
    srcOwner: m.srcOwner,
    destChainId: BigInt(m.destChainId),
    destOwner: m.destOwner,
    to: m.to,
    value: BigInt(m.value),
    data: m.data,
  };
}
