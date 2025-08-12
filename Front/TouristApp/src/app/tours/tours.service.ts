import { HttpClient } from '@angular/common/http';
import { Injectable } from '@angular/core';
import { Tour } from './model/tour.model';
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
}
