import { component$ } from "@builder.io/qwik";
import { routeLoader$, z } from "@builder.io/qwik-city";
import { serverFetch } from "~/helpers/server-fetch";
import { requireNoAuth } from "~/middleware/auth";
import {
  zodForm$,
  useForm,
  formAction$,
  type InitialValues,
} from "@modular-forms/qwik";

export const onGet = requireNoAuth;

const invalidCredentialError = "invalid authentication credentials";

const LoginForm = z.object({
  email: z.string().email("Invalid email address"),
  password: z.string().min(8, "Password must be at least 8 characters"),
});

type LoginForm = z.infer<typeof LoginForm>;

export const useLoginAction = formAction$<LoginForm>(async (data, request) => {
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
    const { error, errors } = (await res.json()) ?? {};

    return {
      message:
        error === invalidCredentialError ? "Invalid email or password" : error,
      errors: {
        email: errors?.email,
        password: errors?.password,
      },
    };
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
}, zodForm$(LoginForm));

export const useFormLoader = routeLoader$<InitialValues<LoginForm>>(() => ({
  email: "",
  password: "",
}));

export default component$(() => {
  const [loginForm, { Form, Field }] = useForm<LoginForm>({
    loader: useFormLoader(),
    action: useLoginAction(),
    validate: zodForm$(LoginForm),
    revalidateOn: "blur",
  });

  return (
    <div class="flex h-screen items-center justify-center">
      <Form class="w-full max-w-sm rounded-lg bg-base-200 p-6 shadow-lg">
        <h2 class="mb-4 text-2xl font-bold">Login</h2>
        {loginForm.response.message && (
          <p class="my-2 text-warning">{loginForm.response.message}</p>
        )}
        <div class="mb-4">
          <label
            class="mb-2 block text-sm font-bold text-base-content"
            for="email"
          >
            Email
          </label>
          <Field name="email">
            {(field, props) => (
              <>
                <input
                  {...props}
                  class={`input input-bordered w-full ${field.error ? "input-error" : ""}`}
                  type="text"
                  id="email"
                  name="email"
                  placeholder="Enter your email"
                  value={field.value}
                />
                {field.error && <p class="my-1 text-error">{field.error}</p>}
              </>
            )}
          </Field>
        </div>
        <div class="mb-6">
          <label
            class="mb-2 block text-sm font-bold text-base-content"
            for="password"
          >
            Password
          </label>
          <Field name="password">
            {(field, props) => (
              <>
                <input
                  {...props}
                  class={`input input-bordered w-full ${field.error ? "input-error" : ""}`}
                  type="password"
                  id="password"
                  name="password"
                  placeholder="Enter your password"
                  value={field.value}
                />
                {field.error && <p class="my-1 text-error">{field.error}</p>}
              </>
            )}
          </Field>
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
