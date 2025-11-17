import type { Route } from "./+types/detail"
import Info from './info'

export function meta({ params }: Route.MetaArgs) {
  return [
    { title: `Pneutrinoutil UI: ${params.title}` },
    { name: "description", content: "Welcome to Pneutrinoutil UI!" },
  ]
}

export default function Component(params: Route.ComponentProps) {
  return <div className="container">
    {Info(params)}
    </div>
}
