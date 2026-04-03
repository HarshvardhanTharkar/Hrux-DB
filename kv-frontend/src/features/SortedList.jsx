import { useState } from "react";
import { kvApi } from "../api/kvApi";
import ResultBox from "../components/ResultBox";

export default function SortedList() {
  const [name, setName] = useState("");
  const [value, setValue] = useState("");
  const [result, setResult] = useState(null);

  return (
    <div>
      <h2 className="text-green-400 mb-4">Sorted List</h2>

      <div className="grid grid-cols-2 gap-4">
        <input
          className="input"
          placeholder="List name"
          onChange={e => setName(e.target.value)}
        />
        <input
          className="input"
          placeholder="Value"
          onChange={e => setValue(e.target.value)}
        />
      </div>

      <div className="mt-4 flex gap-3">
        <button
          className="btn-green"
          onClick={async () => setResult(await kvApi.sortedListAdd(name, value))}
        >
          ADD
        </button>

        <button
          className="btn-gray"
          onClick={async () => setResult(await kvApi.sortedListGet(name))}
        >
          GET
        </button>
      </div>

      <ResultBox data={result} />
    </div>
  );
}
