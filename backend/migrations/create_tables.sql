CREATE TABLE "ships" (
     "id" bigserial PRIMARY KEY,
     "mmsi" bigint NOT NULL UNIQUE,
     "name" varchar,
     "latitude" double precision NOT NULL,
     "longitude" double precision NOT NULL,
     "last_updated" timestamptz NOT NULL DEFAULT (now())
);

CREATE INDEX ON "ships" ("mmsi");

CREATE INDEX ON "ships" ("name");

COMMENT ON COLUMN "ships"."name" IS 'may be empty';
