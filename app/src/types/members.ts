import { z } from "@builder.io/qwik-city";

export const memberStatusEnum = z.enum([
  "owner",
  "pending",
  "accepted",
  "declined",
  "removed",
  "cancelled",
]);

export type MemberStatusEnum = z.infer<typeof memberStatusEnum>;

export const memberSchema = z.object({
  id: z.string().uuid(),
  name: z.string(),
  email: z.string(),
  member_status: memberStatusEnum,
  updated_at: z.string().optional(),
  updated_by: z.string().uuid().optional(),
});

export type Member = z.infer<typeof memberSchema>;
