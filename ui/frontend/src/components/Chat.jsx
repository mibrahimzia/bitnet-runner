import React, { useRef, useEffect, useState } from 'react';
import { useChatStore } from '../stores/chatStore';
import ReactMarkdown from 'react-markdown';
import remarkGfm from 'remark-gfm';
import remarkMath from 'remark-math';
import rehypeKatex from 'rehype-katex';
import { Copy, Check } from 'lucide-react'; // Icons
import 'katex/dist/katex.min.css'; // Math CSS

// Custom Code Block Component with Copy Button
const CodeBlock = ({ node, inline, className, children, ...props }) => {
  const [isCopied, setIsCopied] = useState(false);
  const codeContent = String(children).replace(/\n$/, '');

  const handleCopy = async () => {
    await navigator.clipboard.writeText(codeContent);
    setIsCopied(true);
    setTimeout(() => setIsCopied(false), 2000);
  };

  // Inline code (single backtick) - No copy button
  if (inline) {
    return (
      <code className="bg-black/30 px-1 rounded text-primary font-mono text-sm" {...props}>
        {children}
      </code>
    );
  }

  // Block code (triple backtick)
  return (
    <div className="relative my-4 group rounded-lg overflow-hidden bg-[#1e1e1e] border border-gray-700">
      {/* Header bar with language and copy button */}
      <div className="flex items-center justify-between px-4 py-2 bg-[#2d2d2d] border-b border-gray-700">
        <span className="text-xs text-gray-400 font-mono lowercase">
          {className?.replace('language-', '') || 'text'}
        </span>
        <button
          onClick={handleCopy}
          className="flex items-center gap-1 text-xs text-gray-400 hover:text-white transition-colors"
          title="Copy code"
        >
          {isCopied ? (
            <>
              <Check size={14} className="text-green-500" />
              <span className="text-green-500">Copied!</span>
            </>
          ) : (
            <>
              <Copy size={14} />
              <span>Copy</span>
            </>
          )}
        </button>
      </div>
      
      {/* The Code */}
      <div className="p-4 overflow-x-auto">
        <code className={`!bg-transparent !p-0 font-mono text-sm ${className}`} {...props}>
          {children}
        </code>
      </div>
    </div>
  );
};

export default function Chat() {
  const { messages, isGenerating } = useChatStore();
  const endRef = useRef(null);

  useEffect(() => {
    endRef.current?.scrollIntoView({ behavior: "smooth" });
  }, [messages]);

  return (
    <div className="flex-1 overflow-y-auto p-6 space-y-6 bg-background">
      {messages.length === 0 && (
        <div className="h-full flex flex-col items-center justify-center text-gray-500 opacity-50">
            <div className="text-4xl mb-4">💬</div>
            <p>Select a model and say hello!</p>
        </div>
      )}

      {messages.map((msg, idx) => (
        <div 
          key={idx} 
          className={`flex ${msg.role === 'user' ? 'justify-end' : 'justify-start'}`}
        >
          <div 
            className={`max-w-[85%] rounded-lg p-4 shadow-lg ${
              msg.role === 'user' 
                ? 'bg-primary text-white rounded-br-none' 
                : 'bg-surface text-gray-200 rounded-bl-none'
            }`}
          >
            <div className="prose prose-invert max-w-none prose-sm">
                <ReactMarkdown 
                    // Add Math and GFM plugins
                    remarkPlugins={[remarkGfm, remarkMath]}
                    rehypePlugins={[rehypeKatex]}
                    components={{
                        code: CodeBlock,
                        // Basic styling for other elements
                        p: ({children}) => <p className="mb-2 last:mb-0">{children}</p>,
                        a: ({children, href}) => <a href={href} target="_blank" rel="noopener noreferrer" className="text-blue-400 hover:underline">{children}</a>
                    }}
                >
                    {msg.content}
                </ReactMarkdown>
            </div>
          </div>
        </div>
      ))}
      
      {isGenerating && (
        <div className="flex justify-start">
           <div className="bg-surface p-4 rounded-lg rounded-bl-none animate-pulse">
             <span className="w-2 h-2 bg-gray-400 rounded-full inline-block mr-1"></span>
             <span className="w-2 h-2 bg-gray-400 rounded-full inline-block mr-1"></span>
             <span className="w-2 h-2 bg-gray-400 rounded-full inline-block"></span>
           </div>
        </div>
      )}
      
      <div ref={endRef} />
    </div>
  );
}