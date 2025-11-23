import CodeModal from "../common/modal";

export default function Log(loaderData: string) {
  return CodeModal({ name: "Show Log", code: loaderData });
}
