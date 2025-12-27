# syntax=docker/dockerfile:1
# check=error=true
ARG NODE_VERSION=24
FROM node:${NODE_VERSION}-alpine AS development-dependencies-env
COPY . /app
WORKDIR /app
RUN npm ci

FROM node:${NODE_VERSION}-alpine AS production-dependencies-env
COPY ./package.json package-lock.json /app/
WORKDIR /app
RUN npm ci --omit=dev

ARG ENVIRONMENT=container-dev
FROM node:${NODE_VERSION}-alpine AS build-env
COPY . /app/
COPY --from=development-dependencies-env /app/node_modules /app/node_modules
WORKDIR /app
RUN npm run build -- --mode "$ENVIRONMENT"

FROM node:${NODE_VERSION}-alpine AS final
COPY ./package.json package-lock.json /app/
COPY --from=production-dependencies-env /app/node_modules /app/node_modules
COPY --from=build-env /app/build /app/build
WORKDIR /app
CMD ["npm", "run", "start"]
