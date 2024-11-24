import type { ValueOf } from "~/types/helpers";

export const SERVER_ERROR_MESSAGES = {
  INVALID_CREDENTIALS: "invalid authentication credentials",
  MISSING_CREDENTIALS: "invalid or missing authentication token",
} as const;

export type ServerErrorMessages = ValueOf<typeof SERVER_ERROR_MESSAGES>;
