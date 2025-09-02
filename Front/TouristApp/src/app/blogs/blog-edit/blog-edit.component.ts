import { Component, OnInit } from '@angular/core';
import { FormBuilder, FormGroup, Validators } from '@angular/forms';
import { ActivatedRoute, Router } from '@angular/router';
import { MatSnackBar } from '@angular/material/snack-bar';
import { BlogsService } from '../blogs.service';
import { AuthService } from 'src/app/auth/auth.service';
import { Blog } from '../model/blog.model';

@Component({
  selector: 'app-blog-edit',
  templateUrl: './blog-edit.component.html',
  styleUrls: ['./blog-edit.component.css']
})
export class BlogEditComponent implements OnInit {
  editForm!: FormGroup;
  blog: Blog | null = null;
  selectedFiles: File[] = [];
  loading = true;
  isUpdating = false;
  isAddingImage = false;
  isRemovingImage = false;
  currentUser: any;

  constructor(
    private fb: FormBuilder,
    private route: ActivatedRoute,
    private router: Router,
    private blogService: BlogsService,
    private authService: AuthService,
    private snackBar: MatSnackBar
  ) {}

  ngOnInit(): void {
    this.currentUser = this.authService.getCurrentUser();
    
    if (!this.currentUser) {
      this.snackBar.open("You must be logged in to edit blogs", "Close", {duration: 3000, horizontalPosition: "center"});
      this.router.navigate(['/blogsList']);
      return;
    }

    this.editForm = this.fb.group({
      title: ['', Validators.required],
      description: ['', Validators.required]
    });

    const blogId = this.route.snapshot.paramMap.get('id');
    if (blogId) {
      this.loadBlog(blogId);
    }
  }

  loadBlog(blogId: string): void {
    this.blogService.getBlog(blogId).subscribe({
      next: (blog) => {
        // Proverava da li je korisnik vlasnik bloga
        if (blog.creatorID !== Number(this.currentUser.userId)) {
          this.snackBar.open("You can only edit your own blogs", "Close", {duration: 3000, horizontalPosition: "center"});
          this.router.navigate(['/blogsList']);
          return;
        }

        this.blog = blog;
        this.editForm.patchValue({
          title: blog.title,
          description: blog.description
        });
        this.loading = false;
      },
      error: (err) => {
        console.error('Error loading blog:', err);
        this.snackBar.open("Error loading blog", "Close", {duration: 3000, horizontalPosition: "center"});
        this.router.navigate(['/blogsList']);
      }
    });
  }

  updateBlog(): void {
    if (this.editForm.invalid || !this.blog || this.isUpdating) return;

    this.isUpdating = true;
    
    const updateData = {
      title: this.editForm.value.title,
      description: this.editForm.value.description,
      creatorID: Number(this.currentUser.userId)
    };

    this.blogService.updateBlog(this.blog.id!, updateData).subscribe({
      next: (updatedBlog) => {
        this.blog = updatedBlog;
        this.snackBar.open("Blog updated successfully", "Close", {duration: 3000, horizontalPosition: "center"});
        this.isUpdating = false;
      },
      error: (err) => {
        console.error('Error updating blog:', err);
        this.snackBar.open("Failed to update blog", "Close", {duration: 3000, horizontalPosition: "center"});
        this.isUpdating = false;
      }
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

  removeSelectedFile(index: number): void {
    this.selectedFiles.splice(index, 1);
  }

  getSelectedFilesText(): string {
    if (this.selectedFiles.length === 0) {
      return '';
    }
    return this.selectedFiles.length + ' file(s) selected';
  }

  getFilePreview(file: File): string {
    return URL.createObjectURL(file);
  }

  addSingleImage(file: File, index: number): void {
    if (!this.blog || this.isAddingImage) return;

    this.isAddingImage = true;

    this.blogService.addImageToBlog(this.blog.id!, file, Number(this.currentUser.userId)).subscribe({
      next: (response) => {
        if (!this.blog!.images) {
          this.blog!.images = [];
        }
        this.blog!.images.push(response.imageUrl);
        this.selectedFiles.splice(index, 1);
        this.snackBar.open("Image added successfully", "Close", {duration: 3000, horizontalPosition: "center"});
        this.isAddingImage = false;
      },
      error: (err) => {
        console.error('Error adding image:', err);
        this.snackBar.open("Failed to add image", "Close", {duration: 3000, horizontalPosition: "center"});
        this.isAddingImage = false;
      }
    });
  }

  removeImage(imagePath: string): void {
    if (!this.blog || this.isRemovingImage) return;

    this.isRemovingImage = true;

    this.blogService.removeImageFromBlog(this.blog.id!, imagePath, Number(this.currentUser.userId)).subscribe({
      next: (response) => {
        this.blog!.images = this.blog!.images!.filter(img => img !== imagePath);
        this.snackBar.open("Image removed successfully", "Close", {duration: 3000, horizontalPosition: "center"});
        this.isRemovingImage = false;
      },
      error: (err) => {
        console.error('Error removing image:', err);
        this.snackBar.open("Failed to remove image", "Close", {duration: 3000, horizontalPosition: "center"});
        this.isRemovingImage = false;
      }
    });
  }

  getImageUrl(imagePath: string): string {
    return this.blogService.getImageUrl(imagePath);
  }

  cancel(): void {
    this.router.navigate(['/blogsList']);
  }
}