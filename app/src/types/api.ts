import { z } from "@builder.io/qwik-city";
import { tripSchema } from "./trips";
import { userSchema } from "./user";

export const getTripByIdResponseSchema = z.object({
  trip: tripSchema,
  members: z.array(
    z.object({
      id: z.string().uuid(),
      name: z.string(),
      email: z.string(),
    }),
  ),
});

export const searchUsersResponseSchema = z.object({
  users: z.array(userSchema),
});

export type GetTripByIdResponse = z.infer<typeof getTripByIdResponseSchema>;

export const listTripsResponseSchema = z.object({ trips: z.array(tripSchema) });

export type ListTripsResponse = z.infer<typeof listTripsResponseSchema>;

export type SearchUsersResponse = z.infer<typeof searchUsersResponseSchema>;
