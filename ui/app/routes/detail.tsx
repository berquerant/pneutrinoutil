import type { Route } from "./+type/detail"
import { defaultApi, apiServerUri } from '../api/env'
import Detail from '../detail/detail'
import Config from '../detail/config'
import Log from '../detail/log'
import MusicXML from '../detail/musicxml'
import Wav from '../detail/wav'
import WorldWav from '../detail/worldwav'

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

  const result: Record<string, any> = {
    detail: detailData,
  }
  const isNotFound = err => err.response && err.response.status == 404
  try {
    const x = await defaultApi.procIdConfigGet(params.id)
    result['config'] = x.data.data
  } catch(err) {
    if (!isNotFound(err)) {
      throw err
    }
  }
  try {
    const x = await defaultApi.procIdLogGet(params.id)
    result['log'] = x.data
  } catch(err) {
    if (!isNotFound(err)) {
      throw err
    }
  }

  result['apiServerUri'] = apiServerUri
  return result
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
    log,
    apiServerUri,
  },
}: Route.ComponentProps) {
  return (
    <div className="container">
    {Detail(detail)}
    <hr/>
    <div className="row align-items-start">
    <div className="col d-flex gap-3">
    {config != null && Config(config)}
    {log != null && Log(log)}
    {MusicXML({ apiServerUri: apiServerUri, rid: detail.request_id })}
    {Wav({ apiServerUri: apiServerUri, rid: detail.request_id })}
    {WorldWav({ apiServerUri: apiServerUri, rid: detail.request_id })}
    </div>
    </div>
    </div>
  )
}
