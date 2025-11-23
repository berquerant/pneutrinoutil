import Info, { InfoParams } from "./info";

export default function Component(params: InfoParams) {
  return (
    <div className="container">
      {Info(params)}
    </div>
  );
}
