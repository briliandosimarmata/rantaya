BEGIN;
SELECT pg_advisory_xact_lock(7465822);

DO $$
BEGIN
    IF NOT EXISTS (
        SELECT 1 FROM events e
        JOIN organizers o ON o.id = e.organizer_id
        JOIN accounts a ON a.id = o.account_id
        JOIN event_sessions s ON s.event_id = e.id
        WHERE e.id = 'e1' AND e.slug = 'di-balik-layar'
          AND o.id = 'ruang' AND a.id = 'organizer-demo'
          AND a.role = 'organizer' AND s.id = 's-e1'
    ) THEN
        RAISE EXCEPTION 'Expected fictional demo event was not found';
    END IF;
END $$;

-- Add a new session; booked schedules, orders and tickets stay immutable.
WITH schedule AS (
    SELECT ((CURRENT_TIMESTAMP AT TIME ZONE 'Asia/Jakarta')::date + 30
            + time '19:00') AT TIME ZONE 'Asia/Jakarta' AS starts_at
)
INSERT INTO event_sessions(id, event_id, label, starts_at, ends_at, price, capacity)
SELECT 's-e1-demo-' || to_char(CURRENT_TIMESTAMP AT TIME ZONE 'Asia/Jakarta', 'YYYYMMDD'),
       e.id, 'Pertunjukan malam · demo mendatang',
       schedule.starts_at, schedule.starts_at + interval '90 minutes', s.price, 60
FROM events e
JOIN event_sessions s ON s.event_id = e.id AND s.id = 's-e1'
CROSS JOIN schedule
WHERE e.id = 'e1' AND e.published
  AND NOT EXISTS (
      SELECT 1 FROM event_sessions future
      WHERE future.event_id = e.id
        AND future.starts_at > CURRENT_TIMESTAMP + interval '24 hours'
  )
ON CONFLICT (id) DO NOTHING
RETURNING id, label, starts_at, ends_at, capacity;
COMMIT;
