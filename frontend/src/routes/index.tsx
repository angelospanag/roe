import { createFileRoute } from "@tanstack/react-router";
import { useSuspenseQuery } from "@tanstack/react-query";
import { getGetFeedsQueryOptions } from "../client";
import { Feed } from "../model";

export const Route = createFileRoute("/")({
  loader: ({ context: { queryClient } }) =>
    queryClient.ensureQueryData(getGetFeedsQueryOptions()),
  component: Index,
});

function Index() {
  // Read the data from the cache and subscribe to updates
  const {
    data: { data, isLoading, error },
  } = useSuspenseQuery(getGetFeedsQueryOptions());

  return (
    <div className="p-2">
      {isLoading && <p>Loading...</p>}
      {error && <p>Error: {error.message}</p>}
      <h3>Welcome Home!</h3>
      <ul>
        {data.feeds?.map((feed: Feed) => (
          <li key={feed.id}>
            {feed.id} - {feed.name}
          </li>
        ))}
      </ul>
    </div>
  );
}
