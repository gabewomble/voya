import { z } from "@builder.io/qwik-city";
import { tripSchema } from "./trips";
import { userSchema } from "./user";

const memberStatusEnum = z.enum([
  "owner",
  "pending",
  "accepted",
  "declined",
  "removed",
  "cancelled",
]);

export const getTripByIdResponseSchema = z.object({
  trip: tripSchema,
  members: z.array(
    z.object({
      id: z.string().uuid(),
      name: z.string(),
      email: z.string(),
      member_status: memberStatusEnum,
      updated_at: z.string().optional(),
      updated_by: z.string().uuid().optional(),
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
