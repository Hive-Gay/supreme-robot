-- +migrate Up
CREATE TABLE "public"."mxaliases"
(
    id serial NOT NULL,
    alias character varying NOT NULL,
    mxbox_id integer NOT NULL,
    PRIMARY KEY (id),
    UNIQUE (alias),
    FOREIGN KEY (mxbox_id) REFERENCES mxboxes (id) ON UPDATE RESTRICT ON DELETE RESTRICT
)
;

-- +migrate Down
DROP TABLE "public"."mxaliases";