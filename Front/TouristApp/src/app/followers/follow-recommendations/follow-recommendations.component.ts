import { Component, OnInit } from '@angular/core';
import { FollowerService, FollowRecommendation, User } from '../follower.service';
import { AuthService } from 'src/app/auth/auth.service';
import { StakeholdersService } from 'src/app/stakeholders/stakeholders.service';
import { MatSnackBar } from '@angular/material/snack-bar';
import { Router } from '@angular/router';

@Component({
  selector: 'app-follow-recommendations',
  templateUrl: './follow-recommendations.component.html',
  styleUrls: ['./follow-recommendations.component.css']
})
export class FollowRecommendationsComponent implements OnInit {
  recommendations: FollowRecommendation[] = [];
  loading = true;
  currentUserId = 0;
  followingInProgress = new Set<number>(); // Track ongoing follow operations by userId

  constructor(
    private followerService: FollowerService,
    private authService: AuthService,
    private stakeholdersService: StakeholdersService,
    private snackBar: MatSnackBar,
    private router: Router
  ) {}

  ngOnInit(): void {
    const user = this.authService.getCurrentUser();
    if (!user) {
      this.snackBar.open("You must be logged in", "Close", {duration: 3000});
      this.router.navigate(['/login']);
      return;
    }

    this.currentUserId = Number(user.userId);
    this.loadRecommendations();
  }

  loadRecommendations(): void {
    this.followerService.getFollowRecommendations(this.currentUserId, 20).subscribe({
      next: (response) => {
        this.recommendations = response.recommendations;
        this.loading = false;
      },
      error: (err) => {
        console.error('Error loading recommendations:', err);
        this.snackBar.open("Error loading recommendations", "Close", {duration: 3000});
        this.loading = false;
      }
    });
  }

  followUser(recommendation: FollowRecommendation): void {
    const userId = recommendation.recommendedUser.userId;
    
    if (this.followingInProgress.has(userId)) {
      return; // Prevent multiple simultaneous requests
    }

    this.followingInProgress.add(userId);

    this.followerService.followUser(this.currentUserId, userId).subscribe({
      next: () => {
        const username = recommendation.recommendedUser.username || `User ${userId}`;
        this.snackBar.open(`Now following ${username}`, "Close", {duration: 3000});
        // Remove from recommendations list
        this.recommendations = this.recommendations.filter(
          rec => rec.recommendedUser.userId !== userId
        );
        this.followingInProgress.delete(userId);
      },
      error: (err) => {
        console.error('Error following user:', err);
        const username = recommendation.recommendedUser.username || `User ${userId}`;
        this.snackBar.open(`Failed to follow ${username}`, "Close", {duration: 3000});
        this.followingInProgress.delete(userId);
      }
    });
  }

  isFollowingInProgress(userId: number): boolean {
    return this.followingInProgress.has(userId);
  }

  getMutualFollowersText(mutualFollowers: string[]): string {
    if (mutualFollowers.length === 0) {
      return '';
    }
    
    if (mutualFollowers.length === 1) {
      return `Prati ih ${mutualFollowers[0]}`;
    }
    
    if (mutualFollowers.length === 2) {
      return `Prate ih ${mutualFollowers[0]} i ${mutualFollowers[1]}`;
    }
    
    return `Prate ih ${mutualFollowers.slice(0, 2).join(', ')} i još ${mutualFollowers.length - 2} osoba`;
  }

  goToUserProfile(user: User): void {
    // Navigate to user profile - adjust route as needed
    this.router.navigate(['/profile', user.userId]);
  }

  refreshRecommendations(): void {
    this.loading = true;
    this.loadRecommendations();
  }

  getUserDisplayName(user: User): string {
    return user.username || user.name || `User ${user.userId}`;
  }
}