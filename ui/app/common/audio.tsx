import Download from "./download";

export type AudioParams = {
  url: string;
  name: string;
};

export default function Audio({
  url,
  name,
}: AudioParams) {
  return (
    <div className="card">
      <audio controls src={url} preload="none"></audio>
      {Download({ url: url, name: name })}
    </div>
  );
}
