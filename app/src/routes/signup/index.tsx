import { component$ } from "@builder.io/qwik";
import { Form, zod$, routeAction$, z } from "@builder.io/qwik-city";
import { serverFetch } from "~/helpers/server-fetch";

export const useSignupAction = routeAction$(
  async (formData, request) => {
    const { email, password, confirmPassword, name } = formData;
    // Add your signup logic here
    // For example, you can send a request to your backend to create a new user
    if (password !== confirmPassword) {
      request.fail(400, { error: "Passwords do not match" });
      return;
    }

    const res = await serverFetch(
      "/users",
      {
        method: "POST",
        body: JSON.stringify({
          email,
          password,
          name,
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
    name: z.string().min(1).max(256),
    email: z.string().email(),
    password: z.string().min(8).max(72),
    confirmPassword: z.string().min(8).max(72),
  }),
);

export default component$(() => {
  // TODO: use form errors from useSignupAction
  const signup = useSignupAction();

  return (
    <div class="flex h-screen items-center justify-center">
      <Form
        action={signup}
        class="w-full max-w-sm rounded-lg bg-base-200 p-6 shadow-lg"
      >
        <h2 class="mb-4 text-2xl font-bold">Sign Up</h2>
        <div class="mb-4">
          <label
            class="mb-2 block text-sm font-bold text-base-content"
            for="name"
          >
            Name
          </label>
          <input
            value={signup.formData?.get("name")}
            class="input input-bordered w-full"
            type="text"
            id="name"
            name="name"
            placeholder="Enter your name"
          />
        </div>
        <div class="mb-4">
          <label
            class="mb-2 block text-sm font-bold text-base-content"
            for="email"
          >
            Email
          </label>
          <input
            value={signup.formData?.get("email")}
            class="input input-bordered w-full"
            type="text"
            id="email"
            name="email"
            placeholder="Enter your email"
          />
        </div>
        <div class="mb-4">
          <label
            class="mb-2 block text-sm font-bold text-base-content"
            for="password"
          >
            Password
          </label>
          <input
            value={signup.formData?.get("password")}
            class="input input-bordered w-full"
            type="password"
            id="password"
            name="password"
            placeholder="Enter your password"
          />
        </div>
        <div class="mb-6">
          <label
            class="mb-2 block text-sm font-bold text-base-content"
            for="confirmPassword"
          >
            Confirm Password
          </label>
          <input
            value={signup.formData?.get("confirmPassword")}
            class="input input-bordered w-full"
            type="password"
            id="confirmPassword"
            name="confirmPassword"
            placeholder="Confirm your password"
          />
        </div>
        <div class="flex items-center justify-between">
          <button class="btn btn-primary w-full" type="submit">
            Sign Up
          </button>
        </div>
        <div class="mt-4 text-center">
          <p class="text-sm">
            Already have an account?{" "}
            <a href="/login" class="text-primary">
              Login
            </a>
          </p>
        </div>
      </Form>
    </div>
  );
});
