CREATE TABLE public.restore
(
    Id          serial PRIMARY KEY,
    Username    CHARACTER VARYING(30),
    UuidRestore CHARACTER VARYING(60),
    Exp         timestamp
);

ALTER TABLE public.restore
    ADD CONSTRAINT user_restore FOREIGN KEY (Username) REFERENCES public.user (username) ON DELETE CASCADE;