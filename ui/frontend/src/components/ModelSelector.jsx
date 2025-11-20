import React, { useEffect, useState } from 'react';
import { useChatStore } from '../stores/chatStore';
import { ListModels, LoadModelOnly } from '../wailsjs/go/backend/App'; 
import { EventsOn } from '../wailsjs/runtime/runtime';

export default function ModelSelector() {
  const { models, setModels, selectedModel, setSelectedModel } = useChatStore();
  const [loading, setLoading] = useState(false);
  const [loadedId, setLoadedId] = useState(null);

  useEffect(() => {
    refreshModels();
    
    EventsOn("model_loaded", (id) => {
        setLoading(false);
        setLoadedId(id);
    });
    
    EventsOn("model_load_error", (err) => {
        setLoading(false);
        alert("Load failed: " + err);
    });
  }, []);

  const refreshModels = async () => {
    try {
      const list = await ListModels();
      setModels(list || []);
      if (list && list.length > 0 && !selectedModel) {
        setSelectedModel(list[0].id);
      }
    } catch (e) {
      console.error(e);
    }
  };

  const handleLoad = async () => {
      if (!selectedModel) return;
      setLoading(true);
      await LoadModelOnly(selectedModel);
  };

  return (
    <div className="p-4 border-b border-gray-700 bg-surface">
      <label className="block text-xs text-gray-400 mb-1 uppercase font-bold">
        Model
      </label>
      
      {/* Changed to flex-col for better space management on narrow bars */}
      <div className="flex flex-col gap-2">
          <select 
            className="w-full bg-background text-white p-2 rounded border border-gray-600 focus:border-primary outline-none text-sm text-ellipsis overflow-hidden"
            value={selectedModel || ''}
            onChange={(e) => {
                setSelectedModel(e.target.value);
                setLoadedId(null); 
            }}
            title={selectedModel} // Tooltip for long names
          >
            {models.length === 0 && <option>No models found</option>}
            {models.map((m) => (
              <option key={m.id} value={m.id}>
                {m.name}
              </option>
            ))}
          </select>

          <button 
            onClick={handleLoad}
            disabled={loading || !selectedModel || loadedId === selectedModel}
            className={`w-full py-2 px-3 rounded font-bold text-sm transition-colors flex items-center justify-center gap-2 ${
                loadedId === selectedModel 
                ? "bg-gray-600 text-green-400 cursor-default"
                : "bg-primary hover:bg-green-600 text-white"
            }`}
          >
            {loading ? (
                <span className="animate-spin">↻</span> 
            ) : (
                loadedId === selectedModel ? "✔ Loaded" : "Load Model"
            )}
          </button>
      </div>
    </div>
  );
}