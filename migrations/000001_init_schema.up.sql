-- enable uuid
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

-- users (Clerk)
CREATE TABLE users (
  id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
  clerk_user_id TEXT UNIQUE NOT NULL,
  full_name TEXT,
  avatar_url TEXT,
  created_at TIMESTAMP DEFAULT now()
);

-- campuses
CREATE TABLE campuses (
  id SERIAL PRIMARY KEY,
  name TEXT NOT NULL
);

-- items
CREATE TABLE items (
  id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
  user_id UUID REFERENCES users(id) ON DELETE CASCADE,

  type TEXT CHECK (type IN ('lost', 'found')) NOT NULL,
  name TEXT NOT NULL,
  description TEXT,

  last_seen_location TEXT,
  campus_id INT REFERENCES campuses(id),

  lost_at TIMESTAMP NOT NULL,
  confirmed BOOLEAN DEFAULT false,

  created_at TIMESTAMP DEFAULT now()
);

-- item images
CREATE TABLE item_images (
  id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
  item_id UUID REFERENCES items(id) ON DELETE CASCADE,
  image_url TEXT NOT NULL
);

-- tags
CREATE TABLE tags (
  id SERIAL PRIMARY KEY,
  name TEXT UNIQUE NOT NULL
);

-- item_tags
CREATE TABLE item_tags (
  item_id UUID REFERENCES items(id) ON DELETE CASCADE,
  tag_id INT REFERENCES tags(id) ON DELETE CASCADE,
  PRIMARY KEY (item_id, tag_id)
);
