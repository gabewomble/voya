import { component$ } from "@builder.io/qwik";
import { Link, type DocumentHead } from "@builder.io/qwik-city";
import { useTripData } from "./layout";
import { Card, CardTitle } from "~/components";
import { Members } from "./Members";

export default component$(() => {
  const data = useTripData();
  const { trip } = data.value;

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

      <Members />
    </div>
  );
});

export const head: DocumentHead = {
  title: "View Trip",
};
