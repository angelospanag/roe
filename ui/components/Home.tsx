"use client";

import { useQuery } from "@tanstack/react-query";
import { Rss } from "lucide-react";
import { useState } from "react";
import { countUnreadOptions } from "@/client/@tanstack/react-query.gen";
import { FeedSidebar } from "@/components/FeedSidebar";
import { PostList } from "@/components/PostList";
import { PostReader } from "@/components/PostReader";

export function Home() {
  const [selectedFeedId, setSelectedFeedId] = useState<number | null>(null);
  const [unreadOnly, setUnreadOnly] = useState(false);
  const [offset, setOffset] = useState(0);
  const [selectedPostId, setSelectedPostId] = useState<number | null>(null);
  const { data: unread } = useQuery(countUnreadOptions());

  function selectFeed(feedId: number | null) {
    setSelectedFeedId(feedId);
    setOffset(0);
    setSelectedPostId(null);
  }

  function toggleUnreadOnly() {
    setUnreadOnly((v) => !v);
    setOffset(0);
    setSelectedPostId(null);
  }

  return (
    <div className="flex h-full flex-col">
      <header className="flex items-center gap-2 border-b bg-[var(--color-surface)] px-6 py-4">
        <Rss
          size={20}
          className="text-[var(--color-accent)]"
          strokeWidth={2.5}
        />
        <h1 className="font-serif-content text-xl font-semibold italic">
          Riffle
        </h1>
        {!!unread?.count && (
          <span className="ml-2 rounded-full bg-[var(--color-accent-soft)] px-2 py-0.5 text-xs font-medium text-[var(--color-accent)]">
            {unread.count} unread
          </span>
        )}
      </header>
      <div className="flex min-h-0 flex-1">
        <FeedSidebar
          selectedFeedId={selectedFeedId}
          onSelectFeed={selectFeed}
        />
        <PostList
          feedId={selectedFeedId}
          unreadOnly={unreadOnly}
          onToggleUnreadOnly={toggleUnreadOnly}
          offset={offset}
          onOffsetChange={setOffset}
          selectedPostId={selectedPostId}
          onSelectPost={setSelectedPostId}
        />
        <PostReader postId={selectedPostId} />
      </div>
    </div>
  );
}
