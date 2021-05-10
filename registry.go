package promgo

import (
	"sync"
)

const (
	// WorkerCount 执行协程数量
	WorkerCount = 5
)

// Registry ...
type Registry struct {
	mu         *sync.RWMutex // 读多，写少
	collectors map[string]Collector
}

var defaultRegistry *Registry

func init() {
	defaultRegistry = NewRegistry()
}

// GetDefaultRegistry ...
func GetDefaultRegistry() *Registry {
	return defaultRegistry
}

// NewRegistry ...
func NewRegistry() *Registry {
	return &Registry{
		mu:         new(sync.RWMutex),
		collectors: make(map[string]Collector),
	}
}

// Register 注册
func (r Registry) Register(c Collector) {
	id := c.Describe().ID()
	r.mu.Lock()
	defer r.mu.Unlock()

	if _, ok := r.collectors[id]; ok {
		return
	}

	r.collectors[id] = c
}

// Unregister 取消注册
func (r Registry) Unregister(c Collector) {
	id := c.Describe().ID()
	r.mu.Lock()
	defer r.mu.Unlock()

	if _, ok := r.collectors[id]; !ok {
		return
	}
	delete(r.collectors, id)
}

// Collect 收集
func (r Registry) Collect() []Metric {
	r.mu.RLock()
	defer r.mu.RUnlock()

	chProcess := make(chan struct{}, WorkerCount)
	done := make(chan struct{})
	defer func() {
		close(chProcess)
		close(done)
	}()

	// 遍历所有的采集器，获取对应的指标数据
	ch := make(chan *MetricErr)
	var metrics []Metric

	// 启动协程消费采集到的数据
	go func() {
		for me := range ch {
			if me.Err != nil {
				continue // 这个Err 可以考虑该如何处理
			}

			metrics = append(metrics, *me.Metric)
		}

		done <- struct{}{}
	}()

	wg := new(sync.WaitGroup)
	for _, c := range r.collectors {
		wg.Add(1)
		chProcess <- struct{}{}

		// 启动协程来完成数据收集工作
		go func(c Collector) {
			defer func() {
				wg.Done()
			}()

			c.Collect(ch)
			<-chProcess
		}(c)
	}
	wg.Wait()
	close(ch) // 关闭metric 指标数据通道
	<-done

	return metrics
}
