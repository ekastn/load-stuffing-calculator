"use client"

import { useEffect, useRef, useState } from "react"
import * as THREE from "three"
import { OrbitControls } from "three/addons/controls/OrbitControls.js"

export interface PackedItem {
  itemId: string
  name: string
  dimensions: { length: number; width: number; height: number }
  position: { x: number; y: number; z: number }
  color: string
  step?: number
  rotation?: number
}

interface Canvas3DViewProps {
  items: PackedItem[]
  containerDims: { length: number; width: number; height: number }
  containerName: string
  currentStep?: number
  selectedItemId?: string | null
  onSelect?: (itemId: string | null) => void
}

export function Canvas3DView({ items, containerDims, containerName, currentStep = 0, selectedItemId, onSelect }: Canvas3DViewProps) {
  const containerRef = useRef<HTMLDivElement>(null)
  const sceneRef = useRef<THREE.Scene | null>(null)
  const cameraRef = useRef<THREE.PerspectiveCamera | null>(null)
  const rendererRef = useRef<THREE.WebGLRenderer | null>(null)
  const controlsRef = useRef<any>(null)
  
  // Refs for stable access in event listeners
  const itemGroupsRef = useRef<THREE.Group[]>([])
  const onSelectRef = useRef(onSelect)
  const isDraggingRef = useRef(false)
  const mouseDownPosRef = useRef({ x: 0, y: 0 })

  // Update ref when prop changes
  useEffect(() => {
    onSelectRef.current = onSelect
  }, [onSelect])
  
  const [internalSelected, setInternalSelected] = useState<string | null>(null)
  const activeSelection = selectedItemId !== undefined ? selectedItemId : internalSelected

  const handleSelect = (id: string | null) => {
    setInternalSelected(id)
    if (onSelectRef.current) onSelectRef.current(id)
  }

  // 1. Initialization Effect
  useEffect(() => {
    if (!containerRef.current) return

    // Scene
    const scene = new THREE.Scene()
    scene.background = new THREE.Color(0xffffff) 
    sceneRef.current = scene

    // Camera
    const width = containerRef.current.clientWidth
    const height = containerRef.current.clientHeight
    const camera = new THREE.PerspectiveCamera(25, width / height, 1, 100000)
    camera.position.set(containerDims.length * 2.5, containerDims.height * 2.5, containerDims.width * 2.5)
    cameraRef.current = camera

    // Renderer
    const renderer = new THREE.WebGLRenderer({ antialias: true, alpha: true })
    renderer.setSize(width, height)
    renderer.setPixelRatio(window.devicePixelRatio)
    renderer.shadowMap.enabled = true 
    renderer.shadowMap.type = THREE.PCFSoftShadowMap
    containerRef.current.appendChild(renderer.domElement)
    rendererRef.current = renderer

    // Controls
    const target = new THREE.Vector3(containerDims.length / 2, containerDims.height / 2, containerDims.width / 2)
    const controls = new OrbitControls(camera, renderer.domElement)
    controls.target.copy(target)
    controls.enableDamping = true
    controls.dampingFactor = 0.05
    controlsRef.current = controls

    // Lighting
    const ambientLight = new THREE.AmbientLight(0xffffff, 0.9)
    scene.add(ambientLight)

    const dirLight = new THREE.DirectionalLight(0xffffff, 0.3)
    dirLight.position.set(containerDims.length, containerDims.height * 4, containerDims.width * 2)
    dirLight.castShadow = true
    scene.add(dirLight)

    // Floor
    const floorHeight = 40
    const floorGeo = new THREE.BoxGeometry(containerDims.length, floorHeight, containerDims.width)
    const floorMat = new THREE.MeshLambertMaterial({ color: 0xE3B448 }) 
    const floorMesh = new THREE.Mesh(floorGeo, floorMat)
    floorMesh.position.set(containerDims.length / 2, -floorHeight / 2, containerDims.width / 2)
    floorMesh.receiveShadow = true
    scene.add(floorMesh)

    // Container Wireframe
    const boxGeo = new THREE.BoxGeometry(containerDims.length, containerDims.height, containerDims.width)
    const edges = new THREE.EdgesGeometry(boxGeo)
    const lineMat = new THREE.LineBasicMaterial({ color: 0xdddddd, linewidth: 1 })
    const wireframe = new THREE.LineSegments(edges, lineMat)
    wireframe.position.copy(target)
    scene.add(wireframe)

    // Items
    itemGroupsRef.current = [] 
    
    items.forEach((item, idx) => {
      const group = new THREE.Group()
      group.userData = { 
          itemId: item.itemId, 
          step: item.step !== undefined ? item.step : idx + 1 
      }

      let { length, width, height } = item.dimensions
      const rot = item.rotation || 0
      if (rot === 1 || rot === 3) [length, width] = [width, length]
      
      const geo = new THREE.BoxGeometry(length, height, width)
      
      const mat = new THREE.MeshLambertMaterial({
        color: new THREE.Color(item.color),
        transparent: true,
        opacity: 1.0,
      })

      const mesh = new THREE.Mesh(geo, mat)
      mesh.castShadow = true
      mesh.receiveShadow = true
      mesh.userData.isItemMesh = true
      
      mesh.position.set(
        item.position.x + length / 2,
        item.position.z + height / 2,
        item.position.y + width / 2
      )

      const itemEdges = new THREE.EdgesGeometry(geo)
      const itemEdgeMat = new THREE.LineBasicMaterial({ color: 0xffffff, transparent: true, opacity: 0.4 })
      const itemWireframe = new THREE.LineSegments(itemEdges, itemEdgeMat)
      itemWireframe.position.copy(mesh.position)

      group.add(mesh)
      group.add(itemWireframe)
      scene.add(group)
      
      itemGroupsRef.current.push(group)
    })

    // Interaction Handlers
    const raycaster = new THREE.Raycaster()
    const mouse = new THREE.Vector2()

    const onPointerDown = (e: MouseEvent) => {
        isDraggingRef.current = false
        mouseDownPosRef.current = { x: e.clientX, y: e.clientY }
    }

    const onPointerMove = (e: MouseEvent) => {
        const dx = e.clientX - mouseDownPosRef.current.x
        const dy = e.clientY - mouseDownPosRef.current.y
        if (Math.abs(dx) > 2 || Math.abs(dy) > 2) {
            isDraggingRef.current = true
        }
    }

    const onPointerUp = (e: MouseEvent) => {
        if (isDraggingRef.current) return 

        const rect = containerRef.current!.getBoundingClientRect()
        mouse.x = ((e.clientX - rect.left) / rect.width) * 2 - 1
        mouse.y = -((e.clientY - rect.top) / rect.height) * 2 + 1

        raycaster.setFromCamera(mouse, camera)
        
        // Find visible item meshes
        const visibleGroups = itemGroupsRef.current.filter(g => g.visible)
        const visibleMeshes: THREE.Object3D[] = []
        visibleGroups.forEach(g => {
            g.children.forEach(c => {
                if (c.userData.isItemMesh) visibleMeshes.push(c)
            })
        })

        const intersects = raycaster.intersectObjects(visibleMeshes)
        if (intersects.length > 0) {
            const group = intersects[0].object.parent
            if (group) handleSelect(group.userData.itemId)
        } else {
            handleSelect(null)
        }
    }

    const domEl = renderer.domElement
    domEl.addEventListener("mousedown", onPointerDown)
    domEl.addEventListener("mousemove", onPointerMove)
    domEl.addEventListener("mouseup", onPointerUp)

    let animationId: number
    const animate = () => {
      animationId = requestAnimationFrame(animate)
      controls.update()
      renderer.render(scene, camera)
    }
    animate()

    const handleResize = () => {
      if (!containerRef.current) return
      const w = containerRef.current.clientWidth
      const h = containerRef.current.clientHeight
      camera.aspect = w / h
      camera.updateProjectionMatrix()
      renderer.setSize(w, h)
    }
    window.addEventListener("resize", handleResize)

    return () => {
        window.removeEventListener("resize", handleResize)
        domEl.removeEventListener("mousedown", onPointerDown)
        domEl.removeEventListener("mousemove", onPointerMove)
        domEl.removeEventListener("mouseup", onPointerUp)
        cancelAnimationFrame(animationId)
        controls.dispose()
        containerRef.current?.removeChild(renderer.domElement)
        renderer.dispose()
    }
  }, [items, containerDims]) 

  // 2. Step Visibility Effect
  useEffect(() => {
    itemGroupsRef.current.forEach(group => {
        const itemStep = group.userData.step
        group.visible = itemStep <= currentStep
    })
  }, [currentStep])

  // 3. Selection Effect
  useEffect(() => {
    itemGroupsRef.current.forEach(group => {
        const mesh = group.children.find(c => c.userData.isItemMesh) as THREE.Mesh
        if (mesh && mesh.material instanceof THREE.MeshLambertMaterial) {
            const isSel = group.userData.itemId === activeSelection
            if (activeSelection) {
                mesh.material.opacity = isSel ? 1.0 : 0.2
                mesh.material.emissive.setHex(isSel ? 0x222222 : 0x000000)
            } else {
                mesh.material.opacity = 1.0
                mesh.material.emissive.setHex(0x000000)
            }
        }
    })
  }, [activeSelection])

  return (
    <div ref={containerRef} className="w-full h-full bg-white overflow-hidden cursor-grab active:cursor-grabbing relative" />
  )
}