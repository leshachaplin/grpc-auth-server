CREATE TABLE public.confirmation
(
    Id               serial PRIMARY KEY,
    Username         CHARACTER VARYING(30),
    UuidConfirmation CHARACTER VARYING(60),
    Expiration       timestamp
);

ALTER TABLE public.confirmation
    ADD CONSTRAINT user_confirmed FOREIGN KEY (Username) REFERENCES public.user (username) ON DELETE CASCADE;

