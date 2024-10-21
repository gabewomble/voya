import { useTripData } from "../layout";
import { component$ } from "@builder.io/qwik";

export default component$(() => {
  const trip = useTripData();

  return (
    <>
      <h1>Editing: {trip.value.name}</h1>
      <form class="flex flex-col gap-4">
        <label>
          Name
          <input type="text" name="name" />
        </label>
        <label>
          Description
          <textarea name="name" />
        </label>
      </form>
    </>
  );
});
