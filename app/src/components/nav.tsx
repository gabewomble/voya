import { component$, useContext } from "@builder.io/qwik";
import { Form, globalAction$, zod$ } from "@builder.io/qwik-city";
import { UserContext } from "~/context/user";
import { serverFetch } from "~/helpers/server-fetch";

export const useLogout = globalAction$(async (_, request) => {
  const token = request.cookie.get("token")?.value;

  if (token) {
    await serverFetch(
      "/tokens/current",
      {
        method: "DELETE",
      },
      request,
    );

    request.cookie.delete("token");
    request.sharedMap.delete("user");
  }

  throw request.redirect(303, "/");
}, zod$({}));

export const Nav = component$(() => {
  const user = useContext(UserContext);
  const logout = useLogout();

  return (
    <nav class="navbar bg-base-300">
      <div class="mx-auto flex w-full max-w-screen-2xl items-center justify-between">
        <a class="link text-xl no-underline" href="/trips">
          Voya
        </a>
        {user.value ? (
          <Form action={logout} class="contents">
            <button class="link" type="submit">
              Sign out
            </button>
          </Form>
        ) : (
          <a class="link" href="/login">
            Login
          </a>
        )}
      </div>
    </nav>
  );
});
