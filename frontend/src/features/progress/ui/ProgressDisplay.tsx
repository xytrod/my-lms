import type { ProgressDto } from '../api/types'

export function ProgressDisplay({ progress }: { progress: ProgressDto }) {
  const percentage = Math.min(100, Math.max(0, progress.percentage))
  return <div className="progress-display">
    <div className="progress-copy"><span>{progress.completed_lessons} из {progress.total_lessons} уроков завершено</span><strong>{Math.round(percentage)}%</strong></div>
    <div className="progress-track" role="progressbar" aria-label="Прогресс курса" aria-valuemin={0} aria-valuemax={100} aria-valuenow={Math.round(percentage)}>
      <span style={{ width: `${percentage}%` }} />
    </div>
  </div>
}

