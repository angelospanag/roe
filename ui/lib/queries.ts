import { useMutation, useQuery, useQueryClient } from "@tanstack/react-query";
import {
  countUnread,
  countUnreadByFeed,
  createFeed,
  deleteFeed,
  type Input,
  listFeeds,
  listPosts,
  markAllRead,
  markPostRead,
  refreshFeeds,
} from "@/client";

const feedsKey = ["feeds"] as const;
const unreadCountKey = ["posts", "unread-count"] as const;
const feedUnreadCountKey = (feedId: number) =>
  ["feeds", feedId, "unread-count"] as const;
const postsKey = (params: {
  feedId: number | null;
  unreadOnly: boolean;
  offset: number;
}) => ["posts", params] as const;

export function useFeeds() {
  return useQuery({
    queryKey: feedsKey,
    queryFn: async () => (await listFeeds({ throwOnError: true })).data,
  });
}

export function useCreateFeed() {
  const queryClient = useQueryClient();
  return useMutation({
    mutationFn: async (body: Input) =>
      (await createFeed({ body, throwOnError: true })).data,
    onSuccess: () => queryClient.invalidateQueries({ queryKey: feedsKey }),
  });
}

export function useDeleteFeed() {
  const queryClient = useQueryClient();
  return useMutation({
    mutationFn: async (feedId: number) =>
      (await deleteFeed({ path: { id: feedId }, throwOnError: true })).data,
    onSuccess: () => queryClient.invalidateQueries({ queryKey: feedsKey }),
  });
}

export function useRefreshFeeds() {
  const queryClient = useQueryClient();
  return useMutation({
    mutationFn: async (feedId?: number) =>
      (
        await refreshFeeds({
          body: feedId ? { feed_id: feedId } : {},
          throwOnError: true,
        })
      ).data,
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: feedsKey });
      queryClient.invalidateQueries({ queryKey: ["posts"] });
    },
  });
}

export function useFeedUnreadCount(feedId: number) {
  return useQuery({
    queryKey: feedUnreadCountKey(feedId),
    queryFn: async () =>
      (await countUnreadByFeed({ path: { id: feedId }, throwOnError: true }))
        .data.count,
  });
}

export function useGlobalUnreadCount() {
  return useQuery({
    queryKey: unreadCountKey,
    queryFn: async () => (await countUnread({ throwOnError: true })).data.count,
  });
}

export function useMarkAllRead() {
  const queryClient = useQueryClient();
  return useMutation({
    mutationFn: async (feedId: number) =>
      (await markAllRead({ path: { id: feedId }, throwOnError: true })).data,
    onSuccess: (_data, feedId) => {
      queryClient.invalidateQueries({ queryKey: ["posts"] });
      queryClient.invalidateQueries({ queryKey: unreadCountKey });
      queryClient.invalidateQueries({ queryKey: feedUnreadCountKey(feedId) });
    },
  });
}

export function usePosts(params: {
  feedId: number | null;
  unreadOnly: boolean;
  offset: number;
  limit: number;
}) {
  return useQuery({
    queryKey: postsKey(params),
    queryFn: async () =>
      (
        await listPosts({
          query: {
            feed_id: params.feedId ?? 0,
            unread_only: params.unreadOnly,
            offset: params.offset,
            limit: params.limit,
          },
          throwOnError: true,
        })
      ).data,
  });
}

export function useMarkPostRead() {
  const queryClient = useQueryClient();
  return useMutation({
    mutationFn: async (vars: { postId: number; isRead: boolean }) =>
      (
        await markPostRead({
          path: { id: vars.postId },
          body: { is_read: vars.isRead },
          throwOnError: true,
        })
      ).data,
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ["posts"] });
      queryClient.invalidateQueries({ queryKey: unreadCountKey });
      queryClient.invalidateQueries({ queryKey: ["feeds"] });
    },
  });
}
