CREATE SEQUENCE IF NOT EXISTS expenses_id_seq;

CREATE TABLE "expenses" (
    "id" int4 NOT NULL DEFAULT nextval('expenses_id_seq'::regclass),
    "title" text,
    "amount" float,
    "note" text,
    "tags" text[],
    PRIMARY KEY ("id")
);

INSERT INTO "expenses" ("id", "title", "amount", "note", "tags") VALUES (1, 'test-title', 99.0, 'test-note', '{test}');

SELECT setval('"expenses_id_seq"', (SELECT MAX(id) FROM expenses));