import type { Route } from "./+type/detail"

export default function MusicXML({
  apiServerUri,
  rid,
}: Route.ComponentProps) {
  const url = `${apiServerUri}/proc/${rid}/musicxml`
  return (
    <a href={url} download className="btn btn-primary">
    MusicXML
    </a>
  )
}
