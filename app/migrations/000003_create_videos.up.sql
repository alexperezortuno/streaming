CREATE TABLE IF NOT EXISTS videos (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    name VARCHAR(255) NOT NULL,
    list_id UUID REFERENCES lists(id) ON DELETE SET NULL,
    file_path VARCHAR(500) NOT NULL,
    mime_type VARCHAR(100) NOT NULL DEFAULT 'video/mp4',
    status VARCHAR(20) NOT NULL DEFAULT 'uploading',
    duration DOUBLE PRECISION,
    size BIGINT NOT NULL DEFAULT 0,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_videos_list_id ON videos(list_id);
CREATE INDEX idx_videos_status ON videos(status);
CREATE INDEX idx_videos_name ON videos(name);
