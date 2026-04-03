export default function Header() {
  return (
    <div className="border-b border-zinc-700 px-6 py-4">
      <h1 className="text-xl font-bold text-green-400">
        KV-Distributed DB Console
      </h1>
      <p className="text-sm text-zinc-400">
        In-memory KV store with TTL, Sets, Lists, Queues & Stacks
      </p>
    </div>
  );
}
