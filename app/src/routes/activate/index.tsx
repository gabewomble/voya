import { component$ } from "@builder.io/qwik";
import {
  routeLoader$,
  Form,
  routeAction$,
  type RequestEventBase,
  Link,
} from "@builder.io/qwik-city";
import { serverFetch } from "~/helpers/server-fetch";
import { setCookie } from "~/helpers/set-cookie";

async function handleActivation(requestEvent: RequestEventBase, token: string) {
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
      setCookie("token", newToken, requestEvent);
    }

    return true;
  }

  return false;
}

type ActivateUserData = {
  email: string;
  token?: string;
  activated: boolean;
  attemptedActivation: boolean;
};

export const useActivateUser = routeLoader$<ActivateUserData>(
  async (requestEvent) => {
    const token = requestEvent.query.get("t") ?? undefined;
    const identifier = requestEvent.query.get("i");

    if (!identifier) throw requestEvent.redirect(303, "/login");

    if (token) {
      const ok = await handleActivation(requestEvent, token);

      return {
        email: identifier,
        token,
        attemptedActivation: true,
        activated: ok,
      };
    }

    return {
      email: identifier,
      token,
      activated: false,
      attemptedActivation: false,
    };
  },
);

type ResendActionData = {
  sent: boolean;
};

export const useResendAction = routeAction$<ResendActionData>(
  async (_, request) => {
    const identifier = request.query.get("i");

    if (!identifier) throw request.redirect(303, "/login");

    const res = await serverFetch(
      `/users/resend-activation`,
      {
        method: "POST",
        body: JSON.stringify({ identifier }),
      },
      request,
    );

    return {
      sent: res.ok,
    };
  },
);

const ActivateAccount = component$(() => (
  <>
    <h1 class="card-title mb-4 text-4xl font-bold">Activate Your Account</h1>
    <p class="text-lg">
      Thank you for signing up! We have sent an activation link to your email
      address.
    </p>
    <p class="text-lg">
      Please check your email and click on the activation link to activate your
      account.
    </p>
  </>
));

const ActivationSuccess = component$(() => (
  <>
    <h1 class="card-title mb-4 text-4xl font-bold">Account Activated</h1>
    <p class="text-lg">Your account has been successfully activated.</p>
    <div class="mt-6">
      <Link href="/trips" class="btn btn-primary">
        Go to Trips
      </Link>
    </div>
  </>
));

const ActivationFailure = component$(() => (
  <>
    <h1 class="card-title mb-4 text-4xl font-bold">Activation Failed</h1>
    <p class="text-lg">
      There was an issue activating your account. Please try again.
    </p>
  </>
));

export default component$(() => {
  const { value: data } = useActivateUser();
  const resendAction = useResendAction();

  return (
    <div class="container mx-auto py-8">
      <div class="card bg-base-200 shadow-lg">
        <div class="card-body items-center text-center">
          {data.attemptedActivation ? (
            <>
              {data.activated ? <ActivationSuccess /> : <ActivationFailure />}
            </>
          ) : (
            <ActivateAccount />
          )}
          {!data.activated && (
            <Form action={resendAction}>
              <button type="submit" class="btn btn-secondary mt-6">
                Resend Activation Email
              </button>
            </Form>
          )}
          {resendAction.submitted && (
            <>
              {resendAction.value?.sent ? (
                <p class="mt-4 text-success">
                  Activation link resent successfully. Please check your email.
                </p>
              ) : (
                <p class="mt-4 text-error">
                  There was an issue resending the activation link.
                </p>
              )}
            </>
          )}
        </div>
      </div>
    </div>
  );
});
