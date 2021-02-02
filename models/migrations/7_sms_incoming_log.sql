-- +migrate Up
CREATE TABLE "public"."sms_incoming_log"
(
    id serial NOT NULL,
    account_sid character varying NOT NULL,
    api_version character varying NOT NULL,
    body character varying NOT NULL,
    direction character varying,
    from_id integer NOT NULL,
    message_sid character varying NOT NULL,
    num_media integer NOT NULL,
    num_segments integer NOT NULL,
    sms_message_sid character varying NOT NULL,
    sms_sid character varying NOT NULL,
    sms_status character varying NOT NULL,
    to_id integer NOT NULL,

    created_at timestamp without time zone NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at timestamp without time zone NOT NULL DEFAULT CURRENT_TIMESTAMP,
    PRIMARY KEY (id),
    FOREIGN KEY (from_id) REFERENCES phone_numbers (id) ON DELETE RESTRICT,
    FOREIGN KEY (to_id) REFERENCES phone_numbers (id) ON DELETE RESTRICT
)
;

-- +migrate Down
DROP TABLE "public"."sms_incoming_log";