package elasticsearch

import (
	"strings"
	"testing"
	"time"

	testes "h12.me/realtest/elasticsearch"
)

func TestBulk(t *testing.T) {
	es, err := testes.New()
	if err != nil {
		t.Fatal(err)
	}
	index := testes.RandomIndexName()
	defer es.DeleteIndex(index)
	bulk, err := NewBulk(index, "x", "http://"+es.Addr())
	if err != nil {
		t.Fatal(err)
	}
	bulk.Index(map[string]string{"hello": "es"})
	if err := bulk.Close(); err != nil {
		t.Fatal(err)
	}
	time.Sleep(3 * time.Second)
	content, err := es.DumpIndex(index)
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(content, `"_source":{"hello":"es"}`) {
		t.Fatal("inserted data not found")
	}
}
