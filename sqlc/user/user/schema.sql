create schema "user";
set search_path TO "user";

create table "user"."user" (
  email text not null default ''::text,
  password text not null default ''::text,
  first_name text not null default ''::text,
  last_name text not null default ''::text,
  id uuid not null default gen_random_uuid (),
  verified boolean not null default false,
  avatar uuid not null default '00000000-0000-0000-0000-000000000000'::uuid,
  constraint user_pkey primary key (id),
  constraint user_email_key unique (email)
) TABLESPACE pg_default;

