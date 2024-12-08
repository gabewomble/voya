import type { Signal } from "@builder.io/qwik";
import { $, component$ } from "@builder.io/qwik";
import { useTripData, useInviteUser } from "./layout";
import type { User } from "~/types/user";

export const UserSuggestionsTable = component$(
  ({ userSuggestions }: { userSuggestions: Signal<User[]> }) => {
    if (!userSuggestions.value.length) {
      return null;
    }

    const { trip } = useTripData().value;
    const inviteUser = useInviteUser();

    const onInviteClick = $(async function (user: User) {
      const result = await inviteUser.submit({
        userID: user.id,
        tripID: trip.id,
      });

      if (result.status === 200) {
        userSuggestions.value = userSuggestions.value.filter(
          (u) => u.id !== user.id,
        );
      }
    });

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
                  onClick$={() => onInviteClick(user)}
                  class="btn btn-outline btn-primary btn-sm"
                  type="button"
                >
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
