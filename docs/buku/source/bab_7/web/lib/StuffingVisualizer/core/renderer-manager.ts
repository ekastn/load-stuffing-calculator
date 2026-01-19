import { WebGLRenderer, Scene, Camera } from "three";

/**
 * Manages the Three.js WebGL Renderer
 */
export class RendererManager {
    private renderer: WebGLRenderer;
    private animationId: number | null = null;
    private container: HTMLElement | null = null;
    private resizeObserver: ResizeObserver | null = null;
    private onResizeCallback: ((width: number, height: number, aspect: number) => void) | null = null;

    constructor() {
        this.renderer = new WebGLRenderer({
            antialias: true,
            alpha: false,
            preserveDrawingBuffer: true,
        });
        this.renderer.setPixelRatio(Math.min(window.devicePixelRatio, 2));
        
        this.renderer.domElement.style.display = "block";
        this.renderer.domElement.style.width = "100%";
        this.renderer.domElement.style.height = "100%";
    }

    public attach(
        container: HTMLElement, 
        onResize?: (width: number, height: number, aspect: number) => void
    ): void {
        if (this.container === container) return;
        
        this.container = container;
        this.onResizeCallback = onResize ?? null;
        this.container.appendChild(this.renderer.domElement);
        
        if (this.resizeObserver) {
            this.resizeObserver.disconnect();
        }
        this.resizeObserver = new ResizeObserver(() => this.handleResize());
        this.resizeObserver.observe(this.container);
        
        this.handleResize();
    }

    public startLoop(scene: Scene, camera: Camera, onUpdate?: () => void): void {
        if (this.animationId !== null) return;

        const render = () => {
            if (onUpdate) onUpdate();
            this.renderer.render(scene, camera);
            this.animationId = requestAnimationFrame(render);
        };
        render();
    }

    public stopLoop(): void {
        if (this.animationId !== null) {
            cancelAnimationFrame(this.animationId);
            this.animationId = null;
        }
    }

    public handleResize(): void {
        if (!this.container) return;
        
        const width = this.container.clientWidth;
        const height = this.container.clientHeight;
        
        if (width === 0 || height === 0) return;

        const aspect = width / height;
        this.renderer.setSize(width, height);
        
        if (this.onResizeCallback) {
            this.onResizeCallback(width, height, aspect);
        }
    }

    public getRenderer(): WebGLRenderer {
        return this.renderer;
    }

    public dispose(): void {
        this.stopLoop();
        this.resizeObserver?.disconnect();
        this.renderer.domElement.remove();
        this.renderer.dispose();
    }
}
