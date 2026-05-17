import { useEffect, useState } from "react";

interface ConfirmDialogProps {
  open: boolean;
  title: string;
  description: string;
  /** When provided, the destructive button is disabled until the user types this exact string. */
  confirmString?: string;
  confirmLabel?: string;
  destructive?: boolean;
  onConfirm: () => void;
  onCancel: () => void;
  isLoading?: boolean;
}

export function ConfirmDialog({
  open,
  title,
  description,
  confirmString,
  confirmLabel = "Confirm",
  destructive = false,
  onConfirm,
  onCancel,
  isLoading,
}: ConfirmDialogProps) {
  const [typed, setTyped] = useState("");

  useEffect(() => {
    if (!open) setTyped("");
  }, [open]);

  useEffect(() => {
    if (!open) return;
    const onKey = (e: KeyboardEvent) => {
      if (e.key === "Escape") onCancel();
    };
    document.addEventListener("keydown", onKey);
    return () => document.removeEventListener("keydown", onKey);
  }, [open, onCancel]);

  if (!open) return null;

  const isMatch = !confirmString || typed === confirmString;

  return (
    <div className="fixed inset-0 z-[60] flex items-center justify-center px-4">
      <div className="absolute inset-0 bg-surge-primary/30 backdrop-blur-sm" onClick={onCancel} aria-hidden />
      <div
        role="alertdialog"
        aria-modal="true"
        aria-labelledby="confirm-title"
        className="relative w-full max-w-md rounded-2xl bg-surge-card p-6 shadow-2xl"
      >
        <h3 id="confirm-title" className="text-base font-semibold text-surge-text">
          {title}
        </h3>
        <p className="mt-2 text-sm text-surge-muted">{description}</p>

        {confirmString ? (
          <div className="mt-4">
            <label className="text-xs font-medium text-surge-muted">
              Type{" "}
              <span className="font-mono text-surge-text">{confirmString}</span>{" "}
              to confirm
            </label>
            <input
              type="text"
              value={typed}
              onChange={(e) => setTyped(e.target.value)}
              autoFocus
              className="mt-1 w-full rounded-lg border border-surge-border bg-surge-card-hover/40 px-3 py-2 font-mono text-sm text-surge-text outline-none focus:border-surge-secondary"
            />
          </div>
        ) : null}

        <div className="mt-6 flex items-center justify-end gap-2">
          <button
            type="button"
            onClick={onCancel}
            className="rounded-lg px-3 py-2 text-sm font-medium text-surge-muted hover:bg-surge-card-hover hover:text-surge-text"
          >
            Cancel
          </button>
          <button
            type="button"
            onClick={onConfirm}
            disabled={!isMatch || isLoading}
            className={`rounded-lg px-4 py-2 text-sm font-semibold text-white shadow-sm transition disabled:cursor-not-allowed disabled:opacity-50 ${
              destructive
                ? "bg-red-600 hover:bg-red-700"
                : "bg-surge-primary hover:bg-surge-primary/90"
            }`}
          >
            {isLoading ? "Working…" : confirmLabel}
          </button>
        </div>
      </div>
    </div>
  );
}
