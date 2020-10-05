CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

CREATE TABLE public.claim
(
    Id_Claim uuid PRIMARY KEY,
    Value    CHARACTER VARYING(30) NOT NULL,
    Key      CHARACTER VARYING(30) NOT NULL
);

ALTER TABLE public.claim
    ADD CONSTRAINT claim_user FOREIGN KEY (Id_Claim) REFERENCES public.user ("userid") ON DELETE CASCADE;

