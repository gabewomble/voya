import { routeAction$, routeLoader$, z, zod$ } from "@builder.io/qwik-city";
import { serverFetch } from "~/helpers/server-fetch";
import { getTripByIdResponseSchema } from "~/types/api";
import type { Member } from "~/types/members";
import { memberStatusEnum } from "~/types/members";
import type { InitialValues } from "@modular-forms/qwik";
import type { AddMemberForm } from "./Members";
import { createContextId } from "@builder.io/qwik";

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

export const CurrentMemberContext = createContextId<Member | undefined>(
  "CurrentMember",
);

export const useUpdateMemberStatus = routeAction$(
  async (data, requestEvent) => {
    const response = await serverFetch(
      `/trip/${data.tripID}/members`,
      {
        method: "PATCH",
        body: JSON.stringify({
          user_id: data.userID,
          member_status: data.memberStatus,
        }),
      },
      requestEvent,
    );

    if (!response.ok) {
      return requestEvent.fail(500, {
        error: "Failed to update member status",
      });
    }

    return response.ok;
  },
  zod$(
    z.object({
      tripID: z.string().uuid(),
      userID: z.string().uuid(),
      memberStatus: memberStatusEnum,
    }),
  ),
);
