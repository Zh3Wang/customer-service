package etcd

import "context"

func (s *ServiceData) Put(k, v string) error {
	_, err := s.Client.Put(context.TODO(), k, v)
	if err != nil {
		return err
	}
	return nil
}

func (s *ServiceData) Get(k string) (string, error) {
	res, err := s.Client.Get(context.TODO(), k)
	if err != nil {
		return "", err
	}
	return string(res.Kvs[0].Value), err
}
