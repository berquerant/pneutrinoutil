import type { Route } from "./+types/info";
import { defaultApi } from "../api/env";

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

export default async function Info({
  loaderData: { version, revision },
}: Route.ComponentProps) {
  return (
    <div className="container">
      <table className="table">
        <tbody>
          <tr>
            <td>Version</td>
            <td>{version}</td>
          </tr>
          <tr>
            <td>Revision</td>
            <td>{revision}</td>
          </tr>
        </tbody>
      </table>
    </div>
  );
}
