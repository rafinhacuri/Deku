FROM golang:1.25.3 AS go-builder

WORKDIR /server

COPY ./api/go.mod ./api/go.sum ./
RUN go mod download

COPY ./api ./
RUN CGO_ENABLED=1 GOOS=linux GOARCH=amd64 go build -o api .

FROM oven/bun:1-debian

WORKDIR /main

COPY ./package.json ./bun.lock ./nuxt.config.ts ./tsconfig.json ./

RUN --mount=type=cache,target=/root/.bun/install/cache bun install

COPY ./app ./app
COPY ./public ./public
COPY ./i18n ./i18n

RUN bun run app:build
RUN apt-get update && apt-get install -y rsync
RUN mkdir -p .output/server/node_modules && rsync -a --ignore-existing node_modules/ .output/server/node_modules/

COPY --from=go-builder /server/api /main/api

RUN mkdir /database
VOLUME ["/database"]

EXPOSE 3000

CMD ["sh", "-c", "/main/api & bun /main/.output/server/index.mjs; wait"]