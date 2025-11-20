import React, { useState, useEffect } from 'react';
import { useChatStore } from '../stores/chatStore';
import { StartChat, StopChat } from '../wailsjs/go/backend/App';
import { EventsOn } from '../wailsjs/runtime/runtime';

export default function InputArea() {
  const [input, setInput] = useState('');
  const { 
    addMessage, 
    selectedModel, 
    isGenerating, 
    setGenerating, 
    updateLastMessage,
    config // <--- We need this from the store
  } = useChatStore();

  // Setup Event Listeners
  useEffect(() => {
    // EventsOn returns a 'cancel' function that removes the listener
    const cancelToken = EventsOn("chat_token", (token) => updateLastMessage(token));
    const cancelDone = EventsOn("chat_done", () => setGenerating(false));
    const cancelError = EventsOn("chat_error", (err) => {
        alert(err); 
        setGenerating(false);
    });

    // Cleanup function: This runs when the component unmounts (or re-runs in Strict Mode)
    return () => {
        cancelToken();
        cancelDone();
        cancelError();
    };
  }, []);

  const handleSend = async () => {
    if (!input.trim() || !selectedModel || isGenerating) return;

    const prompt = input;
    setInput('');
    
    // 1. Add User Message
    addMessage('user', prompt);
    
    // 2. Add Empty Assistant Message (placeholder for streaming)
    addMessage('assistant', '');
    setGenerating(true);

    // 3. Call Go Backend with ALL 7 ARGUMENTS
    console.log("Sending config:", config); // Debug log
    
    try {
      await StartChat(
        prompt, 
        selectedModel, 
        Number(config.temperature), // Ensure numbers are numbers
        config.systemPrompt,
        Number(config.topP),
        Number(config.topK),
        Number(config.maxTokens)
      );
    } catch (e) {
      console.error("Failed to start chat:", e);
      setGenerating(false);
    }
  };

  const handleStop = async () => {
      await StopChat();
      setGenerating(false);
  };

  return (
    <div className="p-4 bg-surface border-t border-gray-700">
      <div className="flex gap-2 max-w-4xl mx-auto">
        <input
          type="text"
          className="flex-1 bg-background text-white p-3 rounded-lg border border-gray-600 focus:border-primary outline-none transition-colors"
          placeholder={selectedModel ? "Type a message..." : "Select a model first..."}
          value={input}
          onChange={(e) => setInput(e.target.value)}
          onKeyDown={(e) => e.key === 'Enter' && handleSend()}
          disabled={isGenerating && !selectedModel}
        />
        
        {isGenerating ? (
             <button 
             onClick={handleStop}
             className="bg-red-500 hover:bg-red-600 text-white px-6 rounded-lg font-medium transition-colors"
           >
             Stop
           </button>
        ) : (
            <button 
            onClick={handleSend}
            disabled={!selectedModel}
            className="bg-primary hover:bg-green-600 disabled:opacity-50 disabled:cursor-not-allowed text-white px-6 rounded-lg font-medium transition-colors"
            >
            Send
            </button>
        )}
      </div>
    </div>
  );
}