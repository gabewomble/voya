import { component$, useStore } from "@builder.io/qwik";
import { routeLoader$, Form, routeAction$ } from "@builder.io/qwik-city";
import { serverFetch } from "~/helpers/server-fetch";

export const useActivateUser = routeLoader$(async (requestEvent) => {
  const { token } = requestEvent.params;
  const res = await serverFetch(
    `/users/activated`,
    {
      method: "PUT",
      body: JSON.stringify({ token }),
    },
    requestEvent,
  );

  if (res.ok) {
    const json = await res.json();
    const { token: newToken } = json;

    if (newToken) {
      requestEvent.cookie.set("token", newToken, {
        path: "/",
        httpOnly: true,
        sameSite: true,
        secure: false,
      });
    }
  }
  return {
    success: res.ok,
  };
});

export const useResendAction = routeAction$(() => {
  // TODO: implement activation re-send
});

export default component$(() => {
  const activation = useActivateUser();
  const resendActivation = useResendAction();
  const state = useStore({ resendSuccess: false });

  return (
    <div class="container mx-auto py-8">
      <div class="card bg-base-200 shadow-lg">
        <div class="card-body items-center text-center">
          {activation.value.success ? (
            <>
              <h1 class="card-title mb-4 text-4xl font-bold">
                Account Activated
              </h1>
              <p class="text-lg">
                Your account has been successfully activated.
              </p>
              <div class="mt-6">
                <a href="/trips" class="btn btn-primary">
                  Go to Trips
                </a>
              </div>
            </>
          ) : (
            <>
              <h1 class="card-title mb-4 text-4xl font-bold">
                Activation Failed
              </h1>
              <p class="text-lg">
                There was an issue activating your account. Please try again.
              </p>
              <Form action={resendActivation} class="mt-6">
                <button type="submit" class="btn btn-secondary">
                  Resend Activation Link
                </button>
              </Form>
              {state.resendSuccess && (
                <p class="mt-4 text-green-500">
                  Activation link resent successfully. Please check your email.
                </p>
              )}
            </>
          )}
        </div>
      </div>
    </div>
  );
});
