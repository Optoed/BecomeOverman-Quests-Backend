package services

import (
	"BecomeOverMan/internal/models"
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"
	"time"
)

const (
	url = "https://api.intelligence.io.solutions/api/v1/chat/completions"
)

var apiKey = os.Getenv("API_KEY")

type ChatMessage struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

type ChatRequest struct {
	Model       string        `json:"model"`
	Messages    []ChatMessage `json:"messages"`
	Temperature float64       `json:"temperature,omitempty"`
}

type Choice struct {
	Message ChatMessage `json:"message"`
}

type ChatResponse struct {
	Choices []Choice `json:"choices"`
}

func requestAI(userMessage, systemPrompt, aiModel string) ([]byte, error) {
	if apiKey == "" {
		return nil, fmt.Errorf("API_KEY not found in environment variables")
	}

	if aiModel == "" {
		aiModel = "moonshotai/Kimi-K2-Thinking"
	}

	requestData := ChatRequest{
		Model: aiModel,
		Messages: []ChatMessage{
			{
				Role:    "system",
				Content: systemPrompt,
			},
			{
				Role:    "user",
				Content: userMessage,
			},
		},
		Temperature: 0.7,
	}

	jsonData, err := json.Marshal(requestData)
	if err != nil {
		return nil, fmt.Errorf("error marshaling JSON: %v", err)
	}

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonData))
	if err != nil {
		return nil, fmt.Errorf("error creating request: %v", err)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+apiKey)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("error sending request: %v", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("error reading response: %v", err)
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("API returned error: %s", string(body))
	}

	var chatResponse ChatResponse
	err = json.Unmarshal(body, &chatResponse)
	if err != nil {
		return nil, fmt.Errorf("error parsing chat response: %v", err)
	}

	if len(chatResponse.Choices) == 0 {
		return nil, fmt.Errorf("no choices in response")
	}

	// Очищаем ответ от thinking тегов
	content := chatResponse.Choices[0].Message.Content
	if idx := strings.Index(content, "</think>\n\n"); idx != -1 {
		content = content[idx+11:] // +11 чтобы пропустить "</think>\n\n"
	}

	return []byte(content), nil
}

// GenerateAIQuest - Генерация квеста по запросу пользователя в LLM
func (s *QuestService) GenerateAIQuest(userMessage string) (*models.AIQuestResponse, error) {
	aiModel := "moonshotai/Kimi-K2-Thinking"

	systemPrompt := `
	Ты помощник для генерации квестов в формате строгого JSON. 
	ВОЗВРАЩАЙ ТОЛЬКО JSON БЕЗ ЛЮБЫХ ДОПОЛНИТЕЛЬНЫХ ТЕКСТОВ И КОММЕНТАРИЕВ!

	Структура JSON должна быть такой:
	{
		"quest": {
			"title": "Название квеста",
			"description": "Описание квеста [GENERATED]",
			"category": "health/willpower/intelligence/creativity/social",
			"rarity": "common/rare/epic/legendary",
			"difficulty": 1-5,
			"price": 10-100,
			"tasks_count": 3-7,
			"reward_xp": 50-500,
			"reward_coin": 25-250,
			"time_limit_hours": 24-336
		},
		"tasks": [
			{
				"title": "Название задачи 1",
				"description": "Описание задачи 1",
				"difficulty": 1-3,
				"rarity": "common/rare/epic",
				"category": "health/willpower/intelligence/creativity/social",
				"base_xp_reward": 10-50,
				"base_coin_reward": 5-25,
				"task_order": 1
			}
		]
	}

	Правила:
	- difficulty квеста должен быть средним от difficulty задач
	- price = reward_coin * 1.5 (округлить)
	- tasks_count должно соответствовать количеству задач в массиве
	- time_limit_hours: 24-168 (1-7 дней)
	- reward_xp = сумма base_xp_reward всех задач * 1.5
	- reward_coin = сумма base_coin_reward всех задач * 1.5
	`

	answer, err := requestAI(userMessage, systemPrompt, aiModel)

	// Парсим финальный JSON
	var aiResponse models.AIQuestResponse
	err = json.Unmarshal(answer, &aiResponse)
	if err != nil {
		return nil, fmt.Errorf("error parsing AI quest response: %v", err)
	}

	return &aiResponse, nil
}

// -----------------------------------------------------------
// ----------- GenerateScheduleByAI --------------------------
// -----------------------------------------------------------

type ScheduleTask struct {
	TaskID         int        `json:"task_id"`
	ScheduledStart time.Time  `json:"scheduled_start"`
	ScheduledEnd   time.Time  `json:"scheduled_end"`
	Deadline       *time.Time `json:"deadline"`
	Duration       int        `json:"duration"` // в минутах
}

type AIScheduleResponse struct {
	Schedule []ScheduleTask `json:"schedule"`
}

func (s *QuestService) GenerateScheduleByAI(
	ctx context.Context,
	userID int,
	userMessage string,
) ([]models.Quest, error) {
	quests, err := s.GetMyAllQuestsWithDetails(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("error getting user quests: %v", err)
	}

	info := ""
	for _, q := range quests {
		questStr := fmt.Sprintf("QuestID:%d Title:%s Description:%s ",
			q.ID,
			q.Title,
			q.Description,
		)

		tasksStr := "Tasks:["
		for _, t := range q.Tasks {
			if t.Status != nil && *t.Status != "active" {
				continue
			}

			taskStr := fmt.Sprintf(
				"ID:%d Title:%s Desc:%s Deadline:%v Duration:%v ScheduledStart:%v ScheduledEnd:%v",
				t.ID,
				t.Title,
				t.Description,
				t.Deadline,
				t.Duration,
				t.ScheduledStart,
				t.ScheduledEnd,
			)

			if t.TaskOrder >= 0 {
				taskStr = taskStr + fmt.Sprintf(" TaskOrder:%d; ", t.TaskOrder)
			} else {
				taskStr = taskStr + "; "
			}

			tasksStr += taskStr
		}
		tasksStr += "] "
		questStr += tasksStr

		info += questStr
	}

	userMessageWithInfo := userMessage + "\n\n" + info

	systemPrompt := `
	Ты помощник для генерации расписания задач.
	Верни ТОЛЬКО валидный JSON без комментариев.

	Формат ответа:
	{
	"schedule": [
		{
		"task_id": number,
		"scheduled_start": "RFC3339",
		"scheduled_end": "RFC3339",
		"deadline": "RFC3339|null",
		"duration": number
		}
	]
	}

	Правила:
	- duration только в минутах
	- если у задачи уже есть часть данных — дополни логично
	- распределяй задачи равномерно
	- учитывай taskOrder - задача с меньшим taskOrder должна быть раньше выполнена
	`

	aiModel := "moonshotai/Kimi-K2-Thinking"

	answer, err := requestAI(userMessageWithInfo, systemPrompt, aiModel)
	if err != nil {
		return nil, err
	}

	answer = bytes.TrimSpace(answer)
	if len(answer) == 0 || answer[0] != '{' {
		return nil, fmt.Errorf("AI returned invalid JSON: %s", string(answer))
	}

	// ---------- парсим ----------
	var schedules AIScheduleResponse
	if err := json.Unmarshal(answer, &schedules); err != nil {
		return nil, fmt.Errorf("error parsing AI schedule response: %v\nraw: %s", err, string(answer))
	}

	// ---------- индексируем ----------
	scheduleMap := make(map[int]ScheduleTask)
	for _, sch := range schedules.Schedule {
		scheduleMap[sch.TaskID] = sch
	}

	// ---------- применяем к задачам ----------
	for qi := range quests {
		for ti := range quests[qi].Tasks {
			t := &quests[qi].Tasks[ti]

			sch, ok := scheduleMap[t.ID]
			if !ok {
				continue
			}

			t.ScheduledStart = &sch.ScheduledStart
			t.ScheduledEnd = &sch.ScheduledEnd
			t.Deadline = sch.Deadline
			t.Duration = &sch.Duration
		}
	}

	// ------ Сохраняем в БД -----------
	allTasks := []models.Task{}
	for _, q := range quests {
		allTasks = append(allTasks, q.Tasks...)
	}

	if err := s.questRepo.SetOrUpdateScheduleTasks(ctx, userID, allTasks); err != nil {
		return nil, fmt.Errorf("error updating schedule tasks: %v", err)
	}

	return quests, nil
}
