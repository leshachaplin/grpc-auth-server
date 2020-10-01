CREATE TABLE public.claim
(
    Id serial PRIMARY KEY ,
    Username CHARACTER VARYING(30),
    Value CHARACTER VARYING(30)  NOT NULL,
    Key CHARACTER VARYING(30)  NOT NULL
);

ALTER TABLE public.claim ADD CONSTRAINT user_claim FOREIGN KEY (Username) REFERENCES public.user ("username")  ON DELETE CASCADE;

