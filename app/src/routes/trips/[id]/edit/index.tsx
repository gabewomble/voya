import { Link } from "@builder.io/qwik-city";
import { useTripData } from "../layout";
import { component$ } from "@builder.io/qwik";

export default component$(() => {
  const { trip } = useTripData().value;

  return (
    <div class="container mx-auto py-8">
      <div class="card bg-base-200 shadow-lg">
        <div class="card-body">
          <h1 class="card-title text-4xl font-bold">Editing: {trip.name}</h1>
          <form class="mt-4 flex flex-col gap-4">
            <label class="form-control w-full">
              <div class="label">
                <span class="label-text">Name</span>
              </div>
              <input
                type="text"
                placeholder="A name for your trip"
                value={trip.name}
                class="input input-bordered w-full"
              />
            </label>
            <label class="form-control w-full">
              <div class="label">
                <span class="label-text">Description</span>
              </div>
              <textarea
                class="textarea textarea-bordered h-24"
                placeholder="A description for your trip"
                value={trip.description}
              ></textarea>
            </label>
            <div class="card-actions justify-end">
              <button class="btn btn-primary" type="submit">
                Save changes
              </button>
              <Link href={`/trips/${trip.id}`} class="btn btn-secondary">
                Cancel
              </Link>
            </div>
          </form>
        </div>
      </div>
    </div>
  );
});
