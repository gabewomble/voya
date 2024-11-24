import { component$ } from "@builder.io/qwik";
import { z, routeLoader$ } from "@builder.io/qwik-city";
import { serverFetch } from "~/helpers/server-fetch";
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

const signupFormSchema = z
  .object({
    name: z
      .string()
      .min(1, "Name is required")
      .max(50, "Name cannot be longer than 50 characters"),
    email: z.string().email("Invalid email address"),
    password: z
      .string()
      .min(8, "Password must be at least 8 characters")
      .max(32, "Password cannot be longer than 32 characters"),
    confirmPassword: z
      .string()
      .min(8, "Password must be at least 8 characters"),
  })
  .superRefine(({ password, confirmPassword }, ctx) => {
    if (password && confirmPassword && password !== confirmPassword) {
      ctx.addIssue({
        code: "custom",
        message: "Passwords do not match",
        path: ["confirmPassword"],
      });
    }
  });

type SignupForm = z.infer<typeof signupFormSchema>;

export const useSignupAction = formAction$<SignupForm>(
  async (data, request) => {
    const { name, email, password } = data;

    const res = await serverFetch(
      "/users",
      {
        method: "POST",
        headers: {
          "Content-Type": "application/json",
        },
        body: JSON.stringify({
          name,
          email,
          password,
        }),
      },
      request,
    );

    if (!res.ok) {
      const { errors } = ((await res.json()) ?? {
        errors: [],
      }) as ErrorResponse;

      const { messages, fields } = mapServerErrors({
        errors,
        messages: {},
        fields: {
          name: {},
          email: {
            "duplicate email": "Email is already in use",
          },
          password: {},
        },
      });

      return {
        message: messages[0],
        errors: {
          name: fields.name[0],
          email: fields.email[0],
          password: fields.password[0],
        },
      };
    }

    throw request.redirect(303, "/activate");
  },
  zodForm$(signupFormSchema),
);

export const useSignupFormLoader = routeLoader$<InitialValues<SignupForm>>(
  () => ({
    name: "",
    email: "",
    password: "",
    confirmPassword: "",
  }),
);

export default component$(() => {
  const [signupForm, { Form, Field }] = useForm<SignupForm>({
    loader: useSignupFormLoader(),
    action: useSignupAction(),
    validate: zodForm$(signupFormSchema),
    validateOn: "blur",
    revalidateOn: "blur",
  });

  return (
    <div class="flex h-screen items-center justify-center">
      <Form class="w-full max-w-sm rounded-lg bg-base-200 p-6 shadow-lg">
        <h2 class="mb-4 text-2xl font-bold">Sign Up</h2>
        {signupForm.response.message && (
          <p class="my-2 text-error">{signupForm.response.message}</p>
        )}
        <div class="mb-4">
          <Field name="name">
            {(field, props) => (
              <TextInput
                id="name"
                label="Name"
                name="name"
                placeholder="Enter your name"
                field={field}
                fieldProps={props}
              />
            )}
          </Field>
        </div>
        <div class="mb-4">
          <Field name="email">
            {(field, props) => (
              <TextInput
                id="email"
                label="Email"
                name="email"
                placeholder="Enter your email"
                field={field}
                fieldProps={props}
              />
            )}
          </Field>
        </div>
        <div class="mb-4">
          <Field name="password">
            {(field, props) => (
              <TextInput
                id="password"
                label="Password"
                name="password"
                placeholder="Enter your password"
                type="password"
                field={field}
                fieldProps={props}
              />
            )}
          </Field>
        </div>
        <div class="mb-6">
          <Field name="confirmPassword">
            {(field, props) => (
              <TextInput
                id="confirmPassword"
                label="Confirm Password"
                name="confirmPassword"
                placeholder="Confirm your password"
                type="password"
                field={field}
                fieldProps={props}
              />
            )}
          </Field>
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
