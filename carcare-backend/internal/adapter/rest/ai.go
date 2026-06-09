package rest

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"strings"

	"github.com/google/uuid"

	"github.com/VladimirGrebenev/CarCare-backend/internal/domain/car"
	"github.com/VladimirGrebenev/CarCare-backend/internal/domain/fine"
	"github.com/VladimirGrebenev/CarCare-backend/internal/domain/fuel"
	"github.com/VladimirGrebenev/CarCare-backend/internal/domain/maintenance"
	"github.com/VladimirGrebenev/CarCare-backend/internal/usecase"
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
	ModelURI          string `json:"modelUri"`
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
const systemPrompt = `Ты — AI-помощник приложения CarCare. Приложение для учёта расходов на автомобиль.

Твои задачи:
1. Отвечать на вопросы пользователей о том, как пользоваться приложением.
2. Создавать сущности (автомобили, заправки, ТО, штрафы) по текстовому описанию пользователя.

Правила:
- Отвечай на том же языке, на котором задан вопрос.
- Если данных не хватает — задавай уточняющие вопросы по одному за раз.
- Будь дружелюбным и полезным.

ВАЖНО: Если пользователь просит создать сущность, ты ДОЛЖЕН вернуть JSON-блок с действием в конце своего ответа.
Формат JSON-блока:
\`\`\`json
{"action": "create_car", "data": {"brand": "...", "model": "...", "year": 2024, "vin": "...", "plate": "..."}}
\`\`\`

Доступные действия и их поля:

1. create_car — создать автомобиль
   Поля: brand (обяз), model (обяз), year (обяз, число), vin (опц), plate (опц)

2. create_fuel — создать заправку
   Поля: car_id (обяз), volume (обяз, число), price (обяз, число), type (обяз: АИ-92/АИ-95/АИ-98/Дизель/Газ), date (обяз, формат ГГГГ-ММ-ДД)

3. create_maintenance — создать ТО
   Поля: car_id (обяз), type (обяз: Замена масла/Замена шин/Техосмотр/Ремонт/Страховка/Другое), date (обяз, формат ГГГГ-ММ-ДД), cost (обяз, число), description (опц)

4. create_fine — создать штраф
   Поля: car_id (обяз), amount (обяз, число), type (обяз), date (обяз, формат ГГГГ-ММ-ДД), status (опц: paid/unpaid), description (опц)

Пример ответа с созданием авто:
"Отлично! Добавляю автомобиль Toyota Camry 2020.
\`\`\`json
{"action": "create_car", "data": {"brand": "Toyota", "model": "Camry", "year": 2020}}
\`\`\`"

Если пользователь не указал автомобиль (car_id), спроси какой автомобиль. Список автомобилей пользователя ты не видишь, поэтому спроси название.

Если пользователь не указал обязательные поля — задай уточняющий вопрос.

Информация о приложении:
- Типы топлива: АИ-92, АИ-95, АИ-98, Дизель, Газ
- Типы ТО: Замена масла, Замена шин, Техосмотр, Ремонт, Страховка, Другое
- Статус штрафа: unpaid (не оплачен), paid (оплачен)`

// AIAction — структура действия из JSON
type AIAction struct {
	Action string                 `json:"action"`
	Data   map[string]interface{} `json:"data"`
}

// ChatHandler — обрабатывает запросы к AI
type ChatHandler struct {
	yandexGPTFolderID string
	apiKey            string
	uc                *usecase.UsecaseContainer
}

func NewChatHandler(uc *usecase.UsecaseContainer) *ChatHandler {
	return &ChatHandler{
		yandexGPTFolderID: os.Getenv("YC_FOLDER_ID"),
		apiKey:            os.Getenv("YC_API_KEY"),
		uc:                uc,
	}
}

func (h *ChatHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	if r.Method != http.MethodPost {
		http.Error(w, `{"error":"method not allowed"}`, http.StatusMethodNotAllowed)
		return
	}

	userID := getUserIDFromContext(r)
	if userID == "" {
		http.Error(w, `{"error":"unauthorized"}`, http.StatusUnauthorized)
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
		reply = h.fallbackResponse(req.Message)
	}

	// Пробуем выполнить действие из ответа AI
	result := h.tryExecuteAction(reply, userID)
	if result != "" {
		reply = result
	}

	json.NewEncoder(w).Encode(ChatResponse{Reply: reply})
}

// tryExecuteAction ищет JSON-блок с действием в ответе AI и выполняет его
func (h *ChatHandler) tryExecuteAction(response string, userID string) string {
	// Ищем JSON-блок в ответе
	start := strings.Index(response, "```json\n")
	if start == -1 {
		return "" // нет действия
	}
	start += len("```json\n")
	end := strings.Index(response[start:], "\n```")
	if end == -1 {
		return ""
	}
	jsonStr := response[start : start+end]

	var action AIAction
	if err := json.Unmarshal([]byte(jsonStr), &action); err != nil {
		return ""
	}

	// Выполняем действие
	switch action.Action {
	case "create_car":
		return h.executeCreateCar(action.Data, userID)
	case "create_fuel":
		return h.executeCreateFuel(action.Data, userID)
	case "create_maintenance":
		return h.executeCreateMaintenance(action.Data, userID)
	case "create_fine":
		return h.executeCreateFine(action.Data, userID)
	default:
		return ""
	}
}

func getString(data map[string]interface{}, key string) string {
	if v, ok := data[key]; ok {
		if s, ok := v.(string); ok {
			return s
		}
	}
	return ""
}

func getFloat(data map[string]interface{}, key string) float64 {
	if v, ok := data[key]; ok {
		switch val := v.(type) {
		case float64:
			return val
		case int:
			return float64(val)
		}
	}
	return 0
}

func getInt(data map[string]interface{}, key string) int {
	if v, ok := data[key]; ok {
		switch val := v.(type) {
		case float64:
			return int(val)
		case int:
			return val
		}
	}
	return 0
}

func (h *ChatHandler) executeCreateCar(data map[string]interface{}, userID string) string {
	c := car.Car{
		ID:     uuid.New().String(),
		UserID: userID,
		Brand:  getString(data, "brand"),
		Model:  getString(data, "model"),
		Year:   getInt(data, "year"),
		VIN:    getString(data, "vin"),
		Plate:  getString(data, "plate"),
	}

	if c.Brand == "" || c.Model == "" {
		return "❌ Не указаны марка или модель автомобиля. Пожалуйста, укажите их."
	}

	if err := h.uc.Car.AddCar(c); err != nil {
		return fmt.Sprintf("❌ Ошибка при создании автомобиля: %s", err.Error())
	}

	return fmt.Sprintf("✅ Автомобиль **%s %s** (%d) успешно добавлен! 🚗", c.Brand, c.Model, c.Year)
}

func (h *ChatHandler) executeCreateFuel(data map[string]interface{}, userID string) string {
	e := fuel.FuelEvent{
		ID:     uuid.New().String(),
		CarID:  getString(data, "car_id"),
		Volume: getFloat(data, "volume"),
		Price:  getFloat(data, "price"),
		Type:   getString(data, "type"),
		Date:   getString(data, "date"),
	}

	if e.CarID == "" || e.Volume <= 0 || e.Price <= 0 || e.Type == "" || e.Date == "" {
		return "❌ Не хватает данных для создания заправки. Укажите автомобиль, объём, цену, тип топлива и дату."
	}

	if err := h.uc.Fuel.AddFuelEvent(e); err != nil {
		return fmt.Sprintf("❌ Ошибка при создании заправки: %s", err.Error())
	}

	return fmt.Sprintf("✅ Заправка добавлена: %.1f л %s по %.2f ₽/л на сумму %.2f ₽ 🛢️", e.Volume, e.Type, e.Price, e.Volume*e.Price)
}

func (h *ChatHandler) executeCreateMaintenance(data map[string]interface{}, userID string) string {
	e := maintenance.MaintenanceEvent{
		ID:          uuid.New().String(),
		CarID:       getString(data, "car_id"),
		Type:        getString(data, "type"),
		Date:        getString(data, "date"),
		Cost:        getFloat(data, "cost"),
		Description: getString(data, "description"),
	}

	if e.CarID == "" || e.Type == "" || e.Date == "" || e.Cost <= 0 {
		return "❌ Не хватает данных для создания ТО. Укажите автомобиль, тип работ, дату и стоимость."
	}

	if err := h.uc.Maintenance.AddMaintenanceEvent(e); err != nil {
		return fmt.Sprintf("❌ Ошибка при создании ТО: %s", err.Error())
	}

	return fmt.Sprintf("✅ Техобслуживание добавлено: %s на сумму %.2f ₽ 🔧", e.Type, e.Cost)
}

func (h *ChatHandler) executeCreateFine(data map[string]interface{}, userID string) string {
	status := getString(data, "status")
	if status == "" {
		status = "unpaid"
	}

	f := fine.Fine{
		ID:          uuid.New().String(),
		CarID:       getString(data, "car_id"),
		Amount:      getFloat(data, "amount"),
		Type:        getString(data, "type"),
		Date:        getString(data, "date"),
		Status:      status,
		Description: getString(data, "description"),
	}

	if f.CarID == "" || f.Amount <= 0 || f.Type == "" || f.Date == "" {
		return "❌ Не хватает данных для создания штрафа. Укажите автомобиль, сумму, статью и дату."
	}

	if err := h.uc.Fine.AddFine(f); err != nil {
		return fmt.Sprintf("❌ Ошибка при создании штрафа: %s", err.Error())
	}

	statusText := "не оплачен"
	if f.Status == "paid" {
		statusText = "оплачен"
	}

	return fmt.Sprintf("✅ Штраф добавлен: %.2f ₽ по статье %s (статус: %s) 📋", f.Amount, f.Type, statusText)
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
	httpReq.Header.Set("Authorization", "Api-Key "+h.apiKey)

	resp, err := http.DefaultClient.Do(httpReq)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	respBody, _ := ioutil.ReadAll(resp.Body)

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("YandexGPT API error: status=%d, body=%s", resp.StatusCode, string(respBody))
	}

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