package domain

import (
	"context"
	"fmt"
	"terminal/domain/command"

	"github.com/aws/aws-sdk-go-v2/service/sqs"
	"github.com/aws/aws-sdk-go-v2/service/sqs/types"
)

type Queue struct {
	Conn       *sqs.Client
	Queue_name string
	queue_url  string
	terminal   *Terminal
}

func NewQueue(config *AWSConfig, terminal *Terminal, queue_name string) *Queue {
	queueParams := &sqs.GetQueueUrlInput{
		QueueName: &queue_name,
	}

	queueURL, errURL := config.SQS.GetQueueUrl(context.TODO(), queueParams)
	if errURL != nil {
		panic(errURL)
	}

	return &Queue{
		Conn:       config.SQS,
		Queue_name: queue_name,
		queue_url:  *queueURL.QueueUrl,
		terminal:   terminal,
	}
}

func (q *Queue) DeleteMessage(ctx context.Context, message types.Message) {
	deleteParams := &sqs.DeleteMessageInput{
		QueueUrl:      &q.queue_url,
		ReceiptHandle: message.ReceiptHandle,
	}

	_, errDelete := q.Conn.DeleteMessage(context.TODO(), deleteParams)
	if errDelete != nil {
		fmt.Printf("erro ao tentar deletar uma mensagem: %s\n", errDelete)
	}
}

func (q *Queue) ReadMessage(ctx context.Context, commands *command.Command) {
	queueParam := &sqs.GetQueueUrlInput{
		QueueName: &q.Queue_name,
	}

	sqs_url, errSqs := q.Conn.GetQueueUrl(ctx, queueParam)
	if errSqs != nil {
		panic(fmt.Sprintf("Erro ao estabelecer conexão, %s\n", errSqs))
	}

	for {
		results, resultErr := q.Conn.ReceiveMessage(
			ctx,
			&sqs.ReceiveMessageInput{
				QueueUrl:              sqs_url.QueueUrl,
				MaxNumberOfMessages:   10,
				MessageAttributeNames: []string{"All"},
			},
		)

		if resultErr != nil {
			fmt.Println("Erro ao ler mensagem: ", resultErr.Error())
			continue
		}

		for _, message := range results.Messages {
			messageType, messageTypeOK := message.MessageAttributes["Type"]
			if !messageTypeOK {
				fmt.Println("Nenhum tipo de comando informado")
				q.DeleteMessage(ctx, message)
				continue
			}
			receivedCommand := command.CommandType(*messageType.StringValue)
			if !receivedCommand.IsValid() {
				fmt.Printf("O comando %s não é válido \n", receivedCommand)
				q.DeleteMessage(ctx, message)
				continue
			}

			ipAttr, ipOK := message.MessageAttributes["IP"]
			if !ipOK {
				fmt.Println("Nenhum ip informado")
				q.DeleteMessage(ctx, message)
				continue
			}
			ip := *ipAttr.StringValue
			if ip != q.terminal.Id {
				fmt.Println("O comando não percente a este terminal")
				continue
			}

			if receivedCommand == command.COMMAND_DISABLE_KEYS {
				fmt.Printf("O commando %s foi executado no terminal %s\n", receivedCommand, q.terminal.Slug())
				Terminal.DisableKeys(*q.terminal)
			}

			if receivedCommand == command.COMMAND_ENABLED_KEYS {
				fmt.Printf("O commando %s foi executado no terminal %s\n", receivedCommand, q.terminal.Slug())
				Terminal.EnableKeys(*q.terminal)
			}

			q.DeleteMessage(ctx, message)
		}
	}
}
