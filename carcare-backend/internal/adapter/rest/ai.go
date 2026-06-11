package rest

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"regexp"
	"strings"
	"time"

	"github.com/google/uuid"

	"github.com/VladimirGrebenev/CarCare-backend/internal/domain/car"
	"github.com/VladimirGrebenev/CarCare-backend/internal/domain/fine"
	"github.com/VladimirGrebenev/CarCare-backend/internal/domain/fuel"
	"github.com/VladimirGrebenev/CarCare-backend/internal/domain/maintenance"
	"github.com/VladimirGrebenev/CarCare-backend/internal/usecase"
)

// ChatRequest — запрос от фронтенда
type ChatRequest struct {
	Message string   `json:"message"`
	History []Message `json:"history,omitempty"`
}

type Message struct {
	Role string `json:"role"`
	Text string `json:"text"`
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

// buildSystemPrompt — инструкция для AI с текущей датой
func buildSystemPrompt() string {
	today := time.Now().Format("2006-01-02")
	return fmt.Sprintf(`Сегодняшняя дата: %s. Используй её когда пользователь говорит "сегодня", "вчера", "позавчера" и т.п.

`, today) + systemPromptTemplate
}

// systemPromptTemplate — инструкция для AI
const systemPromptTemplate = "Ты — AI-помощник приложения CarCare для учёта расходов на автомобиль.\n\nТвои задачи:\n1. Отвечать на вопросы пользователей о том, как пользоваться приложением.\n2. Создавать, редактировать и удалять сущности (автомобили, заправки, ТО, штрафы).\n\nВАЖНЫЕ ПРАВИЛА:\n- Не здоровайся в каждом ответе. Отвечай по делу, коротко.\n- НИКОГДА не создавай сущность, пока не собраны ВСЕ обязательные поля.\n- НИКОГДА не используй данные из предыдущих сообщений истории для автоматического заполнения нового запроса — каждое новое \"создай X\" это чистый запрос, спрашивай все поля заново.\n- НИКОГДА не выводи числовые значения (сумму, стоимость, объём) самостоятельно на основе описания или интернета — если пользователь не назвал сумму, спроси; можешь предложить типичные значения как подсказку, но решение за пользователем.\n- ВСЕ даты в JSON передавай строго в формате YYYY-MM-DD. Если пользователь сказал \"сегодня\" — подставь сегодняшнюю дату. Если сказал \"вчера\" — вычти один день. Никогда не передавай в JSON слова \"сегодня\", \"вчера\" и т.п.\n- Если данных не хватает — задавай уточняющие вопросы по одному за раз.\n- НЕ используй ID из базы данных. Спрашивай понятными пользователю названиями.\n- Для привязки к автомобилю достаточно одного госномера. Марку и модель для поиска НЕ нужно спрашивать — только госномер.\n- Для удаления — сначала уточни, какую именно запись удалить, и подтверди у пользователя.\n- Когда собраны все данные для создания — покажи пользователю итог со ВСЕМИ собранными данными и спроси подтверждение простым текстом, БЕЗ JSON. Перечисли каждое поле явно.\n  Пример штрафа: \"Создать штраф для Москвич 3 (Р392РЗ03): сумма 500 ₽, описание: Превышение скорости, дата: 10.06.2026, номер постановления: нет?\"\n  Пример автомобиля: \"Создать автомобиль Москвич 3, 2024, госномер Р392РЗ03, VIN WVGZZZ5NZLW123476?\"\n- Если пользователь ответил \"да\", \"давай\", \"подтверждаю\" или аналогично — сразу выводи JSON с действием.\n- Если пользователь ответил \"нет\" или отказался — не выводи JSON, спроси что исправить.\n- JSON выводи ТОЛЬКО когда пользователь подтвердил действие. Не выводи JSON в вопросе на подтверждение.\n- В JSON для госномера и VIN передавай значения без лишних пробелов.\n- Нормализуй названия марок и моделей автомобилей: \"лада\" → \"Lada\", \"веста\" → \"Vesta\", \"фольксваген\" → \"Volkswagen\", \"тигуан\" → \"Tiguan\", \"москвич\" → \"Москвич\", \"киа\" → \"Kia\", \"хитай\" → \"Hyundai\", \"шкода\" → \"Skoda\", \"одди\" → \"Audi\", \"бмв\" → \"BMW\", \"мерседес\" → \"Mercedes-Benz\", \"ниссан\" → \"Nissan\", \"тойота\" → \"Toyota\", \"хонда\" → \"Honda\", \"мазда\" → \"Mazda\", \"митсубиси\" → \"Mitsubishi\", \"субару\" → \"Subaru\", \"лексус\" → \"Lexus\". Если марка неизвестна — оставь как написал пользователь.\n\nОбязательные поля для создания:\n1. АВТОМОБИЛЬ (create_car): brand, model, year, vin, plate\n2. ЗАПРАВКА (create_fuel): автомобиль, volume, price, type, date\n3. ТЕХОБСЛУЖИВАНИЕ (create_maintenance): автомобиль, type, date, cost\n4. ШТРАФ (create_fine): спрашивай строго в таком порядке — 1) госномер, 2) сумма штрафа (ОБЯЗАТЕЛЬНО спроси явно, никогда не угадывай по описанию), 3) дата нарушения (ОБЯЗАТЕЛЬНО спроси явно, не подставляй сегодняшнюю дату по умолчанию), 4) описание нарушения, 5) номер постановления (необязателен, если нет — создавай без него)\n\nДоступные действия (JSON в конце ответа, только после подтверждения пользователя):\n\nСоздание:\n{\"action\": \"create_car\", \"data\": {\"brand\": \"...\", \"model\": \"...\", \"year\": 2024, \"vin\": \"...\", \"plate\": \"...\"}}\n{\"action\": \"create_fuel\", \"data\": {\"car_plate\": \"...\", \"volume\": 45, \"price\": 55, \"type\": \"АИ-95\", \"date\": \"2026-06-09\"}}\n{\"action\": \"create_maintenance\", \"data\": {\"car_plate\": \"...\", \"type\": \"Замена масла\", \"date\": \"2026-06-09\", \"cost\": 5000}}\n{\"action\": \"create_fine\", \"data\": {\"car_plate\": \"...\", \"amount\": 500, \"date\": \"2026-06-09\", \"description\": \"Превышение скорости на 20 км/ч\", \"bill_number\": \"18810533260608041459\"}}\n\nРедактирование:\n{\"action\": \"update_car\", \"data\": {\"car_plate\": \"...\", \"new_plate\": \"...\", \"new_vin\": \"...\"}}\n{\"action\": \"update_fuel\", \"data\": {\"car_plate\": \"...\", \"date\": \"2026-06-09\", \"new_volume\": 50}}\n{\"action\": \"update_maintenance\", \"data\": {\"car_plate\": \"...\", \"type\": \"Замена масла\", \"date\": \"2026-06-09\", \"new_cost\": 6000}}\n{\"action\": \"update_fine\", \"data\": {\"car_plate\": \"...\", \"amount\": 500, \"date\": \"2026-06-09\", \"new_status\": \"paid\"}}\n\nУдаление (только после подтверждения пользователем):\n{\"action\": \"delete_car\", \"data\": {\"car_plate\": \"...\"}}\n{\"action\": \"delete_fuel\", \"data\": {\"car_plate\": \"...\", \"date\": \"2026-06-09\", \"volume\": 45}}\n{\"action\": \"delete_maintenance\", \"data\": {\"car_plate\": \"...\", \"type\": \"Замена масла\", \"date\": \"2026-06-09\"}}\n{\"action\": \"delete_fine\", \"data\": {\"car_plate\": \"...\", \"amount\": 500, \"date\": \"2026-06-09\"}}"

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

	reply, err := h.callYandexGPT(req.Message, req.History)
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
	start := strings.Index(response, "{\"action\":")
	if start == -1 {
		return ""
	}
	end := strings.LastIndex(response, "}")
	if end == -1 || end < start {
		return ""
	}
	jsonStr := response[start : end+1]

	var action AIAction
	if err := json.Unmarshal([]byte(jsonStr), &action); err != nil {
		return ""
	}

	switch action.Action {
	case "create_car":
		return h.executeCreateCar(action.Data, userID)
	case "create_fuel":
		return h.executeCreateFuel(action.Data, userID)
	case "create_maintenance":
		return h.executeCreateMaintenance(action.Data, userID)
	case "create_fine":
		return h.executeCreateFine(action.Data, userID)
	case "update_car":
		return h.executeUpdateCar(action.Data, userID)
	case "update_fuel":
		return h.executeUpdateFuel(action.Data, userID)
	case "update_maintenance":
		return h.executeUpdateMaintenance(action.Data, userID)
	case "update_fine":
		return h.executeUpdateFine(action.Data, userID)
	case "delete_car":
		return h.executeDeleteCar(action.Data, userID)
	case "delete_fuel":
		return h.executeDeleteFuel(action.Data, userID)
	case "delete_maintenance":
		return h.executeDeleteMaintenance(action.Data, userID)
	case "delete_fine":
		return h.executeDeleteFine(action.Data, userID)
	default:
		return ""
	}
}

func getString(data map[string]interface{}, key string) string {
	if v, ok := data[key]; ok {
		if s, ok := v.(string); ok {
			return strings.TrimSpace(s)
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

// normalizePlate убирает пробелы в госномере
func normalizePlate(plate string) string {
	re := regexp.MustCompile(`\s+`)
	return strings.ToUpper(re.ReplaceAllString(strings.TrimSpace(plate), ""))
}

// findCarByPlate ищет автомобиль пользователя по госномеру
func (h *ChatHandler) findCarByPlate(data map[string]interface{}, userID string) (string, string) {
	plate := normalizePlate(getString(data, "car_plate"))
	if plate == "" {
		return "", ""
	}

	cars, err := h.uc.Car.ListCars(userID)
	if err != nil {
		return "", ""
	}

	for _, c := range cars {
		if strings.EqualFold(normalizePlate(c.Plate), plate) {
			return c.ID, fmt.Sprintf("%s %s (%s)", c.Brand, c.Model, c.Plate)
		}
	}
	return "", ""
}

// ========== CREATE ==========

func (h *ChatHandler) executeCreateCar(data map[string]interface{}, userID string) string {
	brand := getString(data, "brand")
	model := getString(data, "model")
	year := getInt(data, "year")
	vin := getString(data, "vin")
	plate := normalizePlate(getString(data, "plate"))

	missingFields := []string{}
	if brand == "" {
		missingFields = append(missingFields, "марку")
	}
	if model == "" {
		missingFields = append(missingFields, "модель")
	}
	if year <= 0 {
		missingFields = append(missingFields, "год выпуска")
	}
	if vin == "" {
		missingFields = append(missingFields, "VIN-номер")
	}
	if plate == "" {
		missingFields = append(missingFields, "госномер")
	}

	if len(missingFields) > 0 {
		return fmt.Sprintf("❌ Не хватает данных: %s. Пожалуйста, укажите все данные для создания автомобиля.", strings.Join(missingFields, ", "))
	}

	c := car.Car{
		ID:     uuid.New().String(),
		UserID: userID,
		Brand:  brand,
		Model:  model,
		Year:   year,
		VIN:    vin,
		Plate:  plate,
	}

	if err := h.uc.Car.AddCar(c); err != nil {
		return fmt.Sprintf("❌ Ошибка при создании автомобиля: %s", err.Error())
	}

	return fmt.Sprintf("✅ Автомобиль **%s %s** (%d) успешно добавлен! 🚗\nГосномер: %s\nVIN: %s", brand, model, year, plate, vin)
}

func (h *ChatHandler) executeCreateFuel(data map[string]interface{}, userID string) string {
	carID, carName := h.findCarByPlate(data, userID)
	if carID == "" {
		return "❌ Не удалось найти автомобиль. Укажите госномер."
	}

	e := fuel.FuelEvent{
		ID:     uuid.New().String(),
		CarID:  carID,
		Volume: getFloat(data, "volume"),
		Price:  getFloat(data, "price"),
		Type:   getString(data, "type"),
		Date:   getString(data, "date"),
	}

	if e.Volume <= 0 || e.Price <= 0 || e.Type == "" || e.Date == "" {
		return "❌ Не хватает данных для создания заправки. Укажите объём, цену, тип топлива и дату."
	}

	if err := h.uc.Fuel.AddFuelEvent(e); err != nil {
		return fmt.Sprintf("❌ Ошибка при создании заправки: %s", err.Error())
	}

	return fmt.Sprintf("✅ Заправка добавлена для %s: %.1f л %s по %.2f ₽/л 🛢️", carName, e.Volume, e.Type, e.Price)
}

func (h *ChatHandler) executeCreateMaintenance(data map[string]interface{}, userID string) string {
	carID, carName := h.findCarByPlate(data, userID)
	if carID == "" {
		return "❌ Не удалось найти автомобиль. Укажите госномер."
	}

	e := maintenance.MaintenanceEvent{
		ID:          uuid.New().String(),
		CarID:       carID,
		Type:        getString(data, "type"),
		Date:        getString(data, "date"),
		Cost:        getFloat(data, "cost"),
		Description: getString(data, "description"),
	}

	if e.Type == "" || e.Date == "" || e.Cost <= 0 {
		return "❌ Не хватает данных для создания ТО. Укажите тип работ, дату и стоимость."
	}

	if err := h.uc.Maintenance.AddMaintenanceEvent(e); err != nil {
		return fmt.Sprintf("❌ Ошибка при создании ТО: %s", err.Error())
	}

	return fmt.Sprintf("✅ Техобслуживание добавлено для %s: %s на сумму %.2f ₽ 🔧", carName, e.Type, e.Cost)
}

func (h *ChatHandler) executeCreateFine(data map[string]interface{}, userID string) string {
	carID, carName := h.findCarByPlate(data, userID)
	if carID == "" {
		return "❌ Не удалось найти автомобиль. Укажите госномер."
	}

	status := getString(data, "status")
	if status == "" {
		status = "unpaid"
	}

	fineType := getString(data, "type")
	if fineType == "" {
		fineType = "Штраф"
	}

	f := fine.Fine{
		ID:          uuid.New().String(),
		CarID:       carID,
		Amount:      getFloat(data, "amount"),
		Type:        fineType,
		Date:        getString(data, "date"),
		Status:      status,
		Description: getString(data, "description"),
		BillNumber:  getString(data, "bill_number"),
	}

	if f.Amount <= 0 || f.Date == "" {
		return "❌ Не хватает данных для создания штрафа. Укажите сумму и дату."
	}

	if err := h.uc.Fine.AddFine(f); err != nil {
		return fmt.Sprintf("❌ Ошибка при создании штрафа: %s", err.Error())
	}

	statusText := "не оплачен"
	if f.Status == "paid" {
		statusText = "оплачен"
	}

	desc := f.Description
	if desc == "" {
		desc = "без описания"
	}
	return fmt.Sprintf("✅ Штраф добавлен для %s: %.2f ₽, %s (статус: %s) 📋", carName, f.Amount, desc, statusText)
}

// ========== UPDATE ==========

func (h *ChatHandler) executeUpdateCar(data map[string]interface{}, userID string) string {
	carID, carName := h.findCarByPlate(data, userID)
	if carID == "" {
		return "❌ Не удалось найти автомобиль. Укажите госномер."
	}

	existingCar, err := h.uc.Car.GetCar(carID)
	if err != nil {
		return "❌ Автомобиль не найден."
	}

	if plate := normalizePlate(getString(data, "new_plate")); plate != "" {
		existingCar.Plate = plate
	}
	if vin := getString(data, "new_vin"); vin != "" {
		existingCar.VIN = vin
	}
	if brand := getString(data, "new_brand"); brand != "" {
		existingCar.Brand = brand
	}
	if model := getString(data, "new_model"); model != "" {
		existingCar.Model = model
	}
	if year := getInt(data, "new_year"); year > 0 {
		existingCar.Year = year
	}

	if err := h.uc.Car.UpdateCar(existingCar, userID); err != nil {
		return fmt.Sprintf("❌ Ошибка при обновлении: %s", err.Error())
	}

	return fmt.Sprintf("✅ Автомобиль %s обновлён! 🚗", carName)
}

func (h *ChatHandler) executeUpdateFuel(data map[string]interface{}, userID string) string {
	carID, carName := h.findCarByPlate(data, userID)
	if carID == "" {
		return "❌ Не удалось найти автомобиль."
	}

	events, err := h.uc.Fuel.ListFuelEvents(userID)
	if err != nil {
		return "❌ Ошибка при поиске заправки."
	}

	date := getString(data, "date")
	volume := getFloat(data, "volume")

	var found *fuel.FuelEvent
	for _, e := range events {
		if e.CarID == carID && e.Date == date {
			if volume <= 0 || e.Volume == volume {
				found = &e
				break
			}
		}
	}

	if found == nil {
		return "❌ Заправка не найдена. Уточните дату и объём."
	}

	if v := getFloat(data, "new_volume"); v > 0 {
		found.Volume = v
	}
	if p := getFloat(data, "new_price"); p > 0 {
		found.Price = p
	}
	if t := getString(data, "new_type"); t != "" {
		found.Type = t
	}
	if d := getString(data, "new_date"); d != "" {
		found.Date = d
	}

	if err := h.uc.Fuel.UpdateFuelEvent(*found); err != nil {
		return fmt.Sprintf("❌ Ошибка при обновлении: %s", err.Error())
	}

	return fmt.Sprintf("✅ Заправка для %s обновлена! 🛢️", carName)
}

func (h *ChatHandler) executeUpdateMaintenance(data map[string]interface{}, userID string) string {
	carID, carName := h.findCarByPlate(data, userID)
	if carID == "" {
		return "❌ Не удалось найти автомобиль."
	}

	events, err := h.uc.Maintenance.ListMaintenanceEvents(userID)
	if err != nil {
		return "❌ Ошибка при поиске ТО."
	}

	date := getString(data, "date")
	mtype := getString(data, "type")

	var found *maintenance.MaintenanceEvent
	for _, e := range events {
		if e.CarID == carID && e.Date == date && (mtype == "" || e.Type == mtype) {
			found = &e
			break
		}
	}

	if found == nil {
		return "❌ Запись ТО не найдена. Уточните дату и тип работ."
	}

	if c := getFloat(data, "new_cost"); c > 0 {
		found.Cost = c
	}
	if t := getString(data, "new_type"); t != "" {
		found.Type = t
	}
	if d := getString(data, "new_date"); d != "" {
		found.Date = d
	}
	if desc := getString(data, "new_description"); desc != "" {
		found.Description = desc
	}

	if err := h.uc.Maintenance.UpdateMaintenanceEvent(*found); err != nil {
		return fmt.Sprintf("❌ Ошибка при обновлении: %s", err.Error())
	}

	return fmt.Sprintf("✅ Техобслуживание для %s обновлено! 🔧", carName)
}

func (h *ChatHandler) executeUpdateFine(data map[string]interface{}, userID string) string {
	carID, carName := h.findCarByPlate(data, userID)
	if carID == "" {
		return "❌ Не удалось найти автомобиль."
	}

	fines, err := h.uc.Fine.ListFines(userID)
	if err != nil {
		return "❌ Ошибка при поиске штрафа."
	}

	date := getString(data, "date")
	amount := getFloat(data, "amount")

	var found *fine.Fine
	for _, f := range fines {
		if f.CarID == carID && f.Date == date && (amount <= 0 || f.Amount == amount) {
			found = &f
			break
		}
	}

	if found == nil {
		return "❌ Штраф не найден. Уточните дату и сумму."
	}

	if a := getFloat(data, "new_amount"); a > 0 {
		found.Amount = a
	}
	if s := getString(data, "new_status"); s != "" {
		found.Status = s
	}
	if t := getString(data, "new_type"); t != "" {
		found.Type = t
	}
	if d := getString(data, "new_date"); d != "" {
		found.Date = d
	}

	if err := h.uc.Fine.UpdateFine(*found); err != nil {
		return fmt.Sprintf("❌ Ошибка при обновлении: %s", err.Error())
	}

	return fmt.Sprintf("✅ Штраф для %s обновлён! 📋", carName)
}

// ========== DELETE ==========

func (h *ChatHandler) executeDeleteCar(data map[string]interface{}, userID string) string {
	carID, carName := h.findCarByPlate(data, userID)
	if carID == "" {
		return "❌ Не удалось найти автомобиль. Укажите госномер."
	}

	if err := h.uc.Car.DeleteCar(carID, userID); err != nil {
		return fmt.Sprintf("❌ Ошибка при удалении: %s", err.Error())
	}

	return fmt.Sprintf("✅ Автомобиль %s удалён. Все связанные заправки, ТО и штрафы также удалены. 🗑️", carName)
}

func (h *ChatHandler) executeDeleteFuel(data map[string]interface{}, userID string) string {
	carID, carName := h.findCarByPlate(data, userID)
	if carID == "" {
		return "❌ Не удалось найти автомобиль."
	}

	events, err := h.uc.Fuel.ListFuelEvents(userID)
	if err != nil {
		return "❌ Ошибка при поиске заправки."
	}

	date := getString(data, "date")
	volume := getFloat(data, "volume")

	var found *fuel.FuelEvent
	for _, e := range events {
		if e.CarID == carID && e.Date == date && (volume <= 0 || e.Volume == volume) {
			found = &e
			break
		}
	}

	if found == nil {
		return "❌ Заправка не найдена. Уточните дату и объём."
	}

	if err := h.uc.Fuel.DeleteFuelEvent(found.ID); err != nil {
		return fmt.Sprintf("❌ Ошибка при удалении: %s", err.Error())
	}

	return fmt.Sprintf("✅ Заправка для %s от %s удалена. 🗑️", carName, found.Date)
}

func (h *ChatHandler) executeDeleteMaintenance(data map[string]interface{}, userID string) string {
	carID, carName := h.findCarByPlate(data, userID)
	if carID == "" {
		return "❌ Не удалось найти автомобиль."
	}

	events, err := h.uc.Maintenance.ListMaintenanceEvents(userID)
	if err != nil {
		return "❌ Ошибка при поиске ТО."
	}

	date := getString(data, "date")
	mtype := getString(data, "type")

	var found *maintenance.MaintenanceEvent
	for _, e := range events {
		if e.CarID == carID && e.Date == date && (mtype == "" || e.Type == mtype) {
			found = &e
			break
		}
	}

	if found == nil {
		return "❌ Запись ТО не найдена. Уточните дату и тип работ."
	}

	if err := h.uc.Maintenance.DeleteMaintenanceEvent(found.ID); err != nil {
		return fmt.Sprintf("❌ Ошибка при удалении: %s", err.Error())
	}

	return fmt.Sprintf("✅ Техобслуживание для %s от %s удалено. 🗑️", carName, found.Date)
}

func (h *ChatHandler) executeDeleteFine(data map[string]interface{}, userID string) string {
	carID, carName := h.findCarByPlate(data, userID)
	if carID == "" {
		return "❌ Не удалось найти автомобиль."
	}

	fines, err := h.uc.Fine.ListFines(userID)
	if err != nil {
		return "❌ Ошибка при поиске штрафа."
	}

	date := getString(data, "date")
	amount := getFloat(data, "amount")

	var found *fine.Fine
	for _, f := range fines {
		if f.CarID == carID && f.Date == date && (amount <= 0 || f.Amount == amount) {
			found = &f
			break
		}
	}

	if found == nil {
		return "❌ Штраф не найден. Уточните дату и сумму."
	}

	if err := h.uc.Fine.DeleteFine(found.ID); err != nil {
		return fmt.Sprintf("❌ Ошибка при удалении: %s", err.Error())
	}

	return fmt.Sprintf("✅ Штраф для %s на %.2f ₽ удалён. 🗑️", carName, found.Amount)
}

func (h *ChatHandler) callYandexGPT(userMessage string, history []Message) (string, error) {
	if h.apiKey == "" || h.yandexGPTFolderID == "" {
		return "", fmt.Errorf("YandexGPT not configured")
	}

	modelURI := fmt.Sprintf("gpt://%s/yandexgpt/latest", h.yandexGPTFolderID)

	messages := []yandexGPTMessage{
		{Role: "system", Text: buildSystemPrompt()},
	}

	start := 0
	if len(history) > 20 {
		start = len(history) - 20
	}
	for _, msg := range history[start:] {
		messages = append(messages, yandexGPTMessage{Role: msg.Role, Text: msg.Text})
	}

	messages = append(messages, yandexGPTMessage{Role: "user", Text: userMessage})

	reqBody := yandexGPTRequest{}
	reqBody.ModelURI = modelURI
	reqBody.CompletionOptions.Stream = false
	reqBody.CompletionOptions.Temperature = 0.6
	reqBody.CompletionOptions.MaxTokens = 2000
	reqBody.Messages = messages

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

// fallbackResponse — простая за��лушка, если YandexGPT не настроен
func (h *ChatHandler) fallbackResponse(message string) string {
	msg := strings.ToLower(message)

	switch {
	case strings.Contains(msg, "автомобил") || strings.Contains(msg, "машин") || strings.Contains(msg, "добавить авто"):
		return "Чтобы добавить автомобиль, перейдите в раздел \"Мои автомобили\" и нажмите \"Добавить автомобиль\". Укажите марку, модель, год выпуска, VIN и госномер."

	case strings.Contains(msg, "заправк") || strings.Contains(msg, "топлив") || strings.Contains(msg, "бензин"):
		return "Чтобы добавить заправку, перейдите в раздел \"Топливо\" и нажмите \"Добавить заправку\". Укажите дату, объём в литрах, тип топлива, цену и пробег."

	case strings.Contains(msg, "то ") || strings.Contains(msg, "техобслуж") || strings.Contains(msg, "ремонт"):
		return "Чтобы добавить запись о техобслуживании, перейдите в раздел \"Техобслуживание\". Укажите дату, тип работ, пробег и стоимость."

	case strings.Contains(msg, "штраф"):
		return "Чтобы добавить штраф, перейдите в раздел \"Штрафы\" и нажмите \"Добавить штраф\". Укажите дату, номер постановления, сумму и статью нарушения."

	case strings.Contains(msg, "отчёт") || strings.Contains(msg, "статистик") || strings.Contains(msg, "расход"):
		return "В разделе \"Отчёты\" вы можете увидеть сводку по всем расходам: топливо, ТО и штрафы. Доступны фильтры по периоду (день, неделя, месяц, год) и графики распределения расходов."

	case strings.Contains(msg, "привет") || strings.Contains(msg, "здравствуй") || strings.Contains(msg, "здаров"):
		return "Привет! Я AI-помощник CarCare 👋 Я могу помочь с вопросами по приложению, а также создать автомобиль, заправку, ТО или штраф по вашему описанию. Что вас интересует?"

	default:
		return "Я могу помочь вам с вопросами по приложению CarCare! Расскажите, что вас интересует:\n\n• Как добавить автомобиль?\n• Как записать заправку?\n• Как добавить техобслуживание?\n• Как зарегистрировать штраф?\n• Где посмотреть отчёты?\n\nИли просто опишите, что хотите создать, и я помогу собрать все данные! 😊"
	}
}