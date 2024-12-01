import { component$ } from "@builder.io/qwik";
import { Link, type DocumentHead } from "@builder.io/qwik-city";
import { useTripData } from "./layout";
import { Card, CardTitle } from "~/components";

export default component$(() => {
  const data = useTripData();
  const { trip, members } = data.value;

  return (
    <div class="container mx-auto flex flex-col gap-8 py-8">
      <Card>
        <CardTitle level={1}>{trip.name}</CardTitle>
        <p class="py-4">{trip.description}</p>

        <div class="card-actions justify-end">
          <Link href={`/trips/${trip.id}/edit`} class="btn btn-secondary">
            Edit this trip
          </Link>
          <Link href="/trips" class="btn btn-primary">
            Back to trips
          </Link>
        </div>
      </Card>

      <Card>
        <CardTitle level={2}>Members</CardTitle>
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
            <Link
              href={`/trips/${trip.id}/members/add`}
              class="btn btn-primary mt-4"
            >
              Add Members
            </Link>
          </div>
        )}
      </Card>
    </div>
  );
});

export const head: DocumentHead = {
  title: "View Trip",
};
