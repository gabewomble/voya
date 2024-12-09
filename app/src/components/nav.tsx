import { component$, useContext } from "@builder.io/qwik";
import { Form, Link, globalAction$, zod$ } from "@builder.io/qwik-city";
import { ActivityContext } from "~/context/activity";
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
  const activity = useContext(ActivityContext);
  const logout = useLogout();

  const isLoggedIn = !!user.value;
  const unreadActivityCount = activity.unreadCount;

  return (
    <nav class="navbar bg-base-300 shadow-lg">
      <div class="container mx-auto flex items-center justify-between px-4">
        <Link class="text-2xl font-bold text-primary" href="/">
          Voya
        </Link>
        <div class="flex items-center gap-4 space-x-4">
          {isLoggedIn ? (
            <>
              <div class="indicator">
                {unreadActivityCount > 0 && (
                  <span class="badge indicator-item badge-secondary">
                    {unreadActivityCount}
                  </span>
                )}
                <a class="link-hover link p-1" href="/activity">
                  Activity
                </a>
              </div>
              <Link class="link-hover link p-1" href="/trips">
                My Trips
              </Link>
              <div class="dropdown dropdown-end">
                <div
                  tabIndex={0}
                  role="button"
                  class="avatar btn btn-circle btn-ghost"
                >
                  <div class="w-10 rounded-full">
                    <img
                      alt="Tailwind CSS Navbar component"
                      src="https://img.daisyui.com/images/stock/photo-1534528741775-53994a69daeb.webp"
                      height={40}
                      width={40}
                    />
                  </div>
                </div>
                <ul
                  tabIndex={0}
                  class="menu dropdown-content menu-sm z-[1] mt-3 w-52 rounded-box bg-base-100 p-2 shadow"
                >
                  <li>
                    <Link href="/settings/profile">Profile</Link>
                  </li>
                  <Form action={logout} class="contents">
                    <li>
                      <button type="submit">Sign out</button>
                    </li>
                  </Form>
                </ul>
              </div>
            </>
          ) : (
            <Link class="btn btn-primary" href="/login">
              Login
            </Link>
          )}
        </div>
      </div>
    </nav>
  );
});
