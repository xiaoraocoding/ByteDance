package mes

import (
	"ByteDance/model"
	"encoding/json"
	"fmt"
	"time"

	"github.com/nsqio/go-nsq"
)

var (
	Consumer *nsq.Consumer
	address  = "127.0.0.1:4161"
)

type Message struct {
	Followed      model.Followed_sql
	FollowCount   int
	FollowerCount int
}

type ConsumerT struct{}

func (*ConsumerT) HandleMessage(msg *nsq.Message) error {
	fmt.Println("receive", msg.NSQDAddress, "message:", string(msg.Body))
	message := Message{}
	err := json.Unmarshal(msg.Body, &message)
	if err != nil {
		fmt.Println(err)
		return err
	}
	model.Db_write.Table("followed").Create(&message.Followed)
	model.Db_write.Table("user").Where("id = ?", message.Followed.TargetId).UpdateColumn("follow_count", message.FollowCount) //这里的targetId就是之前的userId
	model.Db_write.Table("user").Where("id = ?", message.Followed.UserId).UpdateColumn("follower_count", message.FollowerCount)
	return nil
}

type CancerConsumer struct{}

func (*CancerConsumer) HandleMessage(msg *nsq.Message) error {
	fmt.Println("receive", msg.NSQDAddress, "message:", string(msg.Body))
	message := Message{}
	err := json.Unmarshal(msg.Body, &message)
	if err != nil {
		fmt.Println(err)
		return err
	}
	fmt.Println("解除关注")
	model.Db_write.Table("followed").Where("user_id = ? AND target_id = ?", message.Followed.UserId, message.Followed.TargetId).Delete(&message.Followed)
	model.Db_write.Table("user").Where("id = ?", message.Followed.TargetId).UpdateColumn("follow_count", message.FollowCount)
	model.Db_write.Table("user").Where("id = ?", message.Followed.UserId).UpdateColumn("follower_count", message.FollowerCount)
	return nil
}

//初始化消费者
func InitConsumer(topic string, channel string) {
	cfg := nsq.NewConfig()
	cfg.LookupdPollInterval = time.Second //设置重连时间
	var err error
	Consumer, err = nsq.NewConsumer(topic, channel, cfg) // 新建一个消费者
	if err != nil {
		panic(err)
	}
	Consumer.SetLogger(nil, 0)        //屏蔽系统日志
	Consumer.AddHandler(&ConsumerT{}) // 添加消费者接口
	//建立NSQLookupd连接
	if err := Consumer.ConnectToNSQLookupd(address); err != nil {
		panic(err)
	}

	//建立多个nsqd连接
	// if err := c.ConnectToNSQDs([]string{"127.0.0.1:4150", "127.0.0.1:4152"}); err != nil {
	//  panic(err)
	// }

	// 建立一个nsqd连接
	// if err := c.ConnectToNSQD("127.0.0.1:4150"); err != nil {
	//  panic(err)
	// }
}

// 对应于topic为取消的Consumer
func InitCancerConsumer(topic string, channel string) {
	cfg := nsq.NewConfig()
	cfg.LookupdPollInterval = time.Second //设置重连时间
	var err error
	Consumer, err = nsq.NewConsumer(topic, channel, cfg) // 新建一个消费者
	if err != nil {
		panic(err)
	}
	Consumer.SetLogger(nil, 0)             //屏蔽系统日志
	Consumer.AddHandler(&CancerConsumer{}) // 添加消费者接口
	//建立NSQLookupd连接
	if err := Consumer.ConnectToNSQLookupd(address); err != nil {
		panic(err)
	}
}

func Producer(message Message, topic string) error {
	// 定义nsq生产者
	var producer *nsq.Producer
	// 初始化生产者
	// producer, err := nsq.NewProducer("地址:端口", nsq.*Config )
	producer, err := nsq.NewProducer("127.0.0.1:4150", nsq.NewConfig())
	if err != nil {
		panic(err)
	}

	err = producer.Ping()
	if nil != err {
		// 关闭生产者
		producer.Stop()
		producer = nil
	}
	defer producer.Stop()
	data, err := json.Marshal(message)
	if producer != nil && len(data) != 0 { //不能发布空串，否则会导致error
		err = producer.Publish(topic, data) // 发布消息
		if err != nil {
			fmt.Printf("producer.Publish,err : %v", err)
		}
		fmt.Println(data)

	}
	return err
}
