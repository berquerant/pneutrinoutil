import Download from "../common/download";

export type MusicXMLParams = {
  apiServerUri: string;
  rid: string;
};

export default function MusicXML({
  apiServerUri,
  rid,
}: MusicXMLParams) {
  const url = `${apiServerUri}/proc/${rid}/musicxml`;
  return (
    <div>
      {Download({ url: url, name: "Download MusicXML" })}
    </div>
  );
}
