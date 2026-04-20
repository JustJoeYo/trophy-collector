import clsx from 'clsx'

interface SkeletonProps {
  className?: string
}

// Reusable skeleton loader — use instead of spinners
export default function Skeleton({ className }: SkeletonProps) {
  return <div className={clsx('skeleton', className)} />
}

export function StatCardSkeleton() {
  return (
    <div className="card flex flex-col gap-2">
      <Skeleton className="h-3 w-20" />
      <Skeleton className="h-8 w-28" />
    </div>
  )
}

export function MatchRowSkeleton() {
  return (
    <div className="flex items-center gap-4 p-4 border-b border-surface-600">
      <Skeleton className="h-10 w-10 rounded-lg" />
      <div className="flex flex-col gap-1.5 flex-1">
        <Skeleton className="h-4 w-32" />
        <Skeleton className="h-3 w-20" />
      </div>
      <Skeleton className="h-6 w-16 rounded-full" />
    </div>
  )
}
