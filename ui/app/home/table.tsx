import type { Route } from "./+types/home";

export function Row({
  request_id,
  status,
  created_at,
  updated_at,
  title,
}: Route.ComponentProp) {
  return (
    <tr key={request_id}>
      <td>{title}</td>
      <td>
        <a href={`/detail/${request_id}`}>{request_id}</a>
      </td>
      <td>{status}</td>
      <td>{created_at}</td>
      <td>{updated_at}</td>
    </tr>
  );
}

export default function Table({ data }: Route.ComponentProp) {
  return (
    <div className="container">
      <table className="table">
        <thead>
          <tr>
            <th>Title</th>
            <th>RequestID</th>
            <th>Status</th>
            <th>Created</th>
            <th>Updated</th>
          </tr>
        </thead>
        <tbody>
          {data.map((x) => Row(x))}
        </tbody>
      </table>
    </div>
  );
}
