export function roleColor(color: number): string {
  if (color === 0) {
    return "#99a1af";
  }
  return `#${color.toString(16).padStart(6, "0")}`;
}
