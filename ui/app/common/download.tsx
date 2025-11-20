export default function Download({ url, name }) {
  return <a href={url} download className="btn btn-primary">{name}</a>;
}
