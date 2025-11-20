import type { Route } from "./+types/detail";
import CodeModal from "../common/modal";

export default function Config(loaderData: Route.ComponentProps) {
  const data = JSON.stringify(loaderData, null, "  ");
  return CodeModal({ name: "Show Config", code: data });
}
