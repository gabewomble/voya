import { component$, useContextProvider } from "@builder.io/qwik";
import { Link, type DocumentHead } from "@builder.io/qwik-city";
import {
  CurrentMemberContext,
  useTripData,
  useUpdateMemberStatus,
} from "./layout";
import { Card, CardTitle } from "~/components";
import { Members } from "./Members";
import { useUserData } from "~/routes/layout";
import { getCanMemberEdit } from "~/helpers/members";
import { memberStatusEnum } from "~/types/members";

export default component$(() => {
  const data = useTripData();
  const currentUser = useUserData().value;
  const { trip, members } = data.value;
  const currentMember = members.find(
    (member) => member.id === currentUser?.id,
  )!;

  useContextProvider(CurrentMemberContext, currentMember);

  const canEdit = getCanMemberEdit(currentMember);
  const isPending = currentMember.member_status === "pending";
  const updateMemberStatus = useUpdateMemberStatus();

  return (
    <div class="container mx-auto flex flex-col gap-8 py-8">
      <Card>
        <CardTitle level={1}>{trip.name}</CardTitle>
        <p class="py-4">{trip.description}</p>

        <div class="card-actions justify-end">
          {canEdit && (
            <>
              <Link href={`/trips/${trip.id}/edit`} class="btn btn-secondary">
                Edit this trip
              </Link>
              <Link href="/trips" class="btn btn-primary">
                Back to trips
              </Link>
            </>
          )}
          {isPending && (
            <>
              <button
                class="btn btn-outline btn-secondary btn-sm"
                onClick$={() => {
                  updateMemberStatus.submit({
                    tripID: trip.id,
                    userID: currentMember.id,
                    memberStatus: memberStatusEnum.Values.declined,
                  });
                }}
              >
                Decline invite
              </button>
              <button
                class="btn btn-outline btn-primary btn-sm"
                onClick$={() => {
                  updateMemberStatus.submit({
                    tripID: trip.id,
                    userID: currentMember.id,
                    memberStatus: memberStatusEnum.Values.accepted,
                  });
                }}
              >
                Join this trip
              </button>
            </>
          )}
        </div>
      </Card>

      <Members />
    </div>
  );
});

export const head: DocumentHead = {
  title: "View Trip",
};
