<div class="field">
        <label class="field-label" for="fine-status">Статус *</label>
        <select id="fine-status" class="field-select" bind:value={form.status}>
          <option value="unpaid">Не оплачен</option>
          <option value="paid">Оплачен</option>
        </select>
      </div>
      <div class="col-span-2">
        <Input
          label="Номер постановления"
          placeholder="18810533260608041459"
          bind:value={form.billNumber}
        />
      </div>