# syntax=docker/dockerfile:1
# check=error=true

# https://pnpm.io/docker
ARG NODE_VERSION=24.13.0
FROM node:${NODE_VERSION} AS base
ENV PNPM_HOME="/pnpm"
ENV PATH="$PNPM_HOME:$PATH"
RUN corepack enable
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
