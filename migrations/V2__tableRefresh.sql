CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

CREATE OR REPLACE FUNCTION uuid_refresh_create_tg() RETURNS trigger AS
$$
BEGIN
    IF tg_op = 'INSERT' THEN
        NEW.RefreshId = uuid_generate_v4();
        RETURN NEW;
    END IF;
END;
$$ LANGUAGE plpgsql;

CREATE FUNCTION delete_old_rows() RETURNS trigger
    LANGUAGE plpgsql
AS
$$
BEGIN
    DELETE FROM public.refresh WHERE Expiration < NOW() - INTERVAL '72 hours';
    RETURN NULL;
END;
$$;

CREATE TABLE public.refresh
(
    RefreshId  uuid PRIMARY KEY,
    Expiration timestamp,
    Token      CHARACTER VARYING(256) NOT NULL
);

CREATE INDEX idx_refresh ON public.refresh USING btree (RefreshId);
CREATE INDEX idx_refresh ON public.refresh USING btree (Expiration);

CREATE TRIGGER trigger_delete_old_rows
    AFTER INSERT
    ON public.refresh
EXECUTE PROCEDURE delete_old_rows();

CREATE TRIGGER some_table_uuid_create
    BEFORE INSERT
    ON public.refresh
    FOR EACH ROW
EXECUTE PROCEDURE uuid_refresh_create_tg();

ALTER TABLE public.refresh
    ADD CONSTRAINT refresh_user FOREIGN KEY (RefreshId) REFERENCES public.user (userid) ON DELETE CASCADE;
