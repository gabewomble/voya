import { requireAuth } from "~/middleware/auth";

export const onRequest = requireAuth;
