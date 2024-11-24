import { z } from "@builder.io/qwik-city";

export const errorResponseSchema = z.object({
  errors: z.array(
    z.object({
      message: z.string(),
      field: z.optional(z.string()),
    }),
  ),
});

export type ErrorResponse = z.infer<typeof errorResponseSchema>;
