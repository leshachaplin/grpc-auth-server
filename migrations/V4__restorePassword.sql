CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

CREATE TABLE public.restore
(
    Id_Restore  uuid PRIMARY KEY,
    UuidRestore CHARACTER VARYING(60),
    Exp         timestamp
);

ALTER TABLE public.restore
    ADD CONSTRAINT restore_user FOREIGN KEY (Id_Restore) REFERENCES public.user (userid) ON DELETE CASCADE;