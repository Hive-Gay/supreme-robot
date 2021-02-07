-- +migrate Up
CREATE TABLE "public"."sms_conversation_lines"
(
    id serial NOT NULL,
    account_sid character varying NOT NULL,
    api_version character varying NOT NULL,
    body character varying NOT NULL,
    direction character varying NOT NULL,
    error_code integer,
    error_message character varying,
    from_id integer NOT NULL,
    num_media integer NOT NULL,
    num_segments integer NOT NULL,
    price float,
    price_unit character varying,
    sid character varying NOT NULL,
    sms_conversation_id integer NOT NULL,
    status character varying NOT NULL,
    "timestamp" timestamp without time zone,
    to_id integer NOT NULL,
    created_at timestamp without time zone NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at timestamp without time zone NOT NULL DEFAULT CURRENT_TIMESTAMP,
    PRIMARY KEY (id),
    FOREIGN KEY (sms_conversation_id) REFERENCES sms_conversations (id) ON DELETE CASCADE,
    FOREIGN KEY (from_id) REFERENCES phone_numbers (id) ON DELETE RESTRICT,
    FOREIGN KEY (to_id) REFERENCES phone_numbers (id) ON DELETE RESTRICT
)
;

-- +migrate Down
DROP TABLE "public"."sms_conversation_lines";