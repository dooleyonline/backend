CREATE TABLE public.user (
email text NOT NULL DEFAULT ''::text UNIQUE,
password text NOT NULL DEFAULT ''::text,
liked_items bigint[] NOT NULL DEFAULT '{}'::bigint[],
first_name text not null default ''::text,
last_name text not null default ''::text,
id uuid not null default gen_random_uuid (),
constraint user_pkey primary key (id),
constraint user_email_key unique (email),
constraint user_id_key unique (id)
);
