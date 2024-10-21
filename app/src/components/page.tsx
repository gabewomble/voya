import { component$, Slot } from "@builder.io/qwik";
import { Nav } from "./nav";

export const Page = component$(() => {
  return (
    <div class="flex flex-col gap-4">
      <Nav />
      <div class="flex flex-col gap-4 px-8">
        <Slot />
      </div>
    </div>
  );
});
