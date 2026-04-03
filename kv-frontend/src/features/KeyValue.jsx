import { useState } from "react";
import { kvApi } from "../api/kvApi";
import ResultBox from "../components/ResultBox";

export default function KeyValue() {
  const [bucket, setBucket] = useState("");
  const [key, setKey] = useState("");
  const [value, setValue] = useState("");
  const [ttl, setTTL] = useState("");
  const [result, setResult] = useState(null);

  async function handlePut() {
    setResult(await kvApi.put(bucket, key, value, Number(ttl)));
  }

  async function handleGet() {
    setResult(await kvApi.get(bucket, key));
  }

  async function handleList() {
    setResult(await kvApi.list(bucket));
  }

  async function handleDelete() {
    setResult(await kvApi.del(bucket, key));
  }

  return (
    <div>
      <h2 className="text-lg text-green-400 mb-4">Key-Value Store</h2>

      <div className="grid grid-cols-2 gap-4">
        <input placeholder="Bucket" onChange={e => setBucket(e.target.value)} className="input" />
        <input placeholder="Key" onChange={e => setKey(e.target.value)} className="input" />
        <input placeholder="Value" onChange={e => setValue(e.target.value)} className="input" />
        <input placeholder="TTL (seconds)" onChange={e => setTTL(e.target.value)} className="input" />
      </div>

      <div className="mt-4 flex gap-3">
        <button onClick={handlePut} className="btn-green">PUT</button>
        <button onClick={handleGet} className="btn-blue">GET</button>
        <button onClick={handleList} className="btn-gray">LIST</button>
        <button onClick={handleDelete} className="btn-red">DELETE</button>
      </div>

      <ResultBox data={result} />
    </div>
  );
}
