import { routeAction$, routeLoader$, z, zod$ } from "@builder.io/qwik-city";
import { serverFetch } from "~/helpers/server-fetch";
import { listNotificationsResponseSchema } from "~/types/api";

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
  return listNotificationsResponseSchema.parse(json);
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
