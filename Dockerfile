# Stage 1: build the React client.
FROM node:22-bookworm-slim AS client
WORKDIR /src/frontend
COPY frontend/package.json frontend/package-lock.json ./
RUN npm ci --no-audit --no-fund
COPY frontend/ ./
RUN npm run build

# Stage 2: compile the Go server with the client embedded. Pure Go SQLite, so
# no C toolchain is needed and the binary is fully static.
FROM golang:1.23-bookworm AS server
WORKDIR /src/backend
COPY backend/go.mod backend/go.sum ./
RUN go mod download
COPY backend/ ./
COPY --from=client /src/frontend/dist ./web/dist
RUN CGO_ENABLED=0 go build -trimpath -ldflags "-s -w" -o /out/csit214 .

# Stage 3: a small runtime image holding only the binary and the sample data.
FROM alpine:3.20
RUN adduser -D -H -u 10001 csit214 && mkdir /data && chown csit214 /data
COPY --from=server /out/csit214 /app/csit214
COPY seed/seed.sql /app/seed/seed.sql
ENV CSIT214_DB_PATH=/data/data.db \
    CSIT214_SEED_PATH=/app/seed/seed.sql \
    PORT=8080
USER csit214
VOLUME /data
EXPOSE 8080
HEALTHCHECK --interval=15s --timeout=3s --start-period=5s \
  CMD wget -qO- http://127.0.0.1:8080/api/health || exit 1
ENTRYPOINT ["/app/csit214"]
