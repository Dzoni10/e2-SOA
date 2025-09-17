export enum ExecutionStatus{
    COMPLETED='COMPLETED',
    ABANDONED = 'ABANDONED',
    ONGOING = 'ONGOING'
}

export interface CompletedKeyPoint {
  keyPointId: string;
  reachedAt: string; // ISO datetime string
}

export interface TourExecution {
  id?: string;
  tourId: string;
  touristId: number;
  tourStartDate: string;
  tourEndDate?: string;
  lastActivity: string;
  status: ExecutionStatus;
  completedPercent: number;
  completedKeyPoints: CompletedKeyPoint[];
  currentLatitude: number;
  currentLongitude: number;
}