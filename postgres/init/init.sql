CREATE TABLE IF NOT EXISTS "schedules" (
  "id" VARCHAR(255) PRIMARY KEY,
  "title" VARCHAR(255) NOT NULL,
  "description" TEXT NOT NULL,
  "start_date" TIMESTAMP NOT NULL,
  "end_date" TIMESTAMP NOT NULL,
  "allday" BOOLEAN NOT NULL,
  "tags" TEXT[] NOT NULL,
  "properties" TEXT[] NOT NULL,
  "is_public" BOOLEAN NOT NULL,
  "created_at" TIMESTAMP NOT NULL,
  "updated_at" TIMESTAMP  NOT NULL
);