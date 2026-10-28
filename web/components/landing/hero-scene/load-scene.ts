import * as THREE from "three"

import type { Dimensions, SceneBox } from "./to-scene-boxes"

export type LoadScene = {
  /** Shows boxes up to `step`; the fractional part drops the next box in. */
  setProgress: (step: number) => void
  /** Pointer position in [-1, 1] on both axes; (0, 0) re-centres the view. */
  setPointer: (x: number, y: number) => void
  setVisible: (visible: boolean) => void
  resize: (width: number, height: number) => void
  dispose: () => void
}

type LoadSceneOptions = {
  canvas: HTMLCanvasElement
  container: Dimensions
  boxes: SceneBox[]
  /** One colour per item type, indexed by SceneBox.itemIndex. */
  itemColors: string[]
  highlightColor: string
  frameColor: string
  floorColor: string
  onContextLost: () => void
}

const MAX_PIXEL_RATIO = 1.75
const DROP_HEIGHT_M = 1.1
const VIEW_PADDING = 1.12
const POINTER_ORBIT_RAD = 0.07
const CAMERA_DIRECTION = new THREE.Vector3(0.62, 0.52, 0.86).normalize()

export function createLoadScene(options: LoadSceneOptions): LoadScene {
  const { canvas, boxes, itemColors, highlightColor, frameColor, floorColor, onContextLost } = options
  const length = options.container.length / 1000
  const width = options.container.width / 1000
  const height = options.container.height / 1000

  const renderer = new THREE.WebGLRenderer({ canvas, antialias: true, alpha: true })
  renderer.setPixelRatio(Math.min(window.devicePixelRatio, MAX_PIXEL_RATIO))
  renderer.outputColorSpace = THREE.SRGBColorSpace

  const scene = new THREE.Scene()
  scene.add(new THREE.HemisphereLight(0xffffff, 0x8a8f94, 2.2))
  const sun = new THREE.DirectionalLight(0xffffff, 1.6)
  sun.position.set(-3, 6, 4)
  scene.add(sun)

  const camera = new THREE.OrthographicCamera()
  const disposables: Array<{ dispose: () => void }> = []
  const track = <T extends { dispose: () => void }>(resource: T) => {
    disposables.push(resource)
    return resource
  }

  const frame = buildContainerFrame(length, width, height, frameColor, floorColor, track)
  scene.add(frame)

  const unitBox = track(new THREE.BoxGeometry(1, 1, 1))
  const unitEdges = track(new THREE.EdgesGeometry(unitBox))
  const edgeMaterial = track(new THREE.LineBasicMaterial({ color: frameColor, transparent: true, opacity: 0.45 }))
  const itemMaterials = itemColors.map((color) => track(new THREE.MeshLambertMaterial({ color })))
  const highlightMaterial = track(new THREE.MeshLambertMaterial({ color: highlightColor }))

  const meshes = boxes.map((box) => {
    const mesh = new THREE.Mesh(unitBox, itemMaterials[box.itemIndex % itemMaterials.length])
    mesh.scale.set(box.size[0], box.size[1], box.size[2])
    mesh.position.set(box.center[0], box.center[1], box.center[2])
    mesh.add(new THREE.LineSegments(unitEdges, edgeMaterial))
    mesh.visible = false
    scene.add(mesh)
    return mesh
  })

  const target = new THREE.Vector3(0, -height * 0.08, 0)
  const pointer = { x: 0, y: 0 }
  let viewSize = { width: 1, height: 1 }
  let isVisible = true
  let needsRender = true
  let frameId = 0

  const placeCamera = () => {
    const yaw = pointer.x * POINTER_ORBIT_RAD
    const pitch = -pointer.y * POINTER_ORBIT_RAD * 0.6
    const direction = CAMERA_DIRECTION.clone()
      .applyAxisAngle(new THREE.Vector3(0, 1, 0), yaw)
      .add(new THREE.Vector3(0, pitch, 0))
      .normalize()
    camera.position.copy(target).addScaledVector(direction, 20)
    camera.lookAt(target)
  }

  const fitCamera = () => {
    const aspect = viewSize.width / viewSize.height
    const span = Math.hypot(length, width) * VIEW_PADDING
    const halfWidth = span / 2
    const halfHeight = Math.max(halfWidth / aspect, (height * 1.45 * VIEW_PADDING) / 2)
    const fittedHalfWidth = halfHeight * aspect
    camera.left = -fittedHalfWidth
    camera.right = fittedHalfWidth
    camera.top = halfHeight
    camera.bottom = -halfHeight
    camera.near = 0.1
    camera.far = 60
    camera.updateProjectionMatrix()
  }

  const requestRender = () => {
    needsRender = true
    if (frameId || !isVisible) return
    frameId = requestAnimationFrame(() => {
      frameId = 0
      if (!needsRender) return
      needsRender = false
      placeCamera()
      renderer.render(scene, camera)
    })
  }

  const handleContextLost = (event: Event) => {
    event.preventDefault()
    onContextLost()
  }
  canvas.addEventListener("webglcontextlost", handleContextLost)

  return {
    setProgress(step) {
      const settled = Math.floor(step)
      const fraction = step - settled
      meshes.forEach((mesh, index) => {
        const box = boxes[index]
        const isLanding = box.step === settled + 1 && fraction > 0
        const isLatest = box.step === settled && fraction === 0
        mesh.visible = box.step <= settled || isLanding
        mesh.material = isLanding || isLatest ? highlightMaterial : itemMaterials[box.itemIndex % itemMaterials.length]
        const drop = isLanding ? (1 - easeOutCubic(fraction)) * DROP_HEIGHT_M : 0
        mesh.position.y = box.center[1] + drop
      })
      requestRender()
    },
    setPointer(x, y) {
      pointer.x = THREE.MathUtils.clamp(x, -1, 1)
      pointer.y = THREE.MathUtils.clamp(y, -1, 1)
      requestRender()
    },
    setVisible(visible) {
      isVisible = visible
      if (visible) requestRender()
    },
    resize(nextWidth, nextHeight) {
      if (nextWidth === 0 || nextHeight === 0) return
      viewSize = { width: nextWidth, height: nextHeight }
      renderer.setSize(nextWidth, nextHeight, false)
      fitCamera()
      requestRender()
    },
    dispose() {
      cancelAnimationFrame(frameId)
      canvas.removeEventListener("webglcontextlost", handleContextLost)
      disposables.forEach((resource) => resource.dispose())
      renderer.dispose()
    },
  }
}

function easeOutCubic(t: number) {
  return 1 - Math.pow(1 - t, 3)
}

/**
 * Draws the container as a cutaway: a solid floor, faint far walls, and the
 * full edge outline, so the cargo stays readable from the camera angle.
 */
function buildContainerFrame(
  length: number,
  width: number,
  height: number,
  frameColor: string,
  floorColor: string,
  track: <T extends { dispose: () => void }>(resource: T) => T,
) {
  const group = new THREE.Group()

  const outline = new THREE.LineSegments(
    track(new THREE.EdgesGeometry(track(new THREE.BoxGeometry(length, height, width)))),
    track(new THREE.LineBasicMaterial({ color: frameColor })),
  )
  group.add(outline)

  const floor = new THREE.Mesh(
    track(new THREE.BoxGeometry(length, 0.04, width)),
    track(new THREE.MeshLambertMaterial({ color: floorColor })),
  )
  floor.position.y = -height / 2 - 0.02
  group.add(floor)

  const wallMaterial = track(
    new THREE.MeshBasicMaterial({ color: frameColor, transparent: true, opacity: 0.07, side: THREE.DoubleSide, depthWrite: false }),
  )
  const backWall = new THREE.Mesh(track(new THREE.PlaneGeometry(length, height)), wallMaterial)
  backWall.position.z = -width / 2
  group.add(backWall)

  const endWall = new THREE.Mesh(track(new THREE.PlaneGeometry(width, height)), wallMaterial)
  endWall.rotation.y = Math.PI / 2
  endWall.position.x = -length / 2
  group.add(endWall)

  return group
}
