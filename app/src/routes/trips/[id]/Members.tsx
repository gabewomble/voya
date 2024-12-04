import { component$, useSignal, useTask$ } from "@builder.io/qwik";
import { useTripData } from "./layout";
import { Card, CardTitle, TextInput } from "~/components";
import { server$, z } from "@builder.io/qwik-city";
import { formAction$, useForm, zodForm$ } from "@modular-forms/qwik";
import { useAddMemberLoader } from "./layout";
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

export const useAddMember = formAction$<AddMemberForm>(
  // eslint-disable-next-line @typescript-eslint/no-unused-vars
  async (data, request) => {},
  zodForm$(addMemberFormSchema),
);

export const searchUsers = server$(async function (identifier: string) {
  if (!addMemberFormSchema.safeParse({ identifier }).success) {
    return [];
  }
  const requestEvent = this;
  const response = await serverFetch(
    "/users/search",
    {
      method: "POST",
      body: JSON.stringify({ identifier, limit: 5 }),
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
  const { members } = useTripData().value;
  const [addMemberForm, { Form, Field }] = useForm({
    loader: useAddMemberLoader(),
    action: useAddMember(),
    validateOn: "submit",
    validate: zodForm$(addMemberFormSchema),
  });

  const searchTimeoutId = useSignal<number>();
  const userSuggestions = useSignal<User[]>([]);

  useTask$(async (ctx) => {
    const id = ctx.track(() => addMemberForm.internal.fields.identifier?.value);

    if (id !== undefined && !isServer) {
      window.clearTimeout(searchTimeoutId.value);
      searchTimeoutId.value = window.setTimeout(async () => {
        window.clearTimeout(searchTimeoutId.value);
        userSuggestions.value = await searchUsers(id);
      }, 250);
    }
  });

  return (
    <Card>
      <CardTitle level={2}>Members</CardTitle>
      <Form class="flex flex-col gap-4">
        <div class="dropdown w-full">
          <Field name="identifier">
            {(field, props) => (
              <TextInput
                autocomplete="off"
                id="identifier"
                label="Search for member"
                name="identifier"
                placeholder="Enter username, name, or email"
                field={field}
                fieldProps={props}
              />
            )}
          </Field>
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
                      <button
                        class="btn btn-primary btn-sm"
                        type="button"
                        // onClick={() => {
                        //   addMemberForm.internal.fields.identifier?.setValue("");
                        //   addMemberForm.submit({ identifier: user.id });
                        // }}
                      >
                        Add
                      </button>
                    </td>
                  </tr>
                ))}
              </tbody>
            </table>
          )}
        </div>
      </Form>
      {members.length > 0 ? (
        <ul class="list-disc pl-5">
          {members.map((member) => (
            <li key={member.id} class="py-2">
              {member.name} ({member.email})
            </li>
          ))}
        </ul>
      ) : (
        <div class="py-4">
          <p>No members yet.</p>
        </div>
      )}
    </Card>
  );
});
