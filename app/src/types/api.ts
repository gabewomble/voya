import { z } from "@builder.io/qwik-city";
import { tripSchema } from "./trips";
import { userSchema } from "./users";
import { memberSchema } from "./members";

export const getTripByIdResponseSchema = z.object({
  trip: tripSchema,
  members: z.array(memberSchema),
});

export const searchUsersResponseSchema = z.object({
  users: z.array(userSchema),
});

export type GetTripByIdResponse = z.infer<typeof getTripByIdResponseSchema>;

export const listTripsResponseSchema = z.object({ trips: z.array(tripSchema) });

export type ListTripsResponse = z.infer<typeof listTripsResponseSchema>;

export type SearchUsersResponse = z.infer<typeof searchUsersResponseSchema>;
