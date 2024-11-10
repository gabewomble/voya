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

  const isLoggedIn = !!user.value;

  return (
    <nav class="navbar bg-base-300 shadow-lg">
      <div class="container mx-auto flex items-center justify-between px-4">
        <a class="text-2xl font-bold text-primary" href="/">
          Voya
        </a>
        <div class="flex items-center gap-4 space-x-4">
          {isLoggedIn ? (
            <>
              <a class="btn btn-ghost" href="/trips">
                Trips
              </a>
              <Form action={logout} class="contents">
                <button class="btn btn-ghost" type="submit">
                  Sign out
                </button>
              </Form>
            </>
          ) : (
            <a class="btn btn-primary" href="/login">
              Login
            </a>
          )}
        </div>
      </div>
    </nav>
  );
});
