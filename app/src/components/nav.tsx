import { component$ } from "@builder.io/qwik";

export const Nav = component$(() => {
  return (
    <nav class="navbar bg-base-300">
      <a class="btn btn-ghost text-xl" href="/trips">
        Voya
      </a>
    </nav>
  );
});
