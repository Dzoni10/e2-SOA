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
   updateBlog(blogId: string, blogData: { title: string, description: string, creatorID: number }): Observable<Blog> {
    return this.http.put<Blog>(`${this.apiUrl}/${blogId}`, blogData);
  }

  addImageToBlog(blogId: string, image: File, creatorID: number): Observable<{imageUrl: string, message: string}> {
    const formData = new FormData();
    formData.append('image', image);
    formData.append('creatorID', creatorID.toString());
    
    return this.http.post<{imageUrl: string, message: string}>(`${this.apiUrl}/${blogId}/add-image`, formData);
  }

  removeImageFromBlog(blogId: string, imageUrl: string, creatorID: number): Observable<{message: string}> {
    const requestBody = {
      imageUrl: imageUrl,
      creatorID: creatorID
    };
    
    return this.http.delete<{message: string}>(`${this.apiUrl}/${blogId}/remove-image`, { body: requestBody });
  }

  getImageUrl(imagePath: string): string {
    return `http://localhost:8082${imagePath}`;
  }
}