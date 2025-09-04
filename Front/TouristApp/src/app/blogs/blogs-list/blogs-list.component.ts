import { Component, OnInit } from '@angular/core';
import { Blog } from '../model/blog.model';
import { BlogsService } from '../blogs.service';
import { LikeService } from '../like.service';
import { AuthService } from 'src/app/auth/auth.service';
import { StakeholdersService } from 'src/app/stakeholders/stakeholders.service';
import { MatSnackBar } from '@angular/material/snack-bar';
import { Router } from '@angular/router';
import { FollowerService } from 'src/app/followers/follower.service'; 

@Component({
  selector: 'app-blogs-list',
  templateUrl: './blogs-list.component.html',
  styleUrls: ['./blogs-list.component.css']
})
export class BlogsListComponent implements OnInit {

  blogs: Blog[] = [];
  loading = true;
  currentUser: any;

  constructor(
    private blogService: BlogsService, 
    private authService: AuthService,
    private stakeholdersService: StakeholdersService,
    private likeService: LikeService, 
    private snackBar: MatSnackBar,
    private followerService: FollowerService, 
    private router: Router
  ) {}

  ngOnInit(): void {
    const user = this.authService.getCurrentUser();
    this.currentUser = user;

    if (!user || !user.userId) {
      this.snackBar.open("You must be logged in to see blogs", "Close", {duration: 3000, horizontalPosition: "center"});
      this.loading = false;
      return;
    }

    // Kreiraj korisnika u follower servisu ako ne postoji
    this.createUserInFollowerService(user.userId);

    this.loadBlogs();
  }

  private createUserInFollowerService(userId: number): void {
    this.stakeholdersService.getUserById(userId).subscribe({
      next: (fetchedUser) => {
        this.followerService.createUser(fetchedUser.id, fetchedUser.name, fetchedUser.username).subscribe({
          next: () => console.log('User created/updated in follower service'),
          error: (err) => console.log('User might already exist in follower service:', err)
        });
      },
      error: (err) => console.error('Failed to fetch user info:', err)
    });
  }

  private loadBlogs(): void {
    this.blogService.getAllBlogs().subscribe({
      next: (data) => {
        this.blogs = data.map(blog => {
          this.loadLikeState(blog);
          if (!this.isOwnBlog(blog)) {
            this.checkFollowStatus(blog);
            this.checkCommentPermission(blog);
          } else {
            (blog as any).canViewComments = true;
          }
          return blog;
        });
        this.loading = false;
      },
      error: (err) => {
        console.error('Error during loading blogs', err);
        this.snackBar.open("Error loading blogs", "Close", {duration: 3000, horizontalPosition: "center"});
        this.loading = false;
      }
    });
  }

  checkCommentPermission(blog: Blog): void {
    if (!this.currentUser) return;
    
    this.followerService.canUserComment(
      Number(this.currentUser.userId), 
      blog.creatorID
    ).subscribe({
      next: (response) => {
        (blog as any).canViewComments = response.canComment;
      },
      error: (err) => {
        console.error('Error checking comment permission:', err);
        (blog as any).canViewComments = false;
      }
    });
  }

  checkFollowStatus(blog: Blog): void {
    if (!this.currentUser) return;
    
    this.followerService.isFollowing(
      Number(this.currentUser.userId), 
      blog.creatorID
    ).subscribe({
      next: (response) => {
        (blog as any).isFollowing = response.isFollowing;
      },
      error: (err) => console.error('Error checking follow status:', err)
    });
  }

  toggleFollow(blog: Blog): void {
    if (!this.currentUser) return;

    const currentUserId = Number(this.currentUser.userId);
    const isFollowing = (blog as any).isFollowing;

    if (isFollowing) {
      // Unfollow
      this.followerService.unfollowUser(currentUserId, blog.creatorID).subscribe({
        next: () => {
          this.updateFollowStatusForCreator(blog.creatorID, false);
          this.snackBar.open("Unfollowed user", "Close", {duration: 3000});
        },
        error: (err) => {
          console.error('Error unfollowing user:', err);
          this.snackBar.open("Failed to unfollow user", "Close", {duration: 3000});
        }
      });
    } else {
      // Follow
      this.followerService.followUser(currentUserId, blog.creatorID).subscribe({
        next: () => {
          this.updateFollowStatusForCreator(blog.creatorID, true);
          this.snackBar.open("Now following user", "Close", {duration: 3000});
        },
        error: (err) => {
          console.error('Error following user:', err);
          this.snackBar.open("Failed to follow user", "Close", {duration: 3000});
        }
      });
    }
  }

  private updateFollowStatusForCreator(creatorID: number, isFollowing: boolean): void {
    this.blogs.forEach(blogItem => {
      if (blogItem.creatorID === creatorID) {
        (blogItem as any).isFollowing = isFollowing;
        (blogItem as any).canViewComments = isFollowing;
      }
    });
  }

  getPreview(description: string): string {
    const plainText = description
      .replace(/#{1,6}\s?/g, '') 
      .replace(/\*\*(.*?)\*\*/g, '$1') 
      .replace(/\*(.*?)\*/g, '$1') 
      .replace(/\[(.*?)\]\(.*?\)/g, '$1'); 
    
    return plainText.length > 100 ? plainText.substring(0, 100) + '...' : plainText;
  }

  getImageUrl(imagePath: string): string {
    return this.blogService.getImageUrl(imagePath);
  }

  onImageError(event: any): void {
    event.target.style.display = 'none';
  }

  toggleLike(blog: Blog): void {
    const user = this.authService.getCurrentUser();
    if (!user) return;

    if (blog.likedByUser) {
      this.likeService.unlikeBlog(blog.id!, user.userId).subscribe(() => {
        blog.likedByUser = false;
        if (blog.likesCount && blog.likesCount > 0) {
          blog.likesCount--;
        }
      });
    } else {
      this.likeService.likeBlog(blog.id!, user.userId).subscribe(() => {
        blog.likedByUser = true;
        if (blog.likesCount) {
          blog.likesCount++;
        } else {
          blog.likesCount = 1;
        }
      });
    }
  }

  loadLikeState(blog: Blog): void {
    const user = this.authService.getCurrentUser();
    if (!user) return;

    this.likeService.hasUserLiked(blog.id!, user.userId).subscribe(res => {
      blog.likedByUser = res.liked;
    });

    this.likeService.countLikes(blog.id!).subscribe(res => {
      blog.likesCount = res.count;
    });
  }

  canViewComments(blog: Blog): boolean {
    if (this.isOwnBlog(blog)) {
      return true;
    }
    return (blog as any).canViewComments === true;
  }

  isOwnBlog(blog: Blog): boolean {
    return this.currentUser && 
           blog.creatorID === Number(this.currentUser.userId);
  }

  canEditBlog(blog: Blog): boolean {
    return this.currentUser && 
           this.currentUser.userId && 
           blog.creatorID === Number(this.currentUser.userId);
  }

  editBlog(blogId: string): void {
    this.router.navigate(['/blogs/edit', blogId]);
  }

  getFollowButtonClass(blog: Blog): string {
    const isFollowing = (blog as any).isFollowing;
    return isFollowing ? 'following-btn' : 'follow-btn';
  }

  getFollowIcon(blog: Blog): string {
    const isFollowing = (blog as any).isFollowing;
    return isFollowing ? 'person_remove' : 'person_add';
  }

  getFollowText(blog: Blog): string {
    const isFollowing = (blog as any).isFollowing;
    return isFollowing ? 'Unfollow' : 'Follow';
  }
}