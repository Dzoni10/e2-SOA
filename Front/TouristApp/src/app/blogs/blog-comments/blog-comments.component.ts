import { Component, OnInit } from '@angular/core';
import { ActivatedRoute, Router } from '@angular/router';
import { CommentService } from '../comment.service';
import { StakeholdersService } from 'src/app/stakeholders/stakeholders.service';
import { AuthService } from 'src/app/auth/auth.service'; 
import { MatSnackBar } from '@angular/material/snack-bar';
import { Comment } from '../model/comment.model';

@Component({
  selector: 'app-blog-comments',
  templateUrl: './blog-comments.component.html',
  styleUrls: ['./blog-comments.component.css']
})
export class BlogCommentsComponent implements OnInit {
  blogId!: string;
  comments: Comment[] = [];
  newCommentText: string = '';
  loading = true;
  currentUserId!: number;

  constructor(
    private route: ActivatedRoute,
    private router: Router,
    private commentService: CommentService,
    private authService: AuthService,
    private stakeholdersService: StakeholdersService,
    private snackBar: MatSnackBar
  ) {}

  ngOnInit(): void {
    const user = this.authService.getCurrentUser();
    if (!user) {
      this.snackBar.open("You must be logged in", "Close", {duration: 3000});
      this.router.navigate(['/login']);
      return;
    }
    
    this.currentUserId = Number(user.userId);
    this.blogId = this.route.snapshot.paramMap.get('id')!;
    this.loadComments();
  }

  loadComments(): void {
    this.commentService.getComments(this.blogId, this.currentUserId).subscribe({
      next: (data) => {
        this.comments = data;
        this.loading = false;
      },
      error: (err) => {
        this.loading = false;
        if (err.status === 403) {
          this.snackBar.open("You don't have permission to view these comments. Follow this user first.", "Close", {duration: 5000});
          this.router.navigate(['/blogsList']);
        } else if (err.status === 401) {
          this.snackBar.open("You must be logged in to view comments", "Close", {duration: 3000});
          this.router.navigate(['/login']);
        } else {
          console.error('Error loading comments:', err);
          this.snackBar.open("Error loading comments", "Close", {duration: 3000});
        }
      }
    });
  }

  addComment(): void {
    if (!this.newCommentText.trim()) {
      this.snackBar.open("Comment cannot be empty", "Close", {duration: 3000});
      return;
    }

    const user = this.authService.getCurrentUser();
    if (!user) {
      this.snackBar.open("You must be logged in to comment", "Close", {duration: 3000});
      return;
    }

    // na osnovu userId povuci podatke o korisniku
    this.stakeholdersService.getUserById(user.userId).subscribe({
      next: (fetchedUser) => {
        const newComment: Partial<Comment> = {
          userId: fetchedUser.id,
          username: fetchedUser.username,
          text: this.newCommentText
        };

        this.commentService.addComment(this.blogId, newComment).subscribe({
          next: () => {
            this.loadComments();
            this.newCommentText = '';
            this.snackBar.open("Comment added successfully", "Close", {duration: 3000});
          },
          error: (err) => {
            if (err.status === 403) {
              this.snackBar.open("You must follow this user to comment on their blog", "Close", {duration: 5000});
            } else {
              console.error('Error adding comment:', err);
              this.snackBar.open("Failed to add comment", "Close", {duration: 3000});
            }
          }
        });
      },
      error: (err) => {
        console.error('Failed to fetch user info:', err);
        this.snackBar.open("Failed to fetch user information", "Close", {duration: 3000});
      }
    });
  }
}