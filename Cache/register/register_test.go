package register

import (
	"context"
	"testing"

	clientv3 "go.etcd.io/etcd/client/v3"
)

func TestRegister(t *testing.T) {
	cli, _ := clientv3.New(DefaultEtcdConfig)

	resp, err := cli.Grant(context.Background(), 5)
	if err != nil {
		t.Fatalf(err.Error())
	}

	err = etcdAdd(cli, resp.ID, "test", "localhost:6324")
	if err != nil {
		t.Fatalf(err.Error())
	}
}
