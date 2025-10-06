CREATE TABLE item (
  id bigint GENERATED ALWAYS AS IDENTITY NOT NULL,
  name text NOT NULL DEFAULT ''::text,
  description text NOT NULL DEFAULT ''::text,
  images text[] NOT NULL,
  price double precision NOT NULL DEFAULT '0'::double precision,
  condition smallint NOT NULL DEFAULT '0'::smallint,
  is_negotiable boolean NOT NULL DEFAULT false,
  posted_at timestamp with time zone NOT NULL DEFAULT now(),
  sold_at timestamp with time zone,
  views bigint NOT NULL DEFAULT '0'::bigint,
  category text NOT NULL,
  sub_category text NOT NULL,
  fts tsvector GENERATED ALWAYS AS (to_tsvector('english'::regconfig, ((description || ' '::text) || name))) STORED,
  CONSTRAINT item_pkey PRIMARY KEY (id),
  CONSTRAINT item_category_fkey FOREIGN KEY (category) REFERENCES category(name)
);
