import { Component, OnInit } from '@angular/core';
import { ActivatedRoute } from '@angular/router';
import { CommentService } from '../comment.service';
import { StakeholdersService } from 'src/app/stakeholders/stakeholders.service';
import { AuthService } from 'src/app/auth/auth.service'; 
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

  constructor(
    private route: ActivatedRoute,
    private commentService: CommentService,
    private authService: AuthService,
    private stakeholdersService: StakeholdersService
  ) {}

  ngOnInit(): void {
    this.blogId = this.route.snapshot.paramMap.get('id')!;
    this.loadComments();
  }

  loadComments(): void {
    this.commentService.getComments(this.blogId).subscribe({
      next: (data) => this.comments = data,
      error: (err) => console.error('Error loading comments:', err)
    });
  }
addComment(): void {
  if (!this.newCommentText.trim()) return;

  const user = this.authService.getCurrentUser();
  if (!user) {
    console.error('No logged in user found');
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
      },
        error: (err) => console.error('Error adding comment:', err)
      });
    },
    error: (err) => console.error('Failed to fetch user info:', err)
  });
}
}
