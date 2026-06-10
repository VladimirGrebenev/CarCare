export async function addFine(fine) {
  const body = {
    car_id: fine.carId,
    amount: fine.amount,
    type: fine.type,
    date: fine.date,
    status: fine.status ?? 'unpaid',
    description: fine.description ?? '',
    bill_number: fine.billNumber ?? '',
  };
  const res = await fetch('/api/fines', {
    method: 'POST',
    headers: withAuthHeaders({ 'Content-Type': 'application/json' }),
    credentials: 'include',
    body: JSON.stringify(body)
  });
  if (!res.ok) {
    throw new Error(createErrorMessage(res.status, 'добавление штрафа'));
  }
  const f = await res.json();
  return { id: f.id, carId: f.car_id, amount: f.amount, type: f.type, date: f.date, status: f.status, description: f.description, billNumber: f.bill_number };
}

export async function updateFine(id, fine) {
  const body = {
    car_id: fine.carId,
    amount: fine.amount,
    type: fine.type,
    date: fine.date,
    status: fine.status ?? 'unpaid',
    description: fine.description ?? '',
    bill_number: fine.billNumber ?? '',
  };
  const res = await fetch(`/api/fines/${id}`, {
    method: 'PUT',
    headers: withAuthHeaders({ 'Content-Type': 'application/json' }),
    credentials: 'include',
    body: JSON.stringify(body)
  });
  if (!res.ok) {
    throw new Error(createErrorMessage(res.status, 'обновление штрафа'));
  }
  const f = await res.json();
  return {
    id: f.id,
    carId: f.car_id,
    amount: f.amount,
    type: f.type,
    date: f.date,
    status: f.status,
    description: f.description,
    billNumber: f.bill_number,
  };
}