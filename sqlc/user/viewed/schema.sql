create schema "user";
set search_path TO "user";

create table "user"."viewed" (
  user_id uuid not null,
  item_id bigint not null,
  created_at timestamp with time zone not null default now()
  -- constraint viewed_pkey primary key (user_id, item_id, created_at),
  -- constraint view_user_id_fkey foreign KEY (user_id) references "user"."user" (id) on update CASCADE on delete CASCADE,
  -- constraint viewed_item_id_fkey foreign KEY (item_id) references item.item (id) on update CASCADE on delete CASCADE
) TABLESPACE pg_default;
