import { HttpClient } from '@angular/common/http';
import { Injectable } from '@angular/core';
import { Tour } from './model/tour.model';
import { Review } from './model/review.model'
import { Observable } from 'rxjs';
import { Keypoint } from './model/keypoint.model';
import { Position } from './model/position.model';
import {  HttpParams } from '@angular/common/http';
import { TourStatusInfo, StatusChangeRequest } from './model/tour.model';
@Injectable({
  providedIn: 'root'
})
export class ToursService {

  private apiUrl = 'http://localhost:8070/tours'
  private api = 'http://localhost:8070'
  constructor(private http: HttpClient) { }

  createTour(tour: Tour): Observable<Tour>{
    return this.http.post<Tour>(this.apiUrl,tour);
  }

  /*
  createKeypoint(keypoint: Keypoint): Observable<Keypoint>{
    return this.http.post<Keypoint>(`${this.apiUrl}/keypoints`,keypoint);
  }
    */

  getAllTours(): Observable<Tour[]>{
    return this.http.get<Tour[]>(`${this.apiUrl}/all`);
  }

  getTourByID(tourId: string): Observable<Tour>{
    return this.http.get<Tour>(`${this.apiUrl}/${tourId}`);
  }

  addReviewWithImages(review: Review, files: File[]): Observable<Review> {
    const formData = new FormData();
    formData.append('tourId', review.tourId);
    formData.append('userId', review.userId.toString());
    formData.append('username', review.username);
    formData.append('rating', review.rating.toString());
    formData.append('comment', review.comment);

    files.forEach(file => formData.append('images', file));

    return this.http.post<Review>(`http://localhost:8070/reviews`, formData);
  }

  getReviewsForTour(tourId: string): Observable<Review[]> {
    return this.http.get<Review[]>(`http://localhost:8070/tours/${tourId}/reviews`);
  }

  getKeyPointsForTour(tourId: string): Observable<Keypoint[]>{
    return this.http.get<Keypoint[]>(`${this.apiUrl}/${tourId}/keypoints`);
  }

  // Keypoint operations
  createKeypoint(keypointFormData: FormData): Observable<{ id: string, name: string, order: number }> {
    return this.http.post<{ id: string, name: string, order: number }>(`${this.apiUrl}/keypoints`, keypointFormData);
  }

  updateKeypoint(id: string, formData: FormData): Observable<Keypoint> {
  return this.http.put<Keypoint>(`${this.apiUrl}/keypoints/${id}`, formData);
}
/*
createKeypoint(formData: FormData): Observable<Keypoint> {
  return this.http.post<Keypoint>(`${this.apiUrl}/keypoints`, formData);
}
  

getKeypointsByTourId(tourId: string): Observable<Keypoint[]> {
  return this.http.get<Keypoint[]>(`${this.apiUrl}/tours/${tourId}/keypoints`);
}
  */
  /*
  createKeypoint(keypoint: Keypoint): Observable<{ id: string; name: string; order: number }> {
  let payload: FormData;

  if (keypoint.formData) {
    // Use the provided FormData (e.g., already has files)
    payload = keypoint.formData;
  } else {
    // Build FormData from Keypoint object
    payload = new FormData();
    payload.append('name', keypoint.name);
    payload.append('description', keypoint.description || '');
    payload.append('latitude', keypoint.latitude.toString());
    payload.append('longitude', keypoint.longitude.toString());
    if (keypoint.tourId) payload.append('tourId', keypoint.tourId);

    // If images array contains File objects, append them
    (keypoint.images || []).forEach(img => {
      if (img instanceof File) payload.append('images', img);
    });
  }

  return this.http.post<{ id: string; name: string; order: number }>(
    `${this.apiUrl}/keypoints`,
    payload
  );
  }
  */

  bulkUpdateKeypointsTourId(bulkUpdateData: { keypointIds: string[], tourId: string | undefined}): Observable<{ message: string, updatedCount: number }> {
    return this.http.put<{ message: string, updatedCount: number }>(`${this.apiUrl}/keypoints/bulk-update-tour-id`, bulkUpdateData);
  }

 updateTourLength(tourId: string): Observable<{ status: string,tourLength: number,walkingTime: number,bicycleTime: number,carTime: number}> 
{
  return this.http.put<{status: string,tourLength: number,walkingTime: number,bicycleTime: number, carTime: number}>(`${this.apiUrl}/${tourId}/update-length`, {});
}

calculateTourMetrics(keypoints: Keypoint[]): Observable<{
  tourLength: number,
  walkingTime: number,
  bicycleTime: number,
  carTime: number
}> {
  return this.http.post<{
    tourLength: number,
    walkingTime: number,
    bicycleTime: number,
    carTime: number
  }>(`${this.apiUrl}/calculate-metrics`, keypoints);
}

  getPosition(userId: number): Observable<Position>{
    return this.http.get<Position>(`${this.api}/position/${userId}`);
  }

  initializePosition(userId: number, position: Position): Observable<Position>{
    return this.http.post<Position>(`${this.api}/position/${userId}/create`, position);
  }

  updatePosition(userId: number, position: Position): Observable<Position>{
    return this.http.put<Position>(`${this.api}/position/${userId}/update`, position);
  }

  getTourStatusInfo(tourId: string, creatorId: number): Observable<TourStatusInfo> {
  const params = new HttpParams().set('creatorId', creatorId.toString());
  return this.http.get<TourStatusInfo>(`${this.apiUrl}/${tourId}/status-info`, { params });
}

updateTourStatus(tourId: string, statusRequest: StatusChangeRequest): Observable<any> {
  return this.http.put<any>(`${this.apiUrl}/${tourId}/status`, statusRequest);
}

publishTour(tourId: string, creatorId: number): Observable<any> {
  return this.updateTourStatus(tourId, { status: 1, creatorId });
}

archiveTour(tourId: string, creatorId: number): Observable<any> {
  return this.updateTourStatus(tourId, { status: 2, creatorId });
}

reactivateTour(tourId: string, creatorId: number): Observable<any> {
  return this.updateTourStatus(tourId, { status: 1, creatorId });
}

updateTourCost(tourId: string, cost: number): Observable<any> {
  return this.http.put<any>(`${this.apiUrl}/${tourId}/cost`, { cost });
}

}
