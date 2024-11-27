import { component$ } from "@builder.io/qwik";
import { routeLoader$, z } from "@builder.io/qwik-city";
import {
  useForm,
  zodForm$,
  formAction$,
  type InitialValues,
} from "@modular-forms/qwik";
import { serverFetch } from "~/helpers/server-fetch";
import { TextInput } from "~/components";
import { errorResponseSchema } from "~/types/server-errors";
import { mapServerErrors } from "~/helpers/map-server-errors";

const newTripFormSchema = z.object({
  name: z
    .string()
    .min(1, "Name is required")
    .max(30, "Name must be less than 30 characters"),
  description: z
    .string()
    .min(1, "Description is required")
    .max(500, "Description must be less than 500 characters"),
});

type NewTripForm = z.infer<typeof newTripFormSchema>;

export const useCreateTripAction = formAction$<NewTripForm>(
  async (formData, request) => {
    const { name, description } = formData;

    const res = await serverFetch(
      "/trips",
      {
        method: "POST",
        body: JSON.stringify({
          name,
          description,
        }),
      },
      request,
    );

    const json = await res.json();

    if (!res.ok) {
      const { errors } = errorResponseSchema.parse(await res.json());
      const { messages, fields } = mapServerErrors({
        errors,
        messages: {},
        fields: {
          name: {},
          description: {},
        },
      });

      return {
        message: messages[0],
        errors: {
          name: fields.name[0],
          description: fields.description[0],
        },
      };
    }

    const trip = json.trip;
    throw request.redirect(303, `/trips/${trip.id}`);
  },
  zodForm$(newTripFormSchema),
);

export const useNewTripFormLoader = routeLoader$<InitialValues<NewTripForm>>(
  () => ({
    name: "",
    description: "",
  }),
);

export default component$(() => {
  const [newTripForm, { Form, Field }] = useForm<NewTripForm>({
    loader: useNewTripFormLoader(),
    action: useCreateTripAction(),
    validate: zodForm$(newTripFormSchema),
    validateOn: "blur",
    revalidateOn: "blur",
  });

  return (
    <div class="container mx-auto py-8">
      <div class="card bg-base-200 shadow-lg">
        <div class="card-body">
          <h1 class="card-title mb-4 text-4xl font-bold">New Trip</h1>
          {newTripForm.response.message && (
            <p class="my-2 text-error">{newTripForm.response.message}</p>
          )}
          <Form class="flex flex-col gap-4 space-y-4">
            <Field name="name">
              {(field, props) => (
                <div class="form-control">
                  <TextInput
                    id="name"
                    name="name"
                    label="Trip Name"
                    placeholder="Enter the trip name"
                    field={field}
                    fieldProps={props}
                  />
                </div>
              )}
            </Field>
            <Field name="description">
              {(field, props) => (
                <div class="form-control">
                  <TextInput
                    id="description"
                    name="description"
                    label="Description"
                    placeholder="Enter the trip description"
                    field={field}
                    fieldProps={props}
                    type="textarea"
                  />
                </div>
              )}
            </Field>
            <div class="form-control mt-6">
              <button type="submit" class="btn btn-primary w-full">
                Create Trip
              </button>
            </div>
          </Form>
        </div>
      </div>
    </div>
  );
});
