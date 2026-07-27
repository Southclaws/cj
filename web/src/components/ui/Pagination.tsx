import { ChevronLeft, ChevronRight } from "lucide-react";

import { Button } from "./Button";

interface PaginationProps {
  offset: number;
  pageSize: number;
  total: number;
  onChange: (offset: number) => void;
}

export function Pagination({ offset, pageSize, total, onChange }: PaginationProps) {
  return (
    <div className="flex items-center gap-3 text-sm text-ink-muted">
      <Button variant="secondary" disabled={offset === 0} onClick={() => onChange(Math.max(0, offset - pageSize))}>
        <ChevronLeft className="h-4 w-4" />
      </Button>
      <span>
        {total === 0 ? 0 : offset + 1}-{Math.min(offset + pageSize, total)} of {total}
      </span>
      <Button variant="secondary" disabled={offset + pageSize >= total} onClick={() => onChange(offset + pageSize)}>
        <ChevronRight className="h-4 w-4" />
      </Button>
    </div>
  );
}
