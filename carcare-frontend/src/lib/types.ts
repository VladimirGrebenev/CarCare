export interface Car {
  id: string;
  brand: string;
  model: string;
  year: number;
  plate: string;
  vin: string;
  sts: string;
}
// src/lib/types.ts
export interface Fine {
  id: string;
  date: string;
  amount: number;
  description: string;
  status: 'paid' | 'unpaid';
  carId?: string;
  billNumber?: string;
  [key: string]: unknown;
}