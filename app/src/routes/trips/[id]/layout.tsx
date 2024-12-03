import { routeLoader$ } from "@builder.io/qwik-city";
import { serverFetch } from "~/helpers/server-fetch";
import { getTripByIdResponseSchema } from "~/types/api";
import type { InitialValues } from "@modular-forms/qwik";
import type { AddMemberForm } from "./Members";

// Pretty annoying I need to export this from index.tsx or layout.tsx
export const useAddMemberLoader = routeLoader$<InitialValues<AddMemberForm>>(
  () => ({
    identifier: "",
  }),
);

export const useTripData = routeLoader$(async (requestEvent) => {
  const { id } = requestEvent.params;
  const res = await serverFetch(`/trips/t/${id}`, {}, requestEvent);
  const json = await res.json();
  return getTripByIdResponseSchema.parse(json);
});
