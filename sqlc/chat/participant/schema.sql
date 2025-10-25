create schema "chat";
set search_path to "chat";

create table "chat"."participant" (
  room_id uuid not null,
  user_id uuid not null default gen_random_uuid (),
  last_read_message_id bigint null,
  constraint chat_participant_pkey primary key (room_id, user_id)
  -- constraint participant_room_id_fkey foreign KEY (room_id) references chat.room (id) on update CASCADE on delete CASCADE,
  -- constraint participant_user_id_fkey foreign KEY (user_id) references "user"."user" (id) on update CASCADE on delete set default
) TABLESPACE pg_default;
