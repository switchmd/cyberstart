export type DownloadState = "idle" | "downloading" | "completed" | "error";

export interface DownloadButtonProps {
  state: DownloadState;
  onDownload: () => void;
}

export interface DownloadCardProps {
  state: DownloadState;
  onDownload: () => void;
}
