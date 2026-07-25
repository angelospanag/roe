"use client";

import { ChevronLeft, ChevronRight, ExternalLink } from "lucide-react";
import { useState } from "react";
import type { PostResponse } from "@/client";
import { useMarkPostRead, usePosts } from "@/lib/queries";

const LIMIT = 20;

function formatDate(value?: string) {
  if (!value) return null;
  return new Date(value).toLocaleString(undefined, {
    dateStyle: "medium",
    timeStyle: "short",
  });
}

function PostRow({
  post,
  expanded,
  onToggle,
}: {
  post: PostResponse;
  expanded: boolean;
  onToggle: () => void;
}) {
  const markPostRead = useMarkPostRead();

  function handleToggle() {
    if (!post.is_read) {
      markPostRead.mutate({ postId: post.id, isRead: true });
    }
    onToggle();
  }

  return (
    <li
      className={`border-b border-l-2 ${
        post.is_read ? "border-l-transparent" : "border-l-[var(--color-accent)]"
      } ${expanded ? "bg-[var(--color-surface)]" : ""}`}
    >
      <button
        type="button"
        onClick={handleToggle}
        className="flex w-full flex-col gap-1 px-5 py-3.5 text-left hover:bg-[var(--color-surface-hover)]"
      >
        <p
          className={`font-serif-content text-base leading-snug ${
            post.is_read
              ? "text-[var(--color-text-muted)]"
              : "font-semibold text-[var(--color-text)]"
          }`}
        >
          {post.title}
        </p>
        <p className="text-xs text-[var(--color-text-muted)]">
          {post.author ? `${post.author} · ` : ""}
          {formatDate(post.published_at)}
        </p>
      </button>
      {expanded && (
        <div className="space-y-4 px-5 pb-5">
          {post.description && (
            <p className="text-sm text-[var(--color-text-muted)]">
              {post.description}
            </p>
          )}
          {post.content && (
            <div
              className="font-serif-content prose-content max-w-prose text-sm text-[var(--color-text)]"
              // biome-ignore lint/security/noDangerouslySetInnerHtml: post.content is publisher HTML meant to be rendered
              dangerouslySetInnerHTML={{ __html: post.content }}
            />
          )}
          <div className="flex items-center gap-4 text-sm">
            <a
              href={post.link}
              target="_blank"
              rel="noreferrer"
              className="inline-flex items-center gap-1 font-medium text-[var(--color-accent)] hover:text-[var(--color-accent-hover)]"
            >
              Open original <ExternalLink size={12} />
            </a>
            <button
              type="button"
              onClick={() =>
                markPostRead.mutate({ postId: post.id, isRead: !post.is_read })
              }
              className="text-[var(--color-text-muted)] hover:text-[var(--color-text)]"
            >
              Mark as {post.is_read ? "unread" : "read"}
            </button>
          </div>
        </div>
      )}
    </li>
  );
}

export function PostList({
  feedId,
  unreadOnly,
  onToggleUnreadOnly,
  offset,
  onOffsetChange,
}: {
  feedId: number | null;
  unreadOnly: boolean;
  onToggleUnreadOnly: () => void;
  offset: number;
  onOffsetChange: (offset: number) => void;
}) {
  const { data: posts, isLoading } = usePosts({
    feedId,
    unreadOnly,
    offset,
    limit: LIMIT,
  });
  const [expandedId, setExpandedId] = useState<number | null>(null);

  return (
    <div className="flex h-full flex-col">
      <div className="flex items-center justify-between border-b bg-[var(--color-surface)] px-5 py-2.5">
        <h2 className="text-sm font-medium text-[var(--color-text-muted)]">
          {feedId === null ? "All posts" : "Posts"}
        </h2>
        <label className="flex items-center gap-1.5 text-sm text-[var(--color-text-muted)]">
          <input
            type="checkbox"
            checked={unreadOnly}
            onChange={onToggleUnreadOnly}
            className="accent-[var(--color-accent)]"
          />
          Unread only
        </label>
      </div>

      {isLoading && (
        <p className="px-5 py-6 text-sm text-[var(--color-text-muted)]">
          Loading…
        </p>
      )}
      {!isLoading && posts?.length === 0 && (
        <p className="px-5 py-6 text-sm text-[var(--color-text-muted)]">
          No posts to show.
        </p>
      )}

      <ul>
        {posts?.map((post) => (
          <PostRow
            key={post.id}
            post={post}
            expanded={expandedId === post.id}
            onToggle={() =>
              setExpandedId((id) => (id === post.id ? null : post.id))
            }
          />
        ))}
      </ul>

      <div className="mt-auto flex items-center justify-between border-t bg-[var(--color-surface)] px-5 py-2">
        <button
          type="button"
          disabled={offset === 0}
          onClick={() => onOffsetChange(Math.max(0, offset - LIMIT))}
          className="inline-flex items-center gap-1 rounded px-2 py-1 text-sm hover:bg-[var(--color-surface-hover)] disabled:opacity-40"
        >
          <ChevronLeft size={14} /> Prev
        </button>
        <span className="text-xs text-[var(--color-text-muted)]">
          {offset + 1}–{offset + (posts?.length ?? 0)}
        </span>
        <button
          type="button"
          disabled={(posts?.length ?? 0) < LIMIT}
          onClick={() => onOffsetChange(offset + LIMIT)}
          className="inline-flex items-center gap-1 rounded px-2 py-1 text-sm hover:bg-[var(--color-surface-hover)] disabled:opacity-40"
        >
          Next <ChevronRight size={14} />
        </button>
      </div>
    </div>
  );
}
