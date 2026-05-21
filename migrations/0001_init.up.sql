CREATE EXTENSION IF NOT EXISTS "pgcrypto";

CREATE TABLE IF NOT EXISTS users (
    id          uuid PRIMARY KEY DEFAULT gen_random_uuid(),
    email       text NOT NULL UNIQUE,
    password    text NOT NULL,
    role        text NOT NULL,
    created_at  timestamptz NOT NULL DEFAULT NOW(),
    updated_at  timestamptz NOT NULL DEFAULT NOW()
);

CREATE TABLE IF NOT EXISTS artists (
    id            uuid PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id       uuid NOT NULL UNIQUE REFERENCES users(id) ON DELETE CASCADE,
    display_name  text NOT NULL,
    bio           text,
    avatar_url    text,
    created_at    timestamptz NOT NULL DEFAULT NOW(),
    updated_at    timestamptz NOT NULL DEFAULT NOW()
);

CREATE TABLE IF NOT EXISTS artworks (
    id                          uuid PRIMARY KEY DEFAULT gen_random_uuid(),
    artist_id                   uuid NOT NULL REFERENCES artists(id) ON DELETE CASCADE,
    title                       text NOT NULL,
    description                 text NOT NULL DEFAULT '',
    status                      text NOT NULL DEFAULT 'draft',
    tags                        text[] NOT NULL DEFAULT '{}',
    preview_image_url           text,
    availability_for_exchange   boolean NOT NULL DEFAULT false,
    created_at                  timestamptz NOT NULL DEFAULT NOW(),
    updated_at                  timestamptz NOT NULL DEFAULT NOW(),
    published_at                timestamptz
);

CREATE TABLE IF NOT EXISTS artwork_images (
    id          uuid PRIMARY KEY DEFAULT gen_random_uuid(),
    artwork_id  uuid NOT NULL REFERENCES artworks(id) ON DELETE CASCADE,
    url         text NOT NULL,
    alt_text    text,
    position    integer NOT NULL DEFAULT 0
);

CREATE INDEX IF NOT EXISTS idx_artwork_images_artwork_id ON artwork_images(artwork_id);

CREATE TABLE IF NOT EXISTS likes (
    user_id     uuid NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    artwork_id  uuid NOT NULL REFERENCES artworks(id) ON DELETE CASCADE,
    created_at  timestamptz NOT NULL DEFAULT NOW(),
    PRIMARY KEY (user_id, artwork_id)
);

CREATE TABLE IF NOT EXISTS exchange_requests (
    id            uuid PRIMARY KEY DEFAULT gen_random_uuid(),
    artwork_id    uuid NOT NULL REFERENCES artworks(id) ON DELETE CASCADE,
    requester_id  uuid NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    owner_id      uuid NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    status        text NOT NULL DEFAULT 'pending',
    message       text,
    created_at    timestamptz NOT NULL DEFAULT NOW(),
    updated_at    timestamptz NOT NULL DEFAULT NOW()
);
