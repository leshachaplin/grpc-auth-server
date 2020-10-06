CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

CREATE OR REPLACE FUNCTION uuid_confirmId_create_tg() RETURNS trigger AS
$$
BEGIN
    IF tg_op = 'INSERT' THEN
        NEW.ConfirmationId = uuid_generate_v4();
        RETURN NEW;
    END IF;
END;
$$ LANGUAGE plpgsql;

CREATE OR REPLACE FUNCTION uuid_confirmation_create_tg() RETURNS trigger AS
$$
BEGIN
    IF tg_op = 'INSERT' THEN
        NEW.UuidConfirmation = uuid_generate_v4();
        RETURN NEW;
    END IF;
END;
$$ LANGUAGE plpgsql;

CREATE FUNCTION delete_old_confirmation_rows() RETURNS trigger
    LANGUAGE plpgsql
AS
$$
BEGIN
    DELETE FROM public.confirmation WHERE Expiration < NOW() - INTERVAL '10 minutes';
    RETURN NULL;
END;
$$;

CREATE TABLE public.confirmation
(
    ConfirmationId   uuid PRIMARY KEY,
    UuidConfirmation uuid,
    Expiration       timestamp
);

CREATE INDEX idx_confirmation ON public.confirmation USING btree (ConfirmationId);
CREATE INDEX idx_confirmation ON public.confirmation USING btree (Expiration);

CREATE TRIGGER trigger_delete_old_confirmation_rows
    AFTER INSERT
    ON public.confirmation
EXECUTE PROCEDURE delete_old_confirmation_rows();


CREATE TRIGGER confirmID_uuid_create
    BEFORE INSERT
    ON public.confirmation
    FOR EACH ROW
EXECUTE PROCEDURE uuid_confirmId_create_tg();

CREATE TRIGGER confirm_uuid_create
    BEFORE INSERT
    ON public.confirmation
    FOR EACH ROW
EXECUTE PROCEDURE uuid_confirmation_create_tg();

ALTER TABLE public.confirmation
    ADD CONSTRAINT confirmed_user FOREIGN KEY (ConfirmationId) REFERENCES public.user (userid) ON DELETE CASCADE;

