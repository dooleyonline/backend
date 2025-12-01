create schema if not exists "chat";
set search_path to "chat";

create table "chat"."message" (
  room_id uuid not null,
  sent_by uuid not null,
  body text not null,
  id bigint not null,
  edited boolean not null default false,
  sent_at timestamp with time zone not null default now(),
  constraint chat_message_pkey primary key (id)
  -- constraint message_room_id_fkey foreign KEY (room_id) references chat.room (id) on update CASCADE on delete CASCADE,
  -- constraint message_sender_id_fkey foreign KEY (sender_id) references "user"."user" (id) on update CASCADE on delete set default
) TABLESPACE pg_default;
