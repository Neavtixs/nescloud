CREATE TABLE public_links (
    id UUID PRIMARY KEY NOT NULL,
    file_id UUID NOT NULL UNIQUE,
    token VARCHAR(255) NOT NULL UNIQUE,
    expired_at TIMESTAMPTZ NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    CONSTRAINT fk_public_links_file FOREIGN KEY (file_id) REFERENCES files(id) ON DELETE CASCADE
);

CREATE INDEX idx_public_links_token ON public_links(token);
