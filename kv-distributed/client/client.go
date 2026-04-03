package client

import (
	"net/rpc"

	"kv-distributed/internal/api"
)

type KVClient struct {
	client *rpc.Client
}

func NewKVClient(serverAddr string) (*KVClient, error) {
	client, err := rpc.Dial("tcp", serverAddr)
	if err != nil {
		return nil, err
	}
	return &KVClient{client: client}, nil
}

func (c *KVClient) Put(bucket, key string, value []byte, ttlSeconds int64) error {
	req := api.Request{
		Bucket:     bucket,
		Key:        key,
		Value:      value,
		TTLSeconds: ttlSeconds,
	}
	var resp api.Response
	return c.client.Call("KVServer.Put", req, &resp)
}

func (c *KVClient) Get(bucket, key string) ([]byte, error) {
	req := api.Request{
		Bucket: bucket,
		Key:    key,
	}
	var resp api.Response
	err := c.client.Call("KVServer.Get", req, &resp)
	return resp.Value, err
}

func (c *KVClient) Delete(bucket, key string) error {
	req := api.Request{
		Bucket: bucket,
		Key:    key,
	}
	var resp api.Response
	return c.client.Call("KVServer.Delete", req, &resp)
}

func (c *KVClient) Update(bucket, key string, value []byte) error {
	req := api.Request{
		Bucket: bucket,
		Key:    key,
		Value:  value,
	}
	var resp api.Response
	return c.client.Call("KVServer.Update", req, &resp)
}

func (c *KVClient) List(bucket string) (map[string][]byte, error) {
	req := api.Request{
		Bucket: bucket,
	}
	var resp api.Response
	err := c.client.Call("KVServer.List", req, &resp)
	return resp.Values, err
}

func (c *KVClient) ListKeysByPrefix(bucket, prefix string) ([]string, error) {
	req := api.Request{
		Bucket: bucket,
		Prefix: prefix,
	}
	var resp api.Response
	err := c.client.Call("KVServer.ListKeysByPrefix", req, &resp)
	return resp.Keys, err
}

func (c *KVClient) ListKeysInRange(bucket, start, end string) ([]string, error) {
	req := api.Request{
		Bucket: bucket,
		Start:  start,
		End:    end,
	}
	var resp api.Response
	err := c.client.Call("KVServer.ListKeysInRange", req, &resp)
	return resp.Keys, err
}

func (c *KVClient) SetAdd(setName, value string) error {
	req := api.Request{
		SetName:  setName,
		ValueStr: value,
	}
	var resp api.Response
	return c.client.Call("KVServer.SetAdd", req, &resp)
}

func (c *KVClient) SetRemove(setName, value string) error {
	req := api.Request{
		SetName:  setName,
		ValueStr: value,
	}
	var resp api.Response
	return c.client.Call("KVServer.SetRemove", req, &resp)
}

func (c *KVClient) SetList(setName string) ([]string, error) {
	req := api.Request{
		SetName: setName,
	}
	var resp api.Response
	err := c.client.Call("KVServer.SetList", req, &resp)
	return resp.List, err
}

func (c *KVClient) SortedListAdd(listName, value string) error {
	req := api.Request{
		ListName: listName,
		ValueStr: value,
	}
	var resp api.Response
	return c.client.Call("KVServer.SortedListAdd", req, &resp)
}

func (c *KVClient) SortedListGet(listName string) ([]string, error) {
	req := api.Request{
		ListName: listName,
	}
	var resp api.Response
	err := c.client.Call("KVServer.SortedListGet", req, &resp)
	return resp.List, err
}

func (c *KVClient) MapPut(mapName, key, value string) error {
	req := api.Request{
		MapName:  mapName,
		Key:      key,
		ValueStr: value,
	}
	var resp api.Response
	return c.client.Call("KVServer.MapPut", req, &resp)
}

func (c *KVClient) MapGet(mapName, key string) (interface{}, error) {
	req := api.Request{
		MapName: mapName,
		Key:     key,
	}
	var resp api.Response
	err := c.client.Call("KVServer.MapGet", req, &resp)
	return resp.Data, err
}

func (c *KVClient) QueuePush(name, value string) error {
	req := api.Request{
		QueueName: name,
		ValueStr:  value,
	}
	var resp api.Response
	return c.client.Call("KVServer.QueuePush", req, &resp)
}

func (c *KVClient) QueuePop(name string) (string, error) {
	req := api.Request{
		QueueName: name,
	}
	var resp api.Response
	err := c.client.Call("KVServer.QueuePop", req, &resp)
	return resp.Message, err
}

func (c *KVClient) QueuePeek(name string) (string, error) {
	req := api.Request{
		QueueName: name,
	}
	var resp api.Response
	err := c.client.Call("KVServer.QueuePeek", req, &resp)
	return resp.Message, err
}

func (c *KVClient) StackPush(name, value string) error {
	req := api.Request{
		StackName: name,
		ValueStr:  value,
	}
	var resp api.Response
	return c.client.Call("KVServer.StackPush", req, &resp)
}

func (c *KVClient) StackPop(name string) (string, error) {
	req := api.Request{
		StackName: name,
	}
	var resp api.Response
	err := c.client.Call("KVServer.StackPop", req, &resp)
	return resp.Message, err
}

func (c *KVClient) StackPeek(name string) (string, error) {
	req := api.Request{
		StackName: name,
	}
	var resp api.Response
	err := c.client.Call("KVServer.StackPeek", req, &resp)
	return resp.Message, err
}

func (c *KVClient) ExecuteTransaction(ops []api.TransactionOp) ([]error, error) {
	req := api.Request{
		Transaction: ops,
	}
	var resp api.Response
	err := c.client.Call("KVServer.ExecuteTransaction", req, &resp)
	return resp.Errors, err
}

func (c *KVClient) Close() error {
	return c.client.Close()
}
