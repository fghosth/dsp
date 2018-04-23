package rabbitmq

import (
	"testing"
	"time"

	"jvole.com/dsp/config"
)

func TestRabbitMq_Publish(t *testing.T) {
	config := &Mqserver{
		"amqp://guest:guest@localhost:5672/",
		true,
	}
	exname := "dsp.index.notify"
	extype := "fanout"
	queueName := "index1"
	// queueName2 := "index2"
	// queueName3 := "index3"
	routeKey := "test"
	receiver := NewIndexReceiver(queueName, routeKey)
	// receiver2 := NewIndexReceiver(queueName2, routeKey)
	// receiver3 := NewIndexReceiver(queueName3, routeKey)
	mq := NewRabbitMQ(exname, extype, *config)
	go func() {
		for range time.Tick(time.Duration(1) * time.Second) { //每秒检查一次
			mq.Publish(exname, routeKey, []byte("ddddddddd====="))
		}

	}()

	mq.RegisterReceiver(receiver)
	// mq.RegisterReceiver(receiver2)
	// mq.RegisterReceiver(receiver3)
	mq.Start()
}

func TestIndexMSG(t *testing.T) {
	conf := &Mqserver{
		config.RabbitMQURL,
		true,
	}
	RabbitMQConn = NewRabbitMQ(config.IndexEXName, config.DSPextype, *conf)

	RabbitMQConn.Publish(config.IndexEXName, config.IndexRouteKey, []byte(`{"action":"add","cid":84,"t":1621604501}`))

}

func TestSuccessMSG(t *testing.T) {
	conf := &Mqserver{
		config.RabbitMQURL,
		true,
	}
	RabbitMQConn = NewRabbitMQ(config.IndexEXName, config.DSPextype, *conf)
	RabbitMQConn.Publish(config.IndexEXName, config.IndexRouteKey, []byte(`{"id":"DxU0032U8a","price":20000,"CID":2,"UID":14,"t":1521604501}`))

}
