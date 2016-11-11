package elasticsearch

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"strconv"
	"sync"
	"time"
)

type Bulk struct {
	Client       http.Client
	BulkSize     int
	BulkInterval time.Duration
	startTime    time.Time
	w            *io.PipeWriter
	err          error
	mu           sync.Mutex
	wg           sync.WaitGroup
}

func NewBulk(index, typ, uri string) (*Bulk, error) {
	reader, writer := io.Pipe()
	req, err := http.NewRequest("POST", uri+fmt.Sprintf("/%s/%s/_bulk", index, typ), reader)
	if err != nil {
		return nil, err
	}
	b := &Bulk{
		BulkSize:     5 * 1024 * 1024,
		BulkInterval: 5 * time.Second,
		startTime:    time.Now(),
		w:            writer,
	}
	b.wg.Add(1)
	go b.do(req)
	return b, nil
}

func (b *Bulk) do(req *http.Request) {
	defer b.wg.Done()
	resp, err := b.Client.Do(req)
	if err != nil {
		b.mu.Lock()
		b.err = err
		b.mu.Unlock()
		return
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		errMsg, _ := ioutil.ReadAll(resp.Body)
		err := errors.New(strconv.Itoa(resp.StatusCode) + ": " + string(errMsg))
		b.mu.Lock()
		b.err = err
		b.mu.Unlock()
	}
}

func (b *Bulk) Index(v interface{}) error {
	b.mu.Lock()
	err := b.err
	b.mu.Unlock()
	if err != nil {
		return err
	}
	if _, err := b.w.Write([]byte(`{"index":{}}` + "\n")); err != nil {
		return err
	}
	return json.NewEncoder(b.w).Encode(v)
}

func (b *Bulk) Close() error {
	if err := b.w.Close(); err != nil {
		return err
	}
	b.wg.Wait()
	return b.err // no one will write b.err here
}
