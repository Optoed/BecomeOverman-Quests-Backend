package models

// Search
type RecommendationService_SearchQuest_Request struct {
	Query    string `json:"query" binding:"required"`
	TopK     int    `json:"top_k,omitempty" binding:"omitempty,numeric,min=1,max=100"`
	Category string `json:"category,omitempty"`
}

type RecommendationService_SearchQuest_Result struct {
	ID              int     `json:"id"`
	Title           string  `json:"title"`              // TODO: потом убрать тк ненужная трата трафика
	Description     string  `json:"description"`        // TODO: потом убрать тк ненужная трата трафика
	Category        string  `json:"category,omitempty"` // TODO: потом убрать тк ненужная трата трафика
	SimilarityScore float64 `json:"similarity_score"`
}

type RecommendationService_SearchQuests_Response struct {
	Results []RecommendationService_SearchQuest_Result `json:"results"`
}

// quest/add

type RecommendationService_questToAdd struct {
	ID          string `json:"id"`
	Title       string `json:"title"`
	Description string `json:"description"`
	Category    string `json:"category,omitempty"`
}

type RecommendationService_AddQuests_Request struct {
	Quests []RecommendationService_questToAdd `json:"quests"`
}
