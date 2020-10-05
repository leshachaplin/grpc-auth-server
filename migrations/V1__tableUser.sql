CREATE EXTENSION IF NOT EXISTS pgcrypto;
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

CREATE OR REPLACE FUNCTION trigger_set_timestamp()
    RETURNS TRIGGER AS
$$
BEGIN
    NEW.UpdatedAt = NOW();
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

CREATE OR REPLACE FUNCTION hash_update_tg() RETURNS trigger AS
$$
BEGIN
    IF tg_op = 'INSERT' OR tg_op = 'UPDATE' THEN
        NEW.password = encode(digest(NEW.password, 'sha256'), 'hex');
        RETURN NEW;
    END IF;
END;
$$ LANGUAGE plpgsql;

CREATE TABLE public.user
(
    UserId          uuid PRIMARY KEY,
    Username        CHARACTER VARYING(64),
    Confirmed       boolean     NOT NULL,
    Email           CHARACTER VARYING(64),
    Password        text        NOT NULL,
    CreatedAt       TIMESTAMPTZ NOT NULL DEFAULT (NOW() at time zone 'Europe/Minsk'),
    UpdatedAt       TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    Id_Refresh      uuid,
    Id_Confirmation uuid,
    Id_Restore      uuid,
    Id_Claim        uuid
);

CREATE TRIGGER some_table_hash_insert
    BEFORE INSERT
    ON public.user
    FOR EACH ROW
EXECUTE PROCEDURE hash_update_tg();

CREATE TRIGGER some_table_hash_update
    BEFORE UPDATE
    ON public.user
    FOR EACH ROW
    WHEN ( NEW.password IS DISTINCT FROM OLD.password )
EXECUTE PROCEDURE hash_update_tg();

CREATE TRIGGER set_timestamp
    BEFORE UPDATE
    ON public.user
    FOR EACH ROW
EXECUTE PROCEDURE trigger_set_timestamp();

CREATE TRIGGER check_is_valid_hash
    BEFORE UPDATE
    ON public.user
    FOR EACH ROW
EXECUTE PROCEDURE trigger_set_timestamp();

ALTER TABLE public.user
    ADD CONSTRAINT user_refresh FOREIGN KEY (UserId) REFERENCES public.refresh ("id_refresh");

ALTER TABLE public.user
    ADD CONSTRAINT user_confirmation FOREIGN KEY (UserId) REFERENCES public.confirmation ();

ALTER TABLE public.user
    ADD CONSTRAINT user_refresh FOREIGN KEY (UserId) REFERENCES public.refresh ("id_refresh");

ALTER TABLE public.user
    ADD CONSTRAINT user_refresh FOREIGN KEY (UserId) REFERENCES public.refresh ("id_refresh");