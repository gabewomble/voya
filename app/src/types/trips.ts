import { z } from "@builder.io/qwik-city";

export const tripSchema = z.object({
  id: z.string().uuid(),
  name: z.string(),
  description: z.string(),
  start_date: z.date().nullable(),
  end_date: z.date().nullable(),
  created_at: z.string(),
  updated_at: z.string(),
});

export type Trip = z.infer<typeof tripSchema>;
