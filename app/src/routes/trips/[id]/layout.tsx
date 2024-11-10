import { routeLoader$ } from "@builder.io/qwik-city";
import { serverFetch } from "~/helpers/server-fetch";
import type { Trip } from "~/types/trips";

export const useTripData = routeLoader$(async (requestEvent) => {
  const { id } = requestEvent.params;
  const res = await serverFetch(`/trips/${id}`, {}, requestEvent);
  const json = await res.json();
  return json.trip as Trip;
});
