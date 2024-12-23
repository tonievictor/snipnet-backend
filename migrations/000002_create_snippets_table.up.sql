CREATE TABLE IF NOT EXISTS snippets (
	id TEXT PRIMARY KEY NOT NULL UNIQUE,
	user_id TEXT NOT NULL,
	title TEXT NOT NULL DEFAULT 'Untitled',
	description TEXT NOT NULL DEFAULT '',
	language VARCHAR(20) NOT NULL DEFAULT '',
	code TEXT NOT NULL DEFAULT '',
	document tsvector GENERATED ALWAYS AS (to_tsvector('english', title || ' ' || description || ' ' || code)) STORED,
	is_public BOOLEAN NOT NULL DEFAULT TRUE,
	created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
	updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
	FOREIGN KEY (user_id) REFERENCES users (id)
);

CREATE INDEX IF NOT EXISTS document_idx ON snippets USING GIN(document);
