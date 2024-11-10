import { component$, useContext } from "@builder.io/qwik";
import { UserContext } from "~/context/user";

export const Nav = component$(() => {
  const user = useContext(UserContext);
  return (
    <nav class="navbar bg-base-300">
      <div class="mx-auto flex w-full max-w-screen-2xl items-center justify-between">
        <a class="link text-xl no-underline" href="/trips">
          Voya
        </a>
        {user ? (
          <button>Sign out</button>
        ) : (
          <a class="link" href="/login">
            Login
          </a>
        )}
      </div>
    </nav>
  );
});
