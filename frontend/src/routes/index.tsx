import { createFileRoute } from "@tanstack/react-router";
import { getGetFeedsQueryOptions, useGetFeeds } from "../client.ts";
import type { Feed } from "../model";
import FeedList from "../components/FeedList.tsx";

export const Route = createFileRoute("/")({
  loader: ({ context: { queryClient } }) =>
    queryClient.ensureQueryData(
      getGetFeedsQueryOptions({
        axios: { baseURL: import.meta.env.VITE_API_BASE_URL },
      }),
    ),
  component: App,
});

function App() {
  const { data, isLoading, error } = useGetFeeds({
    axios: { baseURL: import.meta.env.VITE_API_BASE_URL },
  });

  return (
    <div className="p-2">
      {isLoading && <p>Loading...</p>}
      {error && <p>Error: {error.message}</p>}
      <h3>Welcome Home!</h3>
      <ul>
        {data?.data.feeds?.map((feed: Feed) => (
          <li key={feed.id}>
            {feed.id} - {feed.name} - {feed.url}
          </li>
        ))}
      </ul>
      <FeedList feeds={data?.data.feeds} />
    </div>
  );
}
