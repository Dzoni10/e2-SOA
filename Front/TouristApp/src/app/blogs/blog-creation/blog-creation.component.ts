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
  imageList: string[] = [];
  readonly templateImages = signal(this.imageList);

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
      images: [[]],
      creatorID: [Number(user.userId)]
    });
  }

  create(): void {
    if (this.blogCreationForm.invalid) {
      this.snackBar.open("You must fill all required fields!", "Close", {duration: 3000, horizontalPosition: "center"});
      return;
    }

    const newBlog: Blog = {
      ...this.blogCreationForm.value,
      images: this.imageList
    };

    this.blogService.createBlog(newBlog).subscribe({
      next: () => {
        this.snackBar.open("Blog created successfully", "Close", {duration: 3000, horizontalPosition: "center"});
        this.blogCreationForm.reset();
        this.imageList = [];
        this.templateImages.set([]);
        this.router.navigate(['/blogsList']);
      },
      error: (err) => {
        console.error('Error creating blog', err);
        this.snackBar.open("Cannot create blog!", "Close", {duration: 3000, horizontalPosition: "center"});
      }
    });
  }

  addImage(event: MatChipInputEvent): void {
    const value = (event.value || '').trim();

    if (value) {
      this.templateImages.update(images => [...images, value]);
      this.imageList.push(value);
    }
    event.chipInput!.clear();
  }

  removeImage(imageUrl: string) {
    this.templateImages.update(images => {
      const index = images.indexOf(imageUrl);
      if (index >= 0) {
        images.splice(index, 1);
        this.imageList = images;
      }
      return [...images];
    });
  }
}