import type { Route } from "./+types/info"
import { defaultApi } from '../api/env';

export async function loader({}: Route.LoaderArgs) {
  const r = await defaultApi.versionGet()
  return r.data.data
}

export function meta({}: Route.MetaArgs) {
  return [
    { title: "Pneutrinoutil UI Info" },
    { name: "description", content: "Show information" },
  ]
}

export default async function Info({
  loaderData: { version, revision },
}: Route.ComponentProps) {
  return <div class="container">
    <table class="table">
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
}
