-- +migrate Up
CREATE TABLE "public"."sms_conversations"
(
    id serial NOT NULL,
    name character varying NOT NULL,
    created_at timestamp without time zone NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at timestamp without time zone NOT NULL DEFAULT CURRENT_TIMESTAMP,
    PRIMARY KEY (id)
)
;

-- +migrate Down
DROP TABLE "public"."sms_conversations";