"use client";

import { BackgroundDecorations } from "@/components/BackgroundDecorations";
import { Header } from "@/components/Header";
import { DownloadCard } from "@/components/DownloadCard";
import { Footer } from "@/components/Footer";
import { useDownload } from "@/hooks/useDownload";

export default function App() {
  const { state, startDownload } = useDownload();

  return (
    <div className="min-h-screen bg-gradient-to-br from-slate-900 via-purple-900 to-slate-900 flex items-center justify-center p-4 relative overflow-hidden">
      <BackgroundDecorations />

      <div className="relative z-10 text-center mx-auto">
        <Header />
        <DownloadCard state={state} onDownload={startDownload} />
        <Footer />
      </div>
    </div>
  );
}
