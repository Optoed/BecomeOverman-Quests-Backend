package kafka

import (
	"BecomeOverMan/internal/models"
	"context"
	"encoding/json"
	"fmt"
	"os"
	"strings"
	"time"

	kafkago "github.com/segmentio/kafka-go"
)

const (
	defaultTopicUser  = "becomeoverman.user.events"
	defaultTopicQuest = "becomeoverman.quest.events"
)

type Producer struct {
	writer     *kafkago.Writer
	topicUser  string
	topicQuest string
}

type UserQuestPurchasedPayload struct {
	EventType   string `json:"event_type"`
	UserID      int    `json:"user_id"`
	Username    string `json:"username"`
	Email       string `json:"email"`
	XPPoints    int    `json:"xp_points"`
	CoinBalance int    `json:"coin_balance"`
	Level       int    `json:"level"`
	QuestID     int    `json:"quest_id"`
	OccurredAt  string `json:"occurred_at"`
}

type QuestCreatedPayload struct {
	EventType   string `json:"event_type"`
	QuestID     int    `json:"quest_id"`
	Title       string `json:"title"`
	Description string `json:"description"`
	Category    string `json:"category"`
	Rarity      string `json:"rarity"`
	Difficulty  int    `json:"difficulty"`
	Price       int    `json:"price"`
	OccurredAt  string `json:"occurred_at"`
}

func NewProducerFromEnv() *Producer {
	brokers := strings.TrimSpace(os.Getenv("KAFKA_BROKERS"))
	if brokers == "" {
		return nil
	}

	parts := strings.Split(brokers, ",")
	for i := range parts {
		parts[i] = strings.TrimSpace(parts[i])
	}

	topicUser := os.Getenv("KAFKA_TOPIC_USER")
	if topicUser == "" {
		topicUser = defaultTopicUser
	}

	topicQuest := os.Getenv("KAFKA_TOPIC_QUEST")
	if topicQuest == "" {
		topicQuest = defaultTopicQuest
	}

	w := &kafkago.Writer{
		Addr:                   kafkago.TCP(parts...),
		Balancer:               &kafkago.LeastBytes{},
		AllowAutoTopicCreation: true,
		RequiredAcks:           kafkago.RequireOne,
		BatchTimeout:           50 * time.Millisecond,
	}

	return &Producer{
		writer:     w,
		topicUser:  topicUser,
		topicQuest: topicQuest,
	}
}

func (p *Producer) Close() error {
	if p == nil || p.writer == nil {
		return nil
	}
	return p.writer.Close()
}

func (p *Producer) PublishUserQuestPurchased(ctx context.Context, profile models.UserProfile, questID int) error {
	if p == nil || p.writer == nil {
		return nil
	}

	payload := UserQuestPurchasedPayload{
		EventType:   "user.quest_purchased",
		UserID:      profile.ID,
		Username:    profile.Username,
		Email:       profile.Email,
		XPPoints:    profile.XpPoints,
		CoinBalance: profile.CoinBalance,
		Level:       profile.Level,
		QuestID:     questID,
		OccurredAt:  time.Now().UTC().Format(time.RFC3339),
	}

	body, err := json.Marshal(payload)
	if err != nil {
		return err
	}

	msg := kafkago.Message{
		Topic: p.topicUser,
		Key:   fmt.Appendf(nil, "%d", profile.ID),
		Value: body,
		Time:  time.Now(),
	}

	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	return p.writer.WriteMessages(ctx, msg)
}

func (p *Producer) PublishQuestCreated(ctx context.Context, questID int, q *models.Quest) error {
	if p == nil || p.writer == nil || q == nil {
		return nil
	}

	payload := QuestCreatedPayload{
		EventType:   "quest.created",
		QuestID:     questID,
		Title:       q.Title,
		Description: q.Description,
		Category:    q.Category,
		Rarity:      q.Rarity,
		Difficulty:  q.Difficulty,
		Price:       q.Price,
		OccurredAt:  time.Now().UTC().Format(time.RFC3339),
	}

	body, err := json.Marshal(payload)
	if err != nil {
		return err
	}

	msg := kafkago.Message{
		Topic: p.topicQuest,
		Key:   fmt.Appendf(nil, "%d", questID),
		Value: body,
		Time:  time.Now(),
	}

	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	return p.writer.WriteMessages(ctx, msg)
}
