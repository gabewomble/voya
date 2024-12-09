import { createContextId, type Signal } from "@builder.io/qwik";
import type { User } from "~/types/users";

export const UserContext = createContextId<Signal<User | null>>("User");
