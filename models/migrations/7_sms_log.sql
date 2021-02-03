-- +migrate Up
CREATE TABLE "public"."sms_log"
(
    id serial NOT NULL,
    account_sid character varying NOT NULL,
    api_version character varying NOT NULL,
    body character varying NOT NULL,
    date_created timestamp without time zone,
    date_sent timestamp without time zone,
    date_updated timestamp without time zone,
    direction character varying NOT NULL,
    error_code integer,
    error_message character varying,
    from_id integer NOT NULL,
    num_media integer NOT NULL,
    num_segments integer NOT NULL,
    price float,
    price_unit character varying,
    sid character varying NOT NULL,
    status character varying NOT NULL,
    to_id integer NOT NULL,

    created_at timestamp without time zone NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at timestamp without time zone NOT NULL DEFAULT CURRENT_TIMESTAMP,

    PRIMARY KEY (id),
    FOREIGN KEY (from_id) REFERENCES phone_numbers (id) ON DELETE RESTRICT,
    FOREIGN KEY (to_id) REFERENCES phone_numbers (id) ON DELETE RESTRICT
)
;

-- +migrate Down
DROP TABLE "public"."sms_log";