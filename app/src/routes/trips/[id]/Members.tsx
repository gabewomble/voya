import { component$, useSignal, useTask$ } from "@builder.io/qwik";
import { useTripData } from "./layout";
import { Card, CardTitle } from "~/components";
import { server$, z } from "@builder.io/qwik-city";
import { isServer } from "@builder.io/qwik/build";
import { serverFetch } from "~/helpers/server-fetch";
import { searchUsersResponseSchema } from "~/types/api";
import type { User } from "~/types/user";

const addMemberFormSchema = z.object({
  identifier: z
    .string()
    .min(4, "Username or email must be at least 4 characters long."),
});

export type AddMemberForm = z.infer<typeof addMemberFormSchema>;

const searchUsersInputSchema = z.object({
  identifier: addMemberFormSchema.shape.identifier,
  tripID: z.string().uuid(),
});

type SearchUserInput = z.infer<typeof searchUsersInputSchema>;

export const searchUsers = server$(async function (input: SearchUserInput) {
  if (!searchUsersInputSchema.safeParse(input).success) {
    return [];
  }

  const requestEvent = this;
  const response = await serverFetch(
    "/users/search",
    {
      method: "POST",
      body: JSON.stringify({
        identifier: input.identifier,
        trip_id: input.tripID,
        limit: 5,
      }),
    },
    requestEvent,
  );

  if (!response.ok) return [];

  const data = await response.json();
  const result = searchUsersResponseSchema.safeParse(data);

  if (!result.success) return [];

  return result.data.users;
});

export const Members = component$(() => {
  const { trip, members } = useTripData().value;

  const userSearch = useSignal("");
  const searchTimeoutId = useSignal<number>();
  const userSuggestions = useSignal<User[]>([]);

  useTask$(async (ctx) => {
    const searchValue = ctx.track(userSearch);

    if (
      !isServer &&
      addMemberFormSchema.safeParse({ identifier: searchValue }).success
    ) {
      window.clearTimeout(searchTimeoutId.value);
      searchTimeoutId.value = window.setTimeout(async () => {
        window.clearTimeout(searchTimeoutId.value);
        userSuggestions.value = await searchUsers({
          identifier: searchValue,
          tripID: trip.id,
        });
      }, 250);
    } else {
      userSuggestions.value = [];
    }
  });

  return (
    <Card>
      <CardTitle level={2}>Members</CardTitle>
      <label
        class="mb-2 block text-sm font-bold text-base-content"
        for="search-input"
      >
        Search
      </label>
      <input
        autocomplete="off"
        class="input input-bordered w-full"
        type="text"
        id="search-input"
        placeholder="Enter username, name, or email"
        bind:value={userSearch}
      />
      {userSuggestions.value.length > 0 && (
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
                  <button class="btn btn-primary btn-sm" type="button">
                    Add
                  </button>
                </td>
              </tr>
            ))}
          </tbody>
        </table>
      )}
      <div class="divider" />
      {members.length > 0 ? (
        <table class="table">
          <caption>Trip Members</caption>
          <thead>
            <tr>
              <th>Name</th>
              <th>Email</th>
              <th>Status</th>
            </tr>
          </thead>
          <tbody>
            {members.map((user) => (
              <tr key={user.id}>
                <td>{user.name}</td>
                <td>{user.email}</td>
                <td>{user.member_status}</td>
              </tr>
            ))}
          </tbody>
        </table>
      ) : (
        <div class="py-4">
          <p>No members yet.</p>
        </div>
      )}
    </Card>
  );
});
