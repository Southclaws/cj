import { useState } from "react";

export type SortDir = "asc" | "desc";

export function useSort<T extends string>(initialKey: T, initialDir: SortDir = "asc") {
  const [key, setKey] = useState<T>(initialKey);
  const [dir, setDir] = useState<SortDir>(initialDir);

  function toggle(nextKey: T) {
    if (nextKey === key) {
      setDir((d) => (d === "asc" ? "desc" : "asc"));
    } else {
      setKey(nextKey);
      setDir("asc");
    }
  }

  return { key, dir, toggle };
}

export function sortRows<T>(rows: T[], compare: (a: T, b: T) => number, dir: SortDir): T[] {
  const sorted = [...rows].sort(compare);
  return dir === "asc" ? sorted : sorted.reverse();
}
