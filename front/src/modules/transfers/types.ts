/**
 * Transfers Module - Type Definitions
 * Domain: Money transfers, deposits, and withdrawals
 */

export enum TransferType {
  PIX = "pix",
  TED = "ted",
  P2P = "p2p", // Platform user to user
  DEPOSIT = "deposit",
  WITHDRAWAL = "withdrawal",
}

export enum TransferStatus {
  Pending = "pending",
  Processing = "processing",
  Completed = "completed",
  Failed = "failed",
  Cancelled = "cancelled",
  Refunded = "refunded",
}

export enum DepositMethod {
  PIX = "pix",
  BOLETO = "boleto",
  BANK_TRANSFER = "bank_transfer",
}

export enum PIXKeyType {
  CPF = "cpf",
  CNPJ = "cnpj",
  EMAIL = "email",
  PHONE = "phone",
  RANDOM = "random", // Chave aleatória
}

export interface Transfer {
  id: string;
  userId: string;
  type: TransferType;
  status: TransferStatus;
  amount: number;
  currency: string;
  description?: string;
  fromAccount?: string;
  toAccount?: string;
  pixKey?: string;
  pixKeyType?: PIXKeyType;
  recipientName?: string;
  recipientDocument?: string;
  recipientBank?: string;
  recipientBranch?: string;
  recipientAccount?: string;
  recipientAccountType?: "checking" | "savings";
  fee?: number;
  scheduledFor?: string;
  completedAt?: string;
  failureReason?: string;
  createdAt: string;
  updatedAt: string;
}

export interface Deposit {
  id: string;
  userId: string;
  method: DepositMethod;
  status: TransferStatus;
  amount: number;
  currency: string;
  pixQRCode?: string;
  pixKey?: string;
  boletoCode?: string;
  boletoUrl?: string;
  expiresAt?: string;
  completedAt?: string;
  createdAt: string;
  updatedAt: string;
}

export interface PaymentRequest {
  id: string;
  userId: string;
  amount: number;
  description: string;
  pixKey?: string;
  qrCode?: string;
  paymentLink: string;
  expiresAt?: string;
  paidAt?: string;
  status: "pending" | "paid" | "expired" | "cancelled";
  createdAt: string;
  updatedAt: string;
}

export interface TransferReceipt {
  transferId: string;
  type: TransferType;
  amount: number;
  currency: string;
  recipientName: string;
  date: string;
  authenticationCode: string;
  status: TransferStatus;
}
