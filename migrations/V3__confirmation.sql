CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

CREATE TABLE public.confirmation
(
    Id_Confirmation  uuid PRIMARY KEY,
    UuidConfirmation CHARACTER VARYING(60),
    Expiration       timestamp
);

ALTER TABLE public.confirmation
    ADD CONSTRAINT confirmed_user FOREIGN KEY (Id_Confirmation) REFERENCES public.user (userid) ON DELETE CASCADE;

