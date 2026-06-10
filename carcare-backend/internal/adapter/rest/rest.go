func (h *FineHandler) handleAdd(w http.ResponseWriter, r *http.Request) {
	var f fine.Fine
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "invalid body", http.StatusBadRequest)
		return
	}
	if err := json.Unmarshal(body, &f); err != nil {
		http.Error(w, "invalid json", http.StatusBadRequest)
		return
	}
	if f.ID == "" {
		f.ID = uuid.New().String()
	}
	// Проверяем дубликат по номеру постановления
	if f.BillNumber != "" {
		exists, err := h.fineRepo.CheckFineExistsByBillNumber(f.CarID, f.BillNumber)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		if exists {
			http.Error(w, "Штраф с таким номером постановления уже существует", http.StatusConflict)
			return
		}
	}
	if err := h.Add.Execute(f); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(f)
}