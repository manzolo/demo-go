# syntax=docker/dockerfile:1

# Stage 1: Compilazione
FROM golang:1.23-alpine AS builder
WORKDIR /app

# Copia e scarica le dipendenze
COPY go.mod go.sum ./
RUN go mod download

# Copia il codice sorgente
COPY *.go ./

# Compila l'applicazione in modalità statica
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o /demo-go

# Stage 2: Immagine finale minimale
FROM scratch

# Copia il binario dal builder
COPY --from=builder /demo-go /demo-go

EXPOSE 8080

ENTRYPOINT ["/demo-go"]