FROM golang:1.24.2 AS builder

WORKDIR /app

# Copier les fichiers de dépendances
COPY go.mod go.sum* ./

# Télécharger les dépendances
RUN go mod download

# Copier le code source
COPY . .

# Compiler l'application
RUN CGO_ENABLED=0 GOOS=linux go build -o cpu-monitor .

# Image finale
FROM alpine:latest

# Installer les certificats CA pour les requêtes HTTPS
RUN apk --no-cache add ca-certificates

# Create a non-root user and group
RUN addgroup -S appgroup && adduser -S appuser -G appgroup

WORKDIR /app

# Copier l'exécutable compilé depuis l'étape de construction
COPY --from=builder /app/cpu-monitor .

# Copier les templates
COPY --from=builder /app/templates ./templates

# Set ownership and permissions
RUN chown -R appuser:appgroup /app && \
    chmod +x /app/cpu-monitor

# Exposer le port utilisé par l'application
EXPOSE 1323

# Switch to non-root user
USER appuser

# Commande pour exécuter l'application
CMD ["./cpu-monitor"]
