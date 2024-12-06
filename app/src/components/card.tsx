import { component$, Slot } from "@builder.io/qwik";

export const Card = component$(() => (
  <div class="card bg-base-200 shadow-lg">
    <div class="card-body">
      <Slot />
    </div>
  </div>
));

export const CardTitle = component$(({ level = 3 }: { level?: 1 | 2 | 3 }) => {
  if (level === 1) {
    return (
      <h1 class="card-title text-4xl font-bold">
        <Slot />
      </h1>
    );
  }

  if (level === 2) {
    return (
      <h2 class="card-title text-2xl font-bold">
        <Slot />
      </h2>
    );
  }

  return (
    <h3 class="card-title text-2xl font-bold">
      <Slot />
    </h3>
  );
});
