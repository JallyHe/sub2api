/**
 * Credits system API endpoints
 * Handles credit balance, ledger, and admin credit management
 */

import { apiClient } from './client'

// ==================== Types ====================

export interface CreditBalance {
  balance: number
  expires_at: number | null  // Unix timestamp; null = no expiry
  plan_id: number | null
}

export interface CreditLedgerEntry {
  id: number
  delta: number
  reason: string
  model: string | null
  balance_after: number
  created_at: number  // Unix timestamp
}

export interface CreditLedgerResponse {
  items: CreditLedgerEntry[]
  page: number
  limit: number
}

export interface ModelCreditRate {
  id: number
  model_pattern: string
  credits_per_1k_tokens_input: number
  credits_per_1k_tokens_output: number
  priority: number
}

// ==================== User APIs ====================

export async function getBalance(): Promise<CreditBalance> {
  const { data } = await apiClient.get<CreditBalance>('/credits/balance')
  return data
}

export async function getLedger(params: { page?: number; limit?: number } = {}): Promise<CreditLedgerResponse> {
  const { data } = await apiClient.get<CreditLedgerResponse>('/credits/ledger', { params })
  return data
}

// ==================== Admin APIs ====================

export async function getModelRates(): Promise<{ items: ModelCreditRate[] }> {
  const { data } = await apiClient.get<{ items: ModelCreditRate[] }>('/admin/credits/model-rates')
  return data
}

export async function getUserBalance(userId: number): Promise<CreditBalance> {
  const { data } = await apiClient.get<CreditBalance>(`/admin/credits/users/${userId}/balance`)
  return data
}

export async function grantCredits(payload: {
  user_id: number
  credits: number
  notes?: string
}): Promise<{ success: boolean }> {
  const { data } = await apiClient.post<{ success: boolean }>('/admin/credits/grant', payload)
  return data
}

// ==================== Public APIs ====================

export interface PlanItem {
  id: number
  name: string
  description: string
  price: number
  original_price: number | null
  validity_days: number
  features: string
  for_sale: boolean
  sort_order: number
  credits: number
}

export async function getPublicPlans(): Promise<{ items: PlanItem[] }> {
  // /payment/plans is the existing public subscription plans endpoint
  const { data } = await apiClient.get<PlanItem[]>('/payment/plans')
  // Normalize to { items } shape; API returns an array directly
  return { items: Array.isArray(data) ? data : [] }
}

export const creditsApi = {
  getBalance,
  getLedger,
  getModelRates,
  getUserBalance,
  grantCredits,
  getPublicPlans,
}

export default creditsApi
