export default function ResultBox({ data }) {
  if (!data) return null;

  return (
    <pre className="mt-4 bg-black border border-zinc-700 p-4 rounded text-sm overflow-auto max-h-64">
      {JSON.stringify(data, null, 2)}
    </pre>
  );
}
