CREATE TABLE public.user (
email text NOT NULL DEFAULT ''::text UNIQUE,
password text NOT NULL DEFAULT ''::text,
liked_items bigint[] NOT NULL DEFAULT '{}'::bigint[],
CONSTRAINT user_pkey PRIMARY KEY (email)
);
