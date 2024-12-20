import type { Signal } from "@builder.io/qwik";
import { component$ } from "@builder.io/qwik";
import { useInviteUser } from "./layout";
import type { User } from "~/types/users";
import type { Trip } from "~/types/trips";
import type { Member } from "~/types/members";
import { HiUserPlusOutline } from "@qwikest/icons/heroicons";

export const UserSuggestionsTable = component$(
  ({
    trip,
    userSuggestions,
  }: {
    trip: Trip;
    members: Member[];
    userSuggestions: Signal<User[]>;
  }) => {
    if (!userSuggestions.value.length) {
      return null;
    }

    const inviteUser = useInviteUser();

    return (
      <table class="table">
        <caption>Search results</caption>
        <thead>
          <tr>
            <th>Name</th>
            <th>Email</th>
            <th></th>
          </tr>
        </thead>
        <tbody>
          {userSuggestions.value.map((user) => (
            <tr key={user.id}>
              <td>{user.name}</td>
              <td>{user.email}</td>
              <td>
                <button
                  onClick$={() => {
                    inviteUser.submit({
                      userID: user.id,
                      tripID: trip.id,
                    });
                    userSuggestions.value = userSuggestions.value.filter(
                      (u) => u.id !== user.id,
                    );
                  }}
                  class="btn btn-outline btn-primary btn-sm"
                  type="button"
                >
                  <HiUserPlusOutline class="h-4 w-4" />
                  Invite
                </button>
              </td>
            </tr>
          ))}
        </tbody>
      </table>
    );
  },
);
