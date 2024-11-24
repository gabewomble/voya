import { component$ } from "@builder.io/qwik";
import { routeLoader$, z } from "@builder.io/qwik-city";
import { serverFetch } from "~/helpers/server-fetch";
import { setCookie } from "~/helpers/set-cookie";
import { requireNoAuth } from "~/middleware/auth";
import {
  zodForm$,
  useForm,
  formAction$,
  type InitialValues,
} from "@modular-forms/qwik";
import type { ErrorResponse } from "~/types/server-errors";
import { TextInput } from "~/components";
import { mapServerErrors } from "~/helpers/map-server-errors";

export const onGet = requireNoAuth;

const invalidCredentialError = "invalid authentication credentials";

const loginFormSchema = z.object({
  email: z.string().email("Invalid email address"),
  password: z.string().min(8, "Password must be at least 8 characters"),
});

type LoginForm = z.infer<typeof loginFormSchema>;

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
    const { errors } = ((await res.json()) ?? { errors: [] }) as ErrorResponse;

    const { messages, fields } = mapServerErrors({
      errors,
      messages: {
        [invalidCredentialError]: "Invalid email or password",
      },
      fields: {
        email: {},
        password: {},
      },
    });

    return {
      message: messages[0],
      errors: {
        email: fields.email[0],
        password: fields.password[0],
      },
    };
  }

  const json = await res.json();
  const token = json.token;

  setCookie("token", token, request);

  throw request.redirect(303, "/trips");
}, zodForm$(loginFormSchema));

export const useFormLoader = routeLoader$<InitialValues<LoginForm>>(() => ({
  email: "",
  password: "",
}));

export default component$(() => {
  const [loginForm, { Form, Field }] = useForm<LoginForm>({
    loader: useFormLoader(),
    action: useLoginAction(),
    validate: zodForm$(loginFormSchema),
    revalidateOn: "blur",
  });

  return (
    <div class="flex h-screen items-center justify-center">
      <Form class="w-full max-w-sm rounded-lg bg-base-200 p-6 shadow-lg">
        <h2 class="mb-4 text-2xl font-bold">Login</h2>
        {loginForm.response.message && (
          <p class="my-2 text-error">{loginForm.response.message}</p>
        )}
        <div class="mb-4">
          <Field name="email">
            {(field, props) => (
              <TextInput
                id="email"
                label="Email"
                name="email"
                placeholder="Enter your email"
                fieldProps={props}
                field={field}
              />
            )}
          </Field>
        </div>
        <div class="mb-6">
          <Field name="password">
            {(field, props) => (
              <TextInput
                label="Password"
                id="password"
                name="password"
                placeholder="Enter your password"
                type="password"
                fieldProps={props}
                field={field}
              />
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
