import { component$ } from "@builder.io/qwik";

export const Nav = component$(() => {
  return (
    <nav class="navbar bg-base-300">
      <div class="mx-auto flex w-full max-w-screen-2xl items-center justify-between">
        <a class="link text-xl no-underline" href="/trips">
          Voya
        </a>
        <a class="link" href="/login">
          Login
        </a>
      </div>
    </nav>
  );
});
