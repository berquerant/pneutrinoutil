import CodeModal from "../common/modal";

export default function Config(loaderData: unknown) {
  const data = JSON.stringify(loaderData, null, "  ");
  return CodeModal({ name: "Show Config", code: data });
}
