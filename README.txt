--------------------Cambios backend---------------------
se creo el modelo: HistorialPasosPaf, en el cual se consume solamente y se realizan calculos mediante el siguiente controlador ubicado en la linea 108 de "main.go": r.GET("/historialPaso/:id_paf/:run_docente"
historialController.ObtenerHistorialYDuracionesPorIdYRun) se espera que en :run_docente se ingrese el run mediante URL

Se creo un metodo el cual obtiene la cantidad de profesores sin paf apartir de los datos del SAI comprando con pipelsofts y la tabla contratos el controlador es:
r.GET("/estadisticas/obtener-y-comparar-runs", estadisticasController.ObtenerYCompararRunsHandler)  y esta en la linea 111 del archivo "main.go"

Se añadio un parametro extra al modelo "HistorialPafAceptadas.go", el cual es: SemestrePaf (text)  linea: 34 de "HistorialPafAceptadas.go" en la carpeta modelos

Se modifico el servicio de estadistica haciendo que ahora se muestre la descripcion de los estados en lugar del estado carpeta: "Service" nombre "EstadisticaService"
 