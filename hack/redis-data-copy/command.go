package main

import (
	"context"
	"errors"
	"fmt"
	"log"
	"time"

	"go.uber.org/zap"
	"gopkg.in/alecthomas/kingpin.v2"

	"github.com/bucketeer-io/bucketeer/pkg/cli"
	"github.com/bucketeer-io/bucketeer/pkg/metrics"
	v3 "github.com/bucketeer-io/bucketeer/pkg/redis/v3"
)

type command struct {
	*kingpin.CmdClause
	sourceRedisAddress      *string
	destinationRedisAddress *string
	srcPassword             *string
	dstPassword             *string
}

func registerCommand(r cli.CommandRegistry, p cli.ParentCommand) *command {
	cmd := p.Command("copy", "Copy data from source Redis to destination Redis")
	command := &command{
		CmdClause:               cmd,
		sourceRedisAddress:      cmd.Flag("source", "Source Redis address").Required().String(),
		destinationRedisAddress: cmd.Flag("destination", "Destination Redis address").Required().String(),
		srcPassword:             cmd.Flag("src-password", "Source Redis password").String(),
		dstPassword:             cmd.Flag("dst-password", "Destination Redis password").String(),
	}
	r.RegisterCommand(command)
	return command
}

func (c *command) Run(ctx context.Context, metrics metrics.Metrics, logger *zap.Logger) error {
	srcClient, err := v3.NewClient(*c.sourceRedisAddress,
		v3.WithLogger(logger),
		v3.WithPassword(*c.srcPassword),
		v3.WithPoolSize(10),
		v3.WithMinIdleConns(5),
		v3.WithMaxRetries(3),
		v3.WithDialTimeout(5*time.Second),
	)
	if err != nil {
		logger.Error("Error creating source Redis client", zap.Error(err))
		return err
	}
	defer srcClient.Close()

	dstClient, err := v3.NewClient(*c.destinationRedisAddress,
		v3.WithLogger(logger),
		v3.WithPassword(*c.dstPassword),
		v3.WithPoolSize(10),
		v3.WithMinIdleConns(5),
		v3.WithMaxRetries(3),
		v3.WithDialTimeout(5*time.Second),
	)
	if err != nil {
		logger.Error("Error creating destination Redis client", zap.Error(err))
		return err
	}
	defer dstClient.Close()

	if err := c.scanAndCopyBatch(srcClient, dstClient, logger); err != nil {
		logger.Error("Error during scan and copy process", zap.Error(err))
		return err
	}

	logger.Info("Data copy completed")
	return nil
}

func (c *command) scanAndCopyBatch(src, dst v3.Client, logger *zap.Logger) error {
	var cursor uint64
	batchSize := int64(100)

	for {
		nextCursor, keys, err := src.Scan(cursor, "*", batchSize)
		if err != nil {
			return fmt.Errorf("error scanning keys from source Redis: %v", err)
		}

		if err := c.copyBatch(src, dst, keys); err != nil {
			logger.Error("Error copying batch", zap.Error(err))
		}

		if nextCursor == 0 {
			break
		}
		cursor = nextCursor
	}

	return nil
}

func (c *command) copyBatch(src, dst v3.Client, keys []string) error {
	for _, key := range keys {
		value, err := src.Get(key)
		if err != nil {
			if errors.Is(err, v3.ErrNil) {
				log.Printf("Key not found: %s", key)
				continue
			}
			return fmt.Errorf("error getting value for key %s: %v", key, err)
		}

		err = dst.Set(key, value, 0)
		if err != nil {
			return fmt.Errorf("error setting value for key %s: %v", key, err)
		}
	}

	log.Printf("Successfully copied batch of %d keys", len(keys))
	return nil
}
