import Download from "../common/download";

export default function Audio({
  url,
  name,
}) {
  return (
    <div className="card">
      <audio controls src={url} preload="none"></audio>
      {Download({ url: url, name: name })}
    </div>
  );
}
