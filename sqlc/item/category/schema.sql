create schema item;
set search_path to item;

create table item.category (
  name text not null default 'Other'::text,
  subcategory text[] not null,
  icon text not null,
  constraint category_pkey primary key (name)
) TABLESPACE pg_default;
