import { component$, useContext } from "@builder.io/qwik";
import { UserContext } from "~/context/user";

export default component$(() => {
  const user = useContext(UserContext);

  return (
    <div class="container mx-auto py-8">
      <div class="card bg-base-200 shadow-lg">
        <div class="card-body">
          <h1 class="card-title mb-4 text-4xl font-bold">Profile</h1>
          <p class="mb-6 text-lg">Update your profile information</p>
          <div class="grid grid-cols-1 gap-4 md:grid-cols-2">
            <div>
              <label class="block text-sm font-medium text-base-content">
                Username
              </label>
              <p class="mt-1 text-lg font-semibold text-base-content">
                {user.value?.username}
              </p>
            </div>
            <div>
              <label class="block text-sm font-medium text-base-content">
                Name
              </label>
              <p class="mt-1 text-lg font-semibold text-base-content">
                {user.value?.name}
              </p>
            </div>
            <div class="md:col-span-2">
              <label class="block text-sm font-medium text-base-content">
                Email
              </label>
              <p class="mt-1 text-lg font-semibold text-base-content">
                {user.value?.email}
              </p>
            </div>
          </div>
          <div class="mt-6 flex justify-end">
            <a href="/settings/profile/edit" class="btn btn-primary">
              Edit Profile
            </a>
          </div>
        </div>
      </div>
    </div>
  );
});
