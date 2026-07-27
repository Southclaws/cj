import { CheckCircle2, Info, X, XCircle } from "lucide-react";

import { dismissToast, useToasts, type ToastTone } from "../../lib/toast";

const toneClasses: Record<ToastTone, string> = {
  good: "border-good/30 bg-good-soft text-ink",
  bad: "border-bad/30 bg-bad-soft text-ink",
  neutral: "border-border-strong bg-surface-2 text-ink",
};

const toneIcon: Record<ToastTone, typeof Info> = {
  good: CheckCircle2,
  bad: XCircle,
  neutral: Info,
};

const iconToneClasses: Record<ToastTone, string> = {
  good: "text-good",
  bad: "text-bad",
  neutral: "text-ink-faint",
};

export function Toaster() {
  const toasts = useToasts();

  if (toasts.length === 0) return null;

  return (
    <div className="pointer-events-none fixed right-4 top-4 z-[60] flex w-full max-w-sm flex-col gap-2">
      {toasts.map((toast) => {
        const Icon = toneIcon[toast.tone];
        return (
          <div
            key={toast.id}
            className={`animate-toast-in pointer-events-auto flex items-start gap-2.5 rounded-lg border px-3.5 py-3 text-sm shadow-xl backdrop-blur-sm ${toneClasses[toast.tone]}`}
          >
            <Icon className={`mt-0.5 h-4 w-4 shrink-0 ${iconToneClasses[toast.tone]}`} />
            <p className="flex-1">{toast.message}</p>
            <button
              type="button"
              onClick={() => dismissToast(toast.id)}
              className="shrink-0 text-ink-faint transition-colors hover:text-ink"
            >
              <X className="h-3.5 w-3.5" />
            </button>
          </div>
        );
      })}
    </div>
  );
}
