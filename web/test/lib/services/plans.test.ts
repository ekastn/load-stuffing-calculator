import { PlanService } from "@/lib/services/plans"
import { describe, it, expect, beforeEach, vi, Mock } from "vitest"

vi.mock("@/lib/api", () => ({
  apiGet: vi.fn(),
  apiPost: vi.fn(),
  apiPut: vi.fn(),
  apiDelete: vi.fn(),
}))

import { apiGet, apiPost, apiPut, apiDelete } from "@/lib/api"

describe("PlanService", () => {
  beforeEach(() => {
    vi.clearAllMocks()
  })

  it("listPlans calls apiGet correctly", async () => {
    const mockData = [{ plan_id: "1", plan_code: "PL-1" }]
    ;(apiGet as Mock).mockResolvedValue(mockData)

    const result = await PlanService.listPlans(1, 10)
    expect(apiGet).toHaveBeenCalledWith("/plans?page=1&limit=10")
    expect(result).toEqual(mockData)
  })

  it("getPlan calls apiGet correctly", async () => {
    const mockPlan = { plan_id: "1", plan_code: "PL-1" }
    ;(apiGet as Mock).mockResolvedValue(mockPlan)

    const result = await PlanService.getPlan("1")
    expect(apiGet).toHaveBeenCalledWith("/plans/1")
    expect(result).toEqual(mockPlan)
  })

  it("createPlan calls apiPost correctly", async () => {
    const newPlan = { title: "Test Plan", container: { length_mm: 100 }, items: [] }
    const createdPlan = { plan_id: "1", ...newPlan }
    ;(apiPost as Mock).mockResolvedValue(createdPlan)

    const result = await PlanService.createPlan(newPlan as any)
    expect(apiPost).toHaveBeenCalledWith("/plans", newPlan)
    expect(result).toEqual(createdPlan)
  })

  it("createPlan appends workspace_id when provided", async () => {
    const newPlan = { title: "Test Plan", container: { length_mm: 100 }, items: [] }
    const createdPlan = { plan_id: "1", ...newPlan }
    ;(apiPost as Mock).mockResolvedValue(createdPlan)

    const result = await PlanService.createPlan(newPlan as any, "ws-123")
    expect(apiPost).toHaveBeenCalledWith("/plans?workspace_id=ws-123", newPlan)
    expect(result).toEqual(createdPlan)
  })

  it("updatePlan calls apiPut correctly", async () => {
    const updateData = { status: "COMPLETED" }
    ;(apiPut as Mock).mockResolvedValue(undefined)

    await PlanService.updatePlan("1", updateData)
    expect(apiPut).toHaveBeenCalledWith("/plans/1", updateData)
  })

  it("deletePlan calls apiDelete correctly", async () => {
    ;(apiDelete as Mock).mockResolvedValue(undefined)

    await PlanService.deletePlan("1")
    expect(apiDelete).toHaveBeenCalledWith("/plans/1")
  })

  it("addPlanItem calls apiPost correctly", async () => {
    const newItem = { label: "Box", length_mm: 10 }
    const createdItem = { item_id: "2", ...newItem }
    ;(apiPost as Mock).mockResolvedValue(createdItem)

    const result = await PlanService.addPlanItem("1", newItem as any)
    expect(apiPost).toHaveBeenCalledWith("/plans/1/items", newItem)
    expect(result).toEqual(createdItem)
  })

  it("updatePlanItem calls apiPut correctly", async () => {
    const updateData = { quantity: 5 }
    ;(apiPut as Mock).mockResolvedValue(undefined)

    await PlanService.updatePlanItem("1", "2", updateData)
    expect(apiPut).toHaveBeenCalledWith("/plans/1/items/2", updateData)
  })

  it("deletePlanItem calls apiDelete correctly", async () => {
    ;(apiDelete as Mock).mockResolvedValue(undefined)

    await PlanService.deletePlanItem("1", "2")
    expect(apiDelete).toHaveBeenCalledWith("/plans/1/items/2")
  })

  it("calculatePlan calls apiPost correctly", async () => {
    const mockResult = { job_id: "job1", status: "COMPLETED" }
    ;(apiPost as Mock).mockResolvedValue(mockResult)

    const result = await PlanService.calculatePlan("1")
    expect(apiPost).toHaveBeenCalledWith("/plans/1/calculate", {})
    expect(result).toEqual(mockResult)
  })

  it("calculatePlan passes options body", async () => {
    const mockResult = { job_id: "job1", status: "COMPLETED" }
    ;(apiPost as Mock).mockResolvedValue(mockResult)

    const options = { strategy: "parallel", goal: "tightest", gravity: true }
    const result = await PlanService.calculatePlan("1", options)
    expect(apiPost).toHaveBeenCalledWith("/plans/1/calculate", options)
    expect(result).toEqual(mockResult)
  })
})
