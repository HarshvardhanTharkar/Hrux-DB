import { useState } from "react";
import Queue from "./features/Queue";
import Stack from "./features/Stack";
import KeyValue from "./features/KeyValue";
import Sets from "./features/Sets";
import SortedList from "./features/SortedList";
import MapView from "./features/Map";

const tabs = [
  "KeyValue",
  "Sets",
  "SortedList",
  "Map",
  "Queue",
  "Stack",
];


export default function App() {
  const [active, setActive] = useState("KeyValue");

  return (
    <div className="h-screen flex">
      <div className="w-56 border-r border-zinc-700 p-4 space-y-2">
        {tabs.map(t => (
          <button
            key={t}
            onClick={() => setActive(t)}
            className={`w-full text-left px-3 py-2 rounded ${
              active === t ? "bg-zinc-800 text-green-400" : "hover:bg-zinc-800"
            }`}
          >
            {t}
          </button>
        ))}
      </div>

      <div className="flex-1 p-6 overflow-auto">
        {active === "KeyValue" && <KeyValue />}
        {active === "Sets" && <Sets />}
        {active === "Queue" && <Queue />}
        {active === "Stack" && <Stack />} 
        {active === "SortedList" && <SortedList />}
        {active === "Map" && <MapView />}

      </div>
    </div>
  );
}
