import { useQuery } from "@tanstack/react-query";
import { Hash, Megaphone, MessageSquare, Mic, Volume2 } from "lucide-react";
import type { LucideIcon } from "lucide-react";

import { Card } from "../components/ui/Card";
import { CopyableId } from "../components/ui/CopyableId";
import { PageHeader } from "../components/ui/PageHeader";
import { ErrorState, LoadingState } from "../components/ui/States";
import { errorMessage } from "../lib/errors";
import { fetchChannels } from "../lib/api";

const typeIcons: Record<string, LucideIcon> = {
  text: Hash,
  voice: Volume2,
  announcement: Megaphone,
  stage: Mic,
  forum: MessageSquare,
};

export function Channels() {
  const { data, isLoading, isError, error } = useQuery({
    queryKey: ["discord-channels"],
    queryFn: fetchChannels,
  });

  return (
    <div className="flex flex-col gap-6 p-8">
      <PageHeader icon={Hash} title="Channels" description="Every channel in the guild, grouped by category." />

      {isLoading && <LoadingState label="Loading channels" />}

      {isError && <ErrorState title="Channels are unavailable." message={errorMessage(error)} />}

      {data && (
        <div className="flex max-w-2xl flex-col gap-4">
          {data.map((category) => (
            <Card key={category.id || "uncategorised"}>
              <div className="border-b border-border px-4 py-2.5 text-xs font-semibold uppercase tracking-wider text-ink-faint">
                {category.name}
              </div>
              <ul>
                {category.channels.map((channel) => {
                  const Icon = typeIcons[channel.type] ?? Hash;
                  return (
                    <li
                      key={channel.id}
                      className="flex items-center justify-between border-b border-border px-4 py-2.5 text-sm last:border-b-0"
                    >
                      <div className="flex items-center gap-2.5">
                        <Icon className="h-4 w-4 shrink-0 text-ink-faint" strokeWidth={1.75} />
                        <div className="flex flex-col">
                          <span className="text-ink">{channel.name}</span>
                          <CopyableId value={channel.id} className="text-xs" />
                        </div>
                      </div>
                      <span className="text-xs uppercase tracking-wide text-ink-faint">{channel.type}</span>
                    </li>
                  );
                })}
              </ul>
            </Card>
          ))}
        </div>
      )}
    </div>
  );
}
