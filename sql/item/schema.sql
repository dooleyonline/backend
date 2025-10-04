CREATE TABLE public.item (
  id bigint GENERATED ALWAYS AS IDENTITY NOT NULL,
  name text NOT NULL DEFAULT '',
  description text NOT NULL DEFAULT '',
  images text[] NOT NULL,
  price double precision NOT NULL DEFAULT '0'::double precision,
  condition smallint NOT NULL DEFAULT '0'::smallint,
  is_negotiable boolean NOT NULL DEFAULT false,
  posted_at timestamp with time zone NOT NULL DEFAULT now(),
  sold_at timestamp with time zone,
  views bigint NOT NULL DEFAULT '0'::bigint,
  CONSTRAINT item_pkey PRIMARY KEY (id)
);
