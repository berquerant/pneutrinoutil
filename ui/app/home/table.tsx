export type RowParams = {
  request_id: string;
  status: string;
  created_at: string;
  updated_at: string;
  title: string;
};

export function Row({
  request_id,
  status,
  created_at,
  updated_at,
  title,
}: RowParams) {
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

export type TableParams = {
  data: RowParams[];
};

export default function Table({ data }: TableParams) {
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
