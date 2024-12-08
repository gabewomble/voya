import { routeAction$, routeLoader$, z, zod$ } from "@builder.io/qwik-city";
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

const inviteUserInputSchema = z.object({
  userID: z.string().uuid(),
  tripID: z.string().uuid(),
});

export const useInviteUser = routeAction$(async (data, requestEvent) => {
  const response = await serverFetch(
    `/trip/${data.tripID}/members`,
    {
      method: "POST",
      body: JSON.stringify({
        user_id: data.userID,
      }),
    },
    requestEvent,
  );

  if (!response.ok) {
    return requestEvent.fail(500, {
      error: "Failed to invite user",
    });
  }

  return response.ok;
}, zod$(inviteUserInputSchema));

export const useTripData = routeLoader$(async (requestEvent) => {
  const { id } = requestEvent.params;
  const res = await serverFetch(`/trip/${id}`, {}, requestEvent);
  const json = await res.json();
  return getTripByIdResponseSchema.parse(json);
});
