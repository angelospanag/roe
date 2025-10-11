import type { GetFeedsOutputBodyFeeds } from "../model";

interface FeedListProps {
  feeds: GetFeedsOutputBodyFeeds;
}

export default function FeedList({ feeds }: FeedListProps) {
  if (!feeds || feeds.length === 0) {
    return (
      <div className="flex items-center justify-center h-40 text-gray-500">
        No feeds available.
      </div>
    );
  }

  return (
    <div className="grid gap-4 md:grid-cols-2 lg:grid-cols-3 p-4">
      {feeds.map((feed) => (
        <div
          key={feed.id}
          className="bg-white shadow-md rounded-2xl p-4 border border-gray-100 hover:shadow-lg transition-shadow"
        >
          <div className="flex items-center justify-between">
            <h2 className="text-lg font-semibold text-gray-800 truncate">
              {feed.name}
            </h2>
            {feed.unreadItemsCount > 0 && (
              <span className="ml-2 bg-blue-500 text-white text-xs font-medium px-2 py-1 rounded-full">
                {feed.unreadItemsCount} unread
              </span>
            )}
          </div>

          <a
            href={feed.url}
            target="_blank"
            rel="noopener noreferrer"
            className="text-blue-600 hover:underline text-sm mt-2 block truncate"
          >
            {feed.url}
          </a>
        </div>
      ))}
    </div>
  );
}
