import { useState } from "react";
import type { DownloadState } from "@/types";
import { TIMINGS } from "@/constants";
import { downloadFile } from "@/utils/download";

export const useDownload = () => {
  const [state, setState] = useState<DownloadState>("idle");

  const startDownload = async () => {
    setState("downloading");

    try {
      await downloadFile();
      setState("completed");
      setTimeout(() => setState("idle"), TIMINGS.COMPLETED_RESET);
    } catch (error) {
      setState("error");
      setTimeout(() => setState("idle"), TIMINGS.ERROR_RESET);
    }
  };

  return { state, startDownload };
};
