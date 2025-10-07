CREATE TABLE public.category (
  name text NOT NULL,
  subcategory text[] NOT NULL,
  icon text NOT NULL,
  CONSTRAINT category_pkey PRIMARY KEY (name)
);
