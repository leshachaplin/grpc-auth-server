CREATE TABLE public.refresh
(
    Id         serial PRIMARY KEY,
    Username   CHARACTER VARYING(30),
    Expiration timestamp,
    Token      CHARACTER VARYING(256) NOT NULL
);

ALTER TABLE public.refresh
    ADD CONSTRAINT user_refresh FOREIGN KEY (Username) REFERENCES public.user (Username) ON DELETE CASCADE;
