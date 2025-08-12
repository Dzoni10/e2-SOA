export interface Tour
{
    id?: string,
    name: string,
    description: string,
    difficulty: number,
    tags: string,
    status: number,
    cost: number,
    tourLength: number,
    creatorID: number
}