import { HttpClient } from '@angular/common/http';
import { Injectable } from '@angular/core';
import { Blog } from './model/blog.model';
import { Observable } from 'rxjs';

@Injectable({
  providedIn: 'root'
})
export class BlogsService {

  private apiUrl = 'http://localhost:8082/blogs'
  
  constructor(private http: HttpClient) { }

  createBlog(blog: Blog): Observable<Blog>{
    return this.http.post<Blog>(this.apiUrl, blog);
  }

  getAllBlogs(): Observable<Blog[]>{
    return this.http.get<Blog[]>(`${this.apiUrl}/all`);
  }

  getBlogsByCreator(creatorId: number): Observable<Blog[]>{
    return this.http.get<Blog[]>(`${this.apiUrl}/creator/${creatorId}`);
  }

  getBlog(id: string): Observable<Blog>{
    return this.http.get<Blog>(`${this.apiUrl}/${id}`);
  }
}