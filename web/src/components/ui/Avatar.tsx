import { useState } from "react";

import { avatarColor } from "../../lib/avatar";

interface AvatarProps {
  id: string;
  name: string;
  url?: string;
  size?: number;
  className?: string;
}

const sizeText: Record<number, string> = {
  20: "text-[10px]",
  24: "text-xs",
  28: "text-xs",
  36: "text-sm",
  56: "text-xl",
};

export function Avatar({ id, name, url, size = 28, className = "" }: AvatarProps) {
  const [errored, setErrored] = useState(false);
  const initial = name.charAt(0).toUpperCase() || "?";
  const dimension = { width: size, height: size };
  const textSize = sizeText[size] ?? "text-xs";

  if (url && !errored) {
    return (
      <img
        src={url}
        alt=""
        width={size}
        height={size}
        style={dimension}
        className={`shrink-0 rounded-full bg-surface-3 object-cover ${className}`}
        onError={() => setErrored(true)}
      />
    );
  }

  return (
    <div
      style={{ ...dimension, backgroundColor: avatarColor(id) }}
      className={`flex shrink-0 items-center justify-center rounded-full font-bold text-white ${textSize} ${className}`}
    >
      {initial}
    </div>
  );
}
