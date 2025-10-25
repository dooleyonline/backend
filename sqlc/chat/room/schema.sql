create schema "chat";
set search_path to "chat";

create table "chat"."room" (
  id uuid not null default gen_random_uuid (),
  participants uuid[] not null default '{}'::uuid[],
  messages bigint[] not null default '{}'::bigint[],
  constraint chat_room_pkey primary key (id)
) TABLESPACE pg_default;
