import React, { useState } from 'react';
import ModelSelector from './components/ModelSelector';
import ConfigPanel from './components/ConfigPanel'; // <--- IMPORT
import Chat from './components/Chat';
import InputArea from './components/InputArea';

function App() {
  const [sidebarOpen, setSidebarOpen] = useState(true);

  return (
    <div className="flex h-screen w-screen bg-background text-white overflow-hidden relative">
      
      {/* Toggle Button (Floating when closed, Header when open) */}
      <button 
        onClick={() => setSidebarOpen(!sidebarOpen)}
        className={`absolute z-50 top-4 left-4 p-2 rounded bg-surface border border-gray-600 hover:bg-gray-700 transition-all ${sidebarOpen ? 'left-[230px]' : 'left-4'}`}
      >
        {sidebarOpen ? '◀' : '☰'}
      </button>

      {/* Sidebar */}
      <div className={`
        border-r border-gray-700 flex flex-col bg-surface transition-all duration-300 ease-in-out
        ${sidebarOpen ? 'w-64 translate-x-0' : 'w-0 -translate-x-full overflow-hidden border-none'}
      `}>
        <div className="p-4 border-b border-gray-700 font-bold text-lg flex items-center gap-2 pl-12">
            BitNet
        </div>
        
        <ModelSelector />
        <ConfigPanel />  {/* <--- ADDED HERE */}
        
        <div className="flex-1 p-4 overflow-y-auto">
            <div className="text-xs text-gray-500 mt-4">
                <p className="mb-2 font-bold text-gray-400">Instructions:</p>
                <ul className="list-disc pl-4 space-y-2">
                    <li>Select a model</li>
                    <li>Click <strong>Load Model</strong></li>
                    <li>Wait for the ✔</li>
                    <li>Chat freely!</li>
                </ul>
            </div>
        </div>
      </div>

      {/* Main Chat Area */}
      <div className="flex-1 flex flex-col relative">
        <Chat />
        <InputArea />
      </div>
    </div>
  );
}

export default App;