CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

CREATE OR REPLACE FUNCTION uuid_restoreId_create_tg() RETURNS trigger AS
$$
BEGIN
    IF tg_op = 'INSERT' THEN
        NEW.RestoreId = uuid_generate_v4();
        RETURN NEW;
    END IF;
END;
$$ LANGUAGE plpgsql;

CREATE OR REPLACE FUNCTION uuid_restore_create_tg() RETURNS trigger AS
$$
BEGIN
    IF tg_op = 'INSERT' THEN
        NEW.UuidRestore = uuid_generate_v4();
        RETURN NEW;
    END IF;
END;
$$ LANGUAGE plpgsql;

CREATE FUNCTION delete_old_restore_rows() RETURNS trigger
    LANGUAGE plpgsql
AS
$$
BEGIN
    DELETE FROM public.restore WHERE Exp < NOW() - INTERVAL '10 minutes';
    RETURN NULL;
END;
$$;

CREATE TABLE public.restore
(
    RestoreId   uuid PRIMARY KEY,
    UuidRestore CHARACTER VARYING(60),
    Exp         timestamp
);

CREATE INDEX idx_restore ON public.restore USING btree (RestoreId);
CREATE INDEX idx_restore ON public.restore USING btree (Exp);

CREATE TRIGGER trigger_delete_old_restore_rows
    AFTER INSERT
    ON public.restore
EXECUTE PROCEDURE delete_old_restore_rows();

CREATE TRIGGER restoreID_uuid_create
    BEFORE INSERT
    ON public.restore
    FOR EACH ROW
EXECUTE PROCEDURE uuid_restoreId_create_tg();

CREATE TRIGGER restore_uuid_create
    BEFORE INSERT
    ON public.restore
    FOR EACH ROW
EXECUTE PROCEDURE uuid_restore_create_tg();

ALTER TABLE public.restore
    ADD CONSTRAINT restore_user FOREIGN KEY (RestoreId) REFERENCES public.user (userid) ON DELETE CASCADE;