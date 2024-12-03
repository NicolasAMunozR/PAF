
CREATE OR REPLACE FUNCTION fn_eliminar_historial()
RETURNS TRIGGER AS $$
BEGIN
    -- Si se está actualizando el registro y el campo deleted_at no es NULL
    IF NEW.deleted_at IS NOT NULL THEN
        DELETE FROM historial_paf_aceptadas WHERE id = OLD.id;
    END IF;
    
    -- Si se está insertando un registro y el campo deleted_at no es NULL
    IF (TG_OP = 'INSERT') AND (NEW.deleted_at IS NOT NULL) THEN
        DELETE FROM historial_paf_aceptadas WHERE id = NEW.id;
    END IF;

    -- Si se está eliminando un registro, no se necesita eliminar nada, ya que el propio trigger lo maneja
    RETURN NEW; -- Retorna el registro modificado o insertado
END;
$$ LANGUAGE plpgsql;




CREATE TRIGGER trigger_eliminar_historial
AFTER INSERT OR UPDATE OR DELETE ON historial_paf_aceptadas
FOR EACH ROW
EXECUTE FUNCTION fn_eliminar_historial();