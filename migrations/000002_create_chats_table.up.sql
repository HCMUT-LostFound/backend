-- Create chats table
-- Mỗi chat được tạo khi user muốn liên hệ về một item
CREATE TABLE IF NOT EXISTS chats (
  id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  
  -- Item mà chat này liên quan đến
  item_id UUID NOT NULL
    REFERENCES items(id)
    ON DELETE CASCADE,
  
  -- User tạo chat (người muốn liên hệ)
  initiator_id UUID NOT NULL
    REFERENCES users(id)
    ON DELETE CASCADE,
  
  -- User sở hữu item (người được liên hệ)
  item_owner_id UUID NOT NULL
    REFERENCES users(id)
    ON DELETE CASCADE,
  
  -- Đảm bảo mỗi user chỉ có 1 chat với owner về 1 item
  UNIQUE(item_id, initiator_id),
  
  created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
  updated_at TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE INDEX IF NOT EXISTS idx_chats_item_id ON chats(item_id);
CREATE INDEX IF NOT EXISTS idx_chats_initiator_id ON chats(initiator_id);
CREATE INDEX IF NOT EXISTS idx_chats_item_owner_id ON chats(item_owner_id);
CREATE INDEX IF NOT EXISTS idx_chats_updated_at ON chats(updated_at DESC);

-- Create chat_messages table
CREATE TABLE IF NOT EXISTS chat_messages (
  id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  
  chat_id UUID NOT NULL
    REFERENCES chats(id)
    ON DELETE CASCADE,
  
  sender_id UUID NOT NULL
    REFERENCES users(id)
    ON DELETE CASCADE,
  
  content TEXT NOT NULL,
  
  created_at TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE INDEX IF NOT EXISTS idx_chat_messages_chat_id ON chat_messages(chat_id);
CREATE INDEX IF NOT EXISTS idx_chat_messages_created_at ON chat_messages(created_at DESC);
CREATE INDEX IF NOT EXISTS idx_chat_messages_sender_id ON chat_messages(sender_id);

-- Trigger để update updated_at khi có message mới
CREATE OR REPLACE FUNCTION update_chat_updated_at()
RETURNS TRIGGER AS $$
BEGIN
  UPDATE chats
  SET updated_at = NOW()
  WHERE id = NEW.chat_id;
  RETURN NEW;
END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER trigger_update_chat_updated_at
AFTER INSERT ON chat_messages
FOR EACH ROW
EXECUTE FUNCTION update_chat_updated_at();

