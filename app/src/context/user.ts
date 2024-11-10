import { createContextId, type Signal } from "@builder.io/qwik";

type User = Record<string, string> | null;

export const UserContext = createContextId<Signal<User>>("User");
