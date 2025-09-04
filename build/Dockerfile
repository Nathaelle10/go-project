# Étape 1 : Build (Golang 1.25)
FROM golang:1.25 as builder

WORKDIR /app

# Copier les fichiers
COPY src/go.mod src/go.sum ./
RUN go mod download

COPY src/. .

# Compiler un binaire statique pour Linux
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o app

# Étape 2 : Image finale minimale
FROM scratch

# Copier le binaire compilé
COPY --from=builder /app/app /app

# Copier les fichiers statiques embarqués
COPY src/www /www

# Exposer les ports utilisés par l’app (adapter si nécessaire)
EXPOSE 8080 9090

# Point d'entrée
ENTRYPOINT ["/app"]
