CREATE TABLE category (
  name text NOT NULL,
  subcategory text[] NOT NULL,
  icon text,
  CONSTRAINT category_pkey PRIMARY KEY (name)
);
