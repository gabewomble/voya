import { component$ } from "@builder.io/qwik";
import { Form, routeAction$, zod$, z } from "@builder.io/qwik-city";
import { serverFetch } from "~/helpers/server-fetch";

export const useCreateTripAction = routeAction$(
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
      request.fail(res.status, json);
      return;
    }

    const trip = json.trip;

    throw request.redirect(303, `/trips/${trip.id}`);
  },
  zod$({
    name: z.string().min(1, "Name is required").max(255),
    description: z.string().min(1, "Description is required").max(500),
  }),
);

export default component$(() => {
  const createTrip = useCreateTripAction();

  return (
    <div class="container mx-auto py-8">
      <div class="card bg-base-200 shadow-lg">
        <div class="card-body">
          <h1 class="card-title mb-4 text-4xl font-bold">New Trip</h1>
          <Form action={createTrip} class="flex flex-col gap-4 space-y-4">
            <div class="form-control">
              <label for="name" class="label">
                <span class="label-text">Trip Name:</span>
              </label>
              <input
                type="text"
                id="name"
                name="name"
                class="input input-bordered w-full"
                placeholder="Enter the trip name"
                required
              />
            </div>
            <div class="form-control">
              <label for="description" class="label">
                <span class="label-text">Trip Description:</span>
              </label>
              <textarea
                id="description"
                name="description"
                class="textarea textarea-bordered w-full"
                placeholder="Enter the trip description"
              ></textarea>
            </div>
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
