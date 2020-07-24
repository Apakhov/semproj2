DROP FUNCTION IF EXISTS change_deleted;
DROP FUNCTION IF EXISTS check_update_id;
DROP FUNCTION IF EXISTS change_active;

CREATE FUNCTION check_update_id() RETURNS TRIGGER AS $$
    BEGIN
        IF OLD.id <> NEW.id THEN
            RAISE EXCEPTION 'cannot change id';
        END IF;
        RETURN NEW;
    END;
$$ LANGUAGE plpgsql;

CREATE FUNCTION change_deleted() RETURNS TRIGGER AS $$
    BEGIN
         EXECUTE format('UPDATE %I.%I SET deleted=true WHERE id=%s;', TG_TABLE_SCHEMA, TG_TABLE_NAME, OLD.id);
    RETURN NULL;
    END;
$$ LANGUAGE plpgsql;

CREATE FUNCTION change_active() RETURNS TRIGGER AS $$
    BEGIN
         EXECUTE format('UPDATE %I.%I SET active=false WHERE id=%s;', TG_TABLE_SCHEMA, TG_TABLE_NAME, OLD.id);
    RETURN NULL;
    END;
$$ LANGUAGE plpgsql;