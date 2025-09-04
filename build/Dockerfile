# Étape 1 : Build (Golang 1.25)
FROM golang:1.25 as builder

WORKDIR /app

# Copier les fichiers de dépendance
COPY src/go.mod src/go.sum ./
RUN go mod download

# Copier le code source
COPY src/ ./

# Compiler un binaire statique
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o app .

# Étape 2 : Image finale minimale
FROM scratch

# Copier le binaire
COPY --from=builder /app/app /app

# Copier les fichiers statiques (si tu les sers depuis ./www)
COPY src/www /www

# Exposer les ports utilisés
EXPOSE 8080 9090

# Point d'entrée
ENTRYPOINT ["/app"]
