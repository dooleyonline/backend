create schema "user";
set search_path TO "user";

create table "user"."user" (
  email text not null default ''::text,
  password text not null default ''::text,
  first_name text not null default ''::text,
  last_name text not null default ''::text,
  id uuid not null default gen_random_uuid (),
  constraint user_pkey primary key (id),
  constraint user_email_key unique (email)
) TABLESPACE pg_default;
