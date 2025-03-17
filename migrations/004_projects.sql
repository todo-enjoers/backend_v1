CREATE TABLE IF NOT EXISTS projects
(
    "id" UUID NOT NULL UNIQUE,
    "name" VARCHAR NOT NULL,
    "created_by" UUID NOT NULL ,
    PRIMARY KEY ("id"),
    FOREIGN KEY ("created_by") REFERENCES users(id) ON DELETE CASCADE
);
---- create above / drop below ----

-- Write your migrate down statements here. If this migration is irreversible
-- Then delete the separator line above.
