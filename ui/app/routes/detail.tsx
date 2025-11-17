import type { Route } from "./+type/detail"
import { defaultApi } from '../api/env'
import Detail from '../detail/detail'

export async function loader({ params }: Route.LoaderArgs) {
  const r = await defaultApi.procIdDetailGet(params.id)
  const d = r.data.data
  const s = {
    request_id: d.rid,
    title: d.basename,
  }
  return {
    ...d,
    ...s,
  }
}

export function meta({ params }: Route.MetaArgs) {
  return [
    { title: `Pneutrinoutil UI: ${params.id}` },
    { name: "description", content: "Welcome to Pneutrinoutil UI!" },
  ]
}

export default function Component({
  loaderData,
}: Route.ComponentProps) {
  return <div className="container">
    {Detail(loaderData)}
    </div>
}
