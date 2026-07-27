import { useState, type MouseEvent } from "react";
import { Check, Copy } from "lucide-react";

import { copyText } from "../../lib/clipboard";

interface CopyableIdProps {
  value: string;
  className?: string;
}

export function CopyableId({ value, className = "" }: CopyableIdProps) {
  const [copied, setCopied] = useState(false);

  async function handleCopy(e: MouseEvent) {
    e.stopPropagation();
    e.preventDefault();
    const ok = await copyText(value);
    if (ok) {
      setCopied(true);
      setTimeout(() => setCopied(false), 1200);
    }
  }

  return (
    <button
      type="button"
      onClick={handleCopy}
      title="Copy ID"
      className={`inline-flex items-center gap-1 font-mono text-ink-faint transition-colors hover:text-ink ${className}`}
    >
      {value}
      {copied ? <Check className="h-3 w-3 shrink-0 text-good" /> : <Copy className="h-3 w-3 shrink-0" />}
    </button>
  );
}
