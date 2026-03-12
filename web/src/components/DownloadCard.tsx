import type { DownloadCardProps } from "@/types";
import { DownloadButton } from "./DownloadButton";

export const DownloadCard = ({ state, onDownload }: DownloadCardProps) => {
  const hasMessage = state !== "idle";

  return (
    <div className="bg-white/10 backdrop-blur-lg rounded-3xl border border-white/20 p-8 md:p-10 shadow-2xl transition-all duration-500 ease-in-out hover:bg-white/15">
      {/* 설명 텍스트 */}
      <div className="mb-6">
        <div className="text-white/80 text-balance text-lg mb-2 transition-colors duration-300">
          <span>사지방에서 매번</span>
          <span>설치하기 귀찮아.</span>
        </div>
        <p className="text-white/60 text-balance text-sm transition-colors duration-300">
          그래서준비했습니다당신을위한원클릭프로그램
        </p>
      </div>

      {/* 다운로드 버튼 */}
      <DownloadButton state={state} onDownload={onDownload} />

      {/* 상태별 추가 메시지 */}
      <div
        className={`overflow-hidden transition-all duration-500 ease-in-out ${
          hasMessage ? "max-h-20 mt-2" : "max-h-0 mt-0"
        }`}
      >
        <div className="flex items-center justify-center">
          {state === "downloading" && (
            <p className="text-white/60 text-sm animate-in fade-in duration-300">
              잠시만 기다려주세요...
            </p>
          )}
          {state === "completed" && (
            <p className="text-green-400 text-sm animate-in slide-in-from-bottom duration-500">
              다운로드가 완료되었습니다! 🎉
            </p>
          )}
          {state === "error" && (
            <div className="flex flex-col md:flex-row md:gap-1 text-red-400 text-sm text-balance animate-in slide-in-from-bottom duration-300">
              <span>다운로드 중 문제가 발생했습니다.</span>
              <span>다시 시도해주세요.</span>
            </div>
          )}
        </div>
      </div>
    </div>
  );
};
