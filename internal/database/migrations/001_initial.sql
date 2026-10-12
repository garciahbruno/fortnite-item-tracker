CREATE TABLE users (
    id BIGINT GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
    discord_id TEXT NOT NULL UNIQUE CHECK (discord_id ~ '^[0-9]+$'),
    username TEXT NOT NULL,
    alerts_enabled BOOLEAN NOT NULL DEFAULT FALSE,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE TABLE cosmetics (
    id TEXT PRIMARY KEY,
    name TEXT NOT NULL,
    cosmetic_type TEXT NOT NULL,
    rarity TEXT NOT NULL DEFAULT '',
    image_url TEXT,
    updated_at TIMESTAMPTZ NOT NULL DEFAULT now()
);
CREATE INDEX cosmetics_name_lower_idx ON cosmetics (lower(name));

CREATE TABLE shop_snapshots (
    id BIGINT GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
    shop_date DATE NOT NULL,
    content_hash TEXT NOT NULL,
    fetched_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    UNIQUE (shop_date, content_hash)
);
CREATE INDEX shop_snapshots_fetched_idx ON shop_snapshots (fetched_at DESC);

-- Offers can contain multiple cosmetics (bundles).
CREATE TABLE shop_offers (
    snapshot_id BIGINT NOT NULL REFERENCES shop_snapshots(id) ON DELETE CASCADE,
    offer_id TEXT NOT NULL,
    price INTEGER NOT NULL CHECK (price >= 0),
    PRIMARY KEY (snapshot_id, offer_id)
);
CREATE TABLE shop_offer_cosmetics (
    snapshot_id BIGINT NOT NULL,
    offer_id TEXT NOT NULL,
    cosmetic_id TEXT NOT NULL REFERENCES cosmetics(id),
    PRIMARY KEY (snapshot_id, offer_id, cosmetic_id),
    FOREIGN KEY (snapshot_id, offer_id) REFERENCES shop_offers(snapshot_id, offer_id) ON DELETE CASCADE
);
CREATE INDEX shop_offer_cosmetics_cosmetic_idx ON shop_offer_cosmetics (cosmetic_id);

CREATE TABLE watchlist_items (
    user_id BIGINT NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    cosmetic_id TEXT NOT NULL REFERENCES cosmetics(id) ON DELETE CASCADE,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    PRIMARY KEY (user_id, cosmetic_id)
);
CREATE INDEX watchlist_cosmetic_idx ON watchlist_items (cosmetic_id);

-- One notification per user/cosmetic/shop day, even across refreshed snapshots.
CREATE TABLE notifications (
    id BIGINT GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
    user_id BIGINT NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    cosmetic_id TEXT NOT NULL REFERENCES cosmetics(id),
    shop_date DATE NOT NULL,
    status TEXT NOT NULL DEFAULT 'pending' CHECK (status IN ('pending', 'sending', 'sent', 'failed')),
    attempts INTEGER NOT NULL DEFAULT 0 CHECK (attempts >= 0),
    next_attempt_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    sent_at TIMESTAMPTZ,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    UNIQUE (user_id, cosmetic_id, shop_date)
);
CREATE INDEX notifications_pending_idx ON notifications (next_attempt_at) WHERE status IN ('pending', 'failed');
