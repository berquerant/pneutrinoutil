import Audio from "../common/audio";

export type WorldWavParams = {
  apiServerUri: string;
  rid: string;
};

export default function WorldWav({
  apiServerUri,
  rid,
}: WorldWavParams) {
  const url = `${apiServerUri}/proc/${rid}/world_wav`;
  return (
    <div>
      {Audio({ url: url, name: "Download WorldWav (before NEUTRINO v3)" })}
    </div>
  );
}
