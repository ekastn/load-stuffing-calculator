"use client"

import React, { useState, useRef, useCallback } from "react"
import { X, Upload, ImagePlus } from "lucide-react"
import { cn } from "@/lib/utils"

const MAX_IMAGES = 10

export interface UploadedImage {
  file: File
  preview: string
}

interface ImageUploadProps {
  images: UploadedImage[]
  onChange: (images: UploadedImage[]) => void
  maxImages?: number
  disabled?: boolean
}

export function ImageUpload({
  images,
  onChange,
  maxImages = MAX_IMAGES,
  disabled = false,
}: ImageUploadProps) {
  const [isDragging, setIsDragging] = useState(false)
  const inputRef = useRef<HTMLInputElement>(null)

  const addFiles = useCallback(
    (files: FileList | File[]) => {
      const validFiles = Array.from(files).filter((f) => f.type.startsWith("image/"))
      if (validFiles.length === 0) return

      const remaining = maxImages - images.length
      if (remaining <= 0) return

      const filesToAdd = validFiles.slice(0, remaining)
      const newImages: UploadedImage[] = filesToAdd.map((file) => ({
        file,
        preview: URL.createObjectURL(file),
      }))

      onChange([...images, ...newImages])
    },
    [images, onChange, maxImages]
  )

  const removeImage = useCallback(
    (index: number) => {
      const updated = [...images]
      URL.revokeObjectURL(updated[index].preview)
      updated.splice(index, 1)
      onChange(updated)
    },
    [images, onChange]
  )

  const handleDragOver = (e: React.DragEvent) => {
    e.preventDefault()
    if (!disabled) setIsDragging(true)
  }

  const handleDragLeave = (e: React.DragEvent) => {
    e.preventDefault()
    setIsDragging(false)
  }

  const handleDrop = (e: React.DragEvent) => {
    e.preventDefault()
    setIsDragging(false)
    if (disabled) return
    addFiles(e.dataTransfer.files)
  }

  const handleClick = () => {
    if (disabled || images.length >= maxImages) return
    inputRef.current?.click()
  }

  const handleFileChange = (e: React.ChangeEvent<HTMLInputElement>) => {
    if (e.target.files) {
      addFiles(e.target.files)
      e.target.value = ""
    }
  }

  const isFull = images.length >= maxImages

  return (
    <div className="space-y-3">
      <label className="text-sm font-medium">
        Product Images
        <span className="text-muted-foreground font-normal ml-1">
          ({images.length}/{maxImages})
        </span>
      </label>

      {images.length > 0 && (
        <div className="grid grid-cols-5 gap-3">
          {images.map((img, index) => (
            <div
              key={index}
              className="group relative aspect-square rounded-lg border border-border/50 overflow-hidden bg-muted/30"
            >
              <img
                src={img.preview}
                alt={`Product ${index + 1}`}
                className="h-full w-full object-cover"
              />
              {!disabled && (
                <button
                  type="button"
                  onClick={() => removeImage(index)}
                  className="absolute right-1 top-1 flex h-5 w-5 items-center justify-center rounded-full bg-destructive text-destructive-foreground opacity-0 transition-opacity group-hover:opacity-100"
                >
                  <X className="h-3 w-3" />
                </button>
              )}
              <span className="absolute bottom-0 left-0 right-0 bg-black/50 px-1 py-0.5 text-center text-[10px] text-white">
                {index + 1}
              </span>
            </div>
          ))}
        </div>
      )}

      {!isFull && (
        <div
          onClick={handleClick}
          onDragOver={handleDragOver}
          onDragLeave={handleDragLeave}
          onDrop={handleDrop}
          className={cn(
            "flex flex-col items-center justify-center rounded-lg border-2 border-dashed p-6 transition-colors cursor-pointer",
            isDragging
              ? "border-primary bg-primary/5"
              : "border-border/50 hover:border-primary/50 hover:bg-muted/30",
            disabled && "pointer-events-none opacity-50"
          )}
        >
          <div className="flex flex-col items-center gap-2 text-center">
            {images.length === 0 ? (
              <ImagePlus className="h-8 w-8 text-muted-foreground/50" />
            ) : (
              <Upload className="h-6 w-6 text-muted-foreground/50" />
            )}
            <div>
              <p className="text-sm font-medium text-foreground">
                {isDragging ? "Drop images here" : "Drag & drop images or click to browse"}
              </p>
              <p className="text-xs text-muted-foreground mt-1">
                PNG, JPG, WEBP up to 5MB each
              </p>
            </div>
          </div>
        </div>
      )}

      <input
        ref={inputRef}
        type="file"
        accept="image/*"
        multiple
        onChange={handleFileChange}
        className="hidden"
        disabled={disabled}
      />
    </div>
  )
}
