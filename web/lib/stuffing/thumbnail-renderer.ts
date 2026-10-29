import { StuffingVisualizer } from "@/lib/StuffingVisualizer"
import type { StuffingPlanData } from "@/lib/StuffingVisualizer"

export type ThumbnailGeneratorFn = (data: StuffingPlanData) => Promise<string | null>

const THUMBNAIL_WIDTH = 640
const THUMBNAIL_HEIGHT = 360

const thumbnailCache = new Map<string, string>()
let renderQueue: Promise<unknown> = Promise.resolve()
let sharedVisualizer: StuffingVisualizer | null = null
let offscreenHost: HTMLDivElement | null = null
let customGenerator: ThumbnailGeneratorFn | null = null

function getCacheKey(data: StuffingPlanData): string {
  const jobId = data.calculation?.job_id || "no-job"
  const placementCount = data.calculation?.placements?.length ?? 0
  return `${data.plan_id}:${jobId}:${placementCount}`
}

function getOrCreateVisualizer(): StuffingVisualizer {
  if (!sharedVisualizer) {
    offscreenHost = document.createElement("div")
    offscreenHost.style.position = "fixed"
    offscreenHost.style.left = "-9999px"
    offscreenHost.style.top = "-9999px"
    offscreenHost.style.width = `${THUMBNAIL_WIDTH}px`
    offscreenHost.style.height = `${THUMBNAIL_HEIGHT}px`
    offscreenHost.style.pointerEvents = "none"
    offscreenHost.style.visibility = "hidden"
    document.body.appendChild(offscreenHost)

    sharedVisualizer = new StuffingVisualizer({
      cameraFar: 100000,
    })
    sharedVisualizer.attach(offscreenHost)
  }
  return sharedVisualizer
}

async function defaultGenerateThumbnail(data: StuffingPlanData): Promise<string | null> {
  if (typeof window === "undefined" || typeof document === "undefined") {
    return null
  }

  const visualizer = getOrCreateVisualizer()
  visualizer.loadData(data)

  const scene = visualizer.getScene()
  const camera = visualizer.getCamera()

  const aspect = THUMBNAIL_WIDTH / THUMBNAIL_HEIGHT
  const frustumSize = 20
  camera.left = (-frustumSize * aspect) / 2
  camera.right = (frustumSize * aspect) / 2
  camera.top = frustumSize / 2
  camera.bottom = -frustumSize / 2

  camera.position.set(15, 10, -15)
  camera.lookAt(0, 0, 0)

  const maxDim = Math.max(
    data.container.length_mm,
    data.container.width_mm,
    data.container.height_mm
  ) / 1000
  camera.zoom = 20 / Math.max(maxDim, 1)
  camera.updateProjectionMatrix()

  const canvas = offscreenHost?.querySelector("canvas")
  if (!canvas) {
    return null
  }

  // Trigger a synchronous frame capture via captureScreenshot after updating projection
  // We access the scene/camera render through captureStepScreenshot or direct canvas draw
  visualizer.setStep(visualizer.getMaxStep())

  // Wait one microtask / animation frame so the active loop draws the updated camera matrix
  await new Promise<void>((resolve) => requestAnimationFrame(() => resolve()))

  const dataUrl = visualizer.captureScreenshot()
  return dataUrl
}

/**
 * Renders a 3D thumbnail for a given StuffingPlanData using a single shared
 * off-screen WebGL context. Calls are serialized and cached in memory.
 */
export function renderPlanThumbnail(data: StuffingPlanData): Promise<string | null> {
  const key = getCacheKey(data)
  const cached = thumbnailCache.get(key)
  if (cached) {
    return Promise.resolve(cached)
  }

  const generator = customGenerator ?? defaultGenerateThumbnail

  const task = renderQueue
    .catch(() => null)
    .then(async () => {
      const existing = thumbnailCache.get(key)
      if (existing) {
        return existing
      }
      try {
        const result = await generator(data)
        if (result) {
          thumbnailCache.set(key, result)
        }
        return result
      } catch {
        return null
      }
    })

  renderQueue = task
  return task
}

export function getCachedThumbnail(planId: string): string | null {
  for (const [key, value] of thumbnailCache.entries()) {
    if (key.startsWith(`${planId}:`)) {
      return value
    }
  }
  return null
}

export function clearThumbnailCache(): void {
  thumbnailCache.clear()
  renderQueue = Promise.resolve()
}

export function setThumbnailGeneratorForTests(fn: ThumbnailGeneratorFn | null): void {
  customGenerator = fn
}
