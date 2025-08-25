export interface Blog {
    id?: string,
    title: string,
    description: string, // Markdown content
    createdAt?: Date,
    images?: string[],
    creatorID: number
}