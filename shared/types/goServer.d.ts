/* Do not change, this code is generated from Golang structs */

export interface goRes{
  message: string,
}
export interface AuthResponse{
  message: string,
  token?: string,
}
export interface Salary{
  vl: number,
  type: string,
  day?: number,
  month: string,
}
export interface SalaryResponse{
  id: number,
  vl: number,
  type: string,
  day?: number,
  month: string,
}
export interface Expense{
  vl: number,
  type: string,
  description?: string,
  day?: number,
  paymentMethod: string,
  month: string,
}
export interface ExpenseResponse{
  id: number,
  vl: number,
  type: string,
  description?: string,
  day?: number,
  paymentMethod: string,
  month: string,
}
