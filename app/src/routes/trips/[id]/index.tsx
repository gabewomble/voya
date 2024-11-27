import { component$ } from "@builder.io/qwik";
import { type DocumentHead } from "@builder.io/qwik-city";
import { useTripData } from "./layout";

export default component$(() => {
  const data = useTripData();
  const { trip, members } = data.value;

  return (
    <div class="container mx-auto py-8">
      <div class="card mb-8 bg-base-200 shadow-lg">
        <div class="card-body">
          <h1 class="card-title text-4xl font-bold">{trip.name}</h1>
          <p class="py-4">{trip.description}</p>

          <div class="card-actions justify-end">
            <a href={`/trips/${trip.id}/edit`} class="btn btn-secondary">
              Edit this trip
            </a>
            <a href="/trips" class="btn btn-primary">
              Back to trips
            </a>
          </div>
        </div>
      </div>

      <div class="card bg-base-200 shadow-lg">
        <div class="card-body">
          <h2 class="card-title text-2xl font-bold">Members</h2>
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
              <a
                href={`/trips/${trip.id}/members/add`}
                class="btn btn-primary mt-4"
              >
                Add Members
              </a>
            </div>
          )}
        </div>
      </div>
    </div>
  );
});

export const head: DocumentHead = {
  title: "View Trip",
};
