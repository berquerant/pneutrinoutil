import type { Route } from "./+type/detail"
import Download from '../common/download'

export default function MusicXML({
  apiServerUri,
  rid,
}: Route.ComponentProps) {
  const url = `${apiServerUri}/proc/${rid}/musicxml`
  return (
    <div>
      {Download({url: url, name: "Download MusicXML"})}
    </div>
  )
}
