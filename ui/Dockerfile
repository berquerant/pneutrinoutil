# syntax=docker/dockerfile:1
# check=error=true

# https://pnpm.io/docker
ARG NODE_VERSION=24.14.1
FROM node:${NODE_VERSION} AS base
ARG PNPM_VERSION=10.33.0
ENV PNPM_HOME="/pnpm"
ENV PATH="$PNPM_HOME:$PATH"
RUN npm i -g pnpm@$PNPM_VERSION
WORKDIR /app

FROM base AS build
COPY ./package.json ./pnpm-lock.yaml ./pnpm-workspace.yaml /app/
RUN --mount=type=cache,id=pnpm,target=/pnpm/store pnpm install --frozen-lockfile
COPY . .
RUN pnpm run build
RUN pnpm deploy --filter ui --prod deploy/app

FROM base AS final
WORKDIR /app
COPY ./package.json .
COPY --from=build /app/deploy/app/node_modules /app/node_modules
COPY --from=build /app/build /app/build
COPY --from=build /app/deploy/app/public /app/public
COPY --from=build /app/package.json /app/
CMD ["pnpm", "run", "start"]
