import { Button } from "@/components/ui/button";
import type { DownloadButtonProps } from "@/types";
import { COLORS } from "@/constants";
import {
  getButtonStyles,
  getGlowStyles,
  getInlineStyles,
} from "@/utils/download";
import { getButtonContent } from "./ButtonContent";

export const DownloadButton = ({ state, onDownload }: DownloadButtonProps) => {
  const { icon, text } = getButtonContent(state);

  return (
    <div className="flex justify-center">
      <div className="relative">
        <Button
          onClick={onDownload}
          disabled={state === "downloading"}
          className={`${getButtonStyles(state)} group`}
          style={getInlineStyles(state)}
        >
          {/* 버튼 배경 효과 */}
          <div className="absolute inset-0 rounded-full bg-gradient-to-r from-white/20 to-transparent opacity-0 group-hover:opacity-100 transition-opacity duration-300"></div>

          {/* 버튼 내용 */}
          <div
            className="relative flex items-center gap-3 font-light"
            style={{ color: COLORS.white }}
          >
            <div className="transition-all duration-300 ease-in-out">
              {icon}
            </div>
            <span className="transition-all duration-300 ease-in-out">
              {text}
            </span>
          </div>

          {/* 버튼 글로우 효과 */}
          <div className={getGlowStyles(state)}></div>
        </Button>

        {/* 성공/에러 시 추가 효과 */}
        {state === "completed" && (
          <div className="absolute inset-0">
            <div className="absolute inset-0 bg-green-400/20 rounded-full animate-ping"></div>
          </div>
        )}

        {state === "error" && (
          <div className="absolute inset-0">
            <div className="absolute inset-0 bg-red-500/50 rounded-full animate-pulse"></div>
          </div>
        )}
      </div>
    </div>
  );
};
