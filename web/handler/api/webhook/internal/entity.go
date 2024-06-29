package internal

import (
	validation "github.com/go-ozzo/ozzo-validation"
)

type WebhookJson struct {
	NewWebhooks    []NewWebhook    `json:"newWebhooks"`
	UpdateWebhooks []UpdateWebhook `json:"updateWebhooks"`
}

type NewWebhook struct {
	WebhookID        string   `json:"webhookId"`
	SubscriptionType string   `json:"subscriptionType"`
	SubscriptionId   string   `json:"subscriptionId"`
	MentionRoles     []string `json:"mentionRoles"`
	MentionUsers     []string `json:"mentionUsers"`
	NgOrWords        []string `json:"ngOrWords"`
	NgAndWords       []string `json:"ngAndWords"`
	SearchOrWords    []string `json:"searchOrWords"`
	SearchAndWords   []string `json:"searchAndWords"`
	MentionOrWords   []string `json:"mentionOrWords"`
	MentionAndWords  []string `json:"mentionAndWords"`
}

type UpdateWebhook struct {
	WebhookSerialID  int64    `json:"webhookSerialId"`
	WebhookID        string   `json:"webhookId"`
	SubscriptionType string   `json:"subscriptionType"`
	SubscriptionId   string   `json:"subscriptionId"`
	MentionRoles     []string `json:"mentionRoles"`
	MentionUsers     []string `json:"mentionUsers"`
	NgOrWords        []string `json:"ngOrWords"`
	NgAndWords       []string `json:"ngAndWords"`
	SearchOrWords    []string `json:"searchOrWords"`
	SearchAndWords   []string `json:"searchAndWords"`
	MentionOrWords   []string `json:"mentionOrWords"`
	MentionAndWords  []string `json:"mentionAndWords"`
	DeleteFlag       bool     `json:"deleteFlag"`
}

func (g WebhookJson) Validate() error {
	return validation.ValidateStruct(&g,
		validation.Field(&g.NewWebhooks,
			validation.Length(0, 1000),
			validation.Each(validation.By(func(value interface{}) error {
				return validation.ValidateStruct(value,
					validation.Field(value.(NewWebhook).WebhookID, validation.Required),
					validation.Field(value.(NewWebhook).SubscriptionType, validation.In("youtube", "niconico")),
				)
			})),
		),
		validation.Field(&g.UpdateWebhooks,
			validation.Length(0, 1000),
			validation.Each(validation.By(func(value interface{}) error {
				return validation.ValidateStruct(value,
					validation.Field(value.(NewWebhook).WebhookID, validation.Required),
					validation.Field(value.(NewWebhook).SubscriptionType, validation.In("youtube", "niconico")),
				)
			})),
		),
	)
}
