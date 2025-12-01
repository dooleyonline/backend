create schema "chat";
set search_path to "chat";

create table "chat"."room" (
  id uuid not null default gen_random_uuid (),
  message_count bigint not null default 0,
  participants uuid[] not null default '{}'::uuid[],
  constraint chat_room_pkey primary key (id)
) TABLESPACE pg_default;
