import { useCallback, useEffect, useRef } from 'react';
import { queryTxStatus } from '../lib/userOp';
import { useTxStatus } from '../context/TxStatusContext';

/// Poll Catalyst's `surge_txStatus` and translate the returned `UserOpStatus`
/// into the 5-phase overlay (`sequencing → proving → proposing → complete`).
///
/// Shared by both swap pipelines:
///   - L1→L2→L1 (useUserOp) polls by `userOpId`
///   - L2→L1→L2 (useInitiateL2Swap) polls by the L2 tx hash — Catalyst writes
///     the same status transitions under the `by_hash` store tree.
export function useTxStatusPolling() {
  const { setTxStatus } = useTxStatus();
  const pollIntervalRef = useRef<ReturnType<typeof setInterval> | null>(null);
  const txHashRef = useRef<string | undefined>(undefined);

  useEffect(() => {
    return () => {
      if (pollIntervalRef.current) clearInterval(pollIntervalRef.current);
    };
  }, []);

  const pollStatus = useCallback(
    (query: { userOpId: number } | { txHash: string }): Promise<boolean> => {
      return new Promise((resolve) => {
        setTxStatus({ phase: 'sequencing' });

        // Seed txHashRef from the query so `complete` includes it even if we
        // never see a Processing status (rare races on near-instant finality).
        if ('txHash' in query) txHashRef.current = query.txHash;

        // Phase ordering: sequencing(0) < proving(1) < proposing(2) < complete(3)
        // proving comes before proposing because the ZK proof is generated before L1 submission
        const phaseOrder: Record<string, number> = {
          sequencing: 0, proving: 1, proposing: 2, complete: 3, rejected: 3,
        };
        let highestPhase = 0;
        let hasSeenProving = false;
        let pollCount = 0;
        const MAX_POLLS = 300; // 5 minutes at 1s intervals

        pollIntervalRef.current = setInterval(async () => {
          pollCount++;
          if (pollCount > MAX_POLLS) {
            if (pollIntervalRef.current) clearInterval(pollIntervalRef.current);
            pollIntervalRef.current = null;
            setTxStatus({ phase: 'rejected', errorMessage: 'Transaction timed out' });
            resolve(false);
            return;
          }

          const status = await queryTxStatus(query);
          if (!status) return;

          if (status.status === 'Pending') {
            if (highestPhase <= phaseOrder.sequencing) {
              setTxStatus({ phase: 'sequencing', txHash: txHashRef.current });
            }
          } else if (status.status === 'ProvingBlock') {
            hasSeenProving = true;
            if (phaseOrder.proving > highestPhase) {
              highestPhase = phaseOrder.proving;
              setTxStatus({ phase: 'proving', txHash: txHashRef.current });
            }
          } else if (status.status === 'Processing') {
            txHashRef.current = status.tx_hash;
            // Only show "proposing" after proving has been seen
            // Before proving, Processing means "sequencing"
            if (hasSeenProving && phaseOrder.proposing > highestPhase) {
              highestPhase = phaseOrder.proposing;
              setTxStatus({ phase: 'proposing', txHash: txHashRef.current });
            } else if (!hasSeenProving && highestPhase <= phaseOrder.sequencing) {
              setTxStatus({ phase: 'sequencing', txHash: txHashRef.current });
            }
          } else if (status.status === 'Executed') {
            if (pollIntervalRef.current) clearInterval(pollIntervalRef.current);
            pollIntervalRef.current = null;
            setTxStatus({ phase: 'complete', txHash: txHashRef.current });
            resolve(true);
          } else if (status.status === 'Rejected') {
            if (pollIntervalRef.current) clearInterval(pollIntervalRef.current);
            pollIntervalRef.current = null;
            setTxStatus({ phase: 'rejected', errorMessage: status.reason });
            resolve(false);
          }
        }, 1000);
      });
    },
    [setTxStatus]
  );

  return { pollStatus, txHashRef };
}
