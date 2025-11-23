import type { Route } from "./+types/home";
import { defaultApi } from "../api/env";
import { HandlerSearchProcessResponseDataElement } from "../api/client";
import Table, { RowParams } from "../home/table";

export async function loader({ request }: Route.LoaderArgs) {
  const url = new URL(request.url);
  const searchParams = url.searchParams;

  const limitParam = searchParams.get("limit");
  const limit = limitParam && !isNaN(Number(limitParam))
    ? Math.max(10, Number(limitParam))
    : 10;
  const prefix = searchParams.get("prefix") || "";
  const status = searchParams.get("status") || undefined;
  const start = searchParams.get("start") || undefined;
  const end = searchParams.get("end") || undefined;
  const r = await defaultApi.procSearchGet(
    limit,
    prefix,
    status,
    start,
    end,
  );
  return r.data.data;
}

export function meta() {
  return [
    { title: "Pneutrinoutil UI" },
    { name: "description", content: "Welcome to Pneutrinoutil UI!" },
  ];
}

export type HomeParams = {
  loaderData: HandlerSearchProcessResponseDataElement[];
};

export default function Home({
  loaderData,
}: HomeParams) {
  const rp: RowParams = {
    request_id: "",
    status: "",
    created_at: "",
    updated_at: "",
    title: "",
  };
  const data: RowParams[] = loaderData.filter((d) => {
    for (const key in rp) {
      if (
        d[key as keyof HandlerSearchProcessResponseDataElement] === undefined
      ) return false;
    }
    return true;
  }).map((d) => d as RowParams);
  return (
    <div className="container">
      {Table({
        data: data,
      })}
    </div>
  );
}
