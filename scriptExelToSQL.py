import pandas as pd

# Configuración
file_path = 'PAF - PRUEBA.xlsx'  # Reemplaza con la ruta de tu archivo Excel
sheet_name = 'PAF-SAI'  # Reemplaza con el nombre de la hoja que deseas leer
table_name = "profesor_dbs"  # Reemplaza con el nombre de tu tabla
output_file = "script.sql"  # Nombre del archivo de salida
row_limit = 5000  # Límite de filas que deseas incluir en el script

# Leer los datos de la hoja de Excel
df = pd.read_excel(file_path, sheet_name=sheet_name)

# Aplicar el límite de filas
if row_limit > 0:  # Asegúrate de que el límite sea positivo
    df = df.head(row_limit)

# Generar el script SQL
sql_script = ""
for _, row in df.iterrows():
    # Convertir los valores de cada fila a formato de cadena y manejar valores nulos
    values = "', '".join(map(lambda x: str(x) if pd.notnull(x) else "NULL", row.values))
    sql_script += f"INSERT INTO {table_name} VALUES ('{values}');\n"

# Guardar el script en un archivo .sql
with open(output_file, "w") as f:
    f.write(sql_script)

print(f"Script SQL generado con éxito y guardado en {output_file}.")
