# =========================
# Stage 1: Build Go binary
# =========================
FROM golang:1.22-bullseye AS builder
WORKDIR /app

# Ensure updated GPG keys to avoid signature errors
RUN apt-get update || true \
 && apt-get install -y --no-install-recommends gnupg ca-certificates \
 && apt-key adv --keyserver keyserver.ubuntu.com --recv-keys 0xA1F196A8 || true \
 && apt-get update

# Download Go modules
COPY go.mod go.sum ./
RUN go mod download

# Copy full project
COPY . .

# Install cross-compilation dependencies for CGO
RUN dpkg --add-architecture amd64 && \
    apt-get update && \
    apt-get install -y --no-install-recommends \
      g++-x86-64-linux-gnu libc6-dev-amd64-cross \
      gcc libjpeg-dev:amd64 libx11-dev:amd64 xorg-dev:amd64 libxtst-dev:amd64 \
      xsel:amd64 xclip:amd64 libpng++-dev:amd64 \
      libxcb-xkb-dev:amd64 x11-xkb-utils:amd64 libx11-xcb-dev:amd64 \
      libxkbcommon-x11-dev:amd64 libxkbcommon-dev:amd64 \
      dpkg-dev fakeroot \
    && rm -rf /var/lib/apt/lists/*

# Build Go binary
ENV CGO_ENABLED=1 \
    GOOS=linux \
    GOARCH=amd64 \
    CC=x86_64-linux-gnu-gcc
RUN go build -o terminal main.go

# =============================
# Stage 2: Package as .deb file
# =============================
FROM debian:bullseye-slim AS packager
WORKDIR /package

# Copy binary from builder
COPY --from=builder /app/terminal ./terminal-bin

# Copy packaging files
COPY installer/linux-amd64/ .

# Create Debian package folder structure
RUN mkdir -p terminal/DEBIAN \
    && mkdir -p terminal/usr/local/bin \
    && mkdir -p terminal/lib/systemd/system

# Move control scripts and normalize EOLs
RUN mv control terminal/DEBIAN/control \
    && mv postinst terminal/DEBIAN/postinst \
    && mv prerm terminal/DEBIAN/prerm \
    && sed -i 's/\r$//' terminal/DEBIAN/control terminal/DEBIAN/postinst terminal/DEBIAN/prerm \
    && for f in terminal/DEBIAN/control terminal/DEBIAN/postinst terminal/DEBIAN/prerm; do \
         [ -s "$f" ] && tail -c1 "$f" | grep -q $'\n' || echo >> "$f"; \
       done \
    && chmod 644 terminal/DEBIAN/control \
    && chmod 755 terminal/DEBIAN/postinst terminal/DEBIAN/prerm

# Move systemd service file and normalize EOLs
RUN mv terminal.service terminal/lib/systemd/system/terminal.service \
    && sed -i 's/\r$//' terminal/lib/systemd/system/terminal.service \
    && tail -c1 terminal/lib/systemd/system/terminal.service | grep -q $'\n' || echo >> terminal/lib/systemd/system/terminal.service

# Install binary with proper permissions
RUN install -m 0755 terminal-bin terminal/usr/local/bin/terminal

# Build .deb package
RUN dpkg-deb --build terminal terminal_amd64.deb

# ==========================
# Final Stage: Export .deb
# ==========================
FROM debian:bullseye-slim AS final
COPY --from=packager /package/terminal_amd64.deb /terminal.deb
