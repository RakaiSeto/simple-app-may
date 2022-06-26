CREATE TYPE status AS ENUM ('pending', 'success', 'failed')

CREATE TABLE public.product (
   id SERIAL primary key,
   name VARCHAR NOT NULL,
   description VARCHAR,
   price integer NOT NULL,
   created bigint,
   updated bigint,
);

CREATE TABLE public.user (
   id SERIAL PRIMARY KEY,
   uname VARCHAR,
   email VARCHAR,
   password VARCHAR,
   role VARCHAR,
   wallet bigint,
   created bigint,
   updated bigint,
);

CREATE TABLE public.order (
   id SERIAL PRIMARY KEY,
   userid integer not null,
   productid integer not null,
   quantity integer not null,
   totalprice integer default 0,
   created bigint,
   updated bigint,
);

CREATE TABLE public.queue (
   id SERIAL PRIMARY KEY,
   method VARCHAR,
   url VARCHAR,
   reqbody JSON,
   status status DEFAULT 'pending'::status,
   created bigint,
   updated bigint,
)