import { component$, useContextProvider } from "@builder.io/qwik";
import { Link } from "@builder.io/qwik-city";
import { NotificationItem } from "./Notification";
import {
  useIsShowingAll,
  useGetActivity,
  useMarkAllAsRead,
  ActivityUsersContext,
} from "./layout";

export default component$(() => {
  const isShowingAll = useIsShowingAll().value;
  const { notifications, total, users } = useGetActivity().value;
  const markAllAsRead = useMarkAllAsRead();

  useContextProvider(ActivityUsersContext, users);

  return (
    <div class="container mx-auto py-8">
      <div class="mb-8 flex items-center justify-between">
        <h1 class="text-4xl font-bold">Activity</h1>
        <div class="flex gap-4">
          {isShowingAll ? (
            <Link href="/activity?show=unread" class="btn btn-ghost">
              Show Unread
            </Link>
          ) : (
            <Link href="/activity?show=all" class="btn btn-ghost">
              Show All
            </Link>
          )}
          {total > 0 && (
            <button
              class="btn btn-ghost"
              type="button"
              onClick$={() => {
                markAllAsRead.submit({});
              }}
            >
              Mark all as read
            </button>
          )}
        </div>
      </div>

      {total === 0 ? (
        <div class="py-12 text-center">
          <p class="text-xl text-gray-500">Nothing to see here</p>
        </div>
      ) : (
        <div class="flex flex-col gap-4">
          {notifications.map((notification) => (
            <NotificationItem
              notification={notification}
              key={notification.id}
            />
          ))}
        </div>
      )}
    </div>
  );
});
