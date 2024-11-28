import pandas as pd

# Configuración
file_path = 'PAF - PRUEBA.xlsx'  # Ruta del archivo Excel
sheet_name = 'CONTRATOS'  # Nombre de la hoja a leer
table_name = "contratos"  # Nombre de la tabla en la base de datos
output_file = "script.sql"  # Archivo de salida
row_limit = 5000  # Límite de filas a procesar

# Especificar las columnas y tipos
int_columns = ["HORAS"]  # Ajustar nombres
date_columns = []  # Ajustar nombres según Excel

# Leer el archivo Excel
df = pd.read_excel(file_path, sheet_name=sheet_name)

# Aplicar el límite de filas
if row_limit > 0:
    df = df.head(row_limit)

# Generar el script SQL
sql_script = ""
for _, row in df.iterrows():
    formatted_values = []
    for col, value in row.items():
        if pd.isnull(value):  # Valores nulos
            formatted_values.append("NULL")
        elif col in int_columns:  # Columnas INT
            formatted_values.append(str(int(value)))
        elif col in date_columns:  # Columnas DATE
            formatted_values.append(f"'{pd.to_datetime(value).strftime('%Y-%m-%d')}'")
        else:  # Texto y otros
            formatted_values.append(f"'{str(value).replace('\'', '\'\'')}'")  # Escapar comillas simples

    # Asegurar que las columnas coincidan con los datos
    sql_script += f"""INSERT INTO {table_name} (run_docente,unidad_mayor,unidad_menor,planta,horas) VALUES ({', '.join(formatted_values)});\n"""

# Guardar el script en un archivo .sql
with open(output_file, "w") as f:
    f.write(sql_script)

print(f"Script SQL generado con éxito y guardado en {output_file}.")
