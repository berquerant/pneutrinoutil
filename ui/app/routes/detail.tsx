import type { Route } from "./+type/detail"
import { defaultApi } from '../api/env'
import Detail from '../detail/detail'
import Config from '../detail/config'

export async function loader({ params }: Route.LoaderArgs) {
  const detail = await defaultApi.procIdDetailGet(params.id)
  const d = detail.data.data
  const d2 = {
    request_id: d.rid,
    title: d.basename,
  }
  const detailData = {
    ...d,
    ...d2,
  }
  try {
    const config = await defaultApi.procIdConfigGet(params.id)
    const configData = config.data.data
    return {
      detail: detailData,
      config: configData,
    }
  } catch(err) {
    if (err.response && err.response.status == 404) {
      return {
        detail: detailData,
        config: null,
      }
    }
    throw err
  }
}

export function meta({ params }: Route.MetaArgs) {
  return [
    { title: `Pneutrinoutil UI: ${params.id}` },
    { name: "description", content: "Welcome to Pneutrinoutil UI!" },
  ]
}

export default function Component({
  loaderData: {
    detail,
    config,
  },
}: Route.ComponentProps) {
  return <div className="container">
  {Detail(detail)}
  {config != null && Config(config)}
  </div>
}
