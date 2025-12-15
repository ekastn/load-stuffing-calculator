"use client"

import { useEffect, useRef, useState } from "react"
import * as THREE from "three"
import type { PackedItem } from "@/lib/bin-packing"

interface Canvas3DViewProps {
  items: PackedItem[]
  containerDims: { length: number; width: number; height: number }
  containerName: string
}

export function Canvas3DView({ items, containerDims, containerName }: Canvas3DViewProps) {
  const containerRef = useRef<HTMLDivElement>(null)
  const sceneRef = useRef<THREE.Scene | null>(null)
  const cameraRef = useRef<THREE.PerspectiveCamera | null>(null)
  const rendererRef = useRef<THREE.WebGLRenderer | null>(null)
  const controlsRef = useRef<any>(null)
  const [selectedItem, setSelectedItem] = useState<string | null>(null)

  useEffect(() => {
    if (!containerRef.current) return

    // Scene setup
    const scene = new THREE.Scene()
    scene.background = new THREE.Color(0x1a1a2e)
    sceneRef.current = scene

    // Camera setup
    const width = containerRef.current.clientWidth
    const height = containerRef.current.clientHeight
    const camera = new THREE.PerspectiveCamera(75, width / height, 0.1, 10000)
    camera.position.set(containerDims.length * 0.7, containerDims.height * 0.7, containerDims.width * 0.7)
    camera.lookAt(containerDims.length / 2, containerDims.height / 2, containerDims.width / 2)
    cameraRef.current = camera

    // Renderer setup
    const renderer = new THREE.WebGLRenderer({ antialias: true, alpha: true })
    renderer.setSize(width, height)
    renderer.setPixelRatio(window.devicePixelRatio)
    renderer.shadowMap.enabled = true
    containerRef.current.appendChild(renderer.domElement)
    rendererRef.current = renderer

    class SimpleOrbitControls {
      private camera: THREE.PerspectiveCamera
      private target: THREE.Vector3
      private radius: number
      private theta = 0
      private phi: number = Math.PI / 3
      private isDragging = false
      private previousMousePosition = { x: 0, y: 0 }
      private zoomSpeed = 1.2

      constructor(camera: THREE.PerspectiveCamera, container: HTMLElement, target: THREE.Vector3) {
        this.camera = camera
        this.target = target

        const pos = camera.position
        this.radius = Math.sqrt(pos.x ** 2 + pos.y ** 2 + pos.z ** 2)

        container.addEventListener("mousedown", this.onMouseDown.bind(this))
        container.addEventListener("mousemove", this.onMouseMove.bind(this))
        container.addEventListener("mouseup", this.onMouseUp.bind(this))
        container.addEventListener("wheel", this.onMouseWheel.bind(this), { passive: false })
        container.addEventListener("touchstart", this.onTouchStart.bind(this))
        container.addEventListener("touchmove", this.onTouchMove.bind(this))
        container.addEventListener("touchend", this.onTouchEnd.bind(this))
      }

      private onMouseDown(e: MouseEvent) {
        this.isDragging = true
        this.previousMousePosition = { x: e.clientX, y: e.clientY }
      }

      private onMouseMove(e: MouseEvent) {
        if (!this.isDragging) return

        const deltaX = e.clientX - this.previousMousePosition.x
        const deltaY = e.clientY - this.previousMousePosition.y

        this.theta -= deltaX * 0.01
        this.phi -= deltaY * 0.01

        this.phi = Math.max(0.1, Math.min(Math.PI - 0.1, this.phi))

        this.updateCameraPosition()
        this.previousMousePosition = { x: e.clientX, y: e.clientY }
      }

      private onMouseUp() {
        this.isDragging = false
      }

      private onMouseWheel(e: WheelEvent) {
        e.preventDefault()
        this.radius *= e.deltaY > 0 ? this.zoomSpeed : 1 / this.zoomSpeed
        this.radius = Math.max(10, Math.min(500, this.radius))
        this.updateCameraPosition()
      }

      private onTouchStart(e: TouchEvent) {
        if (e.touches.length === 1) {
          this.isDragging = true
          this.previousMousePosition = { x: e.touches[0].clientX, y: e.touches[0].clientY }
        }
      }

      private onTouchMove(e: TouchEvent) {
        if (!this.isDragging || e.touches.length !== 1) return

        const deltaX = e.touches[0].clientX - this.previousMousePosition.x
        const deltaY = e.touches[0].clientY - this.previousMousePosition.y

        this.theta -= deltaX * 0.01
        this.phi -= deltaY * 0.01

        this.phi = Math.max(0.1, Math.min(Math.PI - 0.1, this.phi))

        this.updateCameraPosition()
        this.previousMousePosition = { x: e.touches[0].clientX, y: e.touches[0].clientY }
      }

      private onTouchEnd() {
        this.isDragging = false
      }

      private updateCameraPosition() {
        this.camera.position.x = this.target.x + this.radius * Math.sin(this.phi) * Math.cos(this.theta)
        this.camera.position.y = this.target.y + this.radius * Math.cos(this.phi)
        this.camera.position.z = this.target.z + this.radius * Math.sin(this.phi) * Math.sin(this.theta)

        this.camera.lookAt(this.target)
      }

      dispose() {
        const container = this.camera.userData.container as HTMLElement
        if (container) {
          container.removeEventListener("mousedown", this.onMouseDown.bind(this))
          container.removeEventListener("mousemove", this.onMouseMove.bind(this))
          container.removeEventListener("mouseup", this.onMouseUp.bind(this))
          container.removeEventListener("wheel", this.onMouseWheel.bind(this))
        }
      }
    }

    const controls = new SimpleOrbitControls(
      camera,
      renderer.domElement,
      new THREE.Vector3(containerDims.length / 2, containerDims.height / 2, containerDims.width / 2),
    )
    controlsRef.current = controls

    const ambientLight = new THREE.AmbientLight(0xffffff, 0.6)
    scene.add(ambientLight)

    const directionalLight = new THREE.DirectionalLight(0xffffff, 0.8)
    directionalLight.position.set(100, 150, 100)
    directionalLight.castShadow = true
    directionalLight.shadow.mapSize.width = 2048
    directionalLight.shadow.mapSize.height = 2048
    directionalLight.shadow.camera.far = 500
    scene.add(directionalLight)

    const containerGeometry = new THREE.BoxGeometry(containerDims.length, containerDims.height, containerDims.width)
    const containerMaterial = new THREE.LineBasicMaterial({ color: 0x888888, linewidth: 2 })
    const containerEdges = new THREE.EdgesGeometry(containerGeometry)
    const containerWireframe = new THREE.LineSegments(containerEdges, containerMaterial)
    containerWireframe.position.set(containerDims.length / 2, containerDims.height / 2, containerDims.width / 2)
    scene.add(containerWireframe)

    const floorGeometry = new THREE.PlaneGeometry(containerDims.length * 1.2, containerDims.width * 1.2)
    const floorMaterial = new THREE.MeshStandardMaterial({
      color: 0x2a2a3e,
      metalness: 0.3,
      roughness: 0.4,
    })
    const floor = new THREE.Mesh(floorGeometry, floorMaterial)
    floor.rotation.x = -Math.PI / 2
    floor.receiveShadow = true
    scene.add(floor)

    const itemMeshes: { [key: string]: THREE.Group } = {}

    items.forEach((item, idx) => {
      const group = new THREE.Group()
      group.userData.itemId = item.itemId

      // Item box
      const geometry = new THREE.BoxGeometry(item.dimensions.length, item.dimensions.height, item.dimensions.width)

      // Parse color hex to THREE.Color
      const color = new THREE.Color(item.color)

      const material = new THREE.MeshStandardMaterial({
        color: color,
        metalness: 0.1,
        roughness: 0.6,
        emissive: selectedItem === item.itemId ? new THREE.Color(0xffff00) : new THREE.Color(0x000000),
      })

      const mesh = new THREE.Mesh(geometry, material)
      mesh.castShadow = true
      mesh.receiveShadow = true
      mesh.userData.itemId = item.itemId

      // Position box
      mesh.position.set(
        item.position.x + item.dimensions.length / 2,
        item.position.z + item.dimensions.height / 2,
        item.position.y + item.dimensions.width / 2,
      )

      // Add edges
      const edges = new THREE.EdgesGeometry(geometry)
      const edgeMaterial = new THREE.LineBasicMaterial({ color: 0x000000, linewidth: 1 })
      const wireframe = new THREE.LineSegments(edges, edgeMaterial)
      wireframe.position.copy(mesh.position)

      group.add(mesh)
      group.add(wireframe)

      const canvas = document.createElement("canvas")
      canvas.width = 256
      canvas.height = 128
      const ctx = canvas.getContext("2d")
      if (ctx) {
        ctx.fillStyle = "white"
        ctx.font = "bold 20px Arial"
        ctx.textAlign = "center"
        ctx.fillText(item.name, 128, 40)
        ctx.font = "14px Arial"
        ctx.fillText(`${idx + 1}/${items.length}`, 128, 80)
      }

      const texture = new THREE.CanvasTexture(canvas)
      const labelGeometry = new THREE.PlaneGeometry(item.dimensions.length * 1.2, 10)
      const labelMaterial = new THREE.MeshStandardMaterial({
        map: texture,
        emissive: 0xffffff,
        emissiveIntensity: 0.8,
      })
      const labelMesh = new THREE.Mesh(labelGeometry, labelMaterial)
      labelMesh.position.set(
        item.position.x + item.dimensions.length / 2,
        item.position.z + item.dimensions.height + 5,
        item.position.y + item.dimensions.width / 2,
      )
      labelMesh.lookAt(camera.position)

      group.add(labelMesh)
      itemMeshes[item.itemId] = group
      scene.add(group)
    })

    const raycaster = new THREE.Raycaster()
    const mouse = new THREE.Vector2()

    const onMouseMove = (event: MouseEvent) => {
      if (!containerRef.current) return

      const rect = containerRef.current.getBoundingClientRect()
      mouse.x = ((event.clientX - rect.left) / width) * 2 - 1
      mouse.y = -((event.clientY - rect.top) / height) * 2 + 1

      raycaster.setFromCamera(mouse, camera)

      // Update all item meshes
      Object.values(itemMeshes).forEach((group) => {
        const meshes = group.children.filter((child) => child instanceof THREE.Mesh && !(child as any).isLabel)
        const intersects = raycaster.intersectObjects(meshes)

        if (intersects.length > 0) {
          const material = (meshes[0] as THREE.Mesh).material as THREE.MeshStandardMaterial
          material.emissive.setHex(0x444444)
          material.emissiveIntensity = 0.3
        } else {
          const material = (meshes[0] as THREE.Mesh).material as THREE.MeshStandardMaterial
          if (selectedItem === group.userData.itemId) {
            material.emissive.setHex(0xffff00)
            material.emissiveIntensity = 0.8
          } else {
            material.emissive.setHex(0x000000)
          }
        }
      })
    }

    const onClick = (event: MouseEvent) => {
      if (!containerRef.current) return

      const rect = containerRef.current.getBoundingClientRect()
      mouse.x = ((event.clientX - rect.left) / width) * 2 - 1
      mouse.y = -((event.clientY - rect.top) / height) * 2 + 1

      raycaster.setFromCamera(mouse, camera)

      let clicked = false
      Object.entries(itemMeshes).forEach(([itemId, group]) => {
        const meshes = group.children.filter((child) => child instanceof THREE.Mesh && !(child as any).isLabel)
        const intersects = raycaster.intersectObjects(meshes)

        if (intersects.length > 0 && !clicked) {
          setSelectedItem(itemId)
          clicked = true
        }
      })

      if (!clicked) {
        setSelectedItem(null)
      }
    }

    renderer.domElement.addEventListener("mousemove", onMouseMove)
    renderer.domElement.addEventListener("click", onClick)

    let animationId: number
    const animate = () => {
      animationId = requestAnimationFrame(animate)

      renderer.render(scene, camera)
    }

    animate()

    // Handle window resize
    const handleResize = () => {
      if (!containerRef.current) return
      const newWidth = containerRef.current.clientWidth
      const newHeight = containerRef.current.clientHeight

      camera.aspect = newWidth / newHeight
      camera.updateProjectionMatrix()
      renderer.setSize(newWidth, newHeight)
    }

    window.addEventListener("resize", handleResize)

    // Cleanup
    return () => {
      window.removeEventListener("resize", handleResize)
      renderer.domElement.removeEventListener("mousemove", onMouseMove)
      renderer.domElement.removeEventListener("click", onClick)
      cancelAnimationFrame(animationId)
      containerRef.current?.removeChild(renderer.domElement)
      if (controlsRef.current) {
        controlsRef.current.dispose()
      }
      renderer.dispose()
    }
  }, [items, containerDims])

  // Update selection highlighting
  useEffect(() => {
    if (!sceneRef.current) return

    sceneRef.current.traverse((child) => {
      if (child instanceof THREE.Mesh && child.material instanceof THREE.MeshStandardMaterial) {
        const itemId = (child.parent as any)?.userData?.itemId
        if (itemId === selectedItem) {
          child.material.emissive.setHex(0xffff00)
          child.material.emissiveIntensity = 0.8
        } else if (itemId) {
          child.material.emissive.setHex(0x000000)
          child.material.emissiveIntensity = 0
        }
      }
    })
  }, [selectedItem])

  return (
    <div
      ref={containerRef}
      className="w-full rounded-lg border border-border bg-background overflow-hidden cursor-grab active:cursor-grabbing"
      style={{ height: "600px" }}
    />
  )
}
