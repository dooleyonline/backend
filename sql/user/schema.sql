CREATE TABLE "user" (
  email text NOT NULL DEFAULT ''::text,
  password text NOT NULL DEFAULT ''::text,
  CONSTRAINT User_pkey PRIMARY KEY (email)
);
