import { useCallback, useState } from "react";
import { useAccount, useSignMessage } from "wagmi";
import toast from "react-hot-toast";
import { request, AdminApiError } from "../lib/apiClient";
import { useAuth } from "../context/AuthContext";
import type {
  ChallengeResponse,
  VerifyResponse,
} from "../types/api";

interface SignInState {
  isLoading: boolean;
  error: string | null;
}

export function useSignIn(): SignInState & { signIn: () => Promise<void> } {
  const { address } = useAccount();
  const { signMessageAsync } = useSignMessage();
  const { setSession } = useAuth();
  const [state, setState] = useState<SignInState>({
    isLoading: false,
    error: null,
  });

  const signIn = useCallback(async () => {
    if (!address) {
      toast.error("Connect a wallet first.");
      return;
    }
    setState({ isLoading: true, error: null });
    try {
      const challenge = await request<ChallengeResponse>(
        `/auth/challenge?address=${address}`,
        { anonymous: true },
      );
      const signature = await signMessageAsync({ message: challenge.message });
      const verified = await request<VerifyResponse>("/auth/verify", {
        method: "POST",
        anonymous: true,
        body: { address, signature },
      });
      setSession({
        token: verified.token,
        expiresAt: verified.expires_at,
        eoa: address,
      });
      toast.success("Signed in");
      setState({ isLoading: false, error: null });
    } catch (err) {
      const message =
        err instanceof AdminApiError
          ? err.message
          : err instanceof Error
            ? err.message
            : "sign-in failed";
      setState({ isLoading: false, error: message });
      toast.error(message);
    }
  }, [address, signMessageAsync, setSession]);

  return { ...state, signIn };
}
