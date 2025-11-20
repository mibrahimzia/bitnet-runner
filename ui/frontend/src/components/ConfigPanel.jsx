import React, { useState } from 'react';
import { useChatStore } from '../stores/chatStore';

export default function ConfigPanel() {
  const { config, updateConfig } = useChatStore();
  const [isOpen, setIsOpen] = useState(false);

  return (
    <div className="border-t border-gray-700">
      <button 
        onClick={() => setIsOpen(!isOpen)}
        className="w-full p-4 flex justify-between items-center hover:bg-gray-800 transition-colors text-sm font-bold text-gray-400"
      >
        <span>SETTINGS</span>
        <span>{isOpen ? '▼' : '▶'}</span>
      </button>

      {isOpen && (
        <div className="p-4 bg-surface space-y-4 text-sm animate-in slide-in-from-top-2">
          
          {/* System Prompt */}
          <div>
            <label className="block text-xs text-gray-500 mb-1">System Prompt</label>
            <textarea 
              className="w-full bg-background p-2 rounded border border-gray-600 text-xs h-20 focus:border-primary outline-none resize-none"
              value={config.systemPrompt}
              onChange={(e) => updateConfig('systemPrompt', e.target.value)}
            />
          </div>

          {/* Temperature */}
          <div>
            <div className="flex justify-between mb-1">
              <label className="text-xs text-gray-500">Temperature</label>
              <span className="text-xs text-primary">{config.temperature}</span>
            </div>
            <input 
              type="range" min="0" max="2" step="0.1"
              className="w-full h-1 bg-gray-600 rounded-lg appearance-none cursor-pointer"
              value={config.temperature}
              onChange={(e) => updateConfig('temperature', parseFloat(e.target.value))}
            />
          </div>

          {/* Max Tokens */}
          <div>
            <div className="flex justify-between mb-1">
              <label className="text-xs text-gray-500">Max Tokens</label>
              <span className="text-xs text-primary">{config.maxTokens}</span>
            </div>
            <input 
              type="range" min="128" max="4096" step="128"
              className="w-full h-1 bg-gray-600 rounded-lg appearance-none cursor-pointer"
              value={config.maxTokens}
              onChange={(e) => updateConfig('maxTokens', parseInt(e.target.value))}
            />
          </div>

             {/* Repeat Penalty */}
             <div>
            <div className="flex justify-between mb-1">
              <label className="text-xs text-gray-500">Repeat Penalty</label>
              <span className="text-xs text-primary">{config.repeatPenalty}</span>
            </div>
            <input 
              type="range" min="1.0" max="2.0" step="0.05"
              className="w-full h-1 bg-gray-600 rounded-lg appearance-none cursor-pointer"
              value={config.repeatPenalty}
              onChange={(e) => updateConfig('repeatPenalty', parseFloat(e.target.value))}
            />
          </div>

        </div>
      )}
    </div>
  );
}