import { Component, OnInit, signal } from '@angular/core';
import { FormBuilder, FormGroup, Validators } from '@angular/forms';
import { BlogsService } from '../blogs.service';
import { Router } from '@angular/router';
import { MatSnackBar } from '@angular/material/snack-bar';
import { Blog } from '../model/blog.model';
import { MatChipInputEvent } from '@angular/material/chips';
import { AuthService } from 'src/app/auth/auth.service';

@Component({
  selector: 'app-blog-creation',
  templateUrl: './blog-creation.component.html',
  styleUrls: ['./blog-creation.component.css']
})
export class BlogCreationComponent implements OnInit {

  blogCreationForm!: FormGroup;
  
  selectedFiles: File[] = [];
  isSubmitting = false;

  constructor(
    private fb: FormBuilder, 
    private blogService: BlogsService, 
    private router: Router, 
    private snackBar: MatSnackBar, 
    private authService: AuthService
  ) {}

  ngOnInit(): void {
    const user = this.authService.getCurrentUser();
    
    if (!user || !user.userId) {
        this.snackBar.open("You must be logged in to create a blog", "Close", {duration: 3000, horizontalPosition: "center"});
        return;
    }
    
    this.blogCreationForm = this.fb.group({
      title: ['', Validators.required],
      description: ['', Validators.required],
      creatorID: [Number(user.userId)]
    });
  }

  onFileSelect(event: any): void {
    const files = event.target.files;
    if (files) {
      for (let i = 0; i < files.length; i++) {
        this.selectedFiles.push(files[i]);
      }
    }
  }

  removeFile(index: number): void {
    this.selectedFiles.splice(index, 1);
  }

  getSelectedFilesText(): string {
    if (this.selectedFiles.length === 0) {
      return '';
    }
    return `${this.selectedFiles.length} file(s) selected`;
  }

  getFilePreview(file: File): string {
    return URL.createObjectURL(file);
  }

  create(): void {
    if (this.blogCreationForm.invalid) {
      this.snackBar.open("You must fill all required fields!", "Close", {duration: 3000, horizontalPosition: "center"});
      return;
    }

    if (this.isSubmitting) return;

    this.isSubmitting = true;

    const blogData: Omit<Blog, 'id' | 'createdAt' | 'images'> = {
      title: this.blogCreationForm.value.title,
      description: this.blogCreationForm.value.description,
      creatorID: this.blogCreationForm.value.creatorID
    };

    this.blogService.createBlogWithImages(blogData, this.selectedFiles).subscribe({
      next: (createdBlog) => {
        this.snackBar.open("Blog created successfully", "Close", {duration: 3000, horizontalPosition: "center"});
        this.blogCreationForm.reset();
        this.selectedFiles = [];
        this.isSubmitting = false;
        this.router.navigate(['/blogsList']);
      },
      error: (err) => {
        console.error('Error creating blog', err);
        this.snackBar.open("Cannot create blog!", "Close", {duration: 3000, horizontalPosition: "center"});
        this.isSubmitting = false;
      }
    });
  }

}