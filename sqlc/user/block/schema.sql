create schema "user";
set search_path TO "user";

create table "user"."block" (
  blocker_id uuid not null default gen_random_uuid (),
  blocked_id uuid not null default gen_random_uuid (),
  created_at timestamp without time zone null default now()
  -- constraint block_pkey primary key (blocker_id, blocked_id),
  -- constraint block_blocked_id_fkey foreign KEY (blocked_id) references "user"."user" (id) on update CASCADE on delete CASCADE,
  -- constraint block_blocker_id_fkey foreign KEY (blocker_id) references "user"."user" (id) on update CASCADE on delete CASCADE
);
