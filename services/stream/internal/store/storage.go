package store

import (
	"bytes"
	"io"
	"os"
	"path"
	"sync"
	"time"

	"github.com/anacrolix/missinggo/v2/resource"
)

type Instance struct {
	data     []byte
	name     string
	provider *Provider
	mu       sync.RWMutex
}

func (p *Provider) NewInstance(name string) (resource.Instance, error) {
	p.mu.Lock()
	defer p.mu.Unlock()

	if instance, ok := p.instances[name]; ok {
		return instance, nil
	}

	instance := &Instance{
		name:     name,
		provider: p,
	}

	p.instances[name] = instance

	parent := path.Dir(name)
	base := path.Base(name)

	if _, ok := p.dirs[parent]; !ok {
		p.dirs[parent] = make(map[string]struct{})
	}

	p.dirs[parent][base] = struct{}{}

	return instance, nil
}

func (i *Instance) setData(data []byte) {
	oldSize := int64(len(i.data))
	newSize := int64(len(data))

	i.data = data

	if oldSize == newSize {
		return
	}

	i.provider.mu.Lock()
	i.provider.size += newSize - oldSize
	i.provider.mu.Unlock()
}

func (i *Instance) ReadAt(p []byte, off int64) (int, error) {
	i.mu.RLock()
	defer i.mu.RUnlock()

	if off >= int64(len(i.data)) {
		return 0, io.EOF
	}

	n := copy(p, i.data[off:])

	if n < len(p) {
		return n, io.EOF
	}

	return n, nil
}

func (i *Instance) WriteAt(p []byte, off int64) (int, error) {
	i.mu.Lock()
	defer i.mu.Unlock()

	end := off + int64(len(p))

	if end > int64(len(i.data)) {
		newData := make([]byte, end)
		copy(newData, i.data)

		i.setData(newData)
	}

	copy(i.data[off:], p)

	return len(p), nil
}

func (i *Instance) Get() (io.ReadCloser, error) {
	i.mu.RLock()
	defer i.mu.RUnlock()

	return io.NopCloser(bytes.NewReader(i.data)), nil
}

func (i *Instance) Put(r io.Reader) error {
	data, err := io.ReadAll(r)
	if err != nil {
		return err
	}

	i.mu.Lock()
	defer i.mu.Unlock()

	i.setData(data)

	return nil
}

func (i *Instance) Delete() error {
	i.mu.Lock()
	defer i.mu.Unlock()

	i.setData(nil)

	return nil
}

func (i *Instance) Readdirnames() ([]string, error) {
	i.provider.mu.Lock()
	defer i.provider.mu.Unlock()

	names, ok := i.provider.dirs[i.name]
	if !ok {
		return nil, os.ErrNotExist
	}

	result := make([]string, 0, len(names))

	for name := range names {
		result = append(result, name)
	}

	return result, nil
}

type fileInfo struct {
	size int64
}

func (f fileInfo) Name() string {
	return ""
}

func (f fileInfo) Size() int64 {
	return f.size
}

func (f fileInfo) Mode() os.FileMode {
	return 0
}

func (f fileInfo) ModTime() time.Time {
	return time.Time{}
}

func (f fileInfo) IsDir() bool {
	return false
}

func (f fileInfo) Sys() any {
	return nil
}

func (i *Instance) Stat() (os.FileInfo, error) {
	i.mu.RLock()
	defer i.mu.RUnlock()

	return fileInfo{
		size: int64(len(i.data)),
	}, nil
}

type Provider struct {
	instances map[string]*Instance
	dirs      map[string]map[string]struct{}
	capacity  int64
	size      int64
	mu        sync.Mutex
}

const DefaultCapacity int64 = 512 * 1024 * 1024

func NewProvider(capacity int64) *Provider {
	return &Provider{
		instances: make(map[string]*Instance),
		dirs:      make(map[string]map[string]struct{}),
		capacity:  capacity,
	}
}
