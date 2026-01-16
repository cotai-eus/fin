/**
 * FAQ Section Component
 * Displays frequently asked questions
 */

"use client";

import { useState } from "react";
import { Card } from "@/shared/components/ui/Card";
import type { FAQCategory } from "../types";

interface FAQSectionProps {
  categories: FAQCategory[];
}

export function FAQSection({ categories }: FAQSectionProps) {
  const [expandedItems, setExpandedItems] = useState<Set<string>>(new Set());

  const toggleItem = (id: string) => {
    const newExpanded = new Set(expandedItems);
    if (newExpanded.has(id)) {
      newExpanded.delete(id);
    } else {
      newExpanded.add(id);
    }
    setExpandedItems(newExpanded);
  };

  return (
    <div className="space-y-6">
      {categories.map((category) => (
        <div key={category.name}>
          <h3 className="text-lg font-semibold mb-3 flex items-center gap-2">
            <span>{category.icon}</span>
            {category.name}
          </h3>

          <div className="space-y-2">
            {category.items.map((item) => (
              <Card
                key={item.id}
                className="overflow-hidden cursor-pointer hover:shadow-md transition"
                onClick={() => toggleItem(item.id)}
              >
                <div className="p-4">
                  <div className="flex items-center justify-between">
                    <h4 className="font-medium">{item.question}</h4>
                    <svg
                      className={`w-5 h-5 transition-transform ${
                        expandedItems.has(item.id) ? "rotate-180" : ""
                      }`}
                      fill="none"
                      stroke="currentColor"
                      viewBox="0 0 24 24"
                    >
                      <path
                        strokeLinecap="round"
                        strokeLinejoin="round"
                        strokeWidth={2}
                        d="M19 9l-7 7-7-7"
                      />
                    </svg>
                  </div>

                  {expandedItems.has(item.id) && (
                    <p className="mt-3 text-gray-600 text-sm leading-relaxed">
                      {item.answer}
                    </p>
                  )}
                </div>
              </Card>
            ))}
          </div>
        </div>
      ))}
    </div>
  );
}
