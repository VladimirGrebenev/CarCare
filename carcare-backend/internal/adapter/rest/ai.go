package rest

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"strings"
)

// ChatRequest — запрос от фронтенда
type ChatRequest struct {
	Message string `json:"message"`
}

// ChatResponse — ответ фронтенду
type ChatResponse struct {
	Reply string `json:"reply"`
}

// YandexGPT request/response structures
type yandexGPTMessage struct {
	Role string `json:"role"`
	Text string `json:"text"`
}

type yandexGPTRequest struct {
	ModelURI string `json:"modelUri"`
	CompletionOptions struct {
		Stream      bool    `json:"stream"`
		Temperature float64 `json:"temperature"`
		MaxTokens   int     `json:"maxTokens"`
	} `json:"completionOptions"`
	Messages []yandexGPTMessage `json:"messages"`
}

type yandexGPTResponse struct {
	Result struct {
		Alternatives []struct {
			Message struct {
				Text string `json:"text"`
			} `json:"message"`
		} `json:"alternatives"`
	} `json:"result"`
}

// systemPrompt — инструкция для AI на основе документации
const systemPrompt = `Ты — дружелюбный AI-помощник приложения CarCare. 
Приложение предназначено для учёта расходов на автомобиль: топливо, техобслуживание, штрафы.

Твои задачи:
1. Отвечать на вопросы пользователей о том, как пользоваться приложением.
2. Помогать создавать сущности (автомобили, заправки, ТО, штрафы) по текстовому описанию.
3. Если данных не хватает — задавать уточняющие вопросы.

Правила:
- Отвечай только по вопросам, связанным с CarCare.
- Если вопрос не по теме — вежливо скажи, что можешь помочь только по приложению.
- Если не знаешь ответа — честно признайся.
- Отвечай на том же языке, на котором задан вопрос.
- Будь дружелюбным и полезным.

Информация о приложении:
- Можно добавить несколько автомобилей (марка, модель, год, VIN, госномер)
- Учёт топлива: дата, объём, тип топлива, цена, пробег, АЗС
- Учёт ТО: дата, тип работ, пробег, стоимость, описание, сервис
- Типы ТО: Замена масла, Замена шин, Техосмотр, Ремонт, Страховка, Другое
- Штрафы: дата, номер постановления, сумма, статья, статус оплаты
- Отчёты и статистика по всем расходам
- Регистрация по email или через Яндекс/Google`

// ChatHandler — обрабатывает запросы к AI
type ChatHandler struct {
	yandexGPTFolderID string
	apiKey            string
}

func NewChatHandler() *ChatHandler {
	return &ChatHandler{
		yandexGPTFolderID: os.Getenv("YC_FOLDER_ID"),
		apiKey:            os.Getenv("YC_API_KEY"),
	}
}

func (h *ChatHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	if r.Method != http.MethodPost {
		http.Error(w, `{"error":"method not allowed"}`, http.StatusMethodNotAllowed)
		return
	}

	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		http.Error(w, `{"error":"invalid body"}`, http.StatusBadRequest)
		return
	}

	var req ChatRequest
	if err := json.Unmarshal(body, &req); err != nil {
		http.Error(w, `{"error":"invalid json"}`, http.StatusBadRequest)
		return
	}

	if strings.TrimSpace(req.Message) == "" {
		http.Error(w, `{"error":"message is required"}`, http.StatusBadRequest)
		return
	}

	reply, err := h.callYandexGPT(req.Message)
	if err != nil {
		// Если YandexGPT не настроен, используем заглушку
		reply = h.fallbackResponse(req.Message)
	}

	json.NewEncoder(w).Encode(ChatResponse{Reply: reply})
}

func (h *ChatHandler) callYandexGPT(userMessage string) (string, error) {
	if h.apiKey == "" || h.yandexGPTFolderID == "" {
		return "", fmt.Errorf("YandexGPT not configured")
	}

	modelURI := fmt.Sprintf("gpt://%s/yandexgpt/latest", h.yandexGPTFolderID)

	reqBody := yandexGPTRequest{}
	reqBody.ModelURI = modelURI
	reqBody.CompletionOptions.Stream = false
	reqBody.CompletionOptions.Temperature = 0.6
	reqBody.CompletionOptions.MaxTokens = 2000
	reqBody.Messages = []yandexGPTMessage{
		{Role: "system", Text: systemPrompt},
		{Role: "user", Text: userMessage},
	}

	jsonBody, _ := json.Marshal(reqBody)

	httpReq, _ := http.NewRequest("POST", "https://llm.api.cloud.yandex.net/foundationModels/v1/completion", bytes.NewReader(jsonBody))
	httpReq.Header.Set("Content-Type", "application/json")
	httpReq.Header.Set("Authorization", "Bearer "+h.apiKey)

	resp, err := http.DefaultClient.Do(httpReq)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	respBody, _ := ioutil.ReadAll(resp.Body)
	var yandexResp yandexGPTResponse
	if err := json.Unmarshal(respBody, &yandexResp); err != nil {
		return "", err
	}

	if len(yandexResp.Result.Alternatives) > 0 {
		return yandexResp.Result.Alternatives[0].Message.Text, nil
	}

	return "", fmt.Errorf("no alternatives in response")
}

// fallbackResponse — простая заглушка, если YandexGPT не настроен
func (h *ChatHandler) fallbackResponse(message string) string {
	msg := strings.ToLower(message)

	switch {
	case strings.Contains(msg, "автомобил") || strings.Contains(msg, "машин") || strings.Contains(msg, "добавить авто"):
		return "Чтобы добавить автомобиль, перейдите в раздел \"Мои автомобили\" и нажмите \"Добавить автомобиль\". Укажите марку, модель, год выпуска и при необходимости VIN и госномер. Если хотите, я могу помочь — просто напишите данные автомобиля текстом, например: \"Toyota Camry 2020\"."

	case strings.Contains(msg, "заправк") || strings.Contains(msg, "топлив") || strings.Contains(msg, "бензин"):
		return "Чтобы добавить заправку, перейдите в раздел \"Топливо\" и нажмите \"Добавить заправку\". Укажите дату, объём в литрах, тип топлива, цену и пробег. Если хотите, я могу помочь — просто напишите данные, например: \"Залил 45 литров АИ-95 по 55 рублей\"."

	case strings.Contains(msg, "то ") || strings.Contains(msg, "техобслуж") || strings.Contains(msg, "ремонт"):
		return "Чтобы добавить запись о техобслуживании, перейдите в раздел \"Техобслуживание\". Укажите дату, тип работ, пробег и стоимость. Типы работ: замена масла, замена шин, техосмотр, ремонт, страховка или другое."

	case strings.Contains(msg, "штраф"):
		return "Чтобы добавить штраф, перейдите в раздел \"Штрафы\" и нажмите \"Добавить штраф\". Укажите дату, номер постановления, сумму и статью нарушения. Статус по умолчанию — \"не оплачен\"."

	case strings.Contains(msg, "отчёт") || strings.Contains(msg, "статистик") || strings.Contains(msg, "расход"):
		return "В разделе \"Отчёты\" вы можете увидеть сводку по всем расходам: топливо, ТО и штрафы. Доступны фильтры по периоду (день, неделя, месяц, год) и графики распределения расходов."

	case strings.Contains(msg, "привет") || strings.Contains(msg, "здравствуй") || strings.Contains(msg, "здаров"):
		return "Привет! Я AI-помощник CarCare 👋 Я могу помочь с вопросами по приложению, а также создать автомобиль, заправку, ТО или штраф по вашему описанию. Что вас интересует?"

	default:
		return "Я могу помочь вам с вопросами по приложению CarCare! Расскажите, что вас интересует:\n\n• Как добавить автомобиль?\n• Как записать заправку?\n• Как добавить техобслуживание?\n• Как зарегистрировать штраф?\n• Где посмотреть отчёты?\n\nИли просто опишите, что хотите создать, и я помогу собрать все данные! 😊"
	}
}