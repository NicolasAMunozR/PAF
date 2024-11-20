INSERT INTO pipelsofts(
    run, 
    nombres, 
    primer_apellido, 
    segundo_apellido, 
    correo, 
    codigo_unidad_contratante, 
    nombre_unidad_contratante, 
    nombre_unidad_mayor, 
    codigo_paf, 
    fecha_inicio_contrato, 
    fecha_fin_contrato, 
    codigo_asignatura, 
    nombre_asignatura, 
    cantidad_horas, 
    jerarquia, 
    calidad, 
    estado_proceso, 
    fecha_ultima_modificacion_proceso
) VALUES 
('12345678-9', 'Juan', 'Pérez', 'González', 'juan.perez@example.com', 'U001', 'Unidad de Desarrollo', 'Dirección de Proyectos', 'PAF12345', '2024-01-01', '2024-12-31', 'ASIGN001', 'Matemáticas Avanzadas', 120, 'Profesor Titular', 'Alta', 1, NOW()),
('23456789-0', 'María', 'López', 'Martínez', 'maria.lopez@example.com', 'U002', 'Unidad Administrativa', 'Dirección General', 'PAF12346', '2024-02-01', '2024-11-30', 'ASIGN002', 'Física General', 100, 'Profesor Asistente', 'Media', 2, NOW()),
('34567890-1', 'Carlos', 'García', 'Rodríguez', 'carlos.garcia@example.com', 'U003', 'Unidad de Investigación', 'Unidad de Investigación Avanzada', 'PAF12347', '2024-03-01', '2024-12-31', 'ASIGN003', 'Química Orgánica', 80, 'Profesor Asociado', 'Alta', 1, NOW()),
('45678901-2', 'Laura', 'Martín', 'Hernández', 'laura.martin@example.com', 'U004', 'Unidad de Capacitación', 'Secretaría Académica', 'PAF12348', '2024-04-01', '2024-09-30', 'ASIGN004', 'Biología Molecular', 60, 'Profesor Titular', 'Baja', 3, NOW()),
('56789012-3', 'Pedro', 'Sánchez', 'Fernández', 'pedro.sanchez@example.com', 'U005', 'Unidad de Innovación', 'Consejería Académica', 'PAF12349', '2024-05-01', '2024-12-31', 'ASIGN005', 'Ingeniería de Software', 150, 'Profesor Asistente', 'Media', 2, NOW());



INSERT INTO profesor_dbs (
    run, 
    semestre, 
    codigo_asignatura, 
    nombre_asignatura, 
    seccion, 
    cupo, 
    dia, 
    bloque
) VALUES 
('12345678-9', '2024-1', 'ASIGN001', 'Matemáticas Avanzadas', 'A', 30, 'Lunes', 'M2-M5-V1'),
('23456789-0', '2024-1', 'ASIGN002', 'Física General', 'B', 25, 'Martes', 'M3-M6-V2'),
('34567890-1', '2024-1', 'ASIGN003', 'Química Orgánica', 'C', 20, 'Miércoles', 'M1-M4-V1'),
('45678901-2', '2024-2', 'ASIGN004', 'Biología Molecular', 'A', 15, 'Jueves', 'M2-M5-V3'),
('56789012-3', '2024-2', 'ASIGN005', 'Ingeniería de Software', 'B', 40, 'Viernes', 'M4-M7-V2');




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
