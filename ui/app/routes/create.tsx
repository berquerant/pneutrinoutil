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
      d.get("enhanceBreathiness") as any,
      d.get("formantShift") as any,
      d.get("inference") as any,
      d.get("model") as any,
      d.get("supportModel") as any,
      d.get("transpose") as any,
      d.get("pitchShiftNsf") as any,
      d.get("pitchShiftWorld") as any,
      d.get("smoothFormant") as any,
      d.get("smoothPitch") as any,
      d.get("styleShift") as any,
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
          <label className="form-label" htmlFor="enhanceBreathiness">
            EnhanceBreathiness
          </label>
          <input
            className="form-control"
            id="enhanceBreathiness"
            name="enhanceBreathiness"
            type="text"
            inputMode="numeric"
            defaultValue="0"
          />
          <div className="form-text" id="enhanceBreathiness">
            [0, 100]% (before NEUTRINO v3)
          </div>
        </div>
        <div className="mb-3">
          <label className="form-label" htmlFor="formantShift">
            FormantShift
          </label>
          <input
            className="form-control"
            id="formantShift"
            name="formantShift"
            type="text"
            inputMode="numeric"
            defaultValue="1.0"
          />
          <div className="form-text" id="formantShift">
            Higher values result in a younger tone; lower values result in a
            mature tone (before NEUTRINO v3)
          </div>
        </div>
        <div className="mb-3">
          <label className="form-label" htmlFor="inference">Inference</label>
          <select
            className="form-select"
            id="inference"
            name="inference"
            defaultValue="2"
          >
            <option value="2">Elements</option>
            <option value="3">Standard</option>
            <option value="4">Advanced</option>
          </select>
          <div className="form-text" id="inference">
            Inference quality (before NEUTRINO v3)
          </div>
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
            Support Singer library (NEUTRINO v3)
          </div>
        </div>
        <div className="mb-3">
          <label className="form-label" htmlFor="transpose">Transpose</label>
          <input
            className="form-control"
            id="transpose"
            name="transpose"
            type="text"
            inputMode="numeric"
            pattern="\d*"
            defaultValue="0"
          />
          <div className="form-text" id="transpose">
            Infer by raising the score by the specified key (NEUTRINO v3)
          </div>
        </div>
        <div>
          <label className="form-label" htmlFor="pitchShiftNsf">
            PitchShiftNsf
          </label>
          <input
            className="form-control"
            id="pitchShiftNsf"
            name="pitchShiftNsf"
            type="text"
            inputMode="numeric"
            defaultValue="0"
          />
          <div className="form-text" id="pitchShiftNsf">
            (before NEUTRINO v3)
          </div>
        </div>
        <div className="mb-3">
          <label className="form-label" htmlFor="pitchShiftWorld">
            PitchShiftWorld
          </label>
          <input
            className="form-control"
            id="pitchShiftWorld"
            name="pitchShiftWorld"
            type="text"
            inputMode="numeric"
            defaultValue="0"
          />
          <div className="form-text" id="pitchShiftWorld">
            (before NEUTRINO v3)
          </div>
        </div>
        <div className="mb-3">
          <label className="form-label" htmlFor="smoothFormant">
            SmoothFormant
          </label>
          <input
            className="form-control"
            id="smoothFormant"
            name="smoothFormant"
            type="text"
            inputMode="numeric"
            defaultValue="0"
          />
          <div className="form-text" id="smoothFormant">
            [0, 100]% (before NEUTRINO v3)
          </div>
        </div>
        <div className="mb-3">
          <label className="form-label" htmlFor="smoothPitch">
            SmoothPitch
          </label>
          <input
            className="form-control"
            id="smoothPitch"
            name="smoothPitch"
            type="text"
            inputMode="numeric"
            defaultValue="0"
          />
          <div className="form-text" id="smoothPitch">
            [0, 100]% (before NEUTRINO v3)
          </div>
        </div>
        <div className="mb-3">
          <label className="form-label" htmlFor="styleShift">StyleShift</label>
          <input
            className="form-control"
            id="styleShift"
            name="styleShift"
            type="text"
            inputMode="numeric"
            pattern="\d*"
            defaultValue="0"
          />
          <div className="form-text" id="styleShift">
            Infer by raising the score by the specified key (before NEUTRINO v3)
          </div>
        </div>
        <button className="btn btn-primary" type="submit">
          Create new process
        </button>
      </Form>
    </div>
  );
}
