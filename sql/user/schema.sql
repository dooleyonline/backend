CREATE TABLE public.user (
email text NOT NULL DEFAULT ''::text UNIQUE,
password text NOT NULL DEFAULT ''::text,
CONSTRAINT user_pkey PRIMARY KEY (email)
);
