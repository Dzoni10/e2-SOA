export interface Review {
  id?: string;
  tourId: string;
  userId: number;
  username: string;
  rating: number;
  comment: string;
  images?: string[];   
  createdAt?: Date;
  visitedAt?: Date;
}
