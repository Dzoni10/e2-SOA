export interface Blog {
    id?: string;
    title: string;
    description: string; // Markdown content
    createdAt?: Date;
    images?: string[];
    creatorID: number;
    username?: string;  
    // stiglo iz backenda
    likes?: { userId: number, createdAt: Date }[];

    // frontend computed polja
    likedByUser?: boolean;
    likesCount?: number;
}
