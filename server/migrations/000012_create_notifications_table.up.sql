CREATE TYPE notification_type AS ENUM (
    'trip_cancelled',
    'trip_date_change',
    'trip_invite_pending',
    'trip_invite_accepted',
    'trip_invite_cancelled',
    'trip_invite_declined',
    'trip_member_left',
    'trip_member_removed',
    'trip_ownership_transfer'
);

CREATE TABLE IF NOT EXISTS notifications (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    trip_id UUID REFERENCES trips(id) ON DELETE CASCADE,
    TYPE notification_type NOT NULL,
    message TEXT NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
    read_at TIMESTAMP WITH TIME ZONE,
    metadata JSONB -- flexible field for additional context
);

CREATE INDEX idx_notifications_user_id ON notifications(user_id);

CREATE INDEX idx_notifications_created_at ON notifications(created_at);