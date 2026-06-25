CREATE TABLE files (
    id UUID PRIMARY KEY NOT NULL,
    owner_id UUID NOT NULL,
    folder_id UUID NOT NULL,
    original_name VARCHAR(255) NOT NULL,
    storage_key TEXT NOT NULL,
    mime_type VARCHAR(255) NOT NULL,
    extension VARCHAR(30) NOT NULL,
    size BIGINT NOT NULL,
    checksum VARCHAR(64) NOT NULL,
    upload_status VARCHAR(30) NOT NULL DEFAULT 'uploading',
    deleted_at TIMESTAMPTZ NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    CONSTRAINT fk_files_owner FOREIGN KEY (owner_id) REFERENCES users(id) ON DELETE CASCADE,
    CONSTRAINT fk_files_folder FOREIGN KEY (folder_id) REFERENCES folders(id) ON DELETE CASCADE
);

CREATE INDEX idx_files_owner_id ON files(owner_id);
CREATE INDEX idx_files_folder_id ON files(folder_id);
CREATE INDEX idx_files_deleted_at ON files(deleted_at);
CREATE INDEX idx_files_original_name ON files(original_name);
CREATE INDEX idx_files_storage_key ON files(storage_key);
