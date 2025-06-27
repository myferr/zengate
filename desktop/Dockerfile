FROM rust:latest AS base

ARG TARGETARCH
RUN apt-get update && apt-get install -y --no-install-recommends \
    curl wget unzip build-essential cmake pkg-config libssl-dev git \
    libwebkit2gtk-4.1-dev libgtk-3-dev libayatana-appindicator3-dev librsvg2-dev \
    libsoup-3.0-dev libjavascriptcoregtk-4.1-dev \
    libx11-dev libxrandr-dev libxss-dev libxext-dev libxi-dev \
    xdg-utils \
    && apt-get clean && rm -rf /var/lib/apt/lists/*

# Now this will expand TARGETARCH correctly
RUN rustup target add ${TARGETARCH}-unknown-linux-gnu

RUN curl -fsSL https://bun.sh/install | bash
ENV PATH="/root/.bun/bin:$PATH"

# Set PKG_CONFIG_PATH accordingly
ENV PKG_CONFIG_PATH=/usr/lib/${TARGETARCH}-linux-gnu/pkgconfig

WORKDIR /app
COPY . .

RUN bun install
RUN bunx tauri build

FROM scratch AS artifact
COPY --from=base /app/src-tauri/target/release/bundle/appimage/*.AppImage /Zengate.AppImage
