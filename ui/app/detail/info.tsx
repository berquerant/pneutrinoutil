import type { Route } from "./+types/detail";

export default function Info({
  request_id,
  status,
  started_at,
  completed_at,
  created_at,
  updated_at,
  command,
  title,
}: Route.ComponentProps) {
  return (
    <div className="container">
      <table className="table">
        <tbody>
          <tr>
            <td>Title</td>
            <td>{title}</td>
          </tr>
          <tr>
            <td>RequestID</td>
            <td>{request_id}</td>
          </tr>
          <tr>
            <td>Status</td>
            <td>{status}</td>
          </tr>
          <tr>
            <td>Command</td>
            <td>{command}</td>
          </tr>
          <tr>
            <td>Created</td>
            <td>{created_at}</td>
          </tr>
          <tr>
            <td>Updated</td>
            <td>{updated_at}</td>
          </tr>
          <tr>
            <td>Started</td>
            <td>{started_at}</td>
          </tr>
          <tr>
            <td>Completed</td>
            <td>{completed_at}</td>
          </tr>
        </tbody>
      </table>
    </div>
  );
}
