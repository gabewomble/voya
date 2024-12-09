import { component$, useContext } from "@builder.io/qwik";
import { CurrentMemberContext, useUpdateMemberStatus } from "./layout";
import {
  type Member,
  type MemberStatusEnum,
  memberStatusEnum,
} from "~/types/members";
import { HiUserMinusOutline, HiXCircleOutline } from "@qwikest/icons/heroicons";
import type { Trip } from "~/types/trips";
import { getCanMemberEdit } from "~/helpers/members";
import { Modal, ModalActions, ModalTitle, showModal } from "~/components/modal";

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

const RemoveMemberAction = component$<{ member: Member; tripID: string }>(
  ({ member, tripID }) => {
    const modalId = `remove-${member.id}`;
    const updateStatus = useUpdateMemberStatus();
    const currentMember = useContext(CurrentMemberContext);
    const isSelf = currentMember?.id === member.id;

    return (
      <div class="tooltip" data-tip={isSelf ? "Leave trip" : "Remove member"}>
        <button
          class="hover:text-error"
          type="button"
          onClick$={() => showModal(modalId)}
        >
          <HiUserMinusOutline class="h-6 w-6" />
        </button>
        <Modal id={modalId}>
          <ModalTitle>Confirm</ModalTitle>
          {isSelf ? (
            <p>Are you sure you want to leave the trip?</p>
          ) : (
            <p>Are you sure you want to remove {member.name} from the trip?</p>
          )}
          <ModalActions q:slot="actions">
            <button class="btn btn-outline">Go back</button>
            <button
              class="btn btn-outline btn-error"
              onClick$={() => {
                updateStatus.submit({
                  userID: member.id,
                  tripID: tripID,
                  memberStatus: memberStatusEnum.Values.removed,
                });
              }}
            >
              {isSelf ? "Leave trip" : "Remove member"}
            </button>
          </ModalActions>
        </Modal>
      </div>
    );
  },
);

const CancelInviteAction = component$<{ member: Member; tripID: string }>(
  ({ member, tripID }) => {
    const updateStatus = useUpdateMemberStatus();
    const modalId = `cancel-${member.id}`;
    return (
      <div class="tooltip" data-tip="Cancel invite">
        <button
          class="hover:text-error"
          type="button"
          onClick$={() => showModal(modalId)}
        >
          <HiXCircleOutline class="h-6 w-6" />
        </button>
        <Modal id={modalId}>
          <ModalTitle>Confirm</ModalTitle>
          <p>Are you sure you want to cancel {member.name}'s invitation?</p>
          <ModalActions q:slot="actions">
            <button class="btn btn-outline">Go back</button>
            <button
              class="btn btn-outline btn-error"
              onClick$={() => {
                updateStatus.submit({
                  userID: member.id,
                  tripID: tripID,
                  memberStatus: memberStatusEnum.Values.cancelled,
                });
              }}
            >
              Cancel invitation
            </button>
          </ModalActions>
        </Modal>
      </div>
    );
  },
);

const Actions = component$(
  ({ member, tripID }: { member: Member; tripID: string }) => {
    const actions = [];

    if (member.member_status === memberStatusEnum.Values.accepted) {
      actions.push(<RemoveMemberAction member={member} tripID={tripID} />);
    }

    if (member.member_status === memberStatusEnum.Values.pending) {
      actions.push(<CancelInviteAction member={member} tripID={tripID} />);
    }

    return actions;
  },
);

export const MembersTable = component$(
  ({ trip, members }: { trip: Trip; members: Member[] }) => {
    const currentMember = useContext(CurrentMemberContext);
    const canEdit = getCanMemberEdit(currentMember);

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
              <td>{canEdit && <Actions member={member} tripID={trip.id} />}</td>
            </tr>
          ))}
        </tbody>
      </table>
    );
  },
);
