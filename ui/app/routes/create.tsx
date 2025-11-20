import type { Route } from "./+type/create";
import { defaultApi } from "../api/env";
import { Form } from "react-router";

export async function action({
  request,
}: Route.ActionArgs) {
  const d = await request.formData();
  const r = await defaultApi.procPost(
    d.get("score"),
    d.get("enhanceBreathiness"),
    d.get("formantShift"),
    d.get("inference"),
    d.get("model"),
    d.get("pitchShiftNsf"),
    d.get("pitchShiftWorld"),
    d.get("smoothFormant"),
    d.get("smoothPitch"),
    d.get("styleShift"),
  );
  return r;
}

export default function Create({
  actionData,
}: Route.ComponentProps) {
  const result = actionData
    ? (
      <div className="alert alert-success" role="alert">
        Successfully created process! RequestID={actionData
          .headers["x-request-id"]}
      </div>
    )
    : null;
  return (
    <div>
      {result}
      <Form
        className="form-floating"
        method="post"
        encType="multipart/form-data"
      >
        <div>
          <label className="form-label" htmlFor="score" required>
            Score(musicxml)
          </label>
          <input
            className="form-control"
            id="score"
            name="score"
            type="file"
            accept=".musicxml,application/vnd.recordare.musicxml+xml"
          />
        </div>
        <div>
          <label className="form-label" htmlFor="enhanceBreathiness">
            EnhanceBreathiness(%)
          </label>
          <input
            className="form-control"
            id="enhanceBreathiness"
            name="enhanceBreathiness"
            type="text"
            inputMode="numeric"
            defaultValue="0"
          />
        </div>
        <div>
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
        </div>
        <div>
          <label className="form-label" htmlFor="inference">Inference</label>
          <select
            className="form-select"
            id="inference"
            name="inference"
            defaultValue="2"
          >
            <option value="2">2</option>
            <option value="3">3</option>
            <option value="4">4</option>
          </select>
        </div>
        <div>
          <label className="form-label" htmlFor="model">Model</label>
          <input
            className="form-control"
            id="model"
            name="model"
            type="text"
            defaultValue="KIRITAN"
          />
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
        </div>
        <div>
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
        </div>
        <div>
          <label className="form-label" htmlFor="smoothFormant">
            SmoothFormant(%)
          </label>
          <input
            className="form-control"
            id="smoothFormant"
            name="smoothFormant"
            type="text"
            inputMode="numeric"
            defaultValue="0"
          />
        </div>
        <div>
          <label className="form-label" htmlFor="smoothPitch">
            SmoothPitch(%)
          </label>
          <input
            className="form-control"
            id="smoothPitch"
            name="smoothPitch"
            type="text"
            inputMode="numeric"
            defaultValue="0"
          />
        </div>
        <div>
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
        </div>
        <button className="btn btn-primary" type="submit">Create</button>
      </Form>
    </div>
  );
}
