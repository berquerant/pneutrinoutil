import type { Route } from "./+type/detail"
import Audio from '../common/audio'

export default function WorldWav({
  apiServerUri,
  rid,
}: Route.ComponentProps) {
  const url = `${apiServerUri}/proc/${rid}/world_wav`
  return (
    <div>
      {Audio({url: url, name: "Download WorldWav"})}
    </div>
  )
}
