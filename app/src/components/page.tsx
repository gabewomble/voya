import { component$, Slot } from "@builder.io/qwik";
import { Nav } from "./nav";

export const Page = component$(() => {
  return (
    <div class="flex flex-col">
      <Nav />
      <Slot />
    </div>
  );
});
