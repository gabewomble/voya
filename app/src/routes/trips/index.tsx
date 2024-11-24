import { component$ } from "@builder.io/qwik";
import { routeLoader$, type DocumentHead } from "@builder.io/qwik-city";
import { serverFetch } from "~/helpers/server-fetch";
import { listTripsResponseSchema } from "~/types/api";

export const useGetTrips = routeLoader$(async (request) => {
  const res = await serverFetch("/trips", {}, request);
  const json = await res.json();
  return listTripsResponseSchema.parse(json);
});

export default component$(() => {
  const { trips } = useGetTrips().value;
  return (
    <>
      <div class="container mx-auto py-8">
        <h1 class="mb-8 text-center text-4xl font-bold">Your Trips</h1>
        <div class="mb-8 text-center">
          <a href="/trips/new" class="btn btn-primary">
            Create New Trip
          </a>
        </div>
        <div class="grid grid-cols-1 gap-8 md:grid-cols-2 lg:grid-cols-3">
          {trips.map((trip) => (
            <div key={trip.id} class="card bg-base-200 shadow-lg">
              <div class="card-body">
                <h2 class="card-title">{trip.name}</h2>
                <p>{trip.description}</p>
                <div class="card-actions justify-end">
                  <a href={`/trips/${trip.id}`} class="btn btn-primary">
                    View Details
                  </a>
                </div>
              </div>
            </div>
          ))}
        </div>
      </div>
    </>
  );
});

export const head: DocumentHead = {
  title: "Trips",
};
