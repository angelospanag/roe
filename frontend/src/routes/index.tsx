import {createFileRoute} from "@tanstack/react-router";
import {getGetFeedsQueryOptions, useGetFeeds} from "../client.ts";
import type {Feed} from "../model";


export const Route = createFileRoute("/")({
    loader: ({context: {queryClient}}) =>
        queryClient.ensureQueryData(getGetFeedsQueryOptions()),
    component: Index,
});


function Index() {

    const {data, isLoading, error} = useGetFeeds();

    return (
        <div className="p-2">
            {isLoading && <p>Loading...</p>}
            {error && <p>Error: {error.message}</p>}
            <h3>Welcome Home!</h3>
            <ul>
                {data?.data.feeds?.map((feed: Feed) => (
                    <li key={feed.id}>
                        {feed.id} - {feed.name}
                    </li>
                ))}
            </ul>
        </div>
    );
}
