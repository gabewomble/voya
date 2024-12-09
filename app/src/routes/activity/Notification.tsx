import { component$ } from "@builder.io/qwik";
import { formatDistance } from "~/helpers/date";
import { type Notification } from "~/types/notifications";
import { useMarkAsRead } from "./layout";

const NotificationIcon = component$<{
  type: Notification["notification_type"];
}>(({ type }) => {
  switch (type) {
    case "trip_invite_pending":
      return <span class="text-info">📩</span>;
    case "trip_invite_accepted":
      return <span class="text-success">✅</span>;
    case "trip_invite_declined":
      return <span class="text-error">❌</span>;
    case "trip_member_removed":
      return <span class="text-warning">⚠️</span>;
    case "trip_ownership_transfer":
      return <span class="text-primary">👑</span>;
    default:
      return <span>📢</span>;
  }
});

const NotificationMessage = component$<{ notification: Notification }>(
  ({ notification }) => {
    let message = notification.message;

    switch (notification.notification_type) {
      case "trip_invite_pending":
        message = `${notification.metadata.user_name} invited you to join their trip`;
        break;
      case "trip_invite_accepted":
        message = `${notification.metadata.user_name} accepted your invitation`;
        break;
      case "trip_invite_declined":
        message = `${notification.metadata.user_name} declined your invitation`;
        break;
      case "trip_member_removed":
        message = `${notification.metadata.user_name} removed you from the trip`;
        break;
      case "trip_ownership_transfer":
        message = `${notification.metadata.user_name} transferred the trip ownership to you`;
        break;
    }

    return <p class="text-lg">{message}</p>;
  },
);

export const NotificationItem = component$<{ notification: Notification }>(
  ({ notification }) => {
    const markAsRead = useMarkAsRead();
    return (
      <div
        key={notification.id}
        class={`card bg-base-200 shadow-xl ${
          !notification.read_at ? "border-l-4 border-primary" : ""
        }`}
      >
        <div class="card-body">
          <div class="flex items-start gap-4">
            <div class="text-2xl">
              <NotificationIcon type={notification.notification_type} />
            </div>
            <div class="flex-1">
              <NotificationMessage notification={notification} />
              <div class="mt-2 flex items-center gap-2">
                <span class="text-sm opacity-70">
                  {formatDistance(new Date(notification.created_at))}
                </span>
                {notification.trip_id && (
                  <a
                    href={`/trips/${notification.trip_id}`}
                    class="link-hover link text-sm"
                  >
                    View Trip →
                  </a>
                )}
              </div>
            </div>
            {!notification.read_at && (
              <button
                class="btn btn-ghost btn-sm"
                onClick$={() => {
                  markAsRead.submit({ notificationId: notification.id });
                }}
              >
                Mark as read
              </button>
            )}
          </div>
        </div>
      </div>
    );
  },
);
