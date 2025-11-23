import type { Route } from "./+types/info";
import { defaultApi } from "../api/env";
import { HandlerVersionResponseData } from "../api/client";

export async function loader() {
  const r = await defaultApi.versionGet();
  return r.data.data;
}

export function meta() {
  return [
    { title: "Pneutrinoutil UI Info" },
    { name: "description", content: "Show information" },
  ];
}

export type InfoParams = {
  loaderData: HandlerVersionResponseData;
};

export default async function Info({
  loaderData,
}: InfoParams) {
  return (
    <div className="container">
      <table className="table">
        <tbody>
          <tr>
            <td>Version</td>
            <td>{loaderData.version}</td>
          </tr>
          <tr>
            <td>Revision</td>
            <td>{loaderData.revision}</td>
          </tr>
        </tbody>
      </table>
    </div>
  );
}
