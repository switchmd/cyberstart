import type { DownloadState } from "@/types";
import { DOWNLOAD_FILE, TIMINGS } from "@/constants";

export const getDownloadUrl = (): string =>
  `https://rnseo.kr/${DOWNLOAD_FILE}`;

export const getButtonStyles = (state: DownloadState): string => {
  const baseStyles = `
    relative overflow-hidden
    h-16 px-8 text-lg font-semibold
    border-0 rounded-full
    transition-all duration-500 ease-in-out
    hover:scale-105 hover:shadow-2xl hover:shadow-purple-500/25
    disabled:opacity-70 disabled:cursor-not-allowed disabled:hover:scale-100
  `;

  const stateStyles = {
    idle: "bg-purple-600 hover:bg-purple-700",
    downloading: "bg-purple-600",
    completed: "bg-green-600",
    error: "bg-red-600",
  };

  return `${baseStyles} ${stateStyles[state]}`;
};

export const getGlowStyles = (state: DownloadState): string => {
  const baseGlowClass = "absolute inset-0 rounded-2xl blur-lg opacity-50 group-hover:opacity-75 transition-all duration-500 ease-in-out -z-10";
  
  const glowColors = {
    idle: "bg-purple-600",
    downloading: "bg-purple-600",
    completed: "bg-green-600",
    error: "bg-red-600",
  };

  return `${baseGlowClass} ${glowColors[state]}`;
};

export const getInlineStyles = (state: DownloadState): React.CSSProperties => {
  const gradients = {
    idle: 'linear-gradient(to right, #9333ea, #db2777)',
    downloading: 'linear-gradient(to right, #9333ea, #db2777)',
    completed: 'linear-gradient(to right, #059669, #047857)',
    error: 'linear-gradient(to right, #dc2626, #ea580c)',
  };

  return {
    background: gradients[state],
  };
};

export async function downloadFile(): Promise<void> {
  await new Promise((resolve) => setTimeout(resolve, TIMINGS.DOWNLOAD_DELAY));

  const downloadUrl = getDownloadUrl();
  const response = await fetch(downloadUrl, {
    method: "GET",
    headers: {
      "Cache-Control": "no-cache",
      Pragma: "no-cache",
      Expires: "0",
    },
  });

  // @ts-ignore - IE11 호환성
  if (window.navigator && window.navigator.msSaveOrOpenBlob) {
    const blob = await response.blob();
    // @ts-ignore
    window.navigator.msSaveOrOpenBlob(blob, DOWNLOAD_FILE);
    return;
  }

  if (!response.ok) {
    throw new Error("네트워크 오류.");
  }

  const blob = await response.blob();
  const url = window.URL.createObjectURL(blob);
  const a = document.createElement("a");

  a.href = url;
  a.download = DOWNLOAD_FILE;
  document.body.appendChild(a);
  a.click();
  a.remove();

  window.URL.revokeObjectURL(url);
}
