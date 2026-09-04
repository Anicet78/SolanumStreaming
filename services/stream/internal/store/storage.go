package store

import (
	"io"
	"os"
	"path"
	"sync"
	"time"

	"github.com/anacrolix/missinggo/v2/resource"
)

type Provider struct {
	instances map[string]*Instance
	dirs      map[string]map[string]struct{}
	size      int64
	mu        sync.Mutex
}

func NewProvider() *Provider {
	return &Provider{
		instances: make(map[string]*Instance),
		dirs:      make(map[string]map[string]struct{}),
	}
}

const chunkSize int64 = 1024 * 1024

type Instance struct {
	name     string
	provider *Provider
	chunks   map[int64][]byte
	size     int64
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
		chunks:   make(map[int64][]byte),
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

func (p *Provider) increaseSize(delta int64) {
	p.mu.Lock()
	defer p.mu.Unlock()

	p.size += delta
}

func (i *Instance) ReadAt(p []byte, off int64) (int, error) {
	i.mu.RLock()
	defer i.mu.RUnlock()

	if off < 0 {
		return 0, os.ErrInvalid
	}

	if off >= i.size {
		return 0, io.EOF
	}

	originalLen := len(p)

	if int64(len(p)) > i.size-off {
		p = p[:i.size-off]
	}

	read := 0

	for len(p) > 0 {
		chunkIndex := off / chunkSize
		chunkOffset := off % chunkSize

		n := min(int(chunkSize-chunkOffset), len(p))

		chunk, ok := i.chunks[chunkIndex]
		if !ok {
			return read, io.EOF
		}

		copy(p[:n], chunk[chunkOffset:chunkOffset+int64(n)])

		p = p[n:]
		off += int64(n)
		read += n
	}

	if read < originalLen {
		return read, io.EOF
	}

	return read, nil
}

func (i *Instance) WriteAt(p []byte, off int64) (int, error) {
	i.mu.Lock()
	defer i.mu.Unlock()

	if off < 0 {
		return 0, os.ErrInvalid
	}

	written := 0

	for len(p) > 0 {
		chunkIndex := off / chunkSize
		chunkOffset := off % chunkSize

		n := min(int(chunkSize-chunkOffset), len(p))

		chunk, ok := i.chunks[chunkIndex]
		if !ok {
			chunk = make([]byte, chunkSize)
			i.chunks[chunkIndex] = chunk
			i.provider.increaseSize(chunkSize)
		}

		copy(chunk[chunkOffset:], p[:n])

		p = p[n:]
		off += int64(n)
		written += n
	}

	if off > i.size {
		i.size = off
	}

	return written, nil
}

type instanceReader struct {
	instance *Instance
	offset   int64
}

func (r *instanceReader) Read(p []byte) (int, error) {
	n, err := r.instance.ReadAt(p, r.offset)
	r.offset += int64(n)
	return n, err
}

func (r *instanceReader) Close() error {
	return nil
}

func (i *Instance) Get() (io.ReadCloser, error) {
	return &instanceReader{
		instance: i,
	}, nil
}

func (i *Instance) Put(r io.Reader) error {
	i.mu.Lock()
	defer i.mu.Unlock()

	freed := int64(len(i.chunks)) * chunkSize
	i.provider.increaseSize(-freed)

	i.chunks = make(map[int64][]byte)
	i.size = 0

	buf := make([]byte, chunkSize)
	var off int64

	for {
		n, err := io.ReadFull(r, buf)

		if n > 0 {
			chunk := make([]byte, chunkSize)
			copy(chunk, buf[:n])

			i.chunks[off/chunkSize] = chunk
			i.provider.increaseSize(chunkSize)

			off += int64(n)
		}

		if err == io.EOF {
			break
		}

		if err == io.ErrUnexpectedEOF {
			break
		}

		if err != nil {
			return err
		}
	}

	i.size = off

	return nil
}

func (i *Instance) Delete() error {
	i.mu.Lock()
	defer i.mu.Unlock()

	if len(i.chunks) == 0 {
		i.size = 0
		return nil
	}

	freed := int64(len(i.chunks) * int(chunkSize))

	i.chunks = make(map[int64][]byte)
	i.size = 0

	i.provider.increaseSize(-freed)

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
		size: i.size,
	}, nil
}
