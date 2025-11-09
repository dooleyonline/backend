create schema "user";
set search_path TO "user";

create table "user"."user" (
  email text not null default ''::text,
  password text not null default ''::text,
  liked_items bigint[] not null default '{}'::bigint[],
  first_name text not null default ''::text,
  last_name text not null default ''::text,
  id uuid not null default gen_random_uuid (),
  avatar uuid not null default '92275607-f87d-4a35-8185-27d98c956e93'::uuid,
  constraint user_pkey primary key (id),
  constraint user_email_key unique (email)
) TABLESPACE pg_default;
