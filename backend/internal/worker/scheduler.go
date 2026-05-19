package worker

import (
	"Distributed_Website_monitoring_system/internal/app"
	"Distributed_Website_monitoring_system/internal/kafka/producer"
	"Distributed_Website_monitoring_system/internal/models"
	"fmt"
	"sync"
	"time"

	"github.com/google/uuid"
)

type Scheduler struct {
	monitorRepo app.MonitorRepository
	checksRepo  app.ChecksRepository
	producer    *producer.Producer
	workerCount int
	queue       chan models.Monitor
	lastChecked map[uuid.UUID]time.Time
	mu          sync.Mutex
}

func NewScheduler(monitorRepo app.MonitorRepository, checksRepo app.ChecksRepository, prod *producer.Producer, workerCount int) *Scheduler {
	return &Scheduler{
		monitorRepo: monitorRepo,
		checksRepo:  checksRepo,
		producer:    prod,
		workerCount: workerCount,
		queue:       make(chan models.Monitor, 100),
		lastChecked: make(map[uuid.UUID]time.Time),
	}
}

func (s *Scheduler) Start() {
	for i := 0; i < s.workerCount; i++ {
		go s.work()
	}
	go s.dispatch()
}

func (s *Scheduler) work() {
	for monitor := range s.queue {
		realTime, statusOk, httpErr := app.Ping(monitor.Url)

		for _, region := range Regions {
			regionTime := realTime + region.RandomOffset()
			if err := app.SaveCheck(monitor, s.checksRepo, region.Name, regionTime, statusOk); err != nil {
				fmt.Println("WORKER: failed to save check for region", region.Name, err)
			}
		}

		if err := app.SendKafkaEvent(s.producer, monitor, realTime, statusOk); err != nil {
			fmt.Println("WORKER: kafka error:", err)
		}

		fmt.Println("WORKER: checked", monitor.Url, "err:", httpErr)
	}
}

func (s *Scheduler) dispatch() {
	ticker := time.NewTicker(5 * time.Second)
	defer ticker.Stop()

	for range ticker.C {
		monitors, err := s.monitorRepo.GetAllMonitors()
		if err != nil {
			fmt.Println("SCHEDULER: failed to get monitors:", err)
			continue
		}

		now := time.Now()
		s.mu.Lock()
		for _, m := range monitors {
			last, seen := s.lastChecked[m.Id]
			if !seen || now.Sub(last) >= time.Duration(m.Time_interval)*time.Second {
				s.lastChecked[m.Id] = now
				select {
				case s.queue <- m:
				default:
					fmt.Println("SCHEDULER: queue full, skipping", m.Url)
				}
			}
		}
		s.mu.Unlock()
	}
}
