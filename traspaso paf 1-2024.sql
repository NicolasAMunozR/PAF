CREATE TABLE IF NOT EXISTS historial_paf_aceptadas (
    codigo_asignatura TEXT,
    nombre_asignatura TEXT,
    run TEXT,
    id_paf TEXT,
    cupo TEXT,
    id SERIAL PRIMARY KEY
);

INSERT INTO historial_paf_aceptadas (codigo_asignatura, nombre_asignatura, run, id_paf, cupo) VALUES ('26210', 'TECNICAS DE LA EXPRESION ORAL Y ESCRITA', '9008894', '401', 20);
INSERT INTO historial_paf_aceptadas (codigo_asignatura, nombre_asignatura, run, id_paf, cupo) VALUES ('361445', 'AUDITORÍA III', '15343237', '1978', 0);
INSERT INTO historial_paf_aceptadas (codigo_asignatura, nombre_asignatura, run, id_paf, cupo) VALUES ('2093', '
CONTABILIDAD', '13932477', '2756', 60);
INSERT INTO historial_paf_aceptadas (codigo_asignatura, nombre_asignatura, run, id_paf, cupo) VALUES ('76121', 'DERECHO PENAL I ', '16213713', '2560', 0);
INSERT INTO historial_paf_aceptadas (codigo_asignatura, nombre_asignatura, run, id_paf, cupo) VALUES ('76136', 'DERECHO CIVIL IV', '16382042', '2561', 60);
