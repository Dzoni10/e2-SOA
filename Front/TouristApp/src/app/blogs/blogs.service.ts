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

  createBlogWithImages(blogData: Omit<Blog, 'id' | 'createdAt' | 'images'>, images: File[]): Observable<Blog> {
    const formData = new FormData();
    
    formData.append('title', blogData.title);
    formData.append('description', blogData.description);
    formData.append('creatorID', blogData.creatorID.toString());
    
    for (let i = 0; i < images.length; i++) {
      formData.append('images', images[i]);
    }
    
    return this.http.post<Blog>(this.apiUrl, formData);
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

  getImageUrl(imagePath: string): string {
    return `http://localhost:8082${imagePath}`;
  }
}