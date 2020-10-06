CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

CREATE OR REPLACE FUNCTION uuid_claim_create_tg() RETURNS trigger AS
$$
BEGIN
    IF tg_op = 'INSERT' THEN
        NEW.ClaimId = uuid_generate_v4();
        RETURN NEW;
    END IF;
END;
$$ LANGUAGE plpgsql;

CREATE TABLE public.claim
(
    ClaimId uuid PRIMARY KEY,
    Value   CHARACTER VARYING(30) NOT NULL,
    Key     CHARACTER VARYING(30) NOT NULL
);

CREATE INDEX idx_claim ON public.claim USING btree (ClaimId);

CREATE TRIGGER claimId_uuid_create
    BEFORE INSERT
    ON public.claim
    FOR EACH ROW
EXECUTE PROCEDURE uuid_claim_create_tg();

ALTER TABLE public.claim
    ADD CONSTRAINT claim_user FOREIGN KEY (ClaimId) REFERENCES public.user ("userid") ON DELETE CASCADE;

