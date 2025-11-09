create schema "user";
set search_path TO "user";


create table "user".verify (
  id uuid not null default gen_random_uuid (),
  user_id uuid not null,
  expired_at timestamp with time zone not null default (now() + '00:10:00'::interval),
  constraint verify_pkey primary key (id)
  -- constraint verify_user_id_fkey foreign KEY (user_id) references "user"."user" (id) on update CASCADE on delete CASCADE
) TABLESPACE pg_default;