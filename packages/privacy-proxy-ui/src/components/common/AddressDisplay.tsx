import { useState } from "react";
import { shortenAddress } from "../../lib/format";

interface AddressDisplayProps {
  value: string;
  className?: string;
  fullOnHover?: boolean;
}

export function AddressDisplay({
  value,
  className = "",
  fullOnHover = true,
}: AddressDisplayProps) {
  const [copied, setCopied] = useState(false);

  const copy = async (e: React.MouseEvent) => {
    e.stopPropagation();
    try {
      await navigator.clipboard.writeText(value);
      setCopied(true);
      setTimeout(() => setCopied(false), 1500);
    } catch {
      // ignore
    }
  };

  return (
    <button
      type="button"
      onClick={copy}
      title={fullOnHover ? value : undefined}
      className={`inline-flex items-center gap-1.5 font-mono text-sm text-surge-text hover:text-surge-secondary ${className}`}
    >
      <span>{shortenAddress(value)}</span>
      <span className="text-surge-muted">
        {copied ? (
          <svg
            width="12"
            height="12"
            viewBox="0 0 16 16"
            fill="none"
            stroke="currentColor"
            strokeWidth="2"
            strokeLinecap="round"
            strokeLinejoin="round"
            aria-hidden
          >
            <path d="M3 8l3 3 7-7" />
          </svg>
        ) : (
          <svg
            width="12"
            height="12"
            viewBox="0 0 16 16"
            fill="none"
            stroke="currentColor"
            strokeWidth="1.5"
            strokeLinecap="round"
            strokeLinejoin="round"
            aria-hidden
          >
            <rect x="5" y="5" width="9" height="9" rx="1.5" />
            <path d="M3 11V3a1 1 0 0 1 1-1h8" />
          </svg>
        )}
      </span>
    </button>
  );
}
