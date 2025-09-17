export interface Tour
{
    id?: string,
    name: string,
    description: string,
    difficulty: number,
    tags: string[],
    status: number | string,
    cost: number,
    tourLength: number,
    creatorID: number,
    publishedAt?: Date;
    archivedAt?: Date;
    createdAt?: Date;
    updatedAt?: Date;
    walkingTime?: number;
    bicycleTime?: number;
    carTime?: number;
}



export interface TourStatusInfo {
    currentStatus: string;
    canEdit: boolean;
    publishedAt?: Date;
    archivedAt?: Date;
    createdAt: Date;
    updatedAt: Date;
}

export interface StatusChangeRequest {
    status: number;
    creatorId: number;
}

export enum TourStatus {
    Draft = 0,
    Published = 1,
    Archived = 2
}
