import { component$ } from "@builder.io/qwik";
import { routeLoader$, z } from "@builder.io/qwik-city";
import { serverFetch } from "~/helpers/server-fetch";
import {
  zodForm$,
  useForm,
  formAction$,
  type InitialValues,
} from "@modular-forms/qwik";
import { errorResponseSchema } from "~/types/server-errors";
import { TextInput } from "~/components";
import { mapServerErrors } from "~/helpers/map-server-errors";
import { type User } from "~/types/user";

const updateProfileFormSchema = z.object({
  username: z
    .string()
    .min(1, "Username is required")
    .max(30, "Username must be less than 30 characters"),
  name: z
    .string()
    .min(1, "Name is required")
    .max(30, "Name must be less than 30 characters"),
});

type UpdateProfileForm = z.infer<typeof updateProfileFormSchema>;

export const useUpdateProfileAction = formAction$<UpdateProfileForm>(
  async (data, request) => {
    const user = request.sharedMap.get("user") as User;
    const res = await serverFetch(
      `/users/${user.username}`,
      {
        method: "PATCH",
        headers: {
          "Content-Type": "application/json",
        },
        body: JSON.stringify(data),
      },
      request,
    );

    if (!res.ok) {
      const { errors } = errorResponseSchema.parse(await res.json());

      const { messages, fields } = mapServerErrors({
        errors,
        messages: {},
        fields: {
          username: {
            "duplicate username": "Username is already taken",
          },
          name: {},
        },
      });

      return {
        message: messages[0],
        errors: {
          username: fields.username[0],
          name: fields.name[0],
        },
      };
    }
    throw request.redirect(303, "/settings/profile");
  },
  zodForm$(updateProfileFormSchema),
);

export const useFormLoader = routeLoader$<InitialValues<UpdateProfileForm>>(
  (ctx) => {
    const user = ctx.sharedMap.get("user") as User;

    return {
      username: user.username,
      name: user.name,
    };
  },
);

export default component$(() => {
  const [profileForm, { Form, Field }] = useForm<UpdateProfileForm>({
    loader: useFormLoader(),
    action: useUpdateProfileAction(),
    validate: zodForm$(updateProfileFormSchema),
    revalidateOn: "blur",
  });

  return (
    <div class="container mx-auto py-8">
      <div class="card bg-base-200 shadow-lg">
        <div class="card-body">
          <h1 class="card-title text-4xl font-bold">Edit Profile</h1>
          <p class="py-4">Update your profile information</p>
          <Form class="flex flex-col gap-4">
            {profileForm.response.message && (
              <p class="my-2 text-error">{profileForm.response.message}</p>
            )}
            <div class="form-control">
              <Field name="username">
                {(field, props) => (
                  <TextInput
                    id="username"
                    label="Username"
                    name="username"
                    placeholder="Enter your username"
                    field={field}
                    fieldProps={props}
                  />
                )}
              </Field>
            </div>
            <div class="form-control">
              <Field name="name">
                {(field, props) => (
                  <TextInput
                    id="name"
                    label="Name"
                    name="name"
                    placeholder="Enter your name"
                    field={field}
                    fieldProps={props}
                  />
                )}
              </Field>
            </div>
            <div class="flex justify-end gap-4">
              <a href="/settings/profile" class="btn btn-secondary">
                Cancel
              </a>
              <button type="submit" class="btn btn-primary">
                Save Changes
              </button>
            </div>
          </Form>
        </div>
      </div>
    </div>
  );
});
