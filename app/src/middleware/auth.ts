import type { RequestHandler } from "@builder.io/qwik-city";
import { serverFetch } from "~/helpers/server-fetch";

export const authenticate: RequestHandler = async (requestEvent) => {
  const token = requestEvent.cookie.get("token")?.value;

  if (token) {
    const res = await serverFetch("/users/current", {}, requestEvent);

    if (res.ok) {
      const data = await res.json();
      const user = data?.user ?? null;

      requestEvent.sharedMap.set("user", user);
      return;
    }
  }

  requestEvent.sharedMap.set("user", null);
};

export const requireAuth: RequestHandler = async ({ redirect, sharedMap }) => {
  const user = sharedMap.get("user");

  if (!user) {
    throw redirect(303, "/login");
  }
};
