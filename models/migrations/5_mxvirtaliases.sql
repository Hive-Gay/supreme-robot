-- +migrate Up
CREATE TABLE "public"."mxvirtaliases"
(
    id serial NOT NULL,
    alias character varying NOT NULL,
    mxbox_id integer NOT NULL,
    PRIMARY KEY (id),
    UNIQUE (alias),
    FOREIGN KEY (mxbox_id) REFERENCES public.mxboxes (id) ON UPDATE RESTRICT ON DELETE RESTRICT
)
;

-- +migrate Down
DROP TABLE "public"."mxvirtaliases";