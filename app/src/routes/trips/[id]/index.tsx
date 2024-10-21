import { component$ } from "@builder.io/qwik";
import { type DocumentHead } from "@builder.io/qwik-city";
import { useTripData } from "./layout";

export default component$(() => {
  const trip = useTripData();

  return (
    <>
      <h1>{trip.value.name}</h1>
      <p>{trip.value.description}</p>
      <nav class="flex gap-4">
        <a href={`/trips/${trip.value.id}/edit`}>Edit this trip</a>
        <a href="/trips">Back to trips</a>
      </nav>
    </>
  );
});

export const head: DocumentHead = {
  title: "View Trip",
};
