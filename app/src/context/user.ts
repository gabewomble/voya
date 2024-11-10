import { createContextId } from "@builder.io/qwik";

export const UserContext = createContextId<Record<string, unknown> | null>(
  "User",
);
