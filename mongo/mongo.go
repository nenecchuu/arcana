package mongo

import (
	"context"
	"time"

	"github.com/pkg/errors"
	driver "go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

type (
	Client struct {
		URI            string
		DB             string
		AppName        string
		ConnectTimeout time.Duration
		PingTimeout    time.Duration
	}

	Database struct {
		Database *driver.Database
	}
)

var (
	defaultConnectTimeout = 10 * time.Second
	defaultPingTimeout    = 2 * time.Second
)

func MongoConnectClient(c *Client) *Database {
	client, err := c.MongoConnect()
	if err != nil {
		panic(err)
	}
	return &Database{
		Database: client.Database(c.DB),
	}
}

func (c *Client) MongoConnect() (mc *driver.Client, err error) {
	if c.ConnectTimeout == 0 {
		c.ConnectTimeout = defaultConnectTimeout
	}
	if c.PingTimeout == 0 {
		c.PingTimeout = defaultPingTimeout
	}

	ctx, _ := context.WithTimeout(context.Background(), c.ConnectTimeout)
	opts := []*options.ClientOptions{
		options.Client().SetConnectTimeout(c.ConnectTimeout).ApplyURI(c.URI).SetAppName(c.AppName),
	}
	mc, err = driver.Connect(ctx, opts...)
	if err != nil {
		err = errors.Wrap(err, "failed to create mongodb client")
		return
	}
	ctx, _ = context.WithTimeout(context.Background(), c.PingTimeout)
	if err = mc.Ping(ctx, readpref.Primary()); err != nil {
		err = errors.Wrap(err, "failed to establish connection to mongodb server")
	}
	return
}
