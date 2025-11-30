import React, { useState } from 'react';
import { 
  Plus, 
  Upload, 
  Download, 
  Zap, 
  Edit2, 
  Trash2, 
  Settings, 
  HelpCircle, 
  Box, 
  ShoppingBag, 
  Package, 
  Move3d,
  ArrowRight,
  ArrowLeft,
  Copy,
  X,
  Cylinder,
  Circle,
  Component
} from 'lucide-react';
import { PieChart, Pie, Cell, ResponsiveContainer } from 'recharts';

// --- Data & Constants ---

const PRODUCT_DATA = [
  { id: 1, type: 'box', name: 'Boxes 1', l: 500, w: 400, h: 300, weight: 10, qty: 80, color: '#22c55e' },
  { id: 2, type: 'sack', name: 'Sacks', l: 1000, w: 450, h: 300, weight: 45, qty: 100, color: '#d946ef' },
  { id: 3, type: 'bigbag', name: 'Big bags', l: 1000, w: 1000, h: 1000, weight: 900, qty: 10, color: '#3b82f6' },
];

const CHART_DATA = [
  { name: 'Big bags', value: 10, color: '#3b82f6' },
  { name: 'Sacks', value: 100, color: '#d946ef' },
  { name: 'Boxes 1', value: 80, color: '#22c55e' },
];

const CARGO_TYPES = [
  { id: 'box', label: 'box', icon: Box },
  { id: 'bigbag', label: 'bigbags', icon: ShoppingBag },
  { id: 'sack', label: 'sacks', icon: Package }, // Used Package as generic sack
  { id: 'barrel', label: 'barrels', icon: Cylinder },
  { id: 'roll', label: 'roll', icon: Circle }, // Simplified representation
  { id: 'pipes', label: 'pipes', icon: Component }, // Simplified representation
  { id: 'bulk', label: 'bulk', icon: Move3d },
];

// --- Sub-Components ---

const TypeIcon = ({ type, color }) => {
  const style = { color: color };
  if (type === 'box') return <Box className="w-5 h-5" style={style} />;
  if (type === 'sack') return <ShoppingBag className="w-5 h-5" style={style} />;
  return <Package className="w-5 h-5" style={style} />;
};

const ColorDot = ({ color }) => (
  <div className="w-6 h-6 rounded-full border border-gray-200 shadow-sm" style={{ backgroundColor: color }}></div>
);

// SVG Illustration for Empty Container Outline
const ContainerOutline = () => (
  <svg viewBox="0 0 200 120" className="w-48 h-32 mx-auto opacity-70">
    <path d="M10 30 L60 10 L190 40 L140 60 Z" fill="none" stroke="#64748b" strokeWidth="2" />
    <path d="M10 30 L10 90 L60 70 L60 10" fill="none" stroke="#64748b" strokeWidth="2" />
    <path d="M190 40 L190 100 L140 120 L140 60" fill="none" stroke="#64748b" strokeWidth="2" />
    <path d="M10 90 L140 120" fill="none" stroke="#64748b" strokeWidth="2" />
    <path d="M60 70 L190 100" fill="none" stroke="#64748b" strokeWidth="2" />
    {/* Ribs */}
    <path d="M20 26 L20 86" stroke="#cbd5e1" strokeWidth="1" />
    <path d="M30 22 L30 82" stroke="#cbd5e1" strokeWidth="1" />
    <path d="M40 18 L40 78" stroke="#cbd5e1" strokeWidth="1" />
    <path d="M50 14 L50 74" stroke="#cbd5e1" strokeWidth="1" />
  </svg>
);

// SVG Illustration for Filled Container (Stylized)
const ContainerFilled = () => (
  <svg viewBox="0 0 200 120" className="w-full h-48 drop-shadow-lg">
     {/* Base Floor */}
    <path d="M10 90 L140 120 L190 100 L60 70 Z" fill="#94a3b8" />
    
    {/* Stacked Boxes - Back Row (Blue) */}
    <path d="M60 70 L110 82 L110 52 L60 40 Z" fill="#3b82f6" /> {/* Side */}
    <path d="M60 40 L110 52 L150 42 L100 30 Z" fill="#60a5fa" /> {/* Top */}
    <path d="M110 82 L150 72 L150 42 L110 52 Z" fill="#2563eb" /> {/* Front */}

    {/* Stacked Boxes - Middle Row (Pink) */}
    <path d="M35 65 L85 77 L85 47 L35 35 Z" fill="#d946ef" /> 
    <path d="M35 35 L85 47 L125 37 L75 25 Z" fill="#f0abfc" /> 
    <path d="M85 77 L125 67 L125 37 L85 47 Z" fill="#c026d3" />

    {/* Stacked Boxes - Front Row (Green) */}
    <path d="M10 90 L60 102 L60 72 L10 60 Z" fill="#22c55e" />
    <path d="M10 60 L60 72 L100 62 L50 50 Z" fill="#4ade80" />
    <path d="M60 102 L100 92 L100 62 L60 72 Z" fill="#16a34a" />
    
    {/* Container Frame Overlay (Partial) */}
    <path d="M10 30 L10 90" stroke="#64748b" strokeWidth="2" />
    <path d="M10 90 L140 120" stroke="#64748b" strokeWidth="2" />
    <path d="M140 120 L140 60" stroke="#64748b" strokeWidth="2" />
  </svg>
);

// --- New Add Product Modal Component ---
const AddProductModal = ({ isOpen, onClose }) => {
  const [selectedType, setSelectedType] = useState('box');

  if (!isOpen) return null;

  // Render Cargo Visual based on type (Simplified SVGs)
  const renderCargoVisual = () => {
    return (
      <div className="w-full h-48 flex items-center justify-center relative">
        <svg viewBox="0 0 200 160" className="w-48 h-48 drop-shadow-xl text-blue-100 fill-blue-50 stroke-blue-500">
           {selectedType === 'box' && (
              <>
                 <path d="M50 60 L100 40 L150 60 L100 80 Z" fill="#dbeafe" strokeWidth="2" stroke="#3b82f6"/>
                 <path d="M50 60 L50 120 L100 140 L100 80" fill="#bfdbfe" strokeWidth="2" stroke="#3b82f6"/>
                 <path d="M150 60 L150 120 L100 140" fill="#eff6ff" strokeWidth="2" stroke="#3b82f6"/>
                 {/* Dimension lines */}
                 <text x="20" y="100" className="text-[10px] fill-gray-400 font-sans">Height</text>
                 <line x1="40" y1="60" x2="40" y2="120" stroke="#cbd5e1" strokeWidth="1" />
                 
                 <text x="60" y="150" className="text-[10px] fill-gray-400 font-sans">Width</text>
                 <line x1="50" y1="130" x2="100" y2="150" stroke="#cbd5e1" strokeWidth="1" />
                 
                 <text x="160" y="100" className="text-[10px] fill-gray-400 font-sans">Length</text>
                 <line x1="100" y1="150" x2="150" y2="130" stroke="#cbd5e1" strokeWidth="1" />
              </>
           )}
           {selectedType === 'bigbag' && (
              <>
                <path d="M60 120 C60 135, 140 135, 140 120 L130 60 L70 60 Z" fill="#dbeafe" stroke="#3b82f6" strokeWidth="2"/>
                <path d="M70 60 L60 30 M130 60 L140 30" stroke="#3b82f6" strokeWidth="2" fill="none"/>
                <ellipse cx="100" cy="60" rx="30" ry="10" fill="#eff6ff" stroke="#3b82f6" strokeWidth="2"/>
              </>
           )}
           {(selectedType === 'sack') && (
               <>
                 <path d="M40 80 Q100 20 160 80 Q140 140 100 140 Q60 140 40 80 Z" fill="#dbeafe" stroke="#3b82f6" strokeWidth="2"/>
                 <path d="M40 80 Q100 100 160 80" fill="none" stroke="#3b82f6" strokeWidth="1" strokeDasharray="4"/>
               </>
           )}
            {selectedType === 'barrel' && (
              <>
                <ellipse cx="100" cy="40" rx="40" ry="15" fill="#eff6ff" stroke="#3b82f6" strokeWidth="2"/>
                <path d="M60 40 L60 120 C60 135, 140 135, 140 120 L140 40" fill="#dbeafe" stroke="#3b82f6" strokeWidth="2"/>
                <path d="M60 120 C60 135, 140 135, 140 120" fill="none" stroke="#3b82f6" strokeWidth="2"/>
              </>
           )}
           {/* Fallback for others */}
           {['roll', 'pipes', 'bulk'].includes(selectedType) && (
              <text x="50%" y="50%" textAnchor="middle" fill="#94a3b8">Visual Preview</text>
           )}
        </svg>
      </div>
    )
  }

  return (
    <div className="fixed inset-0 z-50 flex items-center justify-center p-4 bg-black/40 backdrop-blur-sm animate-in fade-in duration-200">
      <div className="bg-white rounded-2xl shadow-2xl w-full max-w-4xl max-h-[90vh] overflow-y-auto overflow-x-hidden flex flex-col">
        
        {/* Modal Header is implied by the layout */}
        
        <div className="p-8 space-y-8">
            {/* --- SECTION 1: CARGO TYPE --- */}
            <div>
                <h3 className="text-blue-900 font-extrabold text-sm uppercase tracking-wide mb-4">1. Select Cargo Type</h3>
                <div className="grid grid-cols-4 md:grid-cols-7 gap-3">
                    {CARGO_TYPES.map((type) => (
                        <button 
                            key={type.id}
                            onClick={() => setSelectedType(type.id)}
                            className={`flex flex-col items-center justify-center py-4 px-2 rounded-xl border-2 transition-all ${selectedType === type.id ? 'border-blue-500 bg-blue-50 text-blue-600' : 'border-gray-100 hover:border-blue-200 text-gray-400'}`}
                        >
                            <type.icon size={28} strokeWidth={1.5} className="mb-2" />
                            <span className="text-xs font-medium capitalize">{type.label}</span>
                        </button>
                    ))}
                </div>
            </div>

            {/* --- SECTION 2: DIMENSIONS --- */}
            <div>
                <h3 className="text-blue-900 font-extrabold text-sm uppercase tracking-wide mb-4">2. Select Cargo Dimensions</h3>
                <div className="flex flex-col md:flex-row gap-8 items-start">
                    {/* Visual */}
                    <div className="w-full md:w-1/3 bg-gray-50 rounded-xl border border-gray-100 p-4 flex items-center justify-center">
                        {renderCargoVisual()}
                    </div>
                    
                    {/* Inputs */}
                    <div className="w-full md:w-2/3 grid grid-cols-1 md:grid-cols-2 gap-x-6 gap-y-4">
                        <div className="col-span-1 md:col-span-2 flex gap-4">
                           <div className="flex-1">
                               <label className="block text-xs font-medium text-blue-300 mb-1">Product Name</label>
                               <input type="text" defaultValue="new product" className="w-full rounded-full border border-gray-300 px-4 py-2 text-sm focus:ring-2 focus:ring-blue-400 outline-none" />
                           </div>
                           <div className="w-1/3">
                               <label className="block text-xs font-medium text-blue-300 mb-1">Color</label>
                               <div className="h-[38px] w-full bg-[#7c3aed] rounded-full cursor-pointer hover:opacity-90 shadow-sm border border-gray-200"></div>
                           </div>
                        </div>

                        <div>
                            <label className="block text-xs font-medium text-blue-300 mb-1">Length</label>
                            <div className="relative">
                                <input type="number" defaultValue="100" className="w-full rounded-full border border-gray-300 px-4 py-2 text-sm focus:ring-2 focus:ring-blue-400 outline-none" />
                                <span className="absolute right-4 top-1/2 -translate-y-1/2 text-xs text-blue-300">mm</span>
                            </div>
                        </div>
                        <div>
                            <label className="block text-xs font-medium text-blue-300 mb-1">Width</label>
                            <div className="relative">
                                <input type="number" defaultValue="100" className="w-full rounded-full border border-gray-300 px-4 py-2 text-sm focus:ring-2 focus:ring-blue-400 outline-none" />
                                <span className="absolute right-4 top-1/2 -translate-y-1/2 text-xs text-blue-300">mm</span>
                            </div>
                        </div>
                        <div>
                            <label className="block text-xs font-medium text-blue-300 mb-1">Height</label>
                            <div className="relative">
                                <input type="number" defaultValue="100" className="w-full rounded-full border border-gray-300 px-4 py-2 text-sm focus:ring-2 focus:ring-blue-400 outline-none" />
                                <span className="absolute right-4 top-1/2 -translate-y-1/2 text-xs text-blue-300">mm</span>
                            </div>
                        </div>
                         {/* Empty placeholder for alignment if needed, or Prediction button span */}
                        <div className="hidden md:block"></div>

                        <div>
                            <label className="block text-xs font-medium text-blue-300 mb-1">Weight</label>
                            <div className="relative">
                                <input type="number" defaultValue="1" className="w-full rounded-full border border-gray-300 px-4 py-2 text-sm focus:ring-2 focus:ring-blue-400 outline-none" />
                                <span className="absolute right-4 top-1/2 -translate-y-1/2 text-xs text-blue-300">kg</span>
                            </div>
                        </div>
                         <div className="flex gap-4 items-end">
                            <div className="flex-1">
                                <label className="block text-xs font-medium text-blue-300 mb-1">Quantity</label>
                                <input type="number" defaultValue="1" className="w-full rounded-full border border-gray-300 px-4 py-2 text-sm focus:ring-2 focus:ring-blue-400 outline-none" />
                            </div>
                            <button className="flex-1 bg-blue-100 text-blue-600 font-semibold py-2 rounded-full text-sm hover:bg-blue-200 transition-colors h-[38px]">
                                Prediction
                            </button>
                        </div>
                    </div>
                </div>
            </div>

            {/* --- SECTIONS 3 & 4 (Split Columns) --- */}
            <div className="flex flex-col md:flex-row gap-8 border-t border-gray-100 pt-6">
                
                {/* 3. SPACING SETTINGS */}
                <div className="flex-1 border-r border-gray-100 pr-0 md:pr-8">
                    <div className="flex items-center gap-2 mb-4">
                        <h3 className="text-blue-900 font-extrabold text-sm uppercase tracking-wide">3. Spacing Settings</h3>
                        <HelpCircle size={16} className="text-gray-400" />
                    </div>
                    
                    <div className="space-y-4">
                        {['Tilt to Length', 'Tilt to Width', 'Tilt to Height'].map((label, idx) => (
                             <div key={idx} className="flex items-start gap-4 group">
                                <input type="checkbox" className="mt-1 rounded border-gray-300 text-blue-600 focus:ring-blue-500 w-4 h-4 cursor-pointer" />
                                <div>
                                    <span className="text-sm text-gray-400 group-hover:text-blue-500 transition-colors cursor-pointer select-none">{label}</span>
                                    {/* Mini Visualization for Tilt */}
                                    <div className="mt-2 opacity-50 flex gap-4">
                                        <div className="w-12 h-8 border border-blue-200 rounded-sm bg-blue-50/50 transform -skew-x-12"></div>
                                        <ArrowRight size={16} className="text-blue-200" />
                                        <div className="w-8 h-12 border border-blue-200 rounded-sm bg-blue-50/50"></div>
                                    </div>
                                </div>
                             </div>
                        ))}
                    </div>
                </div>

                {/* 4. STUFFING SETTINGS */}
                <div className="flex-1 pl-0 md:pl-8">
                     <div className="flex items-center gap-2 mb-4">
                        <h3 className="text-blue-900 font-extrabold text-sm uppercase tracking-wide">4. Stuffing Settings</h3>
                        <HelpCircle size={16} className="text-gray-400" />
                    </div>

                    <div className="grid grid-cols-2 gap-6">
                        {/* Layers Count */}
                        <div className="space-y-2">
                             <div className="flex items-center gap-2">
                                <input type="checkbox" className="rounded border-gray-300 text-blue-600 focus:ring-blue-500 w-4 h-4" />
                                <span className="text-sm text-gray-400">Layers Count</span>
                             </div>
                             <div className="relative">
                                 {/* Simple Stack Icon */}
                                 <div className="absolute -left-10 top-0">
                                     <div className="w-6 h-6 border border-blue-200 bg-blue-50 rounded mb-[-4px]"></div>
                                     <div className="w-6 h-6 border border-blue-200 bg-blue-50 rounded mb-[-4px]"></div>
                                     <div className="w-6 h-6 border border-blue-200 bg-blue-50 rounded"></div>
                                 </div>
                                 <input type="number" placeholder="0" className="w-full bg-gray-50 border border-gray-200 rounded-full px-4 py-1 text-sm text-center" />
                             </div>
                        </div>

                         {/* Mass */}
                        <div className="space-y-2">
                             <div className="flex items-center gap-2">
                                <input type="checkbox" className="rounded border-gray-300 text-blue-600 focus:ring-blue-500 w-4 h-4" />
                                <span className="text-sm text-gray-400">Mass</span>
                             </div>
                             <div className="relative">
                                  {/* Weight Icon */}
                                 <div className="absolute -left-10 top-0 w-8 h-8 border border-blue-200 bg-blue-50 rounded-lg flex items-center justify-center text-[8px] text-blue-400">Kg</div>
                                 <input type="number" placeholder="0" className="w-full bg-gray-50 border border-gray-200 rounded-full px-4 py-1 text-sm text-center" />
                                 <span className="absolute right-3 top-1/2 -translate-y-1/2 text-[10px] text-gray-400">kg</span>
                             </div>
                        </div>

                         {/* Height */}
                        <div className="space-y-2 mt-4">
                             <div className="flex items-center gap-2">
                                <input type="checkbox" className="rounded border-gray-300 text-blue-600 focus:ring-blue-500 w-4 h-4" />
                                <span className="text-sm text-gray-400">Height</span>
                             </div>
                             <div className="relative">
                                 <input type="number" placeholder="0" className="w-full bg-gray-50 border border-gray-200 rounded-full px-4 py-1 text-sm text-center" />
                                 <span className="absolute right-3 top-1/2 -translate-y-1/2 text-[10px] text-gray-400">mm</span>
                             </div>
                        </div>

                        {/* Disable Stacking */}
                        <div className="flex items-end pb-2">
                             <div className="flex items-center gap-2 cursor-pointer hover:bg-gray-50 p-2 rounded-lg transition-colors">
                                <input type="checkbox" className="rounded border-gray-300 text-blue-600 focus:ring-blue-500 w-4 h-4" />
                                <span className="text-sm text-gray-400">Disable stacking</span>
                             </div>
                        </div>
                    </div>
                </div>
            </div>
        </div>

        {/* Modal Footer */}
        <div className="bg-white p-6 border-t border-gray-100 flex justify-end gap-4 rounded-b-2xl sticky bottom-0 z-10">
            <button 
                onClick={onClose}
                className="px-8 py-2.5 rounded-lg bg-blue-50 text-blue-600 font-bold hover:bg-blue-100 transition-colors"
            >
                Cancel
            </button>
            <button 
                onClick={onClose}
                className="px-12 py-2.5 rounded-lg bg-blue-600 text-white font-bold shadow-lg shadow-blue-500/30 hover:bg-blue-700 transition-all"
            >
                Add
            </button>
        </div>
      </div>
    </div>
  );
};


export default function App() {
  const [activeTab, setActiveTab] = useState('products');
  const [showAddProductModal, setShowAddProductModal] = useState(false);

  // Helper to get input style
  const inputClass = "w-full bg-gray-50 border border-gray-200 rounded-full px-4 py-2 text-sm text-gray-700 focus:outline-none focus:ring-2 focus:ring-blue-400 focus:bg-white transition-all text-center";

  return (
    <div className="min-h-screen bg-[#F8F9FC] font-sans text-slate-800 pb-20">
      
      {/* Modal Injection */}
      <AddProductModal 
        isOpen={showAddProductModal} 
        onClose={() => setShowAddProductModal(false)} 
      />

      {/* --- HEADER --- */}
      <header className="pt-10 pb-6 text-center">
        <h1 className="text-4xl font-bold text-slate-800 tracking-tight">Load & Stuffing Calculation</h1>
      </header>

      {/* --- TABS --- */}
      <div className="max-w-6xl mx-auto px-4 mb-8">
        <div className="flex flex-col md:flex-row bg-white rounded-xl shadow-sm overflow-hidden border border-gray-100">
          <button 
            onClick={() => setActiveTab('products')}
            className={`flex-1 py-4 flex items-center justify-center gap-2 font-semibold transition-all ${activeTab === 'products' ? 'text-blue-600 border-b-4 border-blue-500 bg-blue-50/30' : 'text-gray-400 hover:text-gray-600'}`}
          >
            <Box size={18} /> PRODUCTS
          </button>
          <div className="w-px bg-gray-100 hidden md:block"></div>
          <button 
            className="flex-1 py-4 flex items-center justify-center gap-2 font-semibold text-gray-400 cursor-not-allowed"
          >
             <span className="opacity-50">🚚 CONTAINERS & TRUCKS</span>
          </button>
          <div className="w-px bg-gray-100 hidden md:block"></div>
          <button 
            onClick={() => setActiveTab('result')}
            className={`flex-1 py-4 flex items-center justify-center gap-2 font-semibold transition-all ${activeTab === 'result' ? 'text-blue-600 border-b-4 border-blue-500 bg-blue-50/30' : 'text-gray-400 hover:text-gray-600'}`}
          >
            <Package size={18} /> STUFFING RESULT
          </button>
          <div className="w-16 flex items-center justify-center border-l border-gray-100 text-gray-400 hover:text-blue-500 cursor-pointer transition-colors">
            <Settings size={20} />
          </div>
        </div>
      </div>

      {/* --- CONTENT AREA --- */}
      <main className="max-w-6xl mx-auto px-4">
        
        {activeTab === 'products' ? (
          /* ================= PRODUCTS VIEW ================= */
          <div className="space-y-6 animate-in fade-in slide-in-from-bottom-4 duration-500">
            
            {/* Action Bar */}
            <div className="flex flex-wrap items-center justify-between gap-4">
              <button 
                onClick={() => setShowAddProductModal(true)}
                className="bg-blue-600 hover:bg-blue-700 text-white px-6 py-2.5 rounded-xl font-medium shadow-lg shadow-blue-600/20 transition-all flex items-center gap-2"
              >
                <Plus size={18} /> Add Group
              </button>
              
              <div className="flex gap-3">
                <button className="bg-green-100 hover:bg-green-200 text-green-700 px-5 py-2.5 rounded-xl font-medium transition-colors flex items-center gap-2">
                  <Download size={18} /> Import
                </button>
                <button className="bg-blue-50 hover:bg-blue-100 text-blue-600 px-5 py-2.5 rounded-xl font-medium transition-colors flex items-center gap-2">
                  <Upload size={18} /> Export
                </button>
                <button className="bg-orange-500 hover:bg-orange-600 text-white px-5 py-2.5 rounded-xl font-medium shadow-lg shadow-orange-500/20 transition-colors flex items-center gap-2">
                  <Zap size={18} /> Upgrade
                </button>
              </div>
            </div>

            {/* Group Card */}
            <div className="bg-white rounded-xl shadow-sm border border-gray-100 overflow-hidden">
              {/* Card Header */}
              <div className="px-6 py-4 border-b border-gray-100 flex items-center justify-between">
                <div className="flex items-center gap-3">
                  <div className="w-3 h-3 rounded-full bg-black"></div>
                  <h3 className="text-lg font-bold text-gray-800">Group #1</h3>
                  <button className="p-1.5 rounded-full hover:bg-gray-100 text-gray-400 hover:text-blue-500 transition-colors">
                    <Edit2 size={16} />
                  </button>
                </div>
                <div className="flex gap-2">
                    <button className="p-2 text-red-400 hover:bg-red-50 rounded-lg"><Trash2 size={18}/></button>
                    <button className="p-2 text-blue-400 hover:bg-blue-50 rounded-lg"><Upload size={18} className="rotate-180"/></button>
                    <button className="p-2 text-blue-400 hover:bg-blue-50 rounded-lg"><Upload size={18}/></button>
                </div>
              </div>

              {/* Table Area */}
              <div className="p-6 overflow-x-auto">
                <table className="w-full min-w-[900px]">
                  <thead>
                    <tr className="text-left text-xs font-semibold text-gray-400 uppercase tracking-wider">
                      <th className="pb-4 pl-2 w-12">Type</th>
                      <th className="pb-4 w-48">Product Name</th>
                      <th className="pb-4 w-28 text-center">Length <span className="text-[10px] normal-case block text-gray-300">(mm)</span></th>
                      <th className="pb-4 w-28 text-center">Width <span className="text-[10px] normal-case block text-gray-300">(mm)</span></th>
                      <th className="pb-4 w-28 text-center">Height <span className="text-[10px] normal-case block text-gray-300">(mm)</span></th>
                      <th className="pb-4 w-28 text-center">Weight <span className="text-[10px] normal-case block text-gray-300">(kg)</span></th>
                      <th className="pb-4 w-24 text-center">Quantity</th>
                      <th className="pb-4 w-16 text-center">Color</th>
                      <th className="pb-4 w-24 text-center">Stack</th>
                    </tr>
                  </thead>
                  <tbody className="space-y-4">
                    {PRODUCT_DATA.map((item) => (
                      <tr key={item.id} className="group">
                        <td className="py-2 pl-2">
                          <div className="w-10 h-10 rounded-lg bg-gray-50 flex items-center justify-center text-gray-400 border border-gray-100">
                             <TypeIcon type={item.type} color={item.color} />
                          </div>
                        </td>
                        <td className="py-2 pr-4">
                          <input type="text" defaultValue={item.name} className={`${inputClass} !text-left`} />
                        </td>
                        <td className="py-2 px-1">
                          <div className="relative">
                            <input type="number" defaultValue={item.l} className={inputClass} />
                            <span className="absolute right-3 top-1/2 -translate-y-1/2 text-[10px] text-gray-400 pointer-events-none">mm</span>
                          </div>
                        </td>
                        <td className="py-2 px-1">
                          <div className="relative">
                            <input type="number" defaultValue={item.w} className={inputClass} />
                            <span className="absolute right-3 top-1/2 -translate-y-1/2 text-[10px] text-gray-400 pointer-events-none">mm</span>
                          </div>
                        </td>
                        <td className="py-2 px-1">
                          <div className="relative">
                            <input type="number" defaultValue={item.h} className={inputClass} />
                            <span className="absolute right-3 top-1/2 -translate-y-1/2 text-[10px] text-gray-400 pointer-events-none">mm</span>
                          </div>
                        </td>
                        <td className="py-2 px-1">
                          <div className="relative">
                            <input type="number" defaultValue={item.weight} className={inputClass} />
                            <span className="absolute right-3 top-1/2 -translate-y-1/2 text-[10px] text-gray-400 pointer-events-none">kg</span>
                          </div>
                        </td>
                        <td className="py-2 px-1">
                           <input type="number" defaultValue={item.qty} className={inputClass} />
                        </td>
                        <td className="py-2 px-1 text-center">
                          <div className="flex justify-center cursor-pointer hover:scale-110 transition-transform">
                            <ColorDot color={item.color} />
                          </div>
                        </td>
                        <td className="py-2 pl-2 text-center">
                          <div className="flex items-center justify-center gap-2">
                             <button className="p-2 rounded-lg bg-blue-50 text-blue-500 hover:bg-blue-100 transition-colors"><Settings size={16}/></button>
                             <button className="p-2 rounded-lg bg-red-50 text-red-500 hover:bg-red-100 transition-colors"><Trash2 size={16}/></button>
                          </div>
                        </td>
                      </tr>
                    ))}
                  </tbody>
                </table>
              </div>
              
              {/* Card Footer */}
              <div className="px-6 py-4 bg-gray-50/50 border-t border-gray-100 flex items-center justify-between">
                 <button 
                    onClick={() => setShowAddProductModal(true)}
                    className="text-blue-600 hover:text-blue-700 hover:bg-blue-50 px-4 py-2 rounded-lg font-medium transition-colors flex items-center gap-2"
                 >
                    <Plus size={18} /> Add product
                 </button>

                 <div className="flex items-center gap-2">
                    <div className="w-5 h-5 rounded border border-gray-300 bg-white flex items-center justify-center cursor-pointer"></div>
                    <span className="text-sm font-medium text-gray-600">Use pallets</span>
                    <HelpCircle size={14} className="text-gray-400" />
                 </div>
              </div>
            </div>

            {/* Bottom Action */}
            <div className="flex justify-center pt-8">
               <button 
                onClick={() => setActiveTab('result')}
                className="bg-blue-600 hover:bg-blue-700 text-white text-lg px-12 py-3 rounded-xl font-bold shadow-xl shadow-blue-600/30 hover:shadow-blue-600/40 transform hover:-translate-y-0.5 transition-all w-full md:w-auto"
               >
                 Next
               </button>
            </div>

          </div>
        ) : (
          /* ================= RESULT VIEW ================= */
          <div className="space-y-6 animate-in fade-in slide-in-from-right-8 duration-500">
             
             {/* Top Summary Card */}
             <div className="bg-white rounded-xl shadow-sm border border-gray-100 p-6 flex flex-col md:flex-row items-center gap-8">
                <div className="flex-shrink-0 w-32 md:w-48">
                    <ContainerOutline />
                </div>
                <div className="flex-1 space-y-2">
                    <h2 className="text-xl font-bold text-gray-800">20 STANDARD</h2>
                    <div className="grid grid-cols-2 gap-4 text-sm text-gray-600">
                        <div className="bg-gray-50 px-3 py-2 rounded-lg border border-gray-100">
                            <span className="block text-xs text-gray-400 uppercase">Max Weight</span>
                            <span className="font-semibold text-gray-700">14,300 kg</span>
                        </div>
                         <div className="bg-gray-50 px-3 py-2 rounded-lg border border-gray-100">
                            <span className="block text-xs text-gray-400 uppercase">Max Volume</span>
                            <span className="font-semibold text-gray-700">28.30 m³</span>
                        </div>
                    </div>
                </div>
                <div className="flex-shrink-0 text-center px-6 py-4 bg-blue-50 rounded-xl border border-blue-100">
                    <span className="block text-2xl font-bold text-blue-600">1</span>
                    <span className="text-sm font-medium text-blue-400 uppercase">Unit</span>
                </div>
             </div>

             {/* Bottom Grid */}
             <div className="grid grid-cols-1 lg:grid-cols-2 gap-6">
                
                {/* 3D Visualization Card */}
                <div className="bg-white rounded-xl shadow-sm border border-gray-100 p-6 flex flex-col h-full relative overflow-hidden group">
                    <div className="flex justify-between items-start mb-4 z-10">
                        <h3 className="text-xl font-bold text-gray-800">20 Standard #1</h3>
                        <span className="text-blue-600 font-bold bg-blue-50 px-3 py-1 rounded-full text-sm">1 unit</span>
                    </div>

                    <div className="flex-1 flex items-center justify-center min-h-[250px] relative">
                         {/* Fake 3D Content */}
                         <ContainerFilled />
                    </div>

                    <div className="absolute bottom-6 right-6 z-10">
                        <button className="bg-blue-100 hover:bg-blue-200 text-blue-700 px-4 py-2 rounded-lg font-medium flex items-center gap-2 transition-colors shadow-sm">
                            <Move3d size={18} /> 3D VIEW
                        </button>
                    </div>
                </div>

                {/* Statistics Card */}
                <div className="bg-white rounded-xl shadow-sm border border-gray-100 p-6">
                    <div className="mb-6 space-y-3">
                        <div className="flex justify-between items-center border-b border-gray-100 pb-3">
                            <span className="font-bold text-gray-700">Total</span>
                            <span className="font-bold text-gray-900 text-lg">190 packages</span>
                        </div>
                        <div className="flex justify-between items-center">
                            <span className="font-bold text-gray-700">Cargo volume</span>
                            <div className="text-right">
                                <span className="font-bold text-gray-900">28.30 m³</span>
                                <span className="text-sm text-gray-400 ml-2">(85% of volume)</span>
                            </div>
                        </div>
                        <div className="flex justify-between items-center">
                            <span className="font-bold text-gray-700">Cargo weight</span>
                            <div className="text-right">
                                <span className="font-bold text-gray-900">14,300 kg</span>
                                <span className="text-sm text-gray-400 ml-2">(50% of max weight)</span>
                            </div>
                        </div>
                    </div>

                    <div className="flex flex-col md:flex-row items-center gap-6">
                        {/* Donut Chart */}
                        <div className="w-40 h-40 relative">
                            <ResponsiveContainer width="100%" height="100%">
                                <PieChart>
                                    <Pie
                                        data={CHART_DATA}
                                        cx="50%"
                                        cy="50%"
                                        innerRadius={40}
                                        outerRadius={60}
                                        paddingAngle={5}
                                        dataKey="value"
                                        stroke="none"
                                    >
                                        {CHART_DATA.map((entry, index) => (
                                            <Cell key={`cell-${index}`} fill={entry.color} />
                                        ))}
                                    </Pie>
                                </PieChart>
                            </ResponsiveContainer>
                            {/* Center Text Trick */}
                            <div className="absolute inset-0 flex items-center justify-center pointer-events-none">
                                <span className="text-xs text-gray-400 font-medium">MIX</span>
                            </div>
                        </div>

                        {/* Legend Table */}
                        <div className="flex-1 w-full overflow-x-auto">
                            <table className="w-full text-sm">
                                <thead>
                                    <tr className="text-left text-gray-400 text-xs uppercase">
                                        <th className="font-medium pb-2">Name</th>
                                        <th className="font-medium pb-2">Packages</th>
                                        <th className="font-medium pb-2">Volume</th>
                                        <th className="font-medium pb-2 text-right">Weight</th>
                                    </tr>
                                </thead>
                                <tbody className="divide-y divide-gray-50">
                                    {/* Row 1 */}
                                    <tr>
                                        <td className="py-2 flex items-center gap-2 font-bold text-gray-700">
                                            <div className="w-3 h-3 rounded-full bg-blue-500"></div> big bags
                                        </td>
                                        <td className="py-2 text-gray-600"><Package size={14} className="inline mr-1 opacity-50"/> 10</td>
                                        <td className="py-2 text-gray-600">10.00 m³</td>
                                        <td className="py-2 text-right text-gray-600">9000.00 kg</td>
                                    </tr>
                                    {/* Row 2 */}
                                    <tr>
                                        <td className="py-2 flex items-center gap-2 font-bold text-gray-700">
                                            <div className="w-3 h-3 rounded-full bg-pink-500"></div> sacks
                                        </td>
                                        <td className="py-2 text-gray-600"><ShoppingBag size={14} className="inline mr-1 opacity-50"/> 100</td>
                                        <td className="py-2 text-gray-600">13.50 m³</td>
                                        <td className="py-2 text-right text-gray-600">4500.00 kg</td>
                                    </tr>
                                    {/* Row 3 */}
                                    <tr>
                                        <td className="py-2 flex items-center gap-2 font-bold text-gray-700">
                                            <div className="w-3 h-3 rounded-full bg-green-500"></div> boxes 1
                                        </td>
                                        <td className="py-2 text-gray-600"><Box size={14} className="inline mr-1 opacity-50"/> 80</td>
                                        <td className="py-2 text-gray-600">4.80 m³</td>
                                        <td className="py-2 text-right text-gray-600">800.00 kg</td>
                                    </tr>
                                </tbody>
                            </table>
                        </div>
                    </div>
                </div>
             </div>

             {/* Footer Action Buttons */}
             <div className="flex flex-col-reverse sm:flex-row justify-center items-center gap-4 pt-8">
                 <button 
                    onClick={() => setActiveTab('products')}
                    className="w-full sm:w-auto px-8 py-3 bg-blue-50 text-blue-600 hover:bg-blue-100 rounded-xl font-bold transition-colors flex justify-center items-center gap-2"
                 >
                     <ArrowLeft size={18} /> Back
                 </button>
                 <button className="w-full sm:w-auto px-8 py-3 bg-blue-600 text-white hover:bg-blue-700 rounded-xl font-bold shadow-lg shadow-blue-600/20 transition-all flex justify-center items-center gap-2">
                     <Download size={18} /> Export to PDF
                 </button>
                 <button className="w-full sm:w-auto px-8 py-3 bg-[#0EA5E9] text-white hover:bg-[#0284c7] rounded-xl font-bold shadow-lg shadow-sky-500/20 transition-all flex justify-center items-center gap-2">
                     <Copy size={18} /> Copy request
                 </button>
             </div>
          </div>
        )}

      </main>
    </div>
  );
}