FROM node:22-alpine AS build

WORKDIR /app

COPY package.json package-lock.json ./
RUN npm ci

COPY . .
RUN npm run build

FROM alpine:3.20

COPY --from=build /app/dist /dist

CMD ["sh", "-c", "cp -r /dist/. /output/"]
