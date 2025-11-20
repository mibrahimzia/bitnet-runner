import { create } from 'zustand';

export const useChatStore = create((set) => ({
  messages: [],
  isGenerating: false,
  selectedModel: null,
  models: [],
  
  // --- NEW CONFIGURATION STATE ---
  config: {
    systemPrompt: "You are a helpful, precise assistant. Do not repeat words or characters.",
    temperature: 0.7,
    topP: 0.9,
    topK: 40,
    maxTokens: 2048,
    repeatPenalty: 1.1,
    RepeatLastN:   192,
    Mirostat:      2,      // if supported
    MirostatTau:   5.0,
    MirostatEta:   0.1
  },

  // Actions
  setModels: (models) => set({ models }),
  setSelectedModel: (model) => set({ selectedModel: model }),
  setGenerating: (status) => set({ isGenerating: status }),
  
  // New Config Actions
  updateConfig: (key, value) => set((state) => ({
    config: { ...state.config, [key]: value }
  })),

  addMessage: (role, content) => set((state) => ({
    messages: [...state.messages, { role, content, timestamp: new Date() }]
  })),

  updateLastMessage: (content) => set((state) => {
    const newMessages = [...state.messages];
    const lastMsg = newMessages[newMessages.length - 1];
    if (lastMsg && lastMsg.role === 'assistant') {
      lastMsg.content += content;
      return { messages: newMessages };
    }
    return state;
  }),

  clearChat: () => set({ messages: [] }),
}));