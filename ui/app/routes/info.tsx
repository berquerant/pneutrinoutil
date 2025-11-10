import type { Route } from "./+types/info"
import { useLoaderData } from "react-router";
import { defaultApi } from '../api/env';
import type { HandlerSuccessResponseHandlerVersionResponseData } from '../api/client'

export async function loader(): Promise<HandlerSuccessResponseHandlerVersionResponseData> {
  const r = await defaultApi.versionGet()
  return r.data
}

export function meta({ }: Route.MetaArgs) {
  return [
    { title: "Pneutrinoutil UI Info" },
    { name: "description", content: "Show information" },
  ]
}

export default async function Info() {
  const d = useLoaderData<typeof loader>()
  return <div>
    <h1> Info </h1>
    <h2> Server </h2>
    <li> version: { d.data.version } </li>
    <li> revision: { d.data.revision } </li>
    </div>
}
