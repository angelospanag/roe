"use client";

import { useMutation, useQuery, useQueryClient } from "@tanstack/react-query";
import { ExternalLink, Minus, Plus } from "lucide-react";
import { useEffect, useState } from "react";
import {
  getPostOptions,
  markPostReadMutation,
} from "@/client/@tanstack/react-query.gen";

const READER_FONT_SIZE_STORAGE_KEY = "riffle:reader-font-size-rem";
const MIN_FONT_SIZE_REM = 0.75;
const MAX_FONT_SIZE_REM = 1.5;
const FONT_SIZE_STEP_REM = 0.125;
const DEFAULT_FONT_SIZE_REM = 0.875; // matches the old fixed text-sm

function clampFontSize(value: number) {
  return Math.min(MAX_FONT_SIZE_REM, Math.max(MIN_FONT_SIZE_REM, value));
}

// Reads/writes a rem font size for the reading pane's body text from
// localStorage. Starts at the default on both server and first client
// render (localStorage isn't available during SSR) and syncs to the
// stored value in an effect, so there's no hydration mismatch.
function useReaderFontSize() {
  const [fontSize, setFontSize] = useState(DEFAULT_FONT_SIZE_REM);

  useEffect(() => {
    const stored = window.localStorage.getItem(READER_FONT_SIZE_STORAGE_KEY);
    const parsed = stored ? Number.parseFloat(stored) : Number.NaN;
    if (!Number.isNaN(parsed)) {
      setFontSize(clampFontSize(parsed));
    }
  }, []);

  function setAndPersist(next: number) {
    const clamped = clampFontSize(next);
    setFontSize(clamped);
    window.localStorage.setItem(READER_FONT_SIZE_STORAGE_KEY, String(clamped));
  }

  return {
    fontSize,
    increase: () => setAndPersist(fontSize + FONT_SIZE_STEP_REM),
    decrease: () => setAndPersist(fontSize - FONT_SIZE_STEP_REM),
  };
}

function formatDate(value?: string) {
  if (!value) return null;
  return new Date(value).toLocaleString(undefined, {
    dateStyle: "long",
    timeStyle: "short",
  });
}

export function PostReader({ postId }: { postId: number | null }) {
  const queryClient = useQueryClient();
  const { fontSize, increase, decrease } = useReaderFontSize();
  const { data: post } = useQuery({
    ...getPostOptions({ path: { id: postId ?? 0 } }),
    enabled: postId !== null,
  });
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

  if (postId === null || !post) {
    return (
      <div className="flex flex-1 items-center justify-center text-sm text-[var(--color-text-muted)]">
        Select a post to read
      </div>
    );
  }

  // RSS feeds put HTML in either field — some feeds' <description> is just
  // an <img> tag with no <content:encoded> at all. Render whichever is present.
  const body = post.content || post.description;

  return (
    <div className="flex-1 overflow-y-auto">
      <article className="mx-auto max-w-2xl px-8 py-8">
        <h1 className="font-serif-content text-2xl font-semibold text-[var(--color-text)]">
          {post.title}
        </h1>
        <p className="mt-2 text-sm text-[var(--color-text-muted)]">
          {post.author ? `${post.author} · ` : ""}
          {formatDate(post.published_at)}
        </p>
        <div className="mt-4 flex items-center gap-4 text-sm">
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
              markPostRead.mutate({
                path: { id: post.id },
                body: { is_read: !post.is_read },
              })
            }
            className="text-[var(--color-text-muted)] hover:text-[var(--color-text)]"
          >
            Mark as {post.is_read ? "unread" : "read"}
          </button>
          <div className="ml-auto flex items-center gap-1 text-[var(--color-text-muted)]">
            <button
              type="button"
              onClick={decrease}
              disabled={fontSize <= MIN_FONT_SIZE_REM}
              aria-label="Decrease text size"
              className="rounded p-1 hover:bg-[var(--color-surface-hover)] hover:text-[var(--color-text)] disabled:opacity-40 disabled:hover:bg-transparent"
            >
              <Minus size={14} />
            </button>
            <button
              type="button"
              onClick={increase}
              disabled={fontSize >= MAX_FONT_SIZE_REM}
              aria-label="Increase text size"
              className="rounded p-1 hover:bg-[var(--color-surface-hover)] hover:text-[var(--color-text)] disabled:opacity-40 disabled:hover:bg-transparent"
            >
              <Plus size={14} />
            </button>
          </div>
        </div>
        {body && (
          <div
            className="font-serif-content prose-content mt-6 text-[var(--color-text)]"
            style={{ fontSize: `${fontSize}rem` }}
            // biome-ignore lint/security/noDangerouslySetInnerHtml: post content/description is publisher HTML meant to be rendered
            dangerouslySetInnerHTML={{ __html: body }}
          />
        )}
      </article>
    </div>
  );
}
