CREATE TABLE public.user
(
    ID        serial,
    Username  CHARACTER VARYING(64) PRIMARY KEY,
    Confirmed boolean NOT NULL,
    Email     CHARACTER VARYING(64),
    Password  bytea   NOT NULL,
    CreatedAt timestamp,
    UpdatedAt timestamp
)

