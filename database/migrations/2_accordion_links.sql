
-- +migrate Up
CREATE TABLE "public"."accordion_links" (
    id serial NOT NULL,
    accordion_header_id integer,
    title character varying NOT NULL,
    link character varying NOT NULL,
    created_at timestamp without time zone NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at timestamp without time zone NOT NULL DEFAULT CURRENT_TIMESTAMP,
    PRIMARY KEY ("id"),
    unique (accordion_header_id, title),
    FOREIGN KEY (accordion_header_id) REFERENCES accordion_headers (id) ON DELETE CASCADE
)
;

-- +migrate Down
DROP TABLE "public"."accordion_links";