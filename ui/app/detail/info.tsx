export type InfoParams = {
  request_id: string;
  status: string;
  started_at: string;
  completed_at: string;
  created_at: string;
  updated_at: string;
  command: string;
  title: string;
};

export default function Info({
  request_id,
  status,
  started_at,
  completed_at,
  created_at,
  updated_at,
  command,
  title,
}: InfoParams) {
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
