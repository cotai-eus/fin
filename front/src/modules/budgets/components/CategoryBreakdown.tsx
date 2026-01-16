/**
 * Category Breakdown Component
 * Pie/Donut chart showing spending by category
 */

"use client";

import { Card } from "@/shared/components/ui/Card";
import type { CategorySpending } from "../types";
import { CATEGORY_LABELS, CATEGORY_ICONS } from "../validators";
import { PieChart, Pie, Cell, ResponsiveContainer, Legend, Tooltip } from "recharts";

interface CategoryBreakdownProps {
  data: CategorySpending[];
  title?: string;
}

const COLORS = [
  "#3b82f6", // blue
  "#10b981", // green
  "#f59e0b", // amber
  "#ef4444", // red
  "#8b5cf6", // purple
  "#ec4899", // pink
  "#06b6d4", // cyan
  "#6b7280", // gray
];

export function CategoryBreakdown({
  data,
  title = "Gastos por Categoria",
}: CategoryBreakdownProps) {
  // Transform data for pie chart
  const chartData = data.map((item) => ({
    name: CATEGORY_LABELS[item.category],
    value: item.spent,
    percentage: item.percentageOfTotal,
  }));

  return (
    <Card className="p-6">
      <h3 className="text-lg font-semibold mb-4">{title}</h3>

      <ResponsiveContainer width="100%" height={300}>
        <PieChart>
          <Pie
            data={chartData}
            cx="50%"
            cy="50%"
            labelLine={false}
            label={({ name, percent }) =>
              `${name}: ${percent !== undefined ? (percent * 100).toFixed(0) : 0}%`
            }
            outerRadius={80}
            fill="#8884d8"
            dataKey="value"
          >
            {chartData.map((entry, index) => (
              <Cell
                key={`cell-${index}`}
                fill={COLORS[index % COLORS.length]}
              />
            ))}
          </Pie>
          <Tooltip
            formatter={(value: number | undefined) => 
              value !== undefined ? `R$ ${value.toFixed(2)}` : ""
            }
            contentStyle={{
              backgroundColor: "#fff",
              border: "1px solid #e5e7eb",
              borderRadius: "8px",
            }}
          />
        </PieChart>
      </ResponsiveContainer>

      {/* Category List */}
      <div className="mt-4 space-y-2">
        {data.map((item, index) => (
          <div
            key={item.category}
            className="flex items-center justify-between text-sm"
          >
            <div className="flex items-center gap-2">
              <div
                className="w-3 h-3 rounded-full"
                style={{ backgroundColor: COLORS[index % COLORS.length] }}
              />
              <span>{CATEGORY_ICONS[item.category]}</span>
              <span className="font-medium">
                {CATEGORY_LABELS[item.category]}
              </span>
            </div>
            <div className="text-right">
              <span className="font-semibold">R$ {item.spent.toFixed(2)}</span>
              <span className="text-gray-500 ml-2 text-xs">
                ({item.percentageOfTotal.toFixed(1)}%)
              </span>
            </div>
          </div>
        ))}
      </div>
    </Card>
  );
}
