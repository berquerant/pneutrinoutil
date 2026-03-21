import type { Route } from "./+types/create";
import { defaultApi } from "../api/env";
import { Form } from "react-router";

export async function action({
  request,
}: Route.ActionArgs) {
  const d = await request.formData();
  try {
    /* eslint-disable @typescript-eslint/no-explicit-any */
    const r = await defaultApi.procPost(
      d.get("score") as File,
      d.get("model") as any,
      d.get("supportModel") as any,
      d.get("transpose") as any,
    );
    /* eslint-enable @typescript-eslint/no-explicit-any */
    return {
      ok: true,
      data: r.headers["x-request-id"],
    };
  } catch (err) {
    return {
      ok: false,
      err: String(err),
    };
  }
}

export default function Create({
  actionData,
}: Route.ComponentProps) {
  const result = actionData
    ? (
      actionData.ok
        ? (
          <div className="alert alert-success" role="alert">
            Successfully created process! RequestID={actionData
              .data}
          </div>
        )
        : (
          <div className="alert alert-danger" role="alert">
            Failed to create process! {actionData.err}
          </div>
        )
    )
    : null;
  return (
    <div className="container">
      {result}
      <Form
        className="form-floating"
        method="post"
        encType="multipart/form-data"
      >
        <div className="mb-3">
          <label className="form-label" htmlFor="score">
            Score
          </label>
          <input
            className="form-control"
            id="score"
            name="score"
            type="file"
            accept=".musicxml,application/vnd.recordare.musicxml+xml"
            required
          />
          <div className="form-text" id="score">musicxml</div>
        </div>
        <div className="mb-3">
          <label className="form-label" htmlFor="model">Model</label>
          <input
            className="form-control"
            id="model"
            name="model"
            type="text"
            defaultValue="KIRITAN"
          />
          <div className="form-text" id="model">Singer library</div>
        </div>
        <div className="mb-3">
          <label className="form-label" htmlFor="supportModel">
            SupportModel
          </label>
          <input
            className="form-control"
            id="supportModel"
            name="supportModel"
            type="text"
            defaultValue=""
          />
          <div className="form-text" id="supportModel">
            Support Singer library
          </div>
        </div>
        <div className="mb-3">
          <label className="form-label" htmlFor="transpose">Transpose</label>
          <input
            className="form-control"
            id="transpose"
            name="transpose"
            type="number"
            step="1"
            defaultValue="0"
          />
          <div className="form-text" id="transpose">
            Infer by raising the score by the specified key
          </div>
        </div>
        <button className="btn btn-primary" type="submit">
          Create new process
        </button>
      </Form>
    </div>
  );
}
