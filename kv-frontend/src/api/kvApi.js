async function post(url, body = {}) {
  const res = await fetch(`/api${url}`, {
    method: "POST",
    headers: { "Content-Type": "application/json" },
    body: JSON.stringify(body),
  });

  const text = await res.text();
  try {
    return JSON.parse(text);
  } catch {
    return { success: false, error: text || "Server error" };
  }
}

export const kvApi = {
  put: (bucket, key, value, ttl) =>
    post("/put", { bucket, key, value, ttl }),

  get: (bucket, key) =>
    post("/get", { bucket, key }),

  list: (bucket) =>
    post("/list", { bucket }),

  del: (bucket, key) =>
    post("/delete", { bucket, key }),

  setAdd: (setName, valueStr) =>
    post("/set/add", { setName, valueStr }),

  setRemove: (setName, valueStr) =>
    post("/set/remove", { setName, valueStr }),

  setList: (setName) =>
    post("/set/list", { setName }),

  queuePush: (queueName, valueStr) =>
    post("/queue/push", { queueName, valueStr }),

  queuePop: (queueName) =>
    post("/queue/pop", { queueName }),

  stackPush: (stackName, valueStr) =>
    post("/stack/push", { stackName, valueStr }),

  stackPop: (stackName) =>
    post("/stack/pop", { stackName }),

  transaction: (operations) =>
    post("/transaction", { operations }),
  sortedListAdd: (listName, valueStr) =>
  post("/sortedlist/add", { listName, valueStr }),

sortedListGet: (listName) =>
  post("/sortedlist/get", { listName }),

mapPut: (mapName, key, valueStr) =>
  post("/map/put", { mapName, key, valueStr }),

mapGet: (mapName, key) =>
  post("/map/get", { mapName, key }),

};
