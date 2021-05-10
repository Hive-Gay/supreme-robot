-- +migrate Up
CREATE TABLE "public"."phone_numbers"
(
    id serial NOT NULL,
    num character varying UNIQUE NOT NULL,
    city character varying,
    country character varying,
    state character varying,
    zip character varying,
    created_at timestamp without time zone NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at timestamp without time zone NOT NULL DEFAULT CURRENT_TIMESTAMP,
    PRIMARY KEY (id)
)
;

-- +migrate Down
DROP TABLE "public"."phone_numbers";