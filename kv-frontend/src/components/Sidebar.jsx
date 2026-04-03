const items = [
  "KeyValue",
  "Sets",
  "Queue",
  "Stack",
];

export default function Sidebar({ active, onSelect }) {
  return (
    <div className="w-56 border-r border-zinc-700 p-4 space-y-2">
      {items.map(item => (
        <button
          key={item}
          onClick={() => onSelect(item)}
          className={`w-full text-left px-3 py-2 rounded
            ${active === item
              ? "bg-zinc-800 text-green-400"
              : "hover:bg-zinc-800 text-zinc-300"}`}
        >
          {item}
        </button>
      ))}
    </div>
  );
}
