-- Settings table to store user preferences
CREATE TABLE IF NOT EXISTS settings (
    user_id UUID PRIMARY KEY REFERENCES users(id) ON DELETE CASCADE,
    theme TEXT NOT NULL DEFAULT 'light',
    email_notifications BOOLEAN NOT NULL DEFAULT TRUE,
    feedback_alerts BOOLEAN NOT NULL DEFAULT TRUE,
    application_updates BOOLEAN NOT NULL DEFAULT TRUE,
    newsletter BOOLEAN NOT NULL DEFAULT FALSE,
    created_at TIMESTAMP NOT NULL DEFAULT now(),
    updated_at TIMESTAMP NOT NULL DEFAULT now()
);

CREATE INDEX IF NOT EXISTS idx_settings_user_id ON settings(user_id);

CREATE TRIGGER update_settings_updated_at
BEFORE UPDATE ON settings
FOR EACH ROW
EXECUTE FUNCTION update_updated_at_column();

