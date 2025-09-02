import { Injectable } from '@angular/core';
import { HttpClient } from '@angular/common/http';
import { Observable } from 'rxjs';
import { Comment } from './model/comment.model';

@Injectable({
  providedIn: 'root'
})
export class CommentService {
  private baseUrl = 'http://localhost:8082'; // Blog backend

  constructor(private http: HttpClient) {}

  getComments(blogId: string): Observable<Comment[]> {
    return this.http.get<Comment[]>(`${this.baseUrl}/blogs/${blogId}/comments`);
  }

  addComment(blogId: string, comment: Partial<Comment>): Observable<Comment> {
    return this.http.post<Comment>(`${this.baseUrl}/blogs/${blogId}/comments`, comment);
  }
}
