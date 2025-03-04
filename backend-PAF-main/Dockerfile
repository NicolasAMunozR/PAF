FROM golang:1.23-alpine

WORKDIR /app

# Copia los archivos de go.mod y go.sum
COPY go.mod ./
COPY go.sum ./

# Descarga las dependencias
RUN go mod download

# Copia el resto de los archivos de la aplicación
COPY . .

# Asegúrate de que los módulos necesarios se descarguen
RUN go get -d -v ./...

# Construye la aplicación
RUN go build -o /server

# Expone el puerto en el que la aplicación se ejecutará
EXPOSE 3000

# Comando para ejecutar la aplicación
CMD ["/server"]