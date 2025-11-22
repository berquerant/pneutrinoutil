import type { Route } from "./+types/home";
import { defaultApi } from "../api/env";
import Table from "../home/table";

export async function loader({ request }: Route.LoaderArgs) {
  const url = new URL(request.url);
  const searchParams = url.searchParams;

  const limitParam = searchParams.get("limit");
  const limit = limitParam && !isNaN(Number(limitParam))
    ? Math.max(10, Number(limitParam))
    : 10;
  const prefix = searchParams.get("prefix");
  const status = searchParams.get("status");
  const start = searchParams.get("start");
  const end = searchParams.get("end");
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

export default function Home({
  loaderData,
}: Route.ComponentProps) {
  return (
    <div className="container">
      {Table({
        data: loaderData,
      })}
    </div>
  );
}
