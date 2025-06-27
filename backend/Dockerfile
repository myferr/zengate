FROM oven/bun:1.1.13

WORKDIR /app

COPY package.json bun.lock ./
COPY src ./src

RUN bun install

COPY vault.json vault.json
COPY .env .env

ENV PORT=3002

EXPOSE 3002

CMD ["bun", "run", "src/index.ts"]
