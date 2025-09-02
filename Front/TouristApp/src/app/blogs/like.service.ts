import { Injectable } from '@angular/core';
import { HttpClient } from '@angular/common/http';
import { Observable } from 'rxjs';

@Injectable({
  providedIn: 'root'
})
export class LikeService {
  private apiUrl = 'http://localhost:8082/blogs';

  constructor(private http: HttpClient) {}

  likeBlog(blogId: string, userId: number): Observable<any> {
    return this.http.post(`http://localhost:8082/blogs/${blogId}/likes`, { userId });
  }

  unlikeBlog(blogId: string, userId: number): Observable<any> {
    return this.http.request('delete', `http://localhost:8082/blogs/${blogId}/likes`, { body: { userId } });
  }

  countLikes(blogId: string): Observable<{ count: number }> {
    return this.http.get<{ count: number }>(`http://localhost:8082/blogs/${blogId}/likes/count`);
  }

  hasUserLiked(blogId: string, userId: number): Observable<{ liked: boolean }> {
    return this.http.get<{ liked: boolean }>(`http://localhost:8082/blogs/${blogId}/likes/${userId}`);
  }
}
