import { component$, useContext } from "@builder.io/qwik";
import { useCancelInviteUser } from "./layout";
import {
  type Member,
  type MemberStatusEnum,
  memberStatusEnum,
} from "~/types/members";
import { UserContext } from "~/context/user";
import { HiXCircleOutline } from "@qwikest/icons/heroicons";
import type { User } from "~/types/users";
import type { Trip } from "~/types/trips";

function capitalizeFirstLetter(string: string) {
  return string.charAt(0).toUpperCase() + string.slice(1);
}

const MemberStatus = component$(({ status }: { status: MemberStatusEnum }) => {
  const statusString = capitalizeFirstLetter(status);
  if (
    status === memberStatusEnum.Values.owner ||
    status === memberStatusEnum.Values.accepted
  ) {
    return <span class="badge badge-primary">{statusString}</span>;
  }

  if (status === memberStatusEnum.Values.pending) {
    return <span class="badge badge-neutral">{statusString}</span>;
  }
  return <span class="badge">{statusString}</span>;
});

const Actions = component$(
  ({
    canEdit,
    // currentUser,
    member,
    tripID,
  }: {
    canEdit: boolean;
    currentUser: User | null;
    member: Member;
    tripID: string;
  }) => {
    const actions = [];
    const cancelInvite = useCancelInviteUser();

    if (member.member_status === memberStatusEnum.Values.pending) {
      actions.push(
        <div class="tooltip tooltip-error" data-tip="Cancel invite">
          <button
            class="hover:text-error"
            type="button"
            disabled={!canEdit}
            onClick$={() => {
              cancelInvite.submit({
                userID: member.id,
                tripID: tripID,
              });
            }}
          >
            <HiXCircleOutline class="h-6 w-6" />
          </button>
        </div>,
      );
    }

    return actions;
  },
);

export const MembersTable = component$(
  ({ trip, members }: { trip: Trip; members: Member[] }) => {
    const currentUser = useContext(UserContext).value;

    const currentUserCanEdit = members.some((member) => {
      const memberCanEdit =
        member.member_status === memberStatusEnum.Values.owner ||
        member.member_status === memberStatusEnum.Values.accepted;
      return memberCanEdit && member.id === currentUser?.id;
    });

    const membersToRender = members.filter((member) => {
      return (
        member.member_status === memberStatusEnum.Values.owner ||
        member.member_status === memberStatusEnum.Values.accepted ||
        member.member_status === memberStatusEnum.Values.pending
      );
    });

    return (
      <table class="table">
        <caption>Trip Members</caption>
        <thead>
          <tr>
            <th>Name</th>
            <th>Email</th>
            <th>Status</th>
            <th>Action</th>
          </tr>
        </thead>
        <tbody>
          {membersToRender.map((member) => (
            <tr key={member.id}>
              <td>{member.name}</td>
              <td>
                <a class="link-hover" href={`mailto:${member.email}`}>
                  {member.email}
                </a>
              </td>
              <td>
                <MemberStatus status={member.member_status} />
              </td>
              <td>
                <Actions
                  canEdit={currentUserCanEdit}
                  currentUser={currentUser}
                  member={member}
                  tripID={trip.id}
                />
              </td>
            </tr>
          ))}
        </tbody>
      </table>
    );
  },
);
