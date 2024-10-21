import { routeLoader$ } from "@builder.io/qwik-city";
import type { Trip } from "~/types/trips";

export const useTripData = routeLoader$(async (requestEvent) => {
  const { id } = requestEvent.params;
  const res = await fetch(`http://localhost:8080/trips/${id}`);
  const json = await res.json();
  return json.trip as Trip;
});
