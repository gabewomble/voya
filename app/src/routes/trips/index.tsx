import { component$ } from "@builder.io/qwik";
import { routeLoader$, type DocumentHead } from "@builder.io/qwik-city";
import { serverFetch } from "~/helpers/server-fetch";
import type { Trip } from "~/types/trips";

export const useGetTrips = routeLoader$(async (request) => {
  const res = await serverFetch("/trips", {}, request);
  const json = await res.json();
  return (json?.trips ?? []) as Trip[];
});

export default component$(() => {
  const trips = useGetTrips();
  return (
    <>
      <h1>Trips</h1>
      <ul>
        {trips.value.map((trip) => {
          return (
            <li key={trip.id}>
              <a href={`/trips/${trip.id}`}>
                <p>
                  {trip.name}
                  {" - "}
                  <span>{trip.description}</span>
                </p>
              </a>
            </li>
          );
        })}
      </ul>
    </>
  );
});

export const head: DocumentHead = {
  title: "Trips",
};
