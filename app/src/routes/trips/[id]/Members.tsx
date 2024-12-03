import { component$, useSignal, useTask$ } from "@builder.io/qwik";
import { useTripData } from "./layout";
import { Card, CardTitle, TextInput } from "~/components";
import { server$, z } from "@builder.io/qwik-city";
import { formAction$, useForm, zodForm$ } from "@modular-forms/qwik";
import { useAddMemberLoader } from "./layout";
import { isServer } from "@builder.io/qwik/build";

const addMemberFormSchema = z.object({
  identifier: z.string().min(1, "Name or email must be provided"),
});

export type AddMemberForm = z.infer<typeof addMemberFormSchema>;

export const useAddMember = formAction$<AddMemberForm>(
  // eslint-disable-next-line @typescript-eslint/no-unused-vars
  async (data, request) => {},
  zodForm$(addMemberFormSchema),
);

export const searchUsers = server$(async (identifier: string) => {
  if (!addMemberFormSchema.safeParse({ identifier }).success) {
    return [];
  }
  // TODO: implement the search
  return [];
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
  const userSuggestions = useSignal([]);

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
      <Form>
        <Field name="identifier">
          {(field, props) => (
            <TextInput
              id="identifier"
              label="Name or email"
              name="identifier"
              placeholder="Enter name or email"
              field={field}
              fieldProps={props}
            />
          )}
        </Field>
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
