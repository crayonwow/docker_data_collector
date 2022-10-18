package watcher

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"html/template"
	"time"

	"github.com/docker/docker/api/types"
	docker "github.com/docker/docker/client"
	"github.com/docker/go-units"
	"github.com/sirupsen/logrus"
)

const (
	statsTemplate = `Name: {{.Name}}
CPU %:    {{FormatPercentage .CPUPercentage}}
MEM %:    {{FormatPercentage .MemoryPercentage}}
NETWORK : {{.NetworkIO}}
`
)

func (w *watcher) statsHandler() ContainerHandler {
	return func(ctx context.Context, cli *docker.Client, container types.Container) error {
		statsRaw, err := cli.ContainerStatsOneShot(ctx, container.ID)
		if err != nil {
			return fmt.Errorf("container stats: %w", err)
		}
		defer func() {
			cErr := statsRaw.Body.Close()
			if cErr != nil {
				logrus.WithError(cErr).Error("containers stats one shot")
			}
		}()

		stats := &types.StatsJSON{}
		dec := json.NewDecoder(statsRaw.Body)
		err = dec.Decode(stats)
		if err != nil {
			return fmt.Errorf("decode : %w", err)
		}

		previousCPU := stats.PreCPUStats.CPUUsage.TotalUsage
		previousSystem := stats.PreCPUStats.SystemUsage
		cpuPercent := calculateCPUPercentUnix(previousCPU, previousSystem, stats)
		blkRead, blkWrite := calculateBlockIO(stats.BlkioStats)
		mem := calculateMemUsageUnixNoCache(stats.MemoryStats)
		memLimit := float64(stats.MemoryStats.Limit)
		memPercent := calculateMemPercentUnixNoCache(memLimit, mem)
		pidsStatsCurrent := stats.PidsStats.Current
		netRx, netTx := calculateNetwork(stats.Networks)

		se := statsEntry{
			Name:             stats.Name,
			ID:               stats.ID,
			CPUPercentage:    cpuPercent,
			Memory:           mem,
			MemoryPercentage: memPercent,
			MemoryLimit:      memLimit,
			NetworkRx:        netRx,
			NetworkTx:        netTx,
			NetworkIO:        units.HumanSizeWithPrecision(netTx, 3) + " / " + units.HumanSizeWithPrecision(netTx, 3),
			BlockRead:        float64(blkRead),
			BlockWrite:       float64(blkWrite),
			PidsCurrent:      pidsStatsCurrent,
		}
		message, err := prepareStatTemplate(se)
		if err != nil {
			return fmt.Errorf("prepare stats template: %w", err)
		}
		err = w.s.Send(message)
		if err != nil {
			return fmt.Errorf("send mark down: %w", err)
		}
		return nil
	}
}

func (w *watcher) createdHandler() ContainerHandler {
	return func(ctx context.Context, cli *docker.Client, container types.Container) error {
		inspect, err := cli.ContainerInspect(ctx, container.ID)
		if err != nil {
			return fmt.Errorf("container inspect: %w", err)
		}
		created, err := time.Parse(time.RFC3339Nano, inspect.Created)
		if err != nil {
			return fmt.Errorf("time parse: %w", err)
		}
		if created.Day() != time.Now().Day() {
			return nil
		}

		err = w.s.Send(fmt.Sprintf("%s IT's PAYING TIME BABYY", container.Names))
		if err != nil {
			return fmt.Errorf("send: %w", err)
		}

		return nil
	}
}

func prepareStatTemplate(s statsEntry) (string, error) {
	tmlt, err := template.New("stats_message").Funcs(template.FuncMap{
		"FormatPercentage": formatPercentage,
	}).Parse(statsTemplate)
	if err != nil {
		return "", fmt.Errorf("parse template: %w", err)
	}
	buf := bytes.Buffer{}
	err = tmlt.Execute(&buf, s)
	if err != nil {
		return "", fmt.Errorf("execute: %w", err)
	}
	return buf.String(), nil
}
