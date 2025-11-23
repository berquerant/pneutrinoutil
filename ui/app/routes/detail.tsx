import type { Route } from "./+types/detail";
import { apiServerUri, defaultApi } from "../api/env";
import type { InfoParams } from "../detail/info";
import Detail from "../detail/detail";
import Config from "../detail/config";
import Log from "../detail/log";
import MusicXML from "../detail/musicxml";
import Wav from "../detail/wav";
import WorldWav from "../detail/worldwav";
import axios from "axios";

export async function loader({ params }: Route.LoaderArgs) {
  const detail = await defaultApi.procIdDetailGet(params.id);
  const d = detail.data.data;
  if (d === undefined) {
    throw Error("detail: defaultApi.procIdDetailGet returned undefined data!");
  }
  const d2 = {
    request_id: d.rid,
    title: d.basename,
  };
  const detailData = {
    ...d,
    ...d2,
  };

  const result: Record<string, unknown> = {
    detail: detailData,
  };
  const isNotFound = (err: unknown) => {
    if (axios.isAxiosError(err)) {
      return err.response && err.response.status == 404;
    }
    return false;
  };
  try {
    const x = await defaultApi.procIdConfigGet(params.id);
    result["config"] = x.data.data;
  } catch (err) {
    if (!isNotFound(err)) {
      throw err;
    }
  }
  try {
    const x = await defaultApi.procIdLogGet(params.id);
    result["log"] = x.data;
  } catch (err) {
    if (!isNotFound(err)) {
      throw err;
    }
  }

  result["apiServerUri"] = apiServerUri;
  return result;
}

export function meta({ params }: Route.MetaArgs) {
  return [
    { title: `Pneutrinoutil UI: ${params.id}` },
    { name: "description", content: "Welcome to Pneutrinoutil UI!" },
  ];
}

export type ComponentLoaderData = {
  detail: InfoParams;
  config: unknown;
  log: string;
  apiServerUri: string;
};

export type ComponentProps = {
  loaderData: ComponentLoaderData;
};

export default function Component({
  loaderData: {
    detail,
    config,
    log,
    apiServerUri,
  },
}: ComponentProps) {
  return (
    <div className="container">
      {Detail(detail)}
      <hr />
      <div className="row align-items-start">
        <div className="col d-flex gap-3">
          {config != null && Config(config)}
          {log != null && Log(log)}
          {MusicXML({ apiServerUri: apiServerUri, rid: detail.request_id })}
          {Wav({ apiServerUri: apiServerUri, rid: detail.request_id })}
          {WorldWav({ apiServerUri: apiServerUri, rid: detail.request_id })}
        </div>
      </div>
    </div>
  );
}
