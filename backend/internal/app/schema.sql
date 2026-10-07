CREATE TABLE IF NOT EXISTS accounts (
 id text PRIMARY KEY, role text NOT NULL CHECK(role IN ('customer','organizer','admin')),
 google_sub text, email text NOT NULL, name text NOT NULL, city text NOT NULL DEFAULT 'Karawang',
 bio text NOT NULL DEFAULT '', avatar_url text NOT NULL DEFAULT '', interests jsonb NOT NULL DEFAULT '[]',
 preferences jsonb NOT NULL DEFAULT '{"events":true,"replies":true,"reminders":true}',
 onboarding_done boolean NOT NULL DEFAULT false, created_at timestamptz NOT NULL DEFAULT now(),
 UNIQUE(google_sub,role)
);
CREATE TABLE IF NOT EXISTS auth_sessions (token_hash text PRIMARY KEY, account_id text NOT NULL REFERENCES accounts(id) ON DELETE CASCADE, expires_at timestamptz NOT NULL);
CREATE TABLE IF NOT EXISTS oauth_states (id text PRIMARY KEY, role text NOT NULL, verifier text NOT NULL, expires_at timestamptz NOT NULL);
CREATE TABLE IF NOT EXISTS organizers (
 id text PRIMARY KEY, account_id text NOT NULL UNIQUE REFERENCES accounts(id), slug text NOT NULL UNIQUE,
 name text NOT NULL, description text NOT NULL DEFAULT '', about text NOT NULL DEFAULT '',
 city text NOT NULL DEFAULT 'Karawang', category text NOT NULL DEFAULT 'Teater',
 avatar_url text NOT NULL DEFAULT '', cover_url text NOT NULL DEFAULT '', verified boolean NOT NULL DEFAULT false,
 created_at timestamptz NOT NULL DEFAULT now()
);
CREATE TABLE IF NOT EXISTS events (
 id text PRIMARY KEY, slug text NOT NULL UNIQUE, organizer_id text NOT NULL REFERENCES organizers(id),
 title text NOT NULL, description text NOT NULL, category text NOT NULL, city text NOT NULL,
 venue text NOT NULL, address text NOT NULL DEFAULT '', maps_url text NOT NULL DEFAULT '',
 duration text NOT NULL DEFAULT '', language text NOT NULL DEFAULT 'Bahasa Indonesia', age text NOT NULL DEFAULT 'Semua usia',
 flyer_url text NOT NULL DEFAULT '', trailer_url text NOT NULL DEFAULT '', layout_url text NOT NULL DEFAULT '',
 lineup jsonb NOT NULL DEFAULT '[]', published boolean NOT NULL DEFAULT false,
 created_at timestamptz NOT NULL DEFAULT now(), updated_at timestamptz NOT NULL DEFAULT now()
);
CREATE TABLE IF NOT EXISTS event_sessions (
 id text PRIMARY KEY, event_id text NOT NULL REFERENCES events(id), label text NOT NULL,
 starts_at timestamptz NOT NULL, ends_at timestamptz NOT NULL, price bigint NOT NULL CHECK(price>=0),
 capacity integer NOT NULL CHECK(capacity>0), CHECK(ends_at>starts_at)
);
CREATE TABLE IF NOT EXISTS payment_methods (
 id text PRIMARY KEY, organizer_id text NOT NULL REFERENCES organizers(id), kind text NOT NULL CHECK(kind IN ('bank','wallet')),
 provider text NOT NULL, number text NOT NULL, holder text NOT NULL, note text NOT NULL DEFAULT '', enabled boolean NOT NULL DEFAULT true
);
CREATE TABLE IF NOT EXISTS uploads (
 id text PRIMARY KEY, account_id text NOT NULL REFERENCES accounts(id), purpose text NOT NULL CHECK(purpose IN ('media','proof')),
 filename text NOT NULL, content_type text NOT NULL, size bigint NOT NULL, created_at timestamptz NOT NULL DEFAULT now()
);
CREATE TABLE IF NOT EXISTS posts (
 id text PRIMARY KEY, account_id text NOT NULL REFERENCES accounts(id), organizer_id text REFERENCES organizers(id),
 title text NOT NULL DEFAULT '', body text NOT NULL, image_url text NOT NULL DEFAULT '',
 link_url text NOT NULL DEFAULT '', link_label text NOT NULL DEFAULT '', city text NOT NULL,
 pinned boolean NOT NULL DEFAULT false, hidden boolean NOT NULL DEFAULT false,
 created_at timestamptz NOT NULL DEFAULT now()
);
CREATE TABLE IF NOT EXISTS post_mentions (post_id text NOT NULL REFERENCES posts(id) ON DELETE CASCADE, kind text NOT NULL CHECK(kind IN ('event','organizer')), target_id text NOT NULL, PRIMARY KEY(post_id,kind,target_id));
CREATE TABLE IF NOT EXISTS votes (post_id text NOT NULL REFERENCES posts(id) ON DELETE CASCADE, account_id text NOT NULL REFERENCES accounts(id), PRIMARY KEY(post_id,account_id));
CREATE TABLE IF NOT EXISTS comments (id text PRIMARY KEY, post_id text NOT NULL REFERENCES posts(id) ON DELETE CASCADE, account_id text NOT NULL REFERENCES accounts(id), parent_id text REFERENCES comments(id), body text NOT NULL, created_at timestamptz NOT NULL DEFAULT now());
CREATE TABLE IF NOT EXISTS follows (account_id text NOT NULL REFERENCES accounts(id), organizer_id text NOT NULL REFERENCES organizers(id), created_at timestamptz NOT NULL DEFAULT now(), PRIMARY KEY(account_id,organizer_id));
CREATE TABLE IF NOT EXISTS bookmarks (account_id text NOT NULL REFERENCES accounts(id), kind text NOT NULL CHECK(kind IN ('post','event')), target_id text NOT NULL, created_at timestamptz NOT NULL DEFAULT now(), PRIMARY KEY(account_id,kind,target_id));
CREATE TABLE IF NOT EXISTS products (
 id text PRIMARY KEY, organizer_id text NOT NULL REFERENCES organizers(id), event_id text REFERENCES events(id),
 name text NOT NULL, description text NOT NULL, price bigint NOT NULL CHECK(price>=0),
 image_url text NOT NULL DEFAULT '', variants jsonb NOT NULL DEFAULT '["Satu ukuran"]',
 availability text NOT NULL DEFAULT 'Tersedia', purchase_url text NOT NULL DEFAULT '', created_at timestamptz NOT NULL DEFAULT now()
);
CREATE TABLE IF NOT EXISTS orders (
 id text PRIMARY KEY, account_id text NOT NULL REFERENCES accounts(id), session_id text NOT NULL REFERENCES event_sessions(id),
 quantity integer NOT NULL CHECK(quantity BETWEEN 1 AND 6), total bigint NOT NULL CHECK(total>=0),
 status text NOT NULL CHECK(status IN ('awaiting_payment','awaiting_review','correction_requested','approved','cancelled','expired')),
 payment_method_id text REFERENCES payment_methods(id), payment_snapshot jsonb, proof_id text REFERENCES uploads(id),
 note text NOT NULL DEFAULT '', expires_at timestamptz, idempotency_key text NOT NULL,
 created_at timestamptz NOT NULL DEFAULT now(), updated_at timestamptz NOT NULL DEFAULT now(), UNIQUE(account_id,idempotency_key)
);
CREATE TABLE IF NOT EXISTS order_history (id bigserial PRIMARY KEY, order_id text NOT NULL REFERENCES orders(id), status text NOT NULL, actor_id text REFERENCES accounts(id), note text NOT NULL DEFAULT '', created_at timestamptz NOT NULL DEFAULT now());
CREATE TABLE IF NOT EXISTS tickets (
 id text PRIMARY KEY, order_id text NOT NULL REFERENCES orders(id), ordinal integer NOT NULL,
 token text NOT NULL UNIQUE, checked_at timestamptz, checked_by text REFERENCES accounts(id), UNIQUE(order_id,ordinal)
);
CREATE TABLE IF NOT EXISTS reviews (
 id text PRIMARY KEY, event_id text NOT NULL REFERENCES events(id), account_id text NOT NULL REFERENCES accounts(id),
 body text NOT NULL, created_at timestamptz NOT NULL DEFAULT now(), updated_at timestamptz NOT NULL DEFAULT now(), hidden boolean NOT NULL DEFAULT false,
 UNIQUE(event_id,account_id)
);
CREATE TABLE IF NOT EXISTS review_replies (id text PRIMARY KEY, review_id text NOT NULL REFERENCES reviews(id), account_id text NOT NULL REFERENCES accounts(id), body text NOT NULL, created_at timestamptz NOT NULL DEFAULT now());
CREATE TABLE IF NOT EXISTS review_helpful (review_id text NOT NULL REFERENCES reviews(id), account_id text NOT NULL REFERENCES accounts(id), PRIMARY KEY(review_id,account_id));
CREATE TABLE IF NOT EXISTS notifications (id text PRIMARY KEY, account_id text NOT NULL REFERENCES accounts(id), kind text NOT NULL, title text NOT NULL, body text NOT NULL DEFAULT '', url text NOT NULL, read_at timestamptz, created_at timestamptz NOT NULL DEFAULT now());
CREATE TABLE IF NOT EXISTS reports (id text PRIMARY KEY, account_id text NOT NULL REFERENCES accounts(id), kind text NOT NULL CHECK(kind IN ('post','review')), target_id text NOT NULL, reason text NOT NULL, status text NOT NULL DEFAULT 'pending' CHECK(status IN ('pending','reviewed','hidden')), created_at timestamptz NOT NULL DEFAULT now());
CREATE INDEX IF NOT EXISTS event_discovery ON events(city,category,published);
CREATE INDEX IF NOT EXISTS session_schedule ON event_sessions(starts_at);
CREATE INDEX IF NOT EXISTS post_discovery ON posts(city,created_at DESC) WHERE NOT hidden;
CREATE INDEX IF NOT EXISTS mentions_target ON post_mentions(kind,target_id);
CREATE INDEX IF NOT EXISTS order_reservations ON orders(session_id,status,expires_at);
CREATE INDEX IF NOT EXISTS account_orders ON orders(account_id,created_at DESC);
CREATE INDEX IF NOT EXISTS account_notifications ON notifications(account_id,created_at DESC);
CREATE INDEX IF NOT EXISTS account_sessions ON auth_sessions(account_id,expires_at);

CREATE OR REPLACE VIEW organizer_view AS SELECT o.id,o.slug,o.account_id,
 jsonb_build_object('id',o.id,'slug',o.slug,'name',o.name,'description',o.description,'about',o.about,'city',o.city,'category',o.category,'avatar_url',o.avatar_url,'cover_url',o.cover_url,'verified',o.verified,'created_at',o.created_at,'followers',(SELECT count(*) FROM follows f WHERE f.organizer_id=o.id)) doc FROM organizers o;
CREATE OR REPLACE VIEW event_view AS SELECT e.id,e.slug,e.organizer_id,e.city,e.category,e.published,
 jsonb_build_object('id',e.id,'slug',e.slug,'organizer_id',e.organizer_id,'organizer',ov.doc,'title',e.title,'description',e.description,'category',e.category,'city',e.city,'venue',e.venue,'address',e.address,'maps_url',e.maps_url,'duration',e.duration,'language',e.language,'age',e.age,'flyer_url',e.flyer_url,'trailer_url',e.trailer_url,'layout_url',e.layout_url,'lineup',e.lineup,'published',e.published,'created_at',e.created_at,
 'starts_at',(SELECT min(s.starts_at) FROM event_sessions s WHERE s.event_id=e.id),'ends_at',(SELECT max(s.ends_at) FROM event_sessions s WHERE s.event_id=e.id),
 'price',(SELECT min(s.price) FROM event_sessions s WHERE s.event_id=e.id),
 'sessions',COALESCE((SELECT jsonb_agg(to_jsonb(s)||jsonb_build_object('available',greatest(0,s.capacity-COALESCE((SELECT sum(r.quantity) FROM orders r WHERE r.session_id=s.id AND (r.status IN ('awaiting_review','correction_requested','approved') OR (r.status='awaiting_payment' AND r.expires_at>now()))),0))) ORDER BY s.starts_at) FROM event_sessions s WHERE s.event_id=e.id),'[]')) doc
 FROM events e JOIN organizer_view ov ON ov.id=e.organizer_id;
CREATE OR REPLACE VIEW post_view AS SELECT p.id,p.account_id,p.organizer_id,p.city,p.hidden,p.created_at,p.pinned,
 jsonb_build_object('id',p.id,'account_id',p.account_id,'organizer_id',p.organizer_id,'author',COALESCE(o.name,a.name),'avatar_url',COALESCE(o.avatar_url,a.avatar_url),'official',p.organizer_id IS NOT NULL,'title',p.title,'body',p.body,'image_url',p.image_url,'link_url',p.link_url,'link_label',p.link_label,'city',p.city,'pinned',p.pinned,'created_at',p.created_at,'votes',(SELECT count(*) FROM votes v WHERE v.post_id=p.id),'comment_count',(SELECT count(*) FROM comments c WHERE c.post_id=p.id),
 'mentions',COALESCE((SELECT jsonb_agg(jsonb_build_object('kind',m.kind,'id',m.target_id,'name',CASE WHEN m.kind='event' THEN (SELECT title FROM events WHERE id=m.target_id) ELSE (SELECT name FROM organizers WHERE id=m.target_id) END,'slug',CASE WHEN m.kind='event' THEN (SELECT slug FROM events WHERE id=m.target_id) ELSE (SELECT slug FROM organizers WHERE id=m.target_id) END)) FROM post_mentions m WHERE m.post_id=p.id),'[]')) doc
 FROM posts p JOIN accounts a ON a.id=p.account_id LEFT JOIN organizers o ON o.id=p.organizer_id;
CREATE OR REPLACE VIEW review_view AS SELECT r.id,r.event_id,r.account_id,e.organizer_id,r.hidden,
 jsonb_build_object('id',r.id,'event_id',r.event_id,'event_title',e.title,'event_slug',e.slug,'organizer_id',e.organizer_id,'account_id',r.account_id,'author',a.name,'avatar_url',a.avatar_url,'body',r.body,'created_at',r.created_at,'updated_at',r.updated_at,
 'verified',EXISTS(SELECT 1 FROM tickets t JOIN orders x ON x.id=t.order_id JOIN event_sessions s ON s.id=x.session_id WHERE x.account_id=r.account_id AND s.event_id=r.event_id AND t.checked_at IS NOT NULL),
 'helpful',(SELECT count(*) FROM review_helpful h WHERE h.review_id=r.id),
 'replies',COALESCE((SELECT jsonb_agg(jsonb_build_object('id',rr.id,'author',oo.name,'body',rr.body,'created_at',rr.created_at) ORDER BY rr.created_at) FROM review_replies rr JOIN organizers oo ON oo.account_id=rr.account_id WHERE rr.review_id=r.id),'[]')) doc FROM reviews r JOIN events e ON e.id=r.event_id JOIN accounts a ON a.id=r.account_id;
CREATE OR REPLACE VIEW product_view AS SELECT p.id,p.organizer_id,
 to_jsonb(p)||jsonb_build_object('organizer',o.doc,'event_title',e.title,'event_slug',e.slug) doc FROM products p JOIN organizer_view o ON o.id=p.organizer_id LEFT JOIN events e ON e.id=p.event_id;
