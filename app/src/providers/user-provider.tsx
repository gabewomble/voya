import {
  component$,
  createContextId,
  Slot,
  useContextProvider,
  useStore,
} from "@builder.io/qwik";

type UserContext = {
  // TODO: fix type
  user: Record<string, unknown> | null;
};

const UserContextId = createContextId<UserContext>("User");

export const UserProvider = component$(() => {
  const store = useStore<UserContext>({
    user: null,
  });

  useContextProvider(UserContextId, store);

  return <Slot />;
});
