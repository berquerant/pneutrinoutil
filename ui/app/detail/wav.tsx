import Audio from "../common/audio";

export type WavParams = {
  apiServerUri: string;
  rid: string;
};

export default function Wav({
  apiServerUri,
  rid,
}: WavParams) {
  const url = `${apiServerUri}/proc/${rid}/wav`;
  return (
    <div>
      {Audio({ url: url, name: "Download Wav" })}
    </div>
  );
}
