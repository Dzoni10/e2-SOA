export interface Review {
  id?: string;
  tourId: string;
  userId: number;
  username: string;
  rating: number;
  comment: string;
  visitedAt?: Date;
  createdAt?: Date;
  image?: string;
}
