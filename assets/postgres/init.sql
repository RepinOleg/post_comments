CREATE TABLE IF NOT EXISTS posts (
                                     id SERIAL PRIMARY KEY,
                                     title TEXT NOT NULL,
                                     body TEXT NOT NULL,
                                     comments_disabled BOOLEAN DEFAULT FALSE,
                                     created_at TIMESTAMPTZ NOT NULL,
                                     updated_at TIMESTAMPTZ NOT NULL
);

CREATE TABLE IF NOT EXISTS comments (
                                        id SERIAL PRIMARY KEY,
                                        post_id INT NOT NULL REFERENCES posts(id) ON DELETE CASCADE,
                                        parent_id INT,
                                        body TEXT NOT NULL,
                                        created_at TIMESTAMPTZ NOT NULL,
                                        updated_at TIMESTAMPTZ NOT NULL
);