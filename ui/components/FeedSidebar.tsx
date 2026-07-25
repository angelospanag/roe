"use client";

import { CheckCheck, Loader2, Plus, RefreshCw, Trash2 } from "lucide-react";
import { type FormEvent, useState } from "react";
import type { FeedResponse } from "@/client";
import {
  useCreateFeed,
  useDeleteFeed,
  useFeeds,
  useFeedUnreadCount,
  useMarkAllRead,
  useRefreshFeeds,
} from "@/lib/queries";

function FeedRow({
  feed,
  selected,
  onSelect,
}: {
  feed: FeedResponse;
  selected: boolean;
  onSelect: () => void;
}) {
  const { data: unreadCount } = useFeedUnreadCount(feed.id);
  const markAllRead = useMarkAllRead();
  const deleteFeed = useDeleteFeed();

  return (
    <div
      className={`group flex items-center gap-2 rounded-md border-l-2 px-2 py-1.5 text-sm ${
        selected
          ? "border-l-[var(--color-accent)] bg-[var(--color-accent-soft)] text-[var(--color-accent)]"
          : "border-l-transparent hover:bg-[var(--color-surface-hover)]"
      }`}
    >
      <button
        type="button"
        onClick={onSelect}
        className="min-w-0 flex-1 truncate text-left"
        title={feed.title}
      >
        {feed.title}
      </button>
      {!!unreadCount && (
        <span
          className={`shrink-0 rounded-full px-1.5 py-0.5 text-xs font-medium ${
            selected
              ? "bg-[var(--color-accent)] text-white"
              : "bg-[var(--color-accent-soft)] text-[var(--color-accent)]"
          }`}
        >
          {unreadCount}
        </span>
      )}
      <button
        type="button"
        onClick={() => markAllRead.mutate(feed.id)}
        disabled={markAllRead.isPending}
        title="Mark all as read"
        className="hidden shrink-0 opacity-70 hover:opacity-100 group-hover:block"
      >
        <CheckCheck size={14} />
      </button>
      <button
        type="button"
        onClick={() => {
          if (confirm(`Delete "${feed.title}" and all its posts?`)) {
            deleteFeed.mutate(feed.id);
          }
        }}
        disabled={deleteFeed.isPending}
        title="Delete feed"
        className="hidden shrink-0 opacity-70 hover:opacity-100 group-hover:block"
      >
        <Trash2 size={14} />
      </button>
    </div>
  );
}

function AddFeedForm({ onDone }: { onDone: () => void }) {
  const createFeed = useCreateFeed();
  const [url, setUrl] = useState("");
  const [title, setTitle] = useState("");

  function handleSubmit(e: FormEvent) {
    e.preventDefault();
    if (!url.trim()) return;
    createFeed.mutate(
      { url: url.trim(), title: title.trim() || url.trim() },
      { onSuccess: onDone },
    );
  }

  return (
    <form
      onSubmit={handleSubmit}
      className="flex flex-col gap-2 border-b px-3 py-3"
    >
      <input
        value={url}
        onChange={(e) => setUrl(e.target.value)}
        placeholder="Feed URL"
        className="rounded border bg-[var(--color-surface)] px-2 py-1 text-sm outline-none focus:border-[var(--color-accent)]"
      />
      <input
        value={title}
        onChange={(e) => setTitle(e.target.value)}
        placeholder="Title (optional)"
        className="rounded border bg-[var(--color-surface)] px-2 py-1 text-sm outline-none focus:border-[var(--color-accent)]"
      />
      {createFeed.isError && (
        <p className="text-xs text-[var(--color-danger)]">
          Could not add feed — check the URL and try again.
        </p>
      )}
      <div className="flex gap-2">
        <button
          type="submit"
          disabled={createFeed.isPending}
          className="flex-1 rounded bg-[var(--color-accent)] px-2 py-1 text-sm text-white hover:bg-[var(--color-accent-hover)] disabled:opacity-50"
        >
          {createFeed.isPending ? "Adding…" : "Add"}
        </button>
        <button
          type="button"
          onClick={onDone}
          className="rounded border px-2 py-1 text-sm hover:bg-[var(--color-surface-hover)]"
        >
          Cancel
        </button>
      </div>
    </form>
  );
}

export function FeedSidebar({
  selectedFeedId,
  onSelectFeed,
}: {
  selectedFeedId: number | null;
  onSelectFeed: (feedId: number | null) => void;
}) {
  const { data: feeds, isLoading } = useFeeds();
  const refreshFeeds = useRefreshFeeds();
  const [adding, setAdding] = useState(false);

  return (
    <aside className="flex w-64 shrink-0 flex-col border-r bg-[var(--color-surface)]">
      <div className="flex items-center justify-between px-3 py-2">
        <span className="text-xs font-medium tracking-wide text-[var(--color-text-muted)] uppercase">
          Feeds
        </span>
        <div className="flex items-center gap-1">
          <button
            type="button"
            onClick={() => refreshFeeds.mutate(undefined)}
            disabled={refreshFeeds.isPending}
            title="Refresh all feeds"
            className="rounded p-1 hover:bg-[var(--color-surface-hover)]"
          >
            {refreshFeeds.isPending ? (
              <Loader2 size={14} className="animate-spin" />
            ) : (
              <RefreshCw size={14} />
            )}
          </button>
          <button
            type="button"
            onClick={() => setAdding(true)}
            title="Add feed"
            className="rounded p-1 hover:bg-[var(--color-surface-hover)]"
          >
            <Plus size={14} />
          </button>
        </div>
      </div>

      {adding && <AddFeedForm onDone={() => setAdding(false)} />}

      <nav className="flex flex-1 flex-col gap-0.5 overflow-y-auto px-2 pb-2">
        <button
          type="button"
          onClick={() => onSelectFeed(null)}
          className={`rounded-md border-l-2 px-2 py-1.5 text-left text-sm ${
            selectedFeedId === null
              ? "border-l-[var(--color-accent)] bg-[var(--color-accent-soft)] font-medium text-[var(--color-accent)]"
              : "border-l-transparent hover:bg-[var(--color-surface-hover)]"
          }`}
        >
          All feeds
        </button>

        {isLoading && (
          <p className="px-2 py-1.5 text-sm text-[var(--color-text-muted)]">
            Loading…
          </p>
        )}
        {!isLoading && feeds?.length === 0 && (
          <p className="px-2 py-1.5 text-sm text-[var(--color-text-muted)]">
            No feeds yet — add one above.
          </p>
        )}
        {feeds?.map((feed) => (
          <FeedRow
            key={feed.id}
            feed={feed}
            selected={selectedFeedId === feed.id}
            onSelect={() => onSelectFeed(feed.id)}
          />
        ))}
      </nav>
    </aside>
  );
}
