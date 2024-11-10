import { component$ } from "@builder.io/qwik";
import { type DocumentHead } from "@builder.io/qwik-city";
import { useTripData } from "./layout";

export default component$(() => {
  const trip = useTripData();

  return (
    <div class="container mx-auto py-8">
      <div class="card bg-base-200 shadow-lg">
        <div class="card-body">
          <h1 class="card-title text-4xl font-bold">{trip.value.name}</h1>
          <p class="py-4">{trip.value.description}</p>
          <div class="card-actions justify-end">
            <a href={`/trips/${trip.value.id}/edit`} class="btn btn-secondary">
              Edit this trip
            </a>
            <a href="/trips" class="btn btn-primary">
              Back to trips
            </a>
          </div>
        </div>
      </div>
    </div>
  );
});

export const head: DocumentHead = {
  title: "View Trip",
};
