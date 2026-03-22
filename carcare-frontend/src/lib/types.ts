// src/lib/types.ts
export interface Fine {
  id: string;
  date: string;
  amount: number;
  description: string;
  status: 'paid' | 'unpaid';
  carId?: string;
  [key: string]: any;
}
