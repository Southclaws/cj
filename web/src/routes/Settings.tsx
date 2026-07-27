import { useEffect, useState } from "react";
import { useMutation, useQuery, useQueryClient } from "@tanstack/react-query";
import { SlidersHorizontal } from "lucide-react";

import { ChannelPicker } from "../components/ChannelPicker";
import { LazyUserLabel } from "../components/LazyUserLabel";
import { MultiPicker } from "../components/MultiPicker";
import { Button } from "../components/ui/Button";
import { Card } from "../components/ui/Card";
import { PageHeader } from "../components/ui/PageHeader";
import { ErrorState, LoadingState } from "../components/ui/States";
import { errorMessage } from "../lib/errors";
import { useMemberAndUserSearch } from "../lib/entitySearch";
import { pushToast } from "../lib/toast";
import { fetchGuildSettings, saveGuildSettings, type GuildSettings } from "../lib/api";

const emptySettings: GuildSettings = {
  searchMessageChannelId: "",
  errorReportChannelId: "",
  leaderboardChannelId: "",
  adsChannelId: "",
  ltfChannelId: "",
  ltfUserIds: [],
};

const channelFields: { key: keyof GuildSettings; label: string; help: string }[] = [
  {
    key: "searchMessageChannelId",
    label: "Search-message channel",
    help: "Channel /searchmessage must be run in. Also requires roles configured via /config searchmessage.",
  },
  {
    key: "errorReportChannelId",
    label: "Error-report channel",
    help: "Where readme-sync failures are posted. Left empty, failures are simply not announced.",
  },
  {
    key: "leaderboardChannelId",
    label: "Leaderboard channel",
    help: "Where the periodic message-count leaderboard is announced.",
  },
  {
    key: "adsChannelId",
    label: "Ads channel",
    help: "Channel the ad-moderation watcher moderates.",
  },
  {
    key: "ltfChannelId",
    label: "LTF channel",
    help: "Channel for the /ltf command and its heartbeat repost.",
  },
];

export function Settings() {
  const queryClient = useQueryClient();
  const { data, isLoading, isError, error } = useQuery({
    queryKey: ["guild-settings"],
    queryFn: fetchGuildSettings,
  });

  const [form, setForm] = useState<GuildSettings>(emptySettings);

  useEffect(() => {
    if (data) {
      setForm(data);
    }
  }, [data]);

  const mutation = useMutation({
    mutationFn: saveGuildSettings,
    onSuccess: (saved) => {
      queryClient.setQueryData(["guild-settings"], saved);
      pushToast("Settings saved.", "good");
    },
    onError: (error) => {
      pushToast(`Failed to save settings: ${errorMessage(error)}`, "bad");
    },
  });

  function handleSubmit(e: React.FormEvent) {
    e.preventDefault();
    mutation.mutate(form);
  }

  return (
    <div className="flex flex-col gap-6 p-8">
      <PageHeader
        icon={SlidersHorizontal}
        title="Settings"
        description="Guild-specific channels and users. Search by name and pick a result instead of typing an ID. Leave a field empty to keep the feature it backs disabled."
      />

      {isLoading && <LoadingState label="Loading settings" />}

      {isError && <ErrorState title="Settings are unavailable." message={errorMessage(error)} />}

      {data && (
        <Card className="max-w-lg p-6">
          <form onSubmit={handleSubmit} className="flex flex-col gap-5">
            {channelFields.map((field) => (
              <label key={field.key} className="flex flex-col gap-1">
                <span className="text-sm text-ink">{field.label}</span>
                <ChannelPicker
                  value={form[field.key] as string}
                  onChange={(id) => setForm((prev) => ({ ...prev, [field.key]: id }))}
                />
                <span className="text-xs text-ink-faint">{field.help}</span>
              </label>
            ))}

            <label className="flex flex-col gap-1">
              <span className="text-sm text-ink">LTF users</span>
              <MultiPicker
                ids={form.ltfUserIds}
                onChange={(ltfUserIds) => setForm((prev) => ({ ...prev, ltfUserIds }))}
                useOptions={useMemberAndUserSearch}
                placeholder="Search by name or ID..."
                renderChipLabel={(id) => <LazyUserLabel id={id} />}
              />
              <span className="text-xs text-ink-faint">Users eligible for the /ltf command and its heartbeat repost.</span>
            </label>

            <div className="flex items-center gap-3">
              <Button type="submit" variant="primary" disabled={mutation.isPending}>
                {mutation.isPending ? "Saving..." : "Save"}
              </Button>
            </div>
          </form>
        </Card>
      )}
    </div>
  );
}
