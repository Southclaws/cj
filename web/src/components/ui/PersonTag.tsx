import { Link } from "react-router";

import { Avatar } from "./Avatar";
import { CopyableId } from "./CopyableId";
import type { PersonRefRecord } from "../../lib/api";

interface PersonTagProps {
  person: PersonRefRecord;
  fallbackName?: string;
  size?: number;
  showId?: boolean;
  className?: string;
}

export function PersonTag({ person, fallbackName, size = 28, showId = true, className = "" }: PersonTagProps) {
  const name = person.nickname || person.username || fallbackName || "Unknown user";
  const label = person.nickname ? `${person.nickname} (${person.username})` : name;

  return (
    <div className={`flex items-center gap-2.5 ${className}`}>
      <Avatar id={person.id} name={name} url={person.avatarUrl} size={size} />
      <div className="flex min-w-0 flex-col">
        <Link to={`/users/${person.id}`} className="truncate text-ink hover:text-accent">
          {label}
        </Link>
        {showId && <CopyableId value={person.id} className="text-xs" />}
      </div>
    </div>
  );
}
