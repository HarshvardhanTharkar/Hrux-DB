import { useState } from "react";
import { kvApi } from "../api/kvApi";
import ResultBox from "../components/ResultBox";

export default function MapView() {
  const [mapName, setMapName] = useState("");
  const [key, setKey] = useState("");
  const [value, setValue] = useState("");
  const [result, setResult] = useState(null);

  return (
    <div>
      <h2 className="text-green-400 mb-4">Map (Hash)</h2>

      <div className="grid grid-cols-3 gap-4">
        <input
          className="input"
          placeholder="Map name"
          onChange={e => setMapName(e.target.value)}
        />
        <input
          className="input"
          placeholder="Key"
          onChange={e => setKey(e.target.value)}
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
          onClick={async () => setResult(await kvApi.mapPut(mapName, key, value))}
        >
          PUT
        </button>

        <button
          className="btn-blue"
          onClick={async () => setResult(await kvApi.mapGet(mapName, key))}
        >
          GET
        </button>
      </div>

      <ResultBox data={result} />
    </div>
  );
}
