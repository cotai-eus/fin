/**
 * Barcode Scanner Component
 * Uses device camera to scan bill barcodes
 */

"use client";

import { useState, useEffect, useRef } from "react";
import { Button } from "@/shared/components/ui/Button";
import { Card } from "@/shared/components/ui/Card";

interface BarcodeScannerProps {
  onScan: (barcode: string) => void;
  onClose: () => void;
}

export function BarcodeScanner({ onScan, onClose }: BarcodeScannerProps) {
  const [scanning, setScanning] = useState(false);
  const [error, setError] = useState<string | null>(null);
  const [hasPermission, setHasPermission] = useState<boolean | null>(null);
  const videoRef = useRef<HTMLVideoElement>(null);
  const canvasRef = useRef<HTMLCanvasElement>(null);
  const streamRef = useRef<MediaStream | null>(null);

  useEffect(() => {
    // Request camera permission
    const requestCamera = async () => {
      try {
        const stream = await navigator.mediaDevices.getUserMedia({
          video: { facingMode: "environment" }, // Use back camera on mobile
        });
        streamRef.current = stream;
        setHasPermission(true);

        if (videoRef.current) {
          videoRef.current.srcObject = stream;
          videoRef.current.play();
        }
      } catch (err) {
        console.error("Camera access error:", err);
        setError("Não foi possível acessar a câmera");
        setHasPermission(false);
      }
    };

    requestCamera();

    // Cleanup
    return () => {
      if (streamRef.current) {
        streamRef.current.getTracks().forEach((track) => track.stop());
      }
    };
  }, []);

  const captureAndScan = () => {
    if (!videoRef.current || !canvasRef.current) return;

    const video = videoRef.current;
    const canvas = canvasRef.current;
    const context = canvas.getContext("2d");

    if (!context) return;

    // Set canvas size to match video
    canvas.width = video.videoWidth;
    canvas.height = video.videoHeight;

    // Draw current video frame to canvas
    context.drawImage(video, 0, 0, canvas.width, canvas.height);

    // In a real implementation, you would use a barcode detection library here
    // For example: html5-qrcode, zxing-js, or quagga.js
    // This is a placeholder
    
    setScanning(true);
    
    // Simulated barcode detection (replace with actual library)
    setTimeout(() => {
      // Mock barcode result
      const mockBarcode = "34191.79001 01043.510047 91020.150008 1 96610000000100";
      onScan(mockBarcode.replace(/\D/g, ""));
      setScanning(false);
    }, 1000);
  };

  if (hasPermission === false) {
    return (
      <Card className="p-6">
        <div className="text-center">
          <p className="text-red-600 mb-4">{error || "Permissão de câmera negada"}</p>
          <Button variant="outline" onClick={onClose}>
            Fechar
          </Button>
        </div>
      </Card>
    );
  }

  return (
    <Card className="p-6">
      <div className="space-y-4">
        <div className="relative">
          <video
            ref={videoRef}
            className="w-full rounded-lg"
            playsInline
            muted
          />
          <canvas ref={canvasRef} className="hidden" />
          
          {/* Scanning overlay */}
          <div className="absolute inset-0 flex items-center justify-center pointer-events-none">
            <div className="w-3/4 h-32 border-2 border-blue-500 rounded-lg">
              <div className="w-full h-1 bg-blue-500 animate-pulse" />
            </div>
          </div>
        </div>

        <div className="text-center text-sm text-gray-600">
          <p>Posicione o código de barras dentro do quadro</p>
          <p className="text-xs mt-1">O scanner detectará automaticamente</p>
        </div>

        {error && (
          <div className="p-3 bg-red-50 border border-red-200 rounded-lg text-red-700 text-sm">
            {error}
          </div>
        )}

        <div className="flex gap-2">
          <Button
            onClick={captureAndScan}
            disabled={scanning}
            className="flex-1"
          >
            {scanning ? "Escaneando..." : "Capturar"}
          </Button>
          <Button variant="outline" onClick={onClose}>
            Cancelar
          </Button>
        </div>

        <p className="text-xs text-gray-500 text-center">
          Nota: Em produção, use uma biblioteca como html5-qrcode ou quagga.js
          para detecção automática de código de barras
        </p>
      </div>
    </Card>
  );
}
