import { type RouteConfig, index, route } from "@react-router/dev/routes"

export default [
  index("routes/home.tsx"),
  route("detail/:id", "routes/detail.tsx"),
  route("create", "routes/create.tsx"),
  route("info", "routes/info.tsx"),
] satisfies RouteConfig
