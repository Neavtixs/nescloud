CREATE TABLE folders (
    id UUID PRIMARY KEY NOT NULL,
    owner_id UUID NOT NULL,
    parent_folder_id UUID NULL,
    name VARCHAR(255) NOT NULL,
    deleted_at TIMESTAMPTZ NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    CONSTRAINT fk_folders_owner FOREIGN KEY (owner_id) REFERENCES users(id) ON DELETE CASCADE,
    CONSTRAINT fk_folders_parent FOREIGN KEY (parent_folder_id) REFERENCES folders(id) ON DELETE CASCADE
);

CREATE INDEX idx_folders_owner_id ON folders(owner_id);
CREATE INDEX idx_folders_parent_folder_id ON folders(parent_folder_id);
CREATE INDEX idx_folders_deleted_at ON folders(deleted_at);
