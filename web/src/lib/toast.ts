import { useEffect, useState } from "react";

export type ToastTone = "good" | "bad" | "neutral";

export interface Toast {
  id: number;
  tone: ToastTone;
  message: string;
}

let nextId = 1;
let toasts: Toast[] = [];
const listeners = new Set<(toasts: Toast[]) => void>();

function emit() {
  for (const listener of listeners) listener(toasts);
}

export function pushToast(message: string, tone: ToastTone = "neutral", durationMs = 4000) {
  const id = nextId++;
  toasts = [...toasts, { id, tone, message }];
  emit();
  setTimeout(() => dismissToast(id), durationMs);
  return id;
}

export function dismissToast(id: number) {
  toasts = toasts.filter((t) => t.id !== id);
  emit();
}

export function useToasts(): Toast[] {
  const [state, setState] = useState(toasts);
  useEffect(() => {
    listeners.add(setState);
    return () => {
      listeners.delete(setState);
    };
  }, []);
  return state;
}
