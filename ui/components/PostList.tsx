"use client";

import { useMutation, useQuery, useQueryClient } from "@tanstack/react-query";
import { ChevronLeft, ChevronRight } from "lucide-react";
import type { PostResponse } from "@/client";
import {
  listPostsOptions,
  markPostReadMutation,
} from "@/client/@tanstack/react-query.gen";

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
  selected,
  onSelect,
}: {
  post: PostResponse;
  selected: boolean;
  onSelect: () => void;
}) {
  const queryClient = useQueryClient();
  const markPostRead = useMutation({
    ...markPostReadMutation(),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: [{ _id: "listPosts" }] });
      queryClient.invalidateQueries({ queryKey: [{ _id: "countUnread" }] });
      queryClient.invalidateQueries({
        queryKey: [{ _id: "countUnreadByFeed" }],
      });
      queryClient.invalidateQueries({ queryKey: [{ _id: "getPost" }] });
    },
  });

  function handleClick() {
    if (!post.is_read) {
      markPostRead.mutate({ path: { id: post.id }, body: { is_read: true } });
    }
    onSelect();
  }

  return (
    <li
      className={`border-b border-l-2 ${
        post.is_read ? "border-l-transparent" : "border-l-[var(--color-accent)]"
      } ${selected ? "bg-[var(--color-accent-soft)]" : ""}`}
    >
      <button
        type="button"
        onClick={handleClick}
        className={`flex w-full flex-col gap-1 px-4 py-3 text-left ${
          selected ? "" : "hover:bg-[var(--color-surface-hover)]"
        }`}
      >
        <p
          className={`font-serif-content text-sm leading-snug ${
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
    </li>
  );
}

export function PostList({
  feedId,
  unreadOnly,
  onToggleUnreadOnly,
  offset,
  onOffsetChange,
  selectedPostId,
  onSelectPost,
}: {
  feedId: number | null;
  unreadOnly: boolean;
  onToggleUnreadOnly: () => void;
  offset: number;
  onOffsetChange: (offset: number) => void;
  selectedPostId: number | null;
  onSelectPost: (postId: number) => void;
}) {
  const { data: posts, isLoading } = useQuery(
    listPostsOptions({
      query: {
        feed_id: feedId ?? 0,
        unread_only: unreadOnly,
        offset,
        limit: LIMIT,
      },
    }),
  );

  return (
    <div className="flex h-full w-80 shrink-0 flex-col border-r bg-[var(--color-surface)]">
      <div className="flex shrink-0 items-center justify-between border-b px-4 py-2.5">
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

      <div className="flex-1 overflow-y-auto">
        {isLoading && (
          <p className="px-4 py-6 text-sm text-[var(--color-text-muted)]">
            Loading…
          </p>
        )}
        {!isLoading && posts?.length === 0 && (
          <p className="px-4 py-6 text-sm text-[var(--color-text-muted)]">
            No posts to show.
          </p>
        )}

        <ul>
          {posts?.map((post) => (
            <PostRow
              key={post.id}
              post={post}
              selected={selectedPostId === post.id}
              onSelect={() => onSelectPost(post.id)}
            />
          ))}
        </ul>
      </div>

      <div className="flex shrink-0 items-center justify-between border-t px-4 py-2">
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
