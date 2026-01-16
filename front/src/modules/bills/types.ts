/**
 * Bills Module - Type Definitions
 * Domain: Bill payments (water, electricity, internet, etc.)
 */

export enum BillType {
  WATER = "water",
  ELECTRICITY = "electricity",
  INTERNET = "internet",
  PHONE = "phone",
  GAS = "gas",
  OTHER = "other",
}

export enum BillStatus {
  Pending = "pending",
  Processing = "processing",
  Paid = "paid",
  Failed = "failed",
  Cancelled = "cancelled",
  Overdue = "overdue",
}

export interface Bill {
  id: string;
  userId: string;
  type: BillType;
  status: BillStatus;
  barcode: string;
  amount: number;
  currency: string;
  recipientName: string;
  recipientDocument?: string;
  dueDate: string;
  paymentDate?: string;
  description?: string;
  fee?: number;
  discount?: number;
  finalAmount: number;
  createdAt: string;
  updatedAt: string;
}

export interface BarcodeData {
  barcode: string;
  amount?: number;
  dueDate?: string;
  recipientName?: string;
  recipientDocument?: string;
  type?: BillType;
}

export interface BillPaymentReceipt {
  billId: string;
  barcode: string;
  amount: number;
  recipientName: string;
  paymentDate: string;
  authenticationCode: string;
  status: BillStatus;
}
