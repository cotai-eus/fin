/**
 * Live Chat Component
 * Simple chat interface with Server-Sent Events
 */

"use client";

import { useState, useEffect, useRef } from "react";
import { Card } from "@/shared/components/ui/Card";
import { Button } from "@/shared/components/ui/Button";
import type { ChatMessage } from "../types";

export function LiveChat() {
  const [messages, setMessages] = useState<ChatMessage[]>([]);
  const [inputMessage, setInputMessage] = useState("");
  const [isConnected, setIsConnected] = useState(false);
  const [isTyping, setIsTyping] = useState(false);
  const messagesEndRef = useRef<HTMLDivElement>(null);
  const eventSourceRef = useRef<EventSource | null>(null);

  useEffect(() => {
    // Connect to SSE endpoint
    // In production, this would connect to your backend SSE endpoint
    // For now, we skip SSE connection (commented out for future implementation)
    // const eventSource = new EventSource('/api/chat/stream');
    // eventSourceRef.current = eventSource;

    // Cleanup
    return () => {
      const eventSource = eventSourceRef.current;
      if (eventSource) {
        eventSource.close();
      }
    };
  }, []);

  // Separate effect for marking connection ready (no ref dependencies)
  useEffect(() => {
    // Small delay to ensure component is mounted
    const timer = requestAnimationFrame(() => {
      setIsConnected(true);
    });
    return () => cancelAnimationFrame(timer);
  }, []);

  useEffect(() => {
    // Auto-scroll to bottom
    messagesEndRef.current?.scrollIntoView({ behavior: "smooth" });
  }, [messages]);

  const handleSendMessage = async () => {
    if (!inputMessage.trim()) return;

    const newMessage: ChatMessage = {
      id: crypto.randomUUID(),
      userId: "current-user",
      userName: "Você",
      message: inputMessage,
      isStaff: false,
      timestamp: new Date().toISOString(),
    };

    setMessages((prev) => [...prev, newMessage]);
    setInputMessage("");

    // Simulate staff response (in production, this comes from SSE)
    setIsTyping(true);
    setTimeout(() => {
      const staffResponse: ChatMessage = {
        id: crypto.randomUUID(),
        userId: "staff",
        userName: "Suporte LauraTech",
        message:
          "Obrigado pela sua mensagem! Nossa equipe irá responder em breve.",
        isStaff: true,
        timestamp: new Date().toISOString(),
      };
      setMessages((prev) => [...prev, staffResponse]);
      setIsTyping(false);
    }, 1500);
  };

  const handleKeyPress = (e: React.KeyboardEvent) => {
    if (e.key === "Enter" && !e.shiftKey) {
      e.preventDefault();
      handleSendMessage();
    }
  };

  return (
    <Card className="p-0 h-[600px] flex flex-col">
      {/* Header */}
      <div className="p-4 border-b flex items-center justify-between">
        <div className="flex items-center gap-2">
          <div className="w-3 h-3 rounded-full bg-green-500" />
          <h3 className="font-semibold">Chat ao Vivo</h3>
        </div>
        {isConnected && (
          <span className="text-xs text-green-600">Conectado</span>
        )}
      </div>

      {/* Messages */}
      <div className="flex-1 overflow-y-auto p-4 space-y-3">
        {messages.length === 0 && (
          <div className="text-center text-gray-500 mt-8">
            <p className="mb-2">👋 Olá! Como podemos ajudar?</p>
            <p className="text-sm">
              Nossa equipe está online e pronta para atendê-lo
            </p>
          </div>
        )}

        {messages.map((msg) => (
          <div
            key={msg.id}
            className={`flex ${msg.isStaff ? "justify-start" : "justify-end"}`}
          >
            <div
              className={`max-w-[70%] rounded-lg p-3 ${
                msg.isStaff
                  ? "bg-gray-100 text-gray-900"
                  : "bg-blue-600 text-white"
              }`}
            >
              <p className="text-xs font-medium mb-1 opacity-75">
                {msg.userName}
              </p>
              <p className="text-sm whitespace-pre-wrap">{msg.message}</p>
              <p className="text-xs opacity-60 mt-1">
                {new Date(msg.timestamp).toLocaleTimeString("pt-BR", {
                  hour: "2-digit",
                  minute: "2-digit",
                })}
              </p>
            </div>
          </div>
        ))}

        {isTyping && (
          <div className="flex justify-start">
            <div className="bg-gray-100 rounded-lg p-3">
              <div className="flex gap-1">
                <div className="w-2 h-2 bg-gray-400 rounded-full animate-bounce" />
                <div
                  className="w-2 h-2 bg-gray-400 rounded-full animate-bounce"
                  style={{ animationDelay: "0.1s" }}
                />
                <div
                  className="w-2 h-2 bg-gray-400 rounded-full animate-bounce"
                  style={{ animationDelay: "0.2s" }}
                />
              </div>
            </div>
          </div>
        )}

        <div ref={messagesEndRef} />
      </div>

      {/* Input */}
      <div className="p-4 border-t">
        <div className="flex gap-2">
          <textarea
            value={inputMessage}
            onChange={(e) => setInputMessage(e.target.value)}
            onKeyPress={handleKeyPress}
            placeholder="Digite sua mensagem..."
            rows={2}
            className="flex-1 px-3 py-2 border rounded-lg focus:ring-2 focus:ring-blue-500 resize-none"
          />
          <Button
            onClick={handleSendMessage}
            disabled={!inputMessage.trim()}
            className="self-end"
          >
            Enviar
          </Button>
        </div>
        <p className="text-xs text-gray-500 mt-2">
          Pressione Enter para enviar, Shift+Enter para nova linha
        </p>
      </div>
    </Card>
  );
}
