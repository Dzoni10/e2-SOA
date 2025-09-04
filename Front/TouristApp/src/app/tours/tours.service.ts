import { HttpClient } from '@angular/common/http';
import { Injectable } from '@angular/core';
import { Tour } from './model/tour.model';
import { Review } from './model/review.model'
import { Observable } from 'rxjs';

@Injectable({
  providedIn: 'root'
})
export class ToursService {

  private apiUrl = 'http://localhost:8081/tours'
  constructor(private http: HttpClient) { }

  createTour(tour: Tour): Observable<Tour>{
    return this.http.post<Tour>(this.apiUrl,tour);
  }

  getAllTours(): Observable<Tour[]>{
    return this.http.get<Tour[]>(`${this.apiUrl}/all`);
  }

  addReviewWithImages(review: Review, files: File[]): Observable<Review> {
    const formData = new FormData();
    formData.append('tourId', review.tourId);
    formData.append('userId', review.userId.toString());
    formData.append('username', review.username);
    formData.append('rating', review.rating.toString());
    formData.append('comment', review.comment);

    files.forEach(file => formData.append('images', file));

    return this.http.post<Review>(`http://localhost:8081/reviews`, formData);
  }

  getReviewsForTour(tourId: string): Observable<Review[]> {
    return this.http.get<Review[]>(`http://localhost:8081/tours/${tourId}/reviews`);
  }
}
