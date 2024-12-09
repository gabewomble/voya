import { createContextId } from "@builder.io/qwik";
import { routeAction$, routeLoader$, z, zod$ } from "@builder.io/qwik-city";
import { serverFetch } from "~/helpers/server-fetch";
import {
  batchGetUsersResponseSchema,
  listNotificationsResponseSchema,
} from "~/types/api";
import { type User } from "~/types/users";

export const useIsShowingAll = routeLoader$(async (request) => {
  const showAll = request.url.searchParams.get("show") === "all";
  return showAll;
});

export const useGetActivity = routeLoader$(async (request) => {
  const isShowingAll = request.url.searchParams.get("show") === "all";

  const res = await serverFetch(
    isShowingAll ? "/notifications" : "/notifications/unread",
    {},
    request,
  );
  const json = await res.json();
  const result = listNotificationsResponseSchema.parse(json);

  const batchGetUserIds = new Set<string>();

  result.notifications.forEach(
    ({
      created_by: createdBy,
      target_user_id: targetUserID,
      user_id: userID,
    }) => {
      batchGetUserIds.add(createdBy);
      batchGetUserIds.add(targetUserID);
      batchGetUserIds.add(userID);
    },
  );

  const batchResult = await serverFetch(
    `/users/batch`,
    {
      method: "POST",
      body: JSON.stringify({ user_ids: Array.from(batchGetUserIds) }),
    },
    request,
  );

  const usersCache: Record<string, User> = {};
  const batchUsers = batchGetUsersResponseSchema.safeParse(
    await batchResult.json(),
  );

  if (batchUsers.success) {
    batchUsers.data.users.forEach((user) => {
      usersCache[user.id] = user;
    });
  }

  return {
    ...result,
    users: usersCache,
  };
});

export const useMarkAllAsRead = routeAction$(
  async (_, requestEvent) => {
    const res = await serverFetch(
      "/notifications/read",
      {
        method: "POST",
      },
      requestEvent,
    );

    if (!res.ok) {
      return requestEvent.fail(500, {
        error: "Failed to mark all as read",
      });
    }

    return res.ok;
  },
  zod$(z.object({})),
);

export const useMarkAsRead = routeAction$(
  async (data, requestEvent) => {
    const res = await serverFetch(
      `/notification/${data.notificationId}/read`,
      {
        method: "POST",
      },
      requestEvent,
    );

    if (!res.ok) {
      return requestEvent.fail(500, {
        error: "Failed to mark as read",
      });
    }

    return res.ok;
  },
  zod$(z.object({ notificationId: z.string().uuid() })),
);

export const ActivityUsersContext =
  createContextId<Record<string, User>>("ActivityUsers");
