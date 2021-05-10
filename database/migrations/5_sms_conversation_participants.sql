-- +migrate Up
CREATE TABLE "public"."sms_conversation_participants"
(
    id serial NOT NULL,
    phone_number_id integer NOT NULL,
    sms_conversation_id integer NOT NULL,
    created_at timestamp without time zone NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at timestamp without time zone NOT NULL DEFAULT CURRENT_TIMESTAMP,
    PRIMARY KEY (id),
    UNIQUE (phone_number_id, sms_conversation_id),
    FOREIGN KEY (phone_number_id) REFERENCES phone_numbers (id) ON DELETE RESTRICT,
    FOREIGN KEY (sms_conversation_id) REFERENCES sms_conversations (id) ON DELETE CASCADE
)
;

-- +migrate Down
DROP TABLE "public"."sms_conversation_participants";