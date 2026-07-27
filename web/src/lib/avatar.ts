const avatarPalette = ["#6d5ef8", "#3ba55d", "#f0b232", "#ed4245", "#2b90d9", "#d4589e"];

export function avatarColor(id: string): string {
  let hash = 0;
  for (let i = 0; i < id.length; i++) hash = (hash * 31 + id.charCodeAt(i)) >>> 0;
  return avatarPalette[hash % avatarPalette.length];
}
