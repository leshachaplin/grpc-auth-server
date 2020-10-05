CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

CREATE TABLE public.refresh
(
    Id_Refresh uuid PRIMARY KEY,
    Expiration timestamp,
    Token      CHARACTER VARYING(256) NOT NULL
);