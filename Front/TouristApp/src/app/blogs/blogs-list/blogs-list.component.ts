import { Component, OnInit } from '@angular/core';
import { Blog } from '../model/blog.model';
import { BlogsService } from '../blogs.service';
import { AuthService } from 'src/app/auth/auth.service';
import { MatSnackBar } from '@angular/material/snack-bar';

@Component({
  selector: 'app-blogs-list',
  templateUrl: './blogs-list.component.html',
  styleUrls: ['./blogs-list.component.css']
})
export class BlogsListComponent implements OnInit {

  blogs: Blog[] = [];
  loading = true;

  constructor(
    private blogService: BlogsService, 
    private authService: AuthService, 
    private snackBar: MatSnackBar
  ) {}

  ngOnInit(): void {
    const user = this.authService.getCurrentUser();
    
    if (!user || !user.userId) {
        this.snackBar.open("You must be logged in to see blogs", "Close", {duration: 3000, horizontalPosition: "center"});
        this.loading = false;
        return;
    }

    this.blogService.getAllBlogs().subscribe({
      next: (data) => {
        this.blogs = data;
        this.loading = false;
      },
      error: (err) => {
        console.error('Error during loading blogs', err);
        this.snackBar.open("Error loading blogs", "Close", {duration: 3000, horizontalPosition: "center"});
        this.loading = false;
      }
    });
  }

  getPreview(description: string): string {
    // Remove markdown formatting and limit text
    const plainText = description
      .replace(/#{1,6}\s?/g, '') // Remove headers
      .replace(/\*\*(.*?)\*\*/g, '$1') // Remove bold
      .replace(/\*(.*?)\*/g, '$1') // Remove italic
      .replace(/\[(.*?)\]\(.*?\)/g, '$1'); // Remove links
    
    return plainText.length > 100 ? plainText.substring(0, 100) + '...' : plainText;
  }

  getImageUrl(imagePath: string): string {
    return this.blogService.getImageUrl(imagePath);
  }

  onImageError(event: any): void {
    event.target.style.display = 'none';
  }
}