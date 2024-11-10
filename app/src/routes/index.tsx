import { component$ } from "@builder.io/qwik";
import type { DocumentHead } from "@builder.io/qwik-city";

export default component$(() => {
  return (
    <>
      <div class="hero min-h-screen bg-base-200">
        <div class="hero-content text-center">
          <div class="max-w-md">
            <h1 class="text-5xl font-bold">Plan Your Perfect Trip with Voya</h1>
            <p class="py-6">
              Collaborate with friends and family to plan and manage your group
              trips effortlessly.
            </p>
            <a href="/signup" class="btn btn-primary">
              Get Started
            </a>
          </div>
        </div>
      </div>

      <div class="bg-base-100 py-12">
        <div class="container mx-auto">
          <h2 class="mb-8 text-center text-3xl font-bold">Features</h2>
          <div class="grid grid-cols-1 gap-8 md:grid-cols-3">
            <div class="card bg-base-200 shadow-lg">
              <div class="card-body">
                <h3 class="card-title">Collaborative Planning</h3>
                <p>
                  Plan trips together with your friends and family in real-time.
                </p>
              </div>
            </div>
            <div class="card bg-base-200 shadow-lg">
              <div class="card-body">
                <h3 class="card-title">Itinerary Management</h3>
                <p>
                  Keep track of your trip itinerary and make adjustments as
                  needed.
                </p>
              </div>
            </div>
            <div class="card bg-base-200 shadow-lg">
              <div class="card-body">
                <h3 class="card-title">Expense Tracking</h3>
                <p>Manage and split expenses among group members easily.</p>
              </div>
            </div>
          </div>
        </div>
      </div>

      <footer class="footer bg-base-300 p-10 text-base-content">
        <div>
          <p>© 2023 Voya. All rights reserved.</p>
        </div>
      </footer>
    </>
  );
});

export const head: DocumentHead = {
  title: "Welcome to Voya",
  meta: [
    {
      name: "description",
      content: "Plan and manage group trips collaboratively with Voya.",
    },
  ],
};
