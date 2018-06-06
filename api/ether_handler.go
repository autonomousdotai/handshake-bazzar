package api

import (
	"context"
	"encoding/json"
	"errors"
	"log"
	"strconv"

	"cloud.google.com/go/pubsub"
	"github.com/ninjadotorg/handshake-bazzar/utils"
)

type EtherHandler struct {
	BubsubClient       *pubsub.Client
	BubsubSubscription *pubsub.Subscription
}

func NewEthHandler(pubsubClient *pubsub.Client, topicName, subscriptionName string) (*EtherHandler, error) {
	handler := EtherHandler{}

	handler.BubsubClient = pubsubClient

	topic := pubsubClient.Topic(topicName)
	if topic == nil || topic.ID() != topicName {
		var err error
		topic, err = pubsubClient.CreateTopic(context.Background(), topicName)
		if err != nil {
			log.Println("NewEthHandler", err)
			return nil, err
		}
	}

	sub := pubsubClient.Subscription(subscriptionName)
	existed, err := sub.Exists(context.Background())
	if err != nil {
		log.Println("NewEthHandler", err)
		return nil, err
	}
	if !existed {
		var err error
		sub, err = pubsubClient.CreateSubscription(context.Background(), subscriptionName, pubsub.SubscriptionConfig{Topic: topic})
		if err != nil {
			log.Println("NewEthHandler", err)
			return nil, err
		}
	}
	err = sub.Receive(context.Background(), func(ctx context.Context, m *pubsub.Message) {
		log.Printf("Got message : %s", m.Data)
		m.Ack()
		handler.Process(m.Data)
	})
	if err != nil {
		log.Println("NewEthHandler", err)
		return nil, err
	}

	return &handler, nil
}

func (etherHandler *EtherHandler) Process(bytes []byte) error {
	logData := map[string]interface{}{}
	err := json.Unmarshal(bytes, &logData)
	if err != nil {
		log.Println("NewEthHandler.Process()", err)
		return err
	}
	event := logData["event"].(string)
	fromAddress := logData["from_address"].(string)
	data, ok := logData["data"].(map[string]interface{})
	if !ok {
		return errors.New("data is missed")
	}
	switch event {
	case "__init":
		{
			hid := int64(-1)
			val, ok := data["hid"].(float64)
			if !ok {
				return errors.New("hid is invalid")
			}
			hid = int64(val)
			if hid < 0 {
				return errors.New("hid is invalid")
			}
			offchain, ok := data["offchain"].(string)
			offchainType, offchainIdStr, err := utils.ParseOffchain(offchain)
			if err != nil {
				log.Println("NewEthHandler.Process()", err)
				return err
			}
			if offchainType == utils.OFFCHAIN_BAZZAR {
				productShakeId, err := strconv.ParseInt(offchainIdStr, 10, 64)
				if err != nil {
					log.Println("NewEthHandler.Process()", err)
					return err
				}
				err = bazzarService.ProcessEventInit(hid, productShakeId)
				if err != nil {
					log.Println("NewEthHandler.Process()", err)
					return err
				}
			}

			return nil
		}
		break
	case "__shake":
		{
			hid := int64(-1)
			val, ok := data["hid"].(float64)
			if !ok {
				return errors.New("hid is invalid")
			}
			hid = int64(val)
			if hid < 0 {
				return errors.New("hid is invalid")
			}
			offchain, ok := data["offchain"].(string)
			offchainType, offchainIdStr, err := utils.ParseOffchain(offchain)
			if err != nil {
				log.Println("NewEthHandler.Process()", err)
				return err
			}
			if offchainType == utils.OFFCHAIN_BAZZAR {
				productShakeId, err := strconv.ParseInt(offchainIdStr, 10, 64)
				if err != nil {
					log.Println("NewEthHandler.Process()", err)
					return err
				}
				err = bazzarService.ProcessEventShake(hid, productShakeId, fromAddress)
				if err != nil {
					log.Println("NewEthHandler.Process()", err)
					return err
				}
			}

			return nil
		}
		break
	case "__deliver":
		{
			hid := int64(-1)
			val, ok := data["hid"].(float64)
			if !ok {
				return errors.New("hid is invalid")
			}
			hid = int64(val)
			if hid < 0 {
				return errors.New("hid is invalid")
			}
			offchain, ok := data["offchain"].(string)
			offchainType, offchainIdStr, err := utils.ParseOffchain(offchain)
			if err != nil {
				log.Println("NewEthHandler.Process()", err)
				return err
			}
			if offchainType == utils.OFFCHAIN_BAZZAR {
				productShakeId, err := strconv.ParseInt(offchainIdStr, 10, 64)
				if err != nil {
					log.Println("NewEthHandler.Process()", err)
					return err
				}
				err = bazzarService.ProcessEventDeliver(hid, productShakeId)
				if err != nil {
					log.Println("NewEthHandler.Process()", err)
					return err
				}
			}

			return nil
		}
		break
	case "__reject":
		{
			hid := int64(-1)
			val, ok := data["hid"].(float64)
			if !ok {
				return errors.New("hid is invalid")
			}
			hid = int64(val)
			if hid < 0 {
				return errors.New("hid is invalid")
			}
			offchain, ok := data["offchain"].(string)
			offchainType, offchainIdStr, err := utils.ParseOffchain(offchain)
			if err != nil {
				log.Println("NewEthHandler.Process()", err)
				return err
			}
			if offchainType == utils.OFFCHAIN_BAZZAR {
				productShakeId, err := strconv.ParseInt(offchainIdStr, 10, 64)
				if err != nil {
					log.Println("NewEthHandler.Process()", err)
					return err
				}
				err = bazzarService.ProcessEventReject(hid, productShakeId)
				if err != nil {
					log.Println("NewEthHandler.Process()", err)
					return err
				}
			}

			return nil
		}
		break
	case "__accept":
		{
			hid := int64(-1)
			val, ok := data["hid"].(float64)
			if !ok {
				return errors.New("hid is invalid")
			}
			hid = int64(val)
			if hid < 0 {
				return errors.New("hid is invalid")
			}
			offchain, ok := data["offchain"].(string)
			offchainType, offchainIdStr, err := utils.ParseOffchain(offchain)
			if err != nil {
				log.Println("NewEthHandler.Process()", err)
				return err
			}
			if offchainType == utils.OFFCHAIN_BAZZAR {
				productShakeId, err := strconv.ParseInt(offchainIdStr, 10, 64)
				if err != nil {
					log.Println("NewEthHandler.Process()", err)
					return err
				}
				err = bazzarService.ProcessEventAccept(hid, productShakeId)
				if err != nil {
					log.Println("NewEthHandler.Process()", err)
					return err
				}
			}

			return nil
		}
		break
	case "__cancel":
		{
			hid := int64(-1)
			val, ok := data["hid"].(float64)
			if !ok {
				return errors.New("hid is invalid")
			}
			hid = int64(val)
			if hid < 0 {
				return errors.New("hid is invalid")
			}
			offchain, ok := data["offchain"].(string)
			offchainType, offchainIdStr, err := utils.ParseOffchain(offchain)
			if err != nil {
				log.Println("NewEthHandler.Process()", err)
				return err
			}
			if offchainType == utils.OFFCHAIN_BAZZAR {
				productShakeId, err := strconv.ParseInt(offchainIdStr, 10, 64)
				if err != nil {
					log.Println("NewEthHandler.Process()", err)
					return err
				}
				err = bazzarService.ProcessEventCancel(hid, productShakeId)
				if err != nil {
					log.Println("NewEthHandler.Process()", err)
					return err
				}
			}

			return nil
		}
		break
	case "__withdraw":
		{
			hid := int64(-1)
			val, ok := data["hid"].(float64)
			if !ok {
				return errors.New("hid is invalid")
			}
			hid = int64(val)
			if hid < 0 {
				return errors.New("hid is invalid")
			}
			offchain, ok := data["offchain"].(string)
			offchainType, offchainIdStr, err := utils.ParseOffchain(offchain)
			if err != nil {
				log.Println("NewEthHandler.Process()", err)
				return err
			}
			if offchainType == utils.OFFCHAIN_BAZZAR {
				productShakeId, err := strconv.ParseInt(offchainIdStr, 10, 64)
				if err != nil {
					log.Println("NewEthHandler.Process()", err)
					return err
				}
				err = bazzarService.ProcessEventWithdraw(hid, productShakeId)
				if err != nil {
					log.Println("NewEthHandler.Process()", err)
					return err
				}
			}

			return nil
		}
		break
	}
	_ = fromAddress
	return nil
}
