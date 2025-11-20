import type { Route } from "./+type/detail"
import Audio from '../common/audio'

export default function Wav({
  apiServerUri,
  rid,
}: Route.ComponentProps) {
  const url = `${apiServerUri}/proc/${rid}/wav`
  return (
    <div>
      {Audio({url: url, name: "Download Wav"})}
    </div>
  )
}
