-- +migrate Up
CREATE TABLE "public"."mxboxes"
(
    id serial NOT NULL,
    username character varying NOT NULL,
    domain character varying NOT NULL,
    password character varying NOT NULL,
    created_at timestamp without time zone NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at timestamp without time zone NOT NULL DEFAULT CURRENT_TIMESTAMP,
    PRIMARY KEY (id),
    UNIQUE (username, domain)
)
;

-- +migrate Down
DROP TABLE "public"."mxboxes";