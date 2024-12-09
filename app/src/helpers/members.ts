import type { Member } from "~/types/members";

export function getCanMemberEdit(member?: Member): boolean {
  return (
    member?.member_status === "owner" || member?.member_status === "accepted"
  );
}
