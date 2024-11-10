import { component$, Slot, useContextProvider } from "@builder.io/qwik";
import { Page } from "~/components";
import { authenticate } from "~/middleware/auth";
import { routeLoader$, type RequestHandler } from "@builder.io/qwik-city";
import { UserContext } from "~/context/user";

export const onRequest: RequestHandler = async (request) => {
  await authenticate(request);
};

export const onGet: RequestHandler = async ({ cacheControl }) => {
  // Control caching for this request for best performance and to reduce hosting costs:
  // https://qwik.dev/docs/caching/
  cacheControl({
    // Always serve a cached response by default, up to a week stale
    staleWhileRevalidate: 60 * 60 * 24 * 7,
    // Max once every 5 seconds, revalidate on the server to get a fresh version of this page
    maxAge: 5,
  });
};

export const useUserData = routeLoader$(async (request) => {
  const user = request.sharedMap.get("user");
  return user;
});

export default component$(() => {
  const userData = useUserData();

  useContextProvider(UserContext, userData);

  return (
    <Page>
      <Slot />
    </Page>
  );
});
