/**
 * Spending Chart Component
 * Displays spending trends over time using Recharts
 */

"use client";

import { Card } from "@/shared/components/ui/Card";
import type { SpendingTrend } from "../types";
import { LineChart, Line, XAxis, YAxis, CartesianGrid, Tooltip, ResponsiveContainer } from "recharts";

interface SpendingChartProps {
  data: SpendingTrend[];
  title?: string;
}

export function SpendingChart({ data, title = "Gastos ao Longo do Tempo" }: SpendingChartProps) {
  // Transform data for chart
  const chartData = data.map((item) => ({
    date: new Date(item.period).toLocaleDateString("pt-BR", {
      day: "2-digit",
      month: "short",
    }),
    amount: item.amount,
  }));

  return (
    <Card className="p-6">
      <h3 className="text-lg font-semibold mb-4">{title}</h3>
      <ResponsiveContainer width="100%" height={300}>
        <LineChart data={chartData}>
          <CartesianGrid strokeDasharray="3 3" stroke="#e5e7eb" />
          <XAxis
            dataKey="date"
            tick={{ fontSize: 12 }}
            stroke="#6b7280"
          />
          <YAxis
            tick={{ fontSize: 12 }}
            stroke="#6b7280"
            tickFormatter={(value) => `R$ ${value}`}
          />
          <Tooltip
            formatter={(value: number | undefined) => 
              value !== undefined ? [`R$ ${value.toFixed(2)}`, "Gasto"] : ["", ""]
            }
            contentStyle={{
              backgroundColor: "#fff",
              border: "1px solid #e5e7eb",
              borderRadius: "8px",
            }}
          />
          <Line
            type="monotone"
            dataKey="amount"
            stroke="#3b82f6"
            strokeWidth={2}
            dot={{ fill: "#3b82f6", r: 4 }}
            activeDot={{ r: 6 }}
          />
        </LineChart>
      </ResponsiveContainer>
    </Card>
  );
}
