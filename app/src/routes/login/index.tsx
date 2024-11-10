import { component$ } from "@builder.io/qwik";
import { Form, routeAction$, z, zod$ } from "@builder.io/qwik-city";
import { serverFetch } from "~/helpers/server-fetch";

export const useLoginAction = routeAction$(
  async (data, request) => {
    const email = data.email;
    const password = data.password;

    const res = await serverFetch(
      "/tokens/authenticate",
      {
        method: "POST",
        headers: {
          "Content-Type": "application/json",
        },
        body: JSON.stringify({
          email,
          password,
        }),
      },
      request,
    );

    if (!res.ok) {
      request.fail(res.status, await res.json());
      return;
    }

    const json = await res.json();
    const token = json.token;

    request.cookie.set("token", token, {
      path: "/",
      httpOnly: true,
      sameSite: true,
      secure: false,
    });
    throw request.redirect(303, "/trips");
  },
  zod$({
    email: z.string().email(),
    password: z.string(),
  }),
);

export default component$(() => {
  const login = useLoginAction();

  return (
    <div class="flex h-screen items-center justify-center">
      <Form
        action={login}
        class="w-full max-w-sm rounded-lg bg-base-200 p-6 shadow-lg"
      >
        <h2 class="mb-4 text-2xl font-bold">Login</h2>
        <div class="mb-4">
          <label
            class="mb-2 block text-sm font-bold text-base-content"
            for="email"
          >
            Email
          </label>
          <input
            value={login.formData?.get("email")}
            class="input input-bordered w-full"
            type="text"
            id="email"
            name="email"
            placeholder="Enter your email"
          />
        </div>
        <div class="mb-6">
          <label
            class="mb-2 block text-sm font-bold text-base-content"
            for="password"
          >
            Password
          </label>
          <input
            value={login.formData?.get("password")}
            class="input input-bordered w-full"
            type="password"
            id="password"
            name="password"
            placeholder="Enter your password"
          />
        </div>
        <div class="flex items-center justify-between">
          <button class="btn btn-primary w-full" type="submit">
            Login
          </button>
        </div>
        <div class="mt-4 text-center">
          <p class="text-sm">
            Don't have an account?{" "}
            <a href="/signup" class="text-primary">
              Sign up
            </a>
          </p>
        </div>
      </Form>
    </div>
  );
});
