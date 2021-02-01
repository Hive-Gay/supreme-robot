-- +migrate Up
CREATE TABLE "public"."sms_webook_log"
(
    id serial NOT NULL,
    url character varying NOT NULL,
    params character varying NOT NULL,
    i_token uuid UNIQUE NOT NULL,
    signature character varying NOT NULL,
    verified bool,
    PRIMARY KEY (id)
)
;

-- +migrate Down
DROP TABLE "public"."sms_webook_log";