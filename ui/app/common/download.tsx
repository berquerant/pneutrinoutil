export type DownloadParams = {
  url: string;
  name: string;
};

export default function Download({ url, name }: DownloadParams) {
  return <a href={url} download className="btn btn-primary">{name}</a>;
}
