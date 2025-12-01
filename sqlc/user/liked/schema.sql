create schema IF NOT EXISTS "user";
set search_path TO "user";

create table "user"."liked" (
  user_id uuid not null default gen_random_uuid (),
  item_id bigint not null,
  created_at timestamp with time zone not null default now()
);
