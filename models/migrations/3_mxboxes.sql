-- +migrate Up
CREATE TABLE "public"."mxboxes"
(
    id serial NOT NULL,
    username character varying NOT NULL,
    domain character varying NOT NULL,
    password character varying NOT NULL,
    PRIMARY KEY (id),
    UNIQUE (username, domain)
)
;

-- +migrate Down
DROP TABLE "public"."mxboxes";