package service

import "errors"

type TransactionOp struct {
	Action string
	Bucket string
	Key    string
	Value  string
}

func (s *KVService) ExecuteTransaction(ops []TransactionOp) []error {
	var errs []error
	for _, op := range ops {
		switch op.Action {
		case "put":
			s.Put(op.Bucket, op.Key, []byte(op.Value))
			errs = append(errs, nil)
		case "delete":
			errs = append(errs, s.Delete(op.Bucket, op.Key))
		case "update":
			errs = append(errs, s.Update(op.Bucket, op.Key, []byte(op.Value)))
		default:
			errs = append(errs, errors.New("invalid action"))
		}
	}
	return errs
}
