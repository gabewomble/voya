import { createContextId } from "@builder.io/qwik";

type ActivityContextType = {
  unreadCount: number;
};

export const ActivityContext = createContextId<ActivityContextType>("Activity");
