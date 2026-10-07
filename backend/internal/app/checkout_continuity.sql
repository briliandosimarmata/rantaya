ALTER TABLE oauth_states ADD COLUMN next_path text NOT NULL DEFAULT '';
CREATE TABLE payment_reminders (
  order_id text PRIMARY KEY REFERENCES orders(id),
  created_at timestamptz NOT NULL DEFAULT now()
);
ALTER TABLE votes ADD COLUMN created_at timestamptz NOT NULL DEFAULT now();
CREATE TABLE activity_clicks (
  id bigserial PRIMARY KEY,
  organizer_id text NOT NULL REFERENCES organizers(id),
  kind text NOT NULL CHECK(kind IN ('ticket','merch')),
  event_id text REFERENCES events(id),
  product_id text REFERENCES products(id),
  created_at timestamptz NOT NULL DEFAULT now(),
  CHECK((kind='ticket' AND event_id IS NOT NULL AND product_id IS NULL)
    OR (kind='merch' AND product_id IS NOT NULL AND event_id IS NULL))
);
CREATE INDEX activity_clicks_organizer ON activity_clicks(organizer_id,kind,created_at);
