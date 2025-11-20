import type { Route } from "./+types/detail";
import CodeModal from "../common/modal";

export default function Log(loaderData: Route.ComponentProps) {
  return CodeModal({ name: "Show Log", code: loaderData });
}
