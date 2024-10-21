import { useTripData } from "../layout";
import { component$ } from "@builder.io/qwik";

export default component$(() => {
  const trip = useTripData();

  return (
    <>
      <h1>Editing: {trip.value.name}</h1>
      <form class="flex flex-col gap-4">
        <label class="form-control w-full max-w-xs">
          <div class="label">
            <span class="label-text">Name</span>
          </div>
          <input
            type="text"
            placeholder="A name for your trip"
            value={trip.value.name}
            class="input input-bordered w-full max-w-xs"
          />
        </label>
        <label class="form-control">
          <div class="label">
            <span class="label-text">Description</span>
          </div>
          <textarea
            class="textarea textarea-bordered h-24"
            placeholder="A description for your trip"
            value={trip.value.description}
          ></textarea>
        </label>
        <button class="btn btn-primary" type="submit">
          Save changes
        </button>
      </form>
    </>
  );
});
