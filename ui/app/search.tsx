import type { Route } from "./+types/root"
import { useState } from 'react';
import { useSearchParams, useNavigate } from "react-router"

function fromApiDatetime(x: string | null): string | null {
  if (x == null) {
    return null
  }
  const d = new Date(x)
  return d.toISOString().slice(0, 19)
}

function toApiDatetime(x: string | null): string | null {
  if (x == null) {
    return null
  }
  return x + 'Z'
}

export default function Search() {
  const navigate = useNavigate()
  const [searchParams, _] = useSearchParams()
  const params =  Object.fromEntries([
    "limit",
    "prefix",
    "status",
    "start",
    "end",
  ].map(x => [x, searchParams.get(x)]))

  const [formData, setFormData] = useState({
    limit: params.limit,
    status: params.status,
    prefix: params.prefix,
    start: params.start,
    end: params.end,
  })
  const handleChange = (event) => {
    const { name, value } = event.target
    setFormData({...formData, [name]: value})
  }
  const handleSubmit = (event) => {
    event.preventDefault()
    const searchParams = Object.entries(formData)
      .filter(x => x[1] != null)
      .map(x => {
        switch (x[0]) {
          case 'start', 'end':
            return [x[0], toApiDatetime(x[1])]
          default:
            return x
        }
      })
      .map(x => {
        return x[0] + '=' + x[1]
      })
      .join('&')
    navigate(`/?${searchParams}`)
  }

  return <form className="d-flex" role="search" onSubmit={handleSubmit}>
    <input className="form-control me-2" type="number" min="10" max="100" name="limit" placeholder="Limit" aria-label="Search limit" defaultValue={formData.limit} onChange={handleChange} />
    <select className="form-control me-2" type="text" id="status" name="status" placeholder="Status" aria-label="Search status" defaultValue={formData.status} onChange={handleChange} >
    <option value="">--Select Status--</option>
    <option value="running">Running</option>
    <option value="pending">Pending</option>
    <option value="succeed">Succeed</option>
    <option value="failed">Failed</option>
    </select>
    <input className="form-control me-2" type="datetime-local" name="start" placeholder="Start CreatedAt" aria-label="Search Start CreatedAt" step="1" defaultValue={fromApiDatetime(formData.start)} onChange={handleChange} />
    <input className="form-control me-2" type="datetime-local" name="end" placeholder="End CreatedAt" aria-label="Search End CreatedAt" step="1" defaultValue={fromApiDatetime(formData.end)} onChange={handleChange} />
    <input className="form-control me-2" type="search" name="prefix" placeholder="TitlePrefix" aria-label="Search prefix" defaultValue={formData.prefix} onChange={handleChange} />
    <button className="btn btn-outline-success" type="submit">Search</button>
  </form>
}
