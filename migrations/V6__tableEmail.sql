CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

CREATE TABLE public.email
(
    NumberOfAttempts int,
    FromEmail        CHARACTER VARYING(64) NOT NULL PRIMARY KEY ,
    ToEmail          CHARACTER VARYING(64) NOT NULL,
    Status           bool                  NOT NULL,
    CreatedAt        timestamp,
    SentAt           timestamp
);