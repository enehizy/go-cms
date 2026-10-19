CREATE TABLE IF NOT EXISTS media(
    id SERIAL PRIMARY KEY,
    uri VARCHAR(255) NOT NULL UNIQUE,
    mime_type VARCHAR(12) CHECK  (mime_type IN ('image/jpg','image/png','video/mp4','video/webp'))
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP WITH TIME ZONE 
)







