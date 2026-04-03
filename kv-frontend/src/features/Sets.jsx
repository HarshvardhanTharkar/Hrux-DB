import { useState } from "react";
import { kvApi } from "../api/kvApi";
import ResultBox from "../components/ResultBox";

export default function Sets() {
  const [setName, setSetName] = useState("");
  const [value, setValue] = useState("");
  const [result, setResult] = useState(null);

  async function handleAdd() {
    setResult(await kvApi.setAdd(setName, value));
  }

  async function handleRemove() {
    setResult(await kvApi.setRemove(setName, value));
  }

  async function handleList() {
    setResult(await kvApi.setList(setName));
  }

  return (
    <div>
      <h2 className="text-lg text-green-400 mb-4">Sets</h2>

      <div className="grid grid-cols-2 gap-4">
        <input
          placeholder="Set Name"
          className="input"
          onChange={e => setSetName(e.target.value)}
        />
        <input
          placeholder="Value"
          className="input"
          onChange={e => setValue(e.target.value)}
        />
      </div>

      <div className="mt-4 flex gap-3">
        <button onClick={handleAdd} className="btn-green">
          ADD
        </button>
        <button onClick={handleRemove} className="btn-red">
          REMOVE
        </button>
        <button onClick={handleList} className="btn-gray">
          LIST
        </button>
      </div>

      <ResultBox data={result} />
    </div>
  );
}
