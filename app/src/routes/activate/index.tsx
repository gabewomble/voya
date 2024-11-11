import { component$ } from "@builder.io/qwik";

export default component$(() => {
  return (
    <div class="container mx-auto py-8">
      <div class="card bg-base-200 shadow-lg">
        <div class="card-body items-center text-center">
          <h1 class="card-title mb-4 text-4xl font-bold">
            Activate Your Account
          </h1>
          <p class="text-lg">
            Thank you for signing up! We have sent an activation link to your
            email address.
          </p>
          <p class="text-lg">
            Please check your email and click on the activation link to activate
            your account.
          </p>
        </div>
      </div>
    </div>
  );
});
