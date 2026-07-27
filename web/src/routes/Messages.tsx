import { useState } from "react";
import { useQuery } from "@tanstack/react-query";
import { useSearchParams } from "react-router";
import { MessagesSquare } from "lucide-react";

import { MemberUserPicker } from "../components/MemberUserPicker";
import { PageHeader } from "../components/ui/PageHeader";
import { PersonTag } from "../components/ui/PersonTag";
import { SearchInput } from "../components/ui/SearchInput";
import { EmptyState, ErrorState, LoadingState } from "../components/ui/States";
import { Table, Thead, Th, Td, Tr } from "../components/ui/Table";
import { errorMessage } from "../lib/errors";
import { useChannelOptions } from "../lib/entitySearch";
import { useDebounced } from "../lib/useDebounced";
import { searchMessages } from "../lib/api";

export function Messages() {
  const [searchParams] = useSearchParams();
  const presetUserId = searchParams.get("userId") ?? "";

  const [query, setQuery] = useState("");
  const debouncedQuery = useDebounced(query, 300);
  const [selectedUserId, setSelectedUserId] = useState(presetUserId);

  const { all: channels } = useChannelOptions("");
  const channelNames = new Map(channels.map((channel) => [channel.id, channel.label]));

  const {
    data: messages,
    isLoading,
    isError,
    error,
  } = useQuery({
    queryKey: ["message-search", debouncedQuery, selectedUserId],
    queryFn: () => searchMessages(debouncedQuery, selectedUserId || undefined),
    enabled: debouncedQuery.length > 0 || selectedUserId !== "",
  });

  return (
    <div className="flex flex-col gap-6 p-8">
      <PageHeader
        icon={MessagesSquare}
        title="Messages"
        description="Search every message CJ has recorded, optionally narrowed to one user."
      />

      <div className="flex max-w-2xl flex-col gap-3">
        <SearchInput
          placeholder="Search message text..."
          value={query}
          onChange={(e) => setQuery(e.target.value)}
          className="max-w-none"
        />

        <MemberUserPicker value={selectedUserId} onChange={setSelectedUserId} placeholder="Filter by user, ID or name..." />
      </div>

      {debouncedQuery.length === 0 && !selectedUserId && (
        <EmptyState message="Enter a search term, or filter by a user, to find messages." />
      )}

      {isLoading && <LoadingState label="Searching" />}

      {isError && <ErrorState title="Message search is unavailable." message={errorMessage(error)} />}

      {messages && messages.length === 0 && <EmptyState message="No messages match that search." />}

      {messages && messages.length > 0 && (
        <div className="max-w-4xl">
          <Table>
            <Thead>
              <Th>Time</Th>
              <Th>User</Th>
              <Th>Channel</Th>
              <Th>Message</Th>
            </Thead>
            <tbody>
              {messages.map((message) => (
                <Tr key={message.discordMessageId}>
                  <Td className="whitespace-nowrap text-ink-faint">
                    {new Date(message.timestamp * 1000).toLocaleString()}
                  </Td>
                  <Td>
                    <PersonTag person={message.person} size={24} showId={false} />
                  </Td>
                  <Td className="text-ink-muted">{channelNames.get(message.discordChannel) ?? message.discordChannel}</Td>
                  <Td className="text-ink">{message.message}</Td>
                </Tr>
              ))}
            </tbody>
          </Table>
        </div>
      )}
    </div>
  );
}
