import { AlertTriangle, Inbox } from "lucide-react";

export function LoadingState({ label }: { label: string }) {
  return (
    <div className="flex flex-col gap-2 rounded-xl border border-border bg-surface-2 p-6">
      <div className="skeleton h-4 w-1/3 rounded" />
      <div className="skeleton h-4 w-2/3 rounded" />
      <div className="skeleton h-4 w-1/2 rounded" />
      <span className="sr-only">{label}</span>
    </div>
  );
}

export function ErrorState({ title, message }: { title: string; message: string }) {
  return (
    <div className="flex items-start gap-3 rounded-xl border border-bad/25 bg-bad-soft p-5 text-sm">
      <AlertTriangle className="mt-0.5 h-4 w-4 shrink-0 text-bad" />
      <div>
        <p className="font-medium text-bad">{title}</p>
        <p className="mt-1 text-bad/80">{message}</p>
      </div>
    </div>
  );
}

export function EmptyState({ message }: { message: string }) {
  return (
    <div className="flex flex-col items-center gap-2 rounded-xl border border-dashed border-border-strong bg-surface-2/50 p-8 text-center text-sm text-ink-faint">
      <Inbox className="h-5 w-5" />
      {message}
    </div>
  );
}
