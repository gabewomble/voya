import { z } from "@builder.io/qwik-city";

export const notificationTypeEnum = z.enum([
  "trip_cancelled",
  "trip_date_change",
  "trip_invite_pending",
  "trip_invite_accepted",
  "trip_invite_cancelled",
  "trip_invite_declined",
  "trip_member_left",
  "trip_member_removed",
  "trip_ownership_transfer",
]);

export type NotificationTypeEnum = z.infer<typeof notificationTypeEnum>;

export const notificationMetadataSchema = z.object({
  user_id: z.string().uuid(),
  user_name: z.string(),
});

export const notificationSchema = z.object({
  id: z.string().uuid(),
  user_id: z.string().uuid(),
  trip_id: z.string().uuid(),
  message: z.string(),
  notification_type: notificationTypeEnum,
  created_at: z.string(),
  created_by: z.string().uuid(),
  target_user_id: z.string().uuid(),
  read_at: z.string().nullable(),
});

export type Notification = z.infer<typeof notificationSchema>;
