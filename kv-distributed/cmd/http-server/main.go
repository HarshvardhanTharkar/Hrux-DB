package main

import (
	"encoding/json"
	"log"
	"net/http"
	"sort"
	"sync"
	"time"
)

/* ===================== CORS ===================== */

func cors(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type")

		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusOK)
			return
		}
		next(w, r)
	}
}

/* ===================== TYPES ===================== */

type KVRequest struct {
	Bucket string `json:"bucket,omitempty"`
	Key    string `json:"key,omitempty"`
	Value  string `json:"value,omitempty"`
	TTL    int64  `json:"ttl,omitempty"`

	SetName  string `json:"setName,omitempty"`
	ListName string `json:"listName,omitempty"`
	MapName  string `json:"mapName,omitempty"`

	ValueStr string `json:"valueStr,omitempty"`

	QueueName string `json:"queueName,omitempty"`
	StackName string `json:"stackName,omitempty"`
}

type TransactionRequest struct {
	Operations []TransactionOp `json:"operations"`
}

type TransactionOp struct {
	Action string `json:"action"`
	Bucket string `json:"bucket"`
	Key    string `json:"key"`
	Value  string `json:"value,omitempty"`
}

type KVResponse struct {
	Success bool        `json:"success"`
	Data    interface{} `json:"data,omitempty"`
	Error   string      `json:"error,omitempty"`
}

/* ===================== STORAGE ===================== */

type Item struct {
	Value     string
	ExpiresAt time.Time
}

var (
	// KV
	data = make(map[string]map[string]Item)

	// Data structures
	sets        = make(map[string]map[string]bool)
	sortedLists = make(map[string][]string)
	mapsStore   = make(map[string]map[string]string)
	queues      = make(map[string][]string)
	stacks      = make(map[string][]string)

	// Locks
	dataMux   sync.RWMutex
	setsMux   sync.RWMutex
	sortedMux sync.RWMutex
	mapMux    sync.RWMutex
	queueMux  sync.RWMutex
	stackMux  sync.RWMutex
)

/* ===================== HELPERS ===================== */

func expired(item Item) bool {
	return !item.ExpiresAt.IsZero() && time.Now().After(item.ExpiresAt)
}

/* ===================== MAIN ===================== */

func main() {

	/* ---------- HEALTH ---------- */
	http.HandleFunc("/health", cors(func(w http.ResponseWriter, _ *http.Request) {
		json.NewEncoder(w).Encode(map[string]string{"status": "ok"})
	}))

	/* ---------- KV PUT ---------- */
	http.HandleFunc("/put", cors(func(w http.ResponseWriter, r *http.Request) {
		var req KVRequest
		json.NewDecoder(r.Body).Decode(&req)

		if req.Bucket == "" || req.Key == "" {
			json.NewEncoder(w).Encode(KVResponse{Success: false, Error: "bucket & key required"})
			return
		}

		dataMux.Lock()
		if data[req.Bucket] == nil {
			data[req.Bucket] = make(map[string]Item)
		}

		item := Item{Value: req.Value}
		if req.TTL > 0 {
			item.ExpiresAt = time.Now().Add(time.Second * time.Duration(req.TTL))
		}

		data[req.Bucket][req.Key] = item
		dataMux.Unlock()

		json.NewEncoder(w).Encode(KVResponse{Success: true})
	}))

	/* ---------- KV GET ---------- */
	http.HandleFunc("/get", cors(func(w http.ResponseWriter, r *http.Request) {
		var req KVRequest
		json.NewDecoder(r.Body).Decode(&req)

		dataMux.Lock()
		defer dataMux.Unlock()

		item, ok := data[req.Bucket][req.Key]
		if !ok || expired(item) {
			delete(data[req.Bucket], req.Key)
			json.NewEncoder(w).Encode(KVResponse{Success: false, Error: "not found"})
			return
		}

		json.NewEncoder(w).Encode(KVResponse{Success: true, Data: item.Value})
	}))

	/* ---------- KV LIST ---------- */
	http.HandleFunc("/list", cors(func(w http.ResponseWriter, r *http.Request) {
		var req KVRequest
		json.NewDecoder(r.Body).Decode(&req)

		result := map[string]string{}

		dataMux.Lock()
		for k, v := range data[req.Bucket] {
			if !expired(v) {
				result[k] = v.Value
			}
		}
		dataMux.Unlock()

		json.NewEncoder(w).Encode(KVResponse{Success: true, Data: result})
	}))

	/* ---------- KV DELETE ---------- */
	http.HandleFunc("/delete", cors(func(w http.ResponseWriter, r *http.Request) {
		var req KVRequest
		json.NewDecoder(r.Body).Decode(&req)

		dataMux.Lock()
		delete(data[req.Bucket], req.Key)
		dataMux.Unlock()

		json.NewEncoder(w).Encode(KVResponse{Success: true})
	}))

	/* ---------- SET ---------- */
	http.HandleFunc("/set/add", cors(func(w http.ResponseWriter, r *http.Request) {
		var req KVRequest
		json.NewDecoder(r.Body).Decode(&req)

		setsMux.Lock()
		if sets[req.SetName] == nil {
			sets[req.SetName] = make(map[string]bool)
		}
		sets[req.SetName][req.ValueStr] = true
		setsMux.Unlock()

		json.NewEncoder(w).Encode(KVResponse{Success: true})
	}))

	http.HandleFunc("/set/remove", cors(func(w http.ResponseWriter, r *http.Request) {
		var req KVRequest
		json.NewDecoder(r.Body).Decode(&req)

		setsMux.Lock()
		delete(sets[req.SetName], req.ValueStr)
		setsMux.Unlock()

		json.NewEncoder(w).Encode(KVResponse{Success: true})
	}))

	http.HandleFunc("/set/list", cors(func(w http.ResponseWriter, r *http.Request) {
		var req KVRequest
		json.NewDecoder(r.Body).Decode(&req)

		var list []string
		setsMux.RLock()
		for v := range sets[req.SetName] {
			list = append(list, v)
		}
		setsMux.RUnlock()

		json.NewEncoder(w).Encode(KVResponse{Success: true, Data: list})
	}))

	/* ---------- SORTED LIST ---------- */
	http.HandleFunc("/sortedlist/add", cors(func(w http.ResponseWriter, r *http.Request) {
		var req KVRequest
		json.NewDecoder(r.Body).Decode(&req)

		sortedMux.Lock()
		sortedLists[req.ListName] = append(sortedLists[req.ListName], req.ValueStr)
		sort.Strings(sortedLists[req.ListName])
		sortedMux.Unlock()

		json.NewEncoder(w).Encode(KVResponse{Success: true})
	}))

	http.HandleFunc("/sortedlist/get", cors(func(w http.ResponseWriter, r *http.Request) {
		var req KVRequest
		json.NewDecoder(r.Body).Decode(&req)

		sortedMux.RLock()
		list := sortedLists[req.ListName]
		sortedMux.RUnlock()

		json.NewEncoder(w).Encode(KVResponse{Success: true, Data: list})
	}))

	/* ---------- MAP ---------- */
	http.HandleFunc("/map/put", cors(func(w http.ResponseWriter, r *http.Request) {
		var req KVRequest
		json.NewDecoder(r.Body).Decode(&req)

		mapMux.Lock()
		if mapsStore[req.MapName] == nil {
			mapsStore[req.MapName] = make(map[string]string)
		}
		mapsStore[req.MapName][req.Key] = req.ValueStr
		mapMux.Unlock()

		json.NewEncoder(w).Encode(KVResponse{Success: true})
	}))

	http.HandleFunc("/map/get", cors(func(w http.ResponseWriter, r *http.Request) {
		var req KVRequest
		json.NewDecoder(r.Body).Decode(&req)

		mapMux.RLock()
		val := mapsStore[req.MapName][req.Key]
		mapMux.RUnlock()

		json.NewEncoder(w).Encode(KVResponse{Success: true, Data: val})
	}))

	/* ---------- QUEUE ---------- */
	http.HandleFunc("/queue/push", cors(func(w http.ResponseWriter, r *http.Request) {
		var req KVRequest
		json.NewDecoder(r.Body).Decode(&req)

		queueMux.Lock()
		queues[req.QueueName] = append(queues[req.QueueName], req.ValueStr)
		queueMux.Unlock()

		json.NewEncoder(w).Encode(KVResponse{Success: true})
	}))

	http.HandleFunc("/queue/pop", cors(func(w http.ResponseWriter, r *http.Request) {
		var req KVRequest
		json.NewDecoder(r.Body).Decode(&req)

		queueMux.Lock()
		defer queueMux.Unlock()

		q := queues[req.QueueName]
		if len(q) == 0 {
			json.NewEncoder(w).Encode(KVResponse{Success: false, Error: "empty"})
			return
		}

		val := q[0]
		queues[req.QueueName] = q[1:]

		json.NewEncoder(w).Encode(KVResponse{Success: true, Data: val})
	}))

	/* ---------- STACK ---------- */
	http.HandleFunc("/stack/push", cors(func(w http.ResponseWriter, r *http.Request) {
		var req KVRequest
		json.NewDecoder(r.Body).Decode(&req)

		stackMux.Lock()
		stacks[req.StackName] = append(stacks[req.StackName], req.ValueStr)
		stackMux.Unlock()

		json.NewEncoder(w).Encode(KVResponse{Success: true})
	}))

	http.HandleFunc("/stack/pop", cors(func(w http.ResponseWriter, r *http.Request) {
		var req KVRequest
		json.NewDecoder(r.Body).Decode(&req)

		stackMux.Lock()
		defer stackMux.Unlock()

		s := stacks[req.StackName]
		if len(s) == 0 {
			json.NewEncoder(w).Encode(KVResponse{Success: false, Error: "empty"})
			return
		}

		val := s[len(s)-1]
		stacks[req.StackName] = s[:len(s)-1]

		json.NewEncoder(w).Encode(KVResponse{Success: true, Data: val})
	}))

	/* ---------- TRANSACTION ---------- */
	http.HandleFunc("/transaction", cors(func(w http.ResponseWriter, r *http.Request) {
		var req TransactionRequest
		json.NewDecoder(r.Body).Decode(&req)

		dataMux.Lock()
		defer dataMux.Unlock()

		for _, op := range req.Operations {
			switch op.Action {
			case "put":
				if data[op.Bucket] == nil {
					data[op.Bucket] = make(map[string]Item)
				}
				data[op.Bucket][op.Key] = Item{Value: op.Value}
			case "delete":
				delete(data[op.Bucket], op.Key)
			default:
				json.NewEncoder(w).Encode(KVResponse{Success: false, Error: "invalid op"})
				return
			}
		}

		json.NewEncoder(w).Encode(KVResponse{Success: true})
	}))

	log.Println("🚀 HTTP KV Server running on :8081")
	log.Fatal(http.ListenAndServe(":8081", nil))
}
