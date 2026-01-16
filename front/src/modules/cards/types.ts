/**
 * Cards Module - Type Definitions
 * Domain: Card management, virtual/physical cards, controls
 */

export enum CardType {
  PHYSICAL = "physical",
  VIRTUAL = "virtual",
}

export enum CardStatus {
  Active = "active",
  Blocked = "blocked",
  Cancelled = "cancelled",
  Expired = "expired",
  Lost = "lost",
  Stolen = "stolen",
}

export enum CardBrand {
  VISA = "visa",
  MASTERCARD = "mastercard",
  ELO = "elo",
}

export interface Card {
  id: string;
  userId: string;
  type: CardType;
  status: CardStatus;
  brand: CardBrand;
  lastFourDigits: string;
  holderName: string;
  expiryMonth: number;
  expiryYear: number;
  cvv?: string; // Only shown in specific contexts for security
  dailyLimit: number;
  monthlyLimit: number;
  currentDailySpent: number;
  currentMonthlySpent: number;
  isContactless: boolean;
  isInternational: boolean;
  createdAt: string;
  updatedAt: string;
}

export interface CardDetails {
  id: string;
  cardNumber: string; // Masked: **** **** **** 1234
  cvv: string;
  expiryMonth: number;
  expiryYear: number;
  holderName: string;
}

export interface CardTransaction {
  id: string;
  cardId: string;
  amount: number;
  currency: string;
  merchantName: string;
  merchantCategory: string;
  status: "approved" | "declined" | "pending";
  transactionDate: string;
  location?: string;
  isInternational: boolean;
}

export interface CardLimit {
  daily: number;
  monthly: number;
  perTransaction?: number;
}

export interface CardSecuritySettings {
  blockInternational: boolean;
  blockOnline: boolean;
  blockContactless: boolean;
  blockAtm: boolean;
  requireNotification: boolean;
}
