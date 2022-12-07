package mailer

import (
	"fmt"
	"log"

	mailjet "github.com/mailjet/mailjet-apiv3-go"
)

type MailJetConfig struct {
	APIKey    string
	SecretKey string
}

type MailJetClient struct {
	*MailJetConfig
}

func Init(apiKey, secretKey string) *MailJetClient {
	return &MailJetClient{
		&MailJetConfig{
			APIKey:    apiKey,
			SecretKey: secretKey,
		},
	}
}

func (mj *MailJetClient) SendNotifUnidentifiedFace(email, name, plate string) {
	mailjetClient := mailjet.NewMailjetClient(mj.MailJetConfig.APIKey, mj.MailJetConfig.SecretKey)
	messagesInfo := []mailjet.InfoMessagesV31{
		mailjet.InfoMessagesV31{
			From: &mailjet.RecipientV31{
				Email: "vision@mailjet.com",
				Name:  "PT. Vision Aman Sejahtera",
			},
			To: &mailjet.RecipientsV31{
				mailjet.RecipientV31{
					Email: email,
					Name:  name,
				},
			},
			TemplateID:       4413207,
			TemplateLanguage: true,
			Subject:          "Perhatian!",
			Variables: map[string]interface{}{
				"name":  name,
				"plate": plate,
			},
		},
	}
	messages := mailjet.MessagesV31{Info: messagesInfo}
	res, err := mailjetClient.SendMailV31(&messages)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Data: %+v\n", res)
}
