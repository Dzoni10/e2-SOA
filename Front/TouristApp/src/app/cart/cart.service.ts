import { Injectable } from '@angular/core';
import { HttpClient } from '@angular/common/http';
import { Observable } from 'rxjs';

@Injectable({
  providedIn: 'root'
})
export class CartService {
  private apiUrl = 'http://localhost:8070/purchase/cart';

  constructor(private http: HttpClient) {}

  addToCart(userId: string, tour: any): Observable<any> {
    return this.http.post(`${this.apiUrl}/${userId}/add`, {
      tourId: tour.id,
      name: tour.name,
      price: tour.cost
    });
  }

  removeFromCart(userId: string, tourId: string): Observable<any> {
    return this.http.post(`${this.apiUrl}/${userId}/remove`, { tourId });
  }

  getCart(userId: string): Observable<any> {
    return this.http.get(`${this.apiUrl}/${userId}`);
  }

  checkout(userId: string): Observable<any> {
    return this.http.post(`${this.apiUrl}/${userId}/checkout`, {});
  }
  getPurchasedTours(userId: string) {
    return this.http.get<any[]>(`${this.apiUrl}/${userId}/tokens`);
  }
}
