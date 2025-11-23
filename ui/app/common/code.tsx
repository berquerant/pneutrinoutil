export type CodeBlockParams = {
  code: string;
};

export function CodeBlock({ code }: CodeBlockParams) {
  return (
    <pre>
      <code>
        {code}
      </code>
    </pre>
  );
}
