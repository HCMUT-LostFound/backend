CREATE TABLE items (
  id UUID PRIMARY KEY DEFAULT gen_random_uuid(),

  user_id UUID NOT NULL
    REFERENCES users(id)
    ON DELETE CASCADE,

  type TEXT NOT NULL
    CHECK (type IN ('lost', 'found')),

  title TEXT NOT NULL,
  description TEXT,

  image_urls TEXT[],

  location TEXT NOT NULL,
  campus TEXT NOT NULL,

  lost_at TIMESTAMPTZ,

  tags TEXT[],

  is_confirmed BOOLEAN NOT NULL DEFAULT FALSE,

  created_at TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE INDEX idx_items_user_id ON items(user_id);
CREATE INDEX idx_items_type ON items(type);
CREATE INDEX idx_items_is_confirmed ON items(is_confirmed);
