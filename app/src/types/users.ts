import { z } from "@builder.io/qwik-city";

export const userSchema = z.object({
  activated: z.boolean(),
  created_at: z.string(),
  email: z.string(),
  id: z.string().uuid(),
  name: z.string(),
  username: z.string(),
  version: z.number().optional(),
});

export type User = z.infer<typeof userSchema>;
